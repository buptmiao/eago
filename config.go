package eago

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
)

type config struct {

	HttpPort uint16
	// The Crawler Name
	CrawlerName string
	Urls        []string
	Depth       int32
	InSite      bool
	TimeOut     int32
	TTL         int32
	Retry       int32

	Auth           bool
	UserName       string
	Token          string
	//Cluster name identifies your cluster for auto-discovery. If you're running
	//multiple clusters on the same network, make sure you're using unique names.
	//
	ClusterName string
	Local       *NodeInfo
	NodeList    []*NodeInfo
	Redis		map[string]*RedisInstance
}

var Configs = new(config)

func LoadConfig() {
	LogInit()
	var configFile string

	// todo the default config Path
	flag.StringVar(&configFile, "c", "", "the config file path")
	flag.Parse()
	Bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(fmt.Sprintf("reading config file error %s: %v", configFile, err))
	}
	if _, err := toml.Decode(string(Bytes), Configs); err != nil {
		panic(fmt.Sprintf("parse config file error %s: %v", configFile, err))
	}

	ArbitrateConfigs(Configs)
}

func ArbitrateConfigs(c *config) {
	// check the ClusterName, ClusterName is used to Identify the clusters in the Local NetWork
	if c.ClusterName == "" {
		Error.Println("ClusterName should not be empty! please check you config file!")
		os.Exit(1)
	}
	if c.Local == nil || c.Local.NodeName == "" {
		Error.Println("Local node name should not be empty! please check you config file!")
		os.Exit(1)
	} else {
		if c.Local.IP == "" {
			c.Local.IP = "127.0.0.1"
		}
		if c.Local.Port == 0 {
			c.Local.Port = 12001
		}
	}
	if len(c.NodeList) == 0 {
		// If user did not set the NodeList fields, make it a slice
		// with local node in it
		c.NodeList = []*NodeInfo{c.Local}
	}
	if c.HttpPort == 0 {
		c.HttpPort = 12002
	}
	if c.CrawlerName == "" {
		c.CrawlerName = "crawler"
	}
	if c.Depth == 0 {
		c.Depth = 1
	}
	if c.TimeOut == 0 {
		c.TimeOut = 3
	}
	if c.Retry == 0 {
		c.Retry = 2
	}
	Log.Println("Load config file success!")
}
