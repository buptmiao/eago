package main

import (
	"github.com/buptmiao/eago"
	"net/url"
	"log"
	"gopkg.in/iconv.v1"
)

func ByrbbsParser(resp *eago.UrlResponse) (urls []*eago.UrlRequest) {
	cd, err := iconv.Open("utf-8", "gbk")
	if err != nil {
		log.Println("error", err)
		return nil
	}
	defer cd.Close()
	out := make([]byte, len(resp.Body))

	a,_,c := cd.Conv([]byte(resp.Body), out)

	log.Println(string(a),"\n",len(a),c)

	return nil
}

func main() {

	eago.LoadConfig()
	node := eago.GetNodeInstance()
	cluster := eago.GetClusterInstance()
	crawler := eago.NewCrawler("byrbbs", nil, 1, true, 5, 0, 2)
	params := url.Values{}
	params.Set("id", "inaadversity")
	params.Set("passwd", "362622365")
	params.Set("mode", "0")
	params.Set("CookieDate", "0")
	store := eago.NewDefaultStore(eago.GetRedisClient())
	crawler.SetStorage(store)

	node.AddCrawler(crawler)
	eago.NewRpcServer().Start()
	// Discover will Block the execution, until a master node
	// is found, or become master itself.
	cluster.Discover()
	// start the Http Server
	eago.NewHttpServer(node).Serve()
}
