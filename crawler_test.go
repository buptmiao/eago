package eago

import (
	"testing"
	"time"
)

func TestCrawler(t *testing.T) {
	LoadTestConfig()
	c := NewCrawler(Configs.CrawlerName, Configs.Urls, Configs.Depth, Configs.InSite, Configs.TimeOut, Configs.TTL, Configs.Retry)
	req := c.Register("https://www.github.com", "GET", "parse_test", nil)
	c.AddRequest(req)
	reqpop := <-c.req
	AssertEqual(*reqpop[0] == *req)

	c.Start()
	time.Sleep(time.Millisecond*200)

	AssertEqual(c.extract.Status() == RUNNING)
	AssertEqual(c.fetch.Status() == RUNNING)
	AssertEqual(c.report.Status() == RUNNING)

	c.Stop()
	time.Sleep(time.Millisecond*200)

	AssertEqual(c.extract.Status() == STOP)
	AssertEqual(c.fetch.Status() == STOP)
	AssertEqual(c.report.Status() == STOP)

	c.Restart()
	time.Sleep(time.Millisecond*200)

	AssertEqual(c.extract.Status() == RUNNING)
	AssertEqual(c.fetch.Status() == RUNNING)
	AssertEqual(c.report.Status() == RUNNING)
}