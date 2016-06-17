package crawler


// Crawler implements the main work of the node.
// It defines some primitive info.
// If the current node is slave, a node will manage three entities
// fetcher, extractor and reporter. else, if the current node is master
// a distributor is appended.
type Crawler struct {
	// Name stand for a unique identifier of the crawler.
	Name string
	// StartURL is the entrance of website to crawl.
	SeedUrls []string
	//Depth defines the website depth to crawl
	Depth int32
	// If the field is set true, the crawler will only crawl the pages
	// that are of same host address.
	Insite bool
	// the timeout of per request to the target website
	Timeout int32
	// when a request fails, it will retry 'Retry' times.
	Retry int32
	// TTL is the interval of two urls to fetch using by fetch
	TTL int32
	// req is the UrlRequest buffer for current node to fetch the
	// content with minimal blocking
	req RequestChan
	// resp is the UrlResponse buffer for current node to extract
	// the new urls with minimal blocking
	resp ResponseChan
	// upload is the UrlRequest buffer for current node to report
	// the urls to cluster's mster node.
	upload RequestChan
	// three kinds of workers
	fetch *Fetcher
	extract *Extractor
	report *Reporter
}


func NewCrawler(name string, urls []string ,depth int32, insite bool, to int32, ttl int32, retry int32, node *Node) *Crawler{
	res := &Crawler{
		Name: name,
		SeedUrls: urls,
		Depth: depth,
		Retry: retry,
		Timeout:to,
		TTL:ttl,
		req : NewRequestChan(),
		resp : NewResponseChan(),
		upload : NewRequestChan(),
	}
	res.fetch = NewFetcher(res.Timeout, res.TTL, res.Depth, res.Retry, make(chan struct{}), res.req, res.resp)
	res.extract = NewExtractor(res.resp, res.upload)
	res.report = NewReporter(res.upload)
	node.crawl = res
	return res
}


func (c *Crawler)Register(url string, method string, parsename string, p Parser) *UrlRequest{
	c.extract.ParserMap[parsename] = p
	res := NewUrlRequest(url, method, parsename, c.Insite, 0, 0, 0)
	return  res
}

func (c *Crawler)AddRequest(req *UrlRequest){
	Log.Println("add request to fetcher: ",req.url)
	// 统计信息+1
	Stat.AddTotalCount()
	c.req.push(req)
}

func (c *Crawler)Start(){
	Log.Println("Start the crawler...")
	Stat.BeginNow()
	go c.fetch.Run()
	go c.extract.Run()
	go c.report.Run()

}
func (c *Crawler)Stop(){
	Log.Println("Stop the crawler...")
	Stat.Stop()
	c.fetch.Stop()
	c.extract.Stop()
	c.report.Stop()
}

func (c *Crawler)Restart(){
	Log.Println("Restart the crawler...")
	c.fetch.Restart()
	c.extract.Restart()
	c.report.Restart()
}



