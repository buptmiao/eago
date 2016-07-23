package eago

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
	InSite bool
	// the timeout of per request to the target website
	Timeout int32
	// when a request fails, it will retry 'Retry' times.
	Retry int32
	// TTL is the interval of two urls to fetch using by fetch
	TTL int32
	// store defines the store strategy of the response
	store         Storage
	ParserMap     map[string]Parser
	start_request func() []*UrlRequest
	// some extra data
	MetaData map[string]interface{}
}

func NewCrawler(name string, urls []string, depth int32, inSite bool, to int32, ttl int32, retry int32) *Crawler {
	res := &Crawler{
		Name:      name,
		SeedUrls:  urls,
		Depth:     depth,
		InSite:    inSite,
		Retry:     retry,
		Timeout:   to,
		TTL:       ttl,
		ParserMap: make(map[string]Parser),
	}
	return res
}

func (c *Crawler) GetParser(name string) Parser {
	parser, ok := c.ParserMap[name]
	if !ok {
		panic("crawler not found:" + name)
		return nil
	}
	return parser
}

func (c *Crawler) AddParser(name string, p Parser) *Crawler {
	c.ParserMap[name] = p
	return c
}

// To customize the storage strategy.
func (c *Crawler) SetStorage(st Storage) *Crawler {
	c.store = st
	return c
}

func (c *Crawler) StartWith(call func() []*UrlRequest) *Crawler {
	c.start_request = call
	return c
}
