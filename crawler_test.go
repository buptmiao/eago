package eago

import (
	"testing"
)

func TestCrawler(t *testing.T) {
	c := NewCrawler("test")
	c.StartWith(func() []*UrlRequest {
		req := NewUrlRequest("url", "GET", "crawlertest", "parsertest", 0)
		return []*UrlRequest{req}
	})
	c.AddParser("parsertest", func(resp *UrlResponse) (urls []*UrlRequest) {
		return []*UrlRequest{nil}
	})
	p := c.GetParser("parsertest")
	if p(nil)[0] != nil {
		panic("not nil")
	}
}
