package eago

import (
	"consistent"
	"sync"
	"time"
)

const (
	KeeperPeriod        = time.Second * 5
	HeartBeatInterval   = time.Millisecond
	MonitorMasterPeriod = time.Second * 12
)

var ()

type Cluster struct {
	ClusterName string
	Local       *Node
	Master      *NodeInfo
	//Nodes describe the slavers' status, if true, the slaver is active,
	//otherwise, the slaver is down.
	Nodes      map[*NodeInfo]bool
	dis        *Distributor
	hash       *consistent.Consistent
	stopKeeper chan struct{}
	timer      *time.Timer
}

// Constructor of Cluster, init the Cluster, create a new Distributor,
// and set cluster pointer to Node.

var OneCluster sync.Once
var DefaultCluster *Cluster

func GetClusterInstance() *Cluster {
	OneCluster.Do(NewCluster)
	return DefaultCluster
}

func NewCluster() {
	DefaultCluster = &Cluster{
		ClusterName: Configs.ClusterName,
		Local:       GetNodeInstance(),
		Nodes:       make(map[*NodeInfo]bool),
		dis:         NewDistributor(),
		hash:        consistent.New(),
	}
}

func (c *Cluster) PushRequest(req *UrlRequest) {
	c.dis.Requests.push(req)
}

func (c *Cluster) AddNode(node *NodeInfo) {
	// the slaver is added into Nodes, and its status is active by default.
	Log.Println("A new node is joined: ", node.NodeName, node.IP)
	c.Nodes[node] = true
	c.Local.rpc.AddClient(node)
	c.hash.Add(node.NodeName)
	Stat.AddNode(node)
}

// GetNode will return a nodename from the nodelist by hash the url.
func (c *Cluster) GetNode(url string) string {
	res, err := c.hash.Get(url)
	if err != nil {
		Error.Println(err)
	}
	return res
}

// scan nodeList, call Join Rpc Method, if returns error, the remote
// is not the master, or set master to that node. if all the node list
// are not the Master, make itself Master.
func (c *Cluster) Discover() {
	var exist bool
	for _, nodeInfo := range Configs.NodeList {

		if nodeInfo.NodeName == c.Local.Info.NodeName {
			continue
		}
		err := c.Local.rpc.Join(c.Local.Info, nodeInfo)
		// found the master
		if err == nil {
			Log.Println("Join success, Master is: ", nodeInfo.NodeName, nodeInfo.IP)
			exist = true
			c.Master = nodeInfo
			c.Local.rpc.AddClient(c.Master)
			c.BecomeSlaver()
			break
		} else {
			Log.Println("Join failed: ", err, nodeInfo)
		}
	}
	if !exist {
		// make itself Master
		c.BecomeMaster()
	}
}

// Current node becomes Master, and startup tasks belong to master.
func (c *Cluster) BecomeMaster() {
	c.Master = c.Local.Info
	Log.Println("The Master is ", *c.Master)
	Stat.SetClusterName(c.ClusterName).SetMaster(c.Master)
	c.StartKeeper()
	c.StartDistributor()
}

// check the node, if the node has joined in the cluster, return true
func (c *Cluster) IsMember(node *NodeInfo) bool {
	for k, _ := range c.Nodes {
		// node info is equal
		if *node == *k {
			return true
		}
	}
	return false
}

func (c *Cluster) StartDistributor() {
	go c.dis.Run()
}

func (c *Cluster) StopDistributor() {
	c.dis.Stop()
}

func (c *Cluster) RestartDistributor() {
	c.dis.Restart()
}

// Master must detect the slavers, if a slaver is down
// remove it. the func is only invoked by master.
func (c *Cluster) StartKeeper() {
	Log.Println("Keeper is running...")
	c.stopKeeper = make(chan struct{})
	go func() {
		for {
			select {
			case <-c.stopKeeper:
				Log.Println("Keeper is stopped")
				return
			default:
				for node, ok := range c.Nodes {
					if err := c.Local.rpc.KeepAlive(node); err != nil {
						Log.Println("A slaver is down: ", node.NodeName, node.IP)
						// if keep alive failed, don't distribute urls to it any more
						// but keep the rpc client for it, when keep alive success next
						// time, recover it.
						c.UpdateSlaverStatus(node, false)
					} else {
						if !ok {
							Log.Println("A slaver is recovered: ", node.NodeName, node.IP)
							c.UpdateSlaverStatus(node, true)
						}
					}
					//
					time.Sleep(HeartBeatInterval)
				}
				// Keep alive every 5 seconds at least
				time.Sleep(KeeperPeriod - HeartBeatInterval*time.Duration(len(c.Nodes)))
			}
		}
	}()
}

// stop the keeper by closing the chan
func (c *Cluster) StopKeeper() {
	close(c.stopKeeper)
}

func (c *Cluster) UpdateSlaverStatus(node *NodeInfo, v bool) {
	if v {
		c.hash.Add(node.NodeName)
	} else {
		c.hash.Remove(node.NodeName)
	}
	c.Nodes[node] = v
	Stat.UpdateNodeAlive(node, v)
}

/////////////////////////////////////////////////////////////////////
// Functions below are for slavers
/////////////////////////////////////////////////////////////////////
func (c *Cluster) BecomeSlaver() {
	if c.timer == nil {
		c.timer = time.NewTimer(MonitorMasterPeriod)
	} else {
		c.ResetTimer()
	}
	go c.MonitorMaster()
	GetNodeInstance().Restart()
}

func (c *Cluster) ResetTimer() {
	c.timer.Reset(MonitorMasterPeriod)
}

// MonitorMaster check the heart beat package from master, If there is no
// HB package for 12 seconds, stop the world and discover the new master.
// when a new master is selected, restart the world.
func (c *Cluster) MonitorMaster() {
	// slavers block here
	<-c.timer.C
	// Delete Master
	c.Local.rpc.RemClient(c.Master)
	// discover new master
	c.StopTheWorld()
	c.Master = nil
	c.Discover()
}

func (c *Cluster) StopTheWorld() {
	if GetNodeInstance().IsMaster() {
		c.StopDistributor()
		c.StopKeeper()
	}
	GetNodeInstance().Stop()
}
