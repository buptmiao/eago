package main

import (
	"github.com/buptmiao/eago"
)

func main() {

	eago.LoadConfig()
	node := eago.GetNodeInstance()
	cluster := eago.GetClusterInstance()

	store := eago.NewDefaultStore(eago.GetRedisClient())
	node.SetStorage(store)

	eago.NewRpcServer().Start()
	// Descover will Block the execution, until a master node
	// is found, or become master itself.
	cluster.Discover()
	// start the Http Server
	eago.NewHttpServer(node).Serve()
}
