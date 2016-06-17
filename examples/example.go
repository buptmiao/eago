package main

import (
	"crawler"
)

func main() {

	crawler.LoadConfig()
	node := crawler.GetNodeInstance()
	cluster := crawler.GetClusterInstance()

	crawler.NewRpcServer().Start()
	// Descover will Block the execution, until a master node
	// is found, or become master itself.
	cluster.Discover()
	// start the Http Server
	crawler.NewHttpServer(node).Serve()
}
