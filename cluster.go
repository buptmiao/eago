package eago

import (
	"consistent"
	"time"
	"sync"
)

type Cluster struct {
	ClusterName string
	Local *Node
	Master *NodeInfo
	Nodes []*NodeInfo
	dis *Distributor
	hash *consistent.Consistent
}

// Constructor of Cluster, init the Cluster, create a new Distributor,
// and set cluster pointer to Node.

var OneCluster sync.Once
var DefaultCluster *Cluster

func GetClusterInstance() *Cluster{
	OneCluster.Do(NewCluster)
	return DefaultCluster
}

func NewCluster(){
	DefaultCluster = &Cluster{
		ClusterName: Configs.ClusterName,
		Local : GetNodeInstance(),
		dis : NewDistributor(),
		hash : consistent.New(),
	}
}

func (c *Cluster)PushRequest(req *UrlRequest) {
	c.dis.Requests.push(req)
}

func (c *Cluster)AddNode(node *NodeInfo) {
	c.Nodes = append(c.Nodes, node)
	c.hash.Add(node.NodeName)
	Stat.AddNode(node)
}

// GetNode will return a nodename from the nodelist by hash the url.
func (c *Cluster)GetNode(url string) string{
	res, err := c.hash.Get(url)
	if err != nil {
		Error.Println(err)
	}
	return res
}

// scan nodeList, call Join Rpc Method, if returns error, the remote
// is not the master, or set master to that node. if all the node list
// are not the Master, make itself Master.
func (c *Cluster)Discover() {
	var exist bool
	for _, nodeInfo := range Configs.NodeList {

		if nodeInfo.NodeName == c.Local.Info.NodeName {
			continue
		}

		err := c.Local.rpc.Join(c.Local.Info, nodeInfo)
		Log.Println("join done")
		// found the master
		if err == nil {
			exist = true
			c.Master = nodeInfo
			break
		} else {
			Log.Println("Join failed: ", err)
		}
	}
	if !exist {
		// make itself Master
		c.BecomeMaster()
	}
}

func (c *Cluster)BecomeMaster() {
	c.Master = c.Local.Info
	Log.Println("The Master is ", *c.Master)
	Stat.SetClusterName(c.ClusterName).SetMaster(c.Master).SetCrawlerName(Configs.CrawlerName)
	c.StartKeeper()
	c.StartDistributor()
}
// check the node
func (c *Cluster)IsMember(node *NodeInfo) bool {
	for i, _ := range c.Nodes {
		// node info is equal
		if *node == *c.Nodes[i] {
			return true
		}
	}
	return false
}

func (c *Cluster)StartDistributor() {
	go c.dis.Run()
}

// Master must detect the slavers, if a slaver is down
// remove it. the func is only invoked by master.
func (c *Cluster)StartKeeper() {
	go func(){
		for {
			for _, node := range c.Nodes {
				if err := c.Local.rpc.KeepAlive(node); err != nil {
					Error.Println("keepalive failed", err)
					return
				}
			}
			time.Sleep(time.Second * 5)
		}
	}()
}
