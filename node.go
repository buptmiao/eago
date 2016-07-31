package eago

import "sync"

// NodeInfo contains the basic info of a node
type NodeInfo struct {
	NodeName string
	IP       string
	Port     uint16
}

// constructor of NodeInfo
func NewNodeInfo(name string, ip string, port uint16) *NodeInfo {
	res := &NodeInfo{
		NodeName: name,
		IP:       ip,
		Port:     port,
	}
	return res
}

// there is only one Node instance per go process.
//
//
type Node struct {
	Info  *NodeInfo
	rpc   *RpcClient
	crawl map[string]*Crawler

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
	fetch   *Fetcher
	extract *Extractor
	report  *Reporter
}

// Every Node has a DefaultNode.Singleton
var OneNode sync.Once
var DefaultNode *Node

func GetNodeInstance() *Node {
	OneNode.Do(NewNode)
	return DefaultNode
}

func NewNode() {
	res := &Node{
		Info:   Configs.Local,
		rpc:    NewRpcClient(),
		req:    NewRequestChan(),
		resp:   NewResponseChan(),
		upload: NewRequestChan(),
		crawl:  make(map[string]*Crawler),
	}
	res.fetch = NewFetcher(res.req, res.resp)
	res.extract = NewExtractor(res.resp, res.upload)
	res.report = NewReporter(res.upload)
	//	res.crawl[Configs.CrawlerName] = NewCrawler(Configs.CrawlerName, Configs.Urls, Configs.Depth, Configs.InSite, Configs.TimeOut, Configs.TTL, Configs.Retry)
	DefaultNode = res
}

func (n *Node) GetName() string {
	return n.Info.NodeName
}

func (n *Node) IsMaster() bool {
	// If the local node info is equal to master node

	Log.Println(n.Info, GetClusterInstance().Master)
	return *n.Info == *GetClusterInstance().Master
}

func (n *Node) GetStatistic() (*Statistic, error) {
	// for master
	if n.IsMaster() {
		stat := Stat.GetStatistic()
		return stat, nil
	}
	// for slavers
	stat, err := n.rpc.SyncStatistic(n.Info)
	if err != nil {
		Error.Println(err)
		return nil, err
	}
	return stat, nil
}

func (n *Node) AddCrawler(c *Crawler) {
	n.crawl[c.Name] = c
	Stat.AddCrawlerStatistic(c.Name)
}

func (n *Node) RemCrawler(name string) {
	delete(n.crawl, name)
	Stat.RemCrawlerStatistic(name)
}

func (n *Node) GetCrawler(name string) *Crawler {
	crawler, ok := n.crawl[name]
	if !ok {
		panic("crawler not found:" + name)
		return nil
	}
	return crawler
}

func (n *Node) AddRequest(req *UrlRequest) {
	Log.Println("add request to fetcher: ", req.Url)
	//
	Stat.AddTotalCount(req.Crawler)
	n.req.push(req)
}

func (n *Node) Start() {
	Log.Println("Start the crawler...")
	Stat.BeginNow()

	for _, v := range n.crawl {
		for _, req := range v.start_request() {
			n.AddRequest(req)
		}
	}
	go n.fetch.Run()
	go n.extract.Run()
	go n.report.Run()
}

func (n *Node) Stop() {
	Log.Println("Stop the crawler...")
	Stat.Stop()
	n.fetch.Stop()
	n.extract.Stop()
	n.report.Stop()
}

func (n *Node) Restart() {
	Log.Println("Restart the crawler...")
	Stat.BeginNow()
	n.fetch.Restart()
	n.extract.Restart()
	n.report.Restart()
}
