package eago

import (
	"fmt"
	"rpc"
)

//
//
//
type RpcClient struct {
	// client based on rpc
	clients map[string]*rpc.Client
}

func NewRpcClient() *RpcClient {
	res := &RpcClient{
		clients: make(map[string]*rpc.Client),
	}
	return res
}

func (r *RpcClient) AddClient(node *NodeInfo) {
	r.clients[node.NodeName] = rpc.NewClient("tcp", fmt.Sprintf("%s:%d", node.IP, node.Port), 1)
}

func (r *RpcClient) RemClient(node *NodeInfo) {
	err := r.clients[node.NodeName].Close()
	if err != nil {
		Error.Println(err)
	}
	delete(r.clients, node.NodeName)
}

// Invoker should send local NodeInfo to the remote
func (r *RpcClient) Join(local, node *NodeInfo) error {
	Log.Println(local.NodeName, "want to join ", *node)
	var call func(*NodeInfo) error
	client := rpc.NewClient("tcp", fmt.Sprintf("%s:%d", node.IP, node.Port), 1)
	client.MakeRpc("Join", &call)
	return call(local)
}

// Rpc Method at Client side as Master, to distribute the request to
// the slavers.
func (r *RpcClient) Distribute(req *UrlRequest) error {
	var call func(*UrlRequest) error
	r.clients[req.Node].MakeRpc("Distribute", &call)
	return call(req)
}

// Rpc Method at Client side as slavers, Report the new reuests
// to the master.
func (r *RpcClient) ReportRequest(req *UrlRequest) error {
	var call func(*UrlRequest) error
	// Master is down
	if GetClusterInstance().Master == nil {
		return ErrNoneMaster
	}
	r.clients[GetClusterInstance().Master.NodeName].MakeRpc("ReportRequest", &call)
	return call(req)
}

// Rpc Method at Client side as Master, to detect the slavers' status
func (r *RpcClient) KeepAlive(remote *NodeInfo) error {
	var call func(*NodeInfo) error
	r.clients[remote.NodeName].MakeRpc("KeepAlive", &call)
	return call(remote)
}

// Rpc Method at Client side as Slaver, to sync the statistic info
func (r *RpcClient) SyncStatistic(node *NodeInfo) (*Statistic, error) {
	var call func(*NodeInfo) (*Statistic, error)
	// Master is down
	if GetClusterInstance().Master == nil {
		return nil, ErrNoneMaster
	}
	r.clients[GetClusterInstance().Master.NodeName].MakeRpc("SyncStatistic", &call)
	return call(node)
}
