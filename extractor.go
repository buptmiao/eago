package crawler

import (
	"net/url"
)

type Parser func(body *string) (urls []string)

type Extractor struct {
	status string
	stop chan struct{}
	pop ResponseChan
	push RequestChan
	ParserMap map[string]Parser
}

func NewExtractor(in ResponseChan, out RequestChan) *Extractor{
	res := &Extractor{
		status :STOP,
		pop : in,
		push : out,
		ParserMap: make(map[string]Parser),
	}
	return res
}

func (e *Extractor)Run() {
	Log.Println("Extractor is running...")
	for{
		select {
		case <-e.stop:
			Log.Println("the Extractor is stop!")
			e.stop = nil
			return
		case resps := <-e.pop:
			for _, resp := range resps {

				// handle the resp in goroutine
				go e.handle(resp)
			}
		}
	}
}

func (e *Extractor)handle(resp *UrlResponse) {
	if _, ok := e.ParserMap[resp.parser]; !ok {
		Error.Printf("the Parse Method is not defined for %s, url: %s", resp.parser, resp.src.url)
		return
	}
	urls := e.ParserMap[resp.parser](&resp.body)
	// to filter the urls
	urls = e.filter(resp.src, urls)
	if urls != nil {
		newRequests := make([]*UrlRequest,0 ,len(urls))
		for _, url := range urls {
			req := NewUrlRequest(url,resp.src.method, resp.parser, resp.src.insite, resp.src.depth+1, 0, resp.src.cookieJar)
			newRequests = append(newRequests, req)
		}
		e.push.push(newRequests...)
		Log.Println("New Urls: %d, from the src %s", len(newRequests), resp.src.url)
	}
}
// this is the filter to filter the unreasonable and illegal urlrequests
func (e *Extractor)filter(req *UrlRequest, urls []string) []string{
	res := []string{}
	srcurl, _ :=  url.Parse(req.url)
	for _, v := range urls {
		URL, err := url.Parse(v)
		if err != nil {
			Error.Println(err, " bad Url: ", v)
			continue
		}
		if req.insite && srcurl.Host != URL.Host {
			continue
		}
		//remove the url that have been crawled
		client := GetRedisClient().GetClient(v)
		if client.SIsMember(KeyForCrawlByDay(), v).Val() {
			continue
		}
		res = append(res, v)
	}
	return res
}

func (e *Extractor)Stop() {
	defer func(){
		if err := recover(); err != nil {
			Error.Println(err)
		}
	}()
	close(e.stop)
}

// this func will restart the Extractor
func (e *Extractor) Restart() {
	e.stop = make(chan struct{})
	go e.Run()
}