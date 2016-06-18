package eago

import (
	"errors"
	"fmt"
	"rpc"
)

var (
	ErrNotClusterMember = errors.New("not the cluster member")
	ErrNotMaster        = errors.New("I am not the master, thank you!")
)

type RpcServer struct {
	rpc *rpc.Server
}

func NewRpcServer() *RpcServer {
	node := GetNodeInstance()
	res := &RpcServer{
		rpc: rpc.NewServer("tcp", fmt.Sprintf("%s:%d", node.Info.IP, node.Info.Port)),
	}
	return res
}

func (r *RpcServer) Start() {
	Log.Println("Rpc Server starting...")
	r.Register()
	go r.rpc.Start()
}

func (r *RpcServer) Stop() {
	Log.Println("Rpc Server is exit")
	r.rpc.Stop()
}

// Register all the Rpc Service, they may be invoked by either
// the master or the slaver
func (r *RpcServer) Register() {
	r.rpc.Register("Join", r.Join)
	r.rpc.Register("Distribute", r.Distribute)
	r.rpc.Register("ReportRequest", r.ReportRequest)
	r.rpc.Register("KeepAlive", r.KeepAlive)

	r.RegisterType()
}

func (r *RpcServer) RegisterType() {
	rpc.RegisterType(NodeInfo{})
	rpc.RegisterType(UrlRequest{})
}

// Rpc Method at server side as either Master or slave , if it is
// Master, add the remote Node and return nil, otherwise return error.
func (r *RpcServer) Join(remote *NodeInfo) error {
	if !GetNodeInstance().IsMaster() {
		return ErrNotMaster
	}
	GetClusterInstance().AddNode(remote)
	return nil
}

// Rpc Method at server side as Slave, receive the req distributed
// from master, and add it to the crawler to fetch content.
func (r *RpcServer) Distribute(req *UrlRequest) error {
	// the req is send to the wrong node
	if req.node != GetNodeInstance().Info.NodeName {
		return rpc.RpcError{fmt.Sprintf("please check the node info, Node: %s, req:%s", GetNodeInstance().Info.NodeName, req.node)}
	}
	GetNodeInstance().crawl.AddRequest(req)
	return nil
}

// Rpc Method at server side as Master, receive the
// request from slavers and store them to distribute
func (r *RpcServer) ReportRequest(req *UrlRequest) error {
	if !GetNodeInstance().IsMaster() {
		return rpc.RpcError{"I am not the master, thank you!"}
	}
	GetClusterInstance().PushRequest(req)
	return nil
}

// Rpc Method at server side as slaver, response the
// KeepAlive request.
func (r *RpcServer) KeepAlive(remote *NodeInfo) error {

	return nil
}

// Rpc Method at server side as master, response the statistic
// information to the remote. check node's info
func (r *RpcServer) SyncStatistic(node *NodeInfo) (*Statistic, error) {
	// check the id itself
	if !GetNodeInstance().IsMaster() {
		return nil, ErrNotMaster
	}
	if GetClusterInstance().IsMember(node) {
		return Stat.GetStatistic(), nil
	}
	return nil, ErrNotClusterMember
}
