package eago

import (
	"eago/rpc"
	"fmt"
)
//
//
//
type RpcClient struct{
	// client based on rpc
	clients map[string]*rpc.Client
}

func NewRpcClient() *RpcClient {
	res := &RpcClient{
		clients: make(map[string]*rpc.Client),
	}
	return res
}

// Invoker should send local NodeInfo to the remote
func (r *RpcClient)Join(local, node *NodeInfo) error {
	Log.Println(local.NodeName,"want to join ", *node)
	var call func(*NodeInfo) (error)
	client := rpc.NewClient("tcp", fmt.Sprintf("%s:%d", node.IP, node.Port), 0)
	client.MakeRpc("Join", &call)
	return call(local)
}

// RpcClient register RPC methods, about the rpc framework
// see https://github.com/siddontang/golib/tree/master/rpc
//func (r *RpcClient)RegisterForMaster(slavers []*NodeInfo) {
//	for _, v := range slavers {
//		client := rpc.NewClient("tcp", fmt.Sprintf("%s:%d", v.IP, v.Port), 3)
//		client.MakeRpc("Distribute", r.Distribute)
//		client.MakeRpc("KeepAlive", r.KeepAlive)
//		r.clients[v.NodeName] = client
//	}
//}
//func (r *RpcClient)RegisterForSlave(master *NodeInfo) {
//	client := rpc.NewClient("tcp", fmt.Sprintf("%s:%d", master.IP, master.Port), 3)
//	client.MakeRpc("ReportRequest", r.ReportRequest)
//}


// Rpc Method at Client side as Master, to distribute the request to
// the slavers.
func (r *RpcClient)Distribute(req *UrlRequest) error{
	var call func(*UrlRequest) error
	r.clients[req.node].MakeRpc("Distribute", &call)
	return call(req)
}

// Rpc Method at Client side as slavers, Report the new reuests
// to the master.
func (r *RpcClient)ReportRequest(req *UrlRequest) error{
	var call func(*UrlRequest) error
	r.clients[GetClusterInstance().Master.NodeName].MakeRpc("ReportRequest", &call)
	return call(req)
}

// Rpc Method at Client side as Master, to detect the slavers' status
func (r *RpcClient)KeepAlive(remote *NodeInfo) error {
	var call func(*NodeInfo) error
	r.clients[remote.NodeName].MakeRpc("KeepAlive", &call)
	return call(remote)
}

// Rpc Method at Client side as Slaver, to sync the statistic info
func (r *RpcClient)SyncStatistic(node *NodeInfo) (*Statistic, error) {
	var call func() (*Statistic, error)
	r.clients[GetClusterInstance().Master.NodeName].MakeRpc("SyncStatistic", &call)
	return call()
}
