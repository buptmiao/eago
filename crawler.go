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
	// some extra data for http request, such as Header and PostForm
	MetaData map[string]interface{}
}

// NewCrawler return a pointer to a Crawler object.
// by default:
// 		Depth:     1,
//		InSite:    true,
//		Retry:     3,
//		Timeout:   5,
//		TTL:       0,
// these params can be set by methods of the crawler.

func NewCrawler(name string) *Crawler {
	res := &Crawler{
		Name:      name,
		SeedUrls:  make([]string, 1),
		Depth:     1,
		InSite:    true,
		Retry:     3,
		Timeout:   5,
		TTL:       0,
		ParserMap: make(map[string]Parser),
	}
	return res
}

// Set the Urls of the crawler
func (c *Crawler) AddSeedUrls(urls ...string) *Crawler {
	c.SeedUrls = append(c.SeedUrls, urls...)
	return c
}

// Set the depth of the crawler
func (c *Crawler) SetDepth(depth int32) *Crawler {
	c.Depth = depth
	return c
}

// Set the InsSite of the crawler
func (c *Crawler) SetInSite(inSite bool) *Crawler {
	c.InSite = inSite
	return c
}

// Set the Timeout of the crawler
func (c *Crawler) SetTimeout(to int32) *Crawler {
	c.Timeout = to
	return c
}

// Set the TTL of the crawler
func (c *Crawler) SetTTL(ttl int32) *Crawler {
	c.TTL = ttl
	return c
}

// Set the Retry of the crawler
func (c *Crawler) SetRetry(retry int32) *Crawler {
	c.Retry = retry
	return c
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

func (c *Crawler) SetParam(key string, value interface{}) *Crawler {
	c.MetaData[key] = value
	return c
}

func (c *Crawler) GetParam(key string) interface{} {
	res, ok := c.MetaData[key]
	if !ok {
		panic("key not found")
		return nil
	}
	return res
}
