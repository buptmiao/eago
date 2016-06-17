package crawler

import "sync"

// NodeInfo contains the basic info of a node
type NodeInfo struct {
	NodeName string
	IP string
	Port uint16
}
// constructor of NodeInfo
func NewNodeInfo(name string, ip string, port uint16) *NodeInfo{
	res := &NodeInfo{
		NodeName: name,
		IP: ip,
		Port:port,
	}
	return res
}

// there is only one Node instance per go process.
//
//
type Node struct {
	Info *NodeInfo
	rpc *RpcClient
	crawl *Crawler
}

// Every Node has a DefaultNode.Singleton
var OneNode sync.Once
var DefaultNode *Node

func GetNodeInstance() *Node{
	OneNode.Do(NewNode)
	return DefaultNode
}

func NewNode(){
	res := &Node{}
	res.Info = Configs.Local
	res.rpc = NewRpcClient()
	res.crawl = NewCrawler(Configs.CrawlerName, Configs.Urls, Configs.Depth, Configs.InSite, Configs.TimeOut, Configs.TTL, Configs.Retry, res)
	DefaultNode = res
}


func (n *Node)IsMaster() bool{
	// If the local node info is equal to master node
	return *n.Info == *GetClusterInstance().Master
}

func (n *Node)GetStatistic() (*Statistic, error){
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