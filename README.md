# Eago 
[![Build Status](https://travis-ci.org/buptmiao/eago.svg?branch=master)](https://travis-ci.org/buptmiao/eago)
![License](https://img.shields.io/dub/l/vibe-d.svg)

An easy distributed and restful crawler framework

## Installation
Before install eago, you should install these dependencies

    go get github.com/gin-gonic/gin
    go get gopkg.in/redis.v3
    go get github.com/BurntSushi/toml
    go get gopkg.in/iconv.v1

Install:

    go get github.com/buptmiao/eago
    
## Features
* Eago uses [Toml](https://github.com/BurntSushi/toml) to configure the parameters, for details: [config.toml](https://github.com/buptmiao/eago/blob/master/config.toml)

* Make sure redis-server is correctly installed and launched on your system. Eago filters the duplicate urls by Redis, and the urls is sharded with a configurable number of redis shards 

* You can customize the storage strategy in your application by implementing the interface [Storage](https://github.com/buptmiao/eago/blob/master/storage.go)

* Eago supports RESTful API, through which users can monitor eago's statistic information, add new crawler job, control the crawler and so on.

* Eago can be deployed as clusters. An eago cluster consist of one master and multiple slavers, and the master node is auto-discovered

## QuickStart

An Example:

```go
import (
	"github.com/buptmiao/eago"
)

func main() {

	eago.LoadConfig()
	node := eago.GetNodeInstance()
	cluster := eago.GetClusterInstance()
    // store the web page into redis by default.
	store := eago.NewDefaultStore(eago.GetRedisClient())
	node.SetStorage(store)

	eago.NewRpcServer().Start()
	// Descover will Block the execution, until a master node
	// is found, or become master itself.
	cluster.Discover()
	// start the Http Server
	eago.NewHttpServer(node).Serve()
}

```

Run the example:
    
    go run demo.go -c /yourpath/config.toml 

Monitor the eago's status by REST API:

    curl -XGET localhost:12002?pretty

Response:

```
    {
        "ClusterName": "eagles",
        "Running": false,
        "Begin At": "",
        "Elapse": "",
        "NodeNumber": 1,
        "Master": {
            "NodeName": "miao_shard1",
            "IP": "192.168.1.12",
            "Port": 12001
        },
        "slavers": null,
        "CrawlerStatistics": {},
        "Message": "You know, for data"
    }
```

## More updates will come
