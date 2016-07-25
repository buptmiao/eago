package eago

import (
	"net/http"
	"net/url"
	"sync"
	"time"
)

// Fetcher is an executer doing some kind of job
type Fetcher struct {
	status string

	stop chan struct{}
	// It is a stacked channel from which the Fetcher pop UrlRequests
	// see http://gowithconfidence.tumblr.com/post/31426832143/stacked-channels
	pop RequestChan
	// It is a stacked channel to which the Fetcher push UrlResponses
	push ResponseChan

	// clientMap record the cookie:client pairs to deal with the requests
	// need cookies
	clientMap map[int]*http.Client
	// to protect clientMap
	cookieMu *sync.Mutex
}

func NewFetcher(in chan []*UrlRequest, out chan []*UrlResponse) *Fetcher {
	res := &Fetcher{
		status:    STOP,
		stop:      make(chan struct{}),
		pop:       in,
		push:      out,
		clientMap: make(map[int]*http.Client),
		cookieMu:  new(sync.Mutex),
	}
	return res
}

// It is a dead loop until the stop signal is received.
// every request is handled per 'ttl' seconds
func (f *Fetcher) Run() {
	Log.Println("Fetcher is running...")
	f.status = RUNNING
	for {
		select {
		case <-f.stop:
			Log.Println("the Fetcher is stop!")
			f.stop = nil
			f.status = STOP
			return
		case reqs := <-f.pop:
			for _, req := range reqs {
				crawler := GetNodeInstance().GetCrawler(req.Crawler)
				ttl := time.After(time.Second * time.Duration(crawler.TTL))
				// handle the req in goroutine
				go f.handle(req)
				<-ttl
			}
		}
	}
}

// do the handle in goroutine, and push the response to the responsechan
func (f *Fetcher) handle(req *UrlRequest) {
	// get the crawler of the req, if not found return.
	crawler := GetNodeInstance().GetCrawler(req.Crawler)
	client := f.getClient(req)
	// set proxy for this client
	if req.Proxy != "" {
		url, _ := url.Parse(req.Url)
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(url),
		}
	} else {
		// use the default Transport
		client.Transport = nil
	}
	// Set node name
	req.Node = GetNodeInstance().GetName()
	request, err := req.ToRequest()
	if err != nil {
		Error.Println("create request failed, ", err, req.Url)
		return
	}
	response, err := client.Do(request)
	if err != nil {
		Error.Println("http request error, ", err)
		if req.Retry < crawler.Retry {
			req.Incr()
			f.pop.push(req)
			return
		} else {
			Log.Println("dropped url: ", req.Url)
		}
	}
	if response.StatusCode != 200 {
		Log.Println("status of the response: ", response.StatusCode)
		return
	}
	Stat.AddCrawledCount(req.Crawler)
	resp := NewResponse(req, response)

	// store the resp body
	if crawler.store != nil {
		go crawler.store.Store(resp)
	}
	//Add the url to Redis, to mark as crawled
	redisCli := GetRedisClient().GetClient(req.Url)
	redisCli.SAdd(KeyForCrawlByDay(), req.Url)
	if req.Depth >= crawler.Depth {
		Log.Println("this request is reach the Max depth, so stop creating new requests")
		return
	}
	f.push.push(resp)
}

// this func will get the client by cookie of the request,
// when not found, create one.
func (f *Fetcher) getClient(req *UrlRequest) *http.Client {
	cookie := req.CookieJar
	var client *http.Client
	f.cookieMu.Lock()
	if _, ok := f.clientMap[cookie]; !ok {
		jar := NewJar()
		crawler := GetNodeInstance().GetCrawler(req.Crawler)
		f.clientMap[cookie] = &http.Client{
			Jar:     jar,
			Timeout: time.Second * time.Duration(crawler.Timeout),
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
	defer func() {
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

func (f *Fetcher) Status() string {
	return f.status
}
