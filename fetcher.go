package eago

import (
	"time"
	"net/http"
	"net/http/cookiejar"
	"sync"
)

// Fetcher is an executer doing some kind of job
type Fetcher struct {
	status string
	// timeout defines the max waiting time for per request
	timeout int32
	// ttl is the interval between two requests
	ttl int32
	// this channel can stop the Fetcher
	stop chan struct{}
	depth int32
	retry int32
	// It is a stacked channel from which the Fetcher pop UrlRequests
	// see http://gowithconfidence.tumblr.com/post/31426832143/stacked-channels
	pop  RequestChan
	// It is a stacked channel to which the Fetcher push UrlResponses
	push ResponseChan

	// clientMap record the cookie:client pairs to deal with the requests
	// need cookies
	clientMap map[int]*http.Client
	// to protect clientMap
	cookieMu *sync.Mutex
	// store defines the store strategy of the response
	store Storer
}

func NewFetcher(to, ttl,depth,retry int32, stop chan struct{}, in chan []*UrlRequest, out chan []*UrlResponse) *Fetcher{
	res := &Fetcher{
		status : STOP,
		timeout: to,
		ttl:ttl,
		depth:depth,
		retry:retry,
		stop : make(chan struct{}),
		pop : in,
		push : out,
		clientMap: make(map[int]*http.Client),
		cookieMu: new(sync.Mutex),
	}
	return res
}

// It is a dead loop until the stop signal is received.
// every request is handled per 'ttl' seconds
func (f *Fetcher) Run() {
	Log.Println("Fetcher is running...")
	f.status = RUNNING
	for{
		select {
		case <-f.stop:
			Log.Println("the Fetcher is stop!")
			f.stop = nil
			f.status = STOP
			return
		case reqs := <-f.pop:
			for _, req := range reqs {
				ttl := time.After(time.Second * time.Duration(f.ttl))
				// handle the req in goroutine
				go f.handle(req)
				<-ttl
			}
		}
	}
}
// do the handle in goroutine, and push the response to the responsechan
func (f *Fetcher) handle(req *UrlRequest ) {
	client := f.getClient(req)
	request, err := http.NewRequest(req.method, req.url, nil)
	if err != nil {
		Error.Println("create request failed, ",err, req.url)
		return
	}
	response, err := client.Do(request)
	if err != nil {
		Error.Println("http request error, ",err)
		if req.retry < f.retry {
			req.Incr()
			f.pop.push(req)
			return
		} else {
			Log.Println("dropped url: ", req.url)
		}
	}
	if response.StatusCode != 200 {
		Log.Println("status of the response: ",response.StatusCode)
		return
	}
	resp := NewResponse(req, response)

	// store the resp
	if f.store != nil {
		go f.store.Store(resp)
		//Add the url to Redis, to mark as crawled
		rediscli := GetRedisClient().GetClient(req.url)
		rediscli.SAdd(KeyForCrawlByDay(), req.url)
	}
	if req.depth >= req.depth {
		Log.Println("this request is reach the Max depth, so stop creating new requests")
		return
	}
	f.push.push(resp)
}

// this func will get the client by cookie of the request,
// when not found, create one.
func (f *Fetcher) getClient(req *UrlRequest)  *http.Client{
	cookie := req.cookieJar
	var client *http.Client

	f.cookieMu.Lock()
	if _, ok := f.clientMap[cookie]; !ok {
		jar, _ := cookiejar.New(nil)
		f.clientMap[cookie] = &http.Client{
			Jar: jar,
			Timeout: time.Second * time.Duration(f.timeout),
		}
	}
	client = f.clientMap[cookie]
	f.cookieMu.Unlock()
	return client
}
func (f *Fetcher) Add(req *UrlRequest) {

}

// stop the Fetcher, in fact, it just send the STOP signal to
// the Fetcher itself, it is invoked by the up-level in general
func (f *Fetcher) Stop() {
	defer func(){
		if err := recover(); err != nil {
			Error.Println(err)
		}
	}()
	close(f.stop)
}

// this func will restart the Fetcher
func (f *Fetcher) Restart() {
	f.stop = make(chan struct{})
	go f.Run()
}
