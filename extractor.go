package eago

import (
	"net/url"
)

type Parser func(resp *UrlResponse) (urls []*UrlRequest)

type Extractor struct {
	status string
	stop   chan struct{}
	pop    ResponseChan
	push   RequestChan
}

func NewExtractor(in ResponseChan, out RequestChan) *Extractor {
	res := &Extractor{
		status: STOP,
		stop:   make(chan struct{}),
		pop:    in,
		push:   out,
	}
	return res
}

func (e *Extractor) Run() {
	Log.Println("Extractor is running...")
	e.status = RUNNING
	for {
		select {
		case <-e.stop:
			Log.Println("the Extractor is stop!")
			e.status = STOP
			e.stop = nil
			return
		case resps := <-e.pop:
			for _, resp := range resps {

				// handle the resp in go routine
				go e.handle(resp)
			}
		}
	}
}

func (e *Extractor) handle(resp *UrlResponse) {
	crawler := GetNodeInstance().GetCrawler(resp.Src.Crawler)
	parser := crawler.GetParser(resp.Src.Parser)
	urls := parser(resp)
	// to filter the urls
	urls = e.filter(resp.Src, urls)
	if len(urls) > 0 {
		e.push.push(urls...)
		Log.Printf("New Urls: %d, from the src %s", len(urls), resp.Src.Url)
	}
}

// this is the filter to filter the unreasonable and illegal urlrequests
func (e *Extractor) filter(req *UrlRequest, urls []*UrlRequest) []*UrlRequest {
	res := []*UrlRequest{}
	src, _ := url.Parse(req.Url)
	crawler := GetNodeInstance().GetCrawler(req.Crawler)
	for _, v := range urls {
		URL, err := url.Parse(v.Url)
		if err != nil {
			Error.Println(err, " bad Url: ", v)
			continue
		}
		if crawler.InSite && src.Host != URL.Host {
			continue
		}
		//remove the url that have been crawled
		client := GetRedisClient().GetClient(v.Url)
		if client.SIsMember(KeyForCrawlByDay(), v.Url).Val() {
			continue
		}
		res = append(res, v)
	}
	return res
}

func (e *Extractor) Stop() {
	defer func() {
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

func (e *Extractor) Status() string {
	return e.status
}
