# Eago 
[![Build Status](https://travis-ci.org/buptmiao/eago.svg?branch=master)](https://travis-ci.org/buptmiao/eago)
![License](https://img.shields.io/dub/l/vibe-d.svg)

An easy distributed and restful crawler framework

## Installation

    go get -u github.com/buptmiao/eago
    
## Features

* Eago works like [scrapy](https://github.com/scrapy/scrapy), but it is more lightweight and effective.

* Eago supports RESTful API, through which users can monitor eago's statistic information, add new crawler job, control the crawler and so on.

* Eago can be deployed as clusters. An eago cluster consist of one master and multiple slavers, and the master node is auto-discovered


## QuickStart

You can run eago like this without any crawlers:

```go
import (
	"github.com/buptmiao/eago"
)

func main() {

	eago.LoadConfig()
	node := eago.GetNodeInstance()
	cluster := eago.GetClusterInstance()

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
## Write a crawler

see the [demo](https://github.com/buptmiao/eago/blob/master/examples/byrbbs/byrbbs.go), which crawls the bbs web pages. The crawler implements multiple Parsers that analysis web page and extract key information.

## More updates will come
