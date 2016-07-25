package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/buptmiao/eago"
)

func ByrBBSCrawler() *eago.Crawler {
	userid := ""
	passwd := ""
	crawler := eago.NewCrawler("byrbbs")
	crawler.SetDepth(3).AddSeedUrls(
		//"https://bbs.byr.cn/section/ajax_list.json?uid=" + userid + "&root=sec-0",
		//"https://bbs.byr.cn/section/ajax_list.json?uid="+userid+"&root=sec-1",
		//"https://bbs.byr.cn/section/ajax_list.json?uid=" + userid + "&root=sec-2",
		//"https://bbs.byr.cn/section/ajax_list.json?uid="+userid+"&root=sec-3",
		//"https://bbs.byr.cn/section/ajax_list.json?uid="+userid+"&root=sec-4",
		//"https://bbs.byr.cn/section/ajax_list.json?uid="+userid+"&root=sec-5",
		//"https://bbs.byr.cn/section/ajax_list.json?uid="+userid+"&root=sec-6",
		//"https://bbs.byr.cn/section/ajax_list.json?uid="+userid+"&root=sec-7",
		//"https://bbs.byr.cn/section/ajax_list.json?uid="+userid+"&root=sec-8",
		//"https://bbs.byr.cn/section/ajax_list.json?uid="+userid+"&root=sec-9",
		"https://bbs.byr.cn/board/ACM_ICPC",
	)
	crawler.StartWith(func() []*eago.UrlRequest {
		params := url.Values{}
		params.Set("id", userid)
		params.Set("passwd", passwd)
		//params.Set("mode", "0")
		//params.Set("CookieDate", "0")
		header := http.Header{}
		header.Set("Host", "bbs.byr.cn")
		header.Set("X-Requested-With", "XMLHttpRequest")
		header.Set("Connection", "keep-alive")
		header.Set("content-type", "application/x-www-form-urlencoded; charset=UTF-8")
		req := eago.NewUrlRequest("https://bbs.byr.cn/user/ajax_login.json", "POST", "byrbbs", "login", 1)
		req.SetParams(params.Encode())
		req.SetHeader(header)
		return []*eago.UrlRequest{req}
	})
	crawler.AddParser("login", func(resp *eago.UrlResponse) (urls []*eago.UrlRequest) {
		log.Println("Url:", resp.Src.Url)
		unicode, _ := eago.GBKtoUTF(resp.Body)
		log.Println("Body:", unicode)
		log.Println("Content-Length:", len(resp.Body))
		log.Println("cookies***********", resp.Resp.Header)
		log.Println("cookies***********", resp.Resp.Header["Set-Cookie"])
		res := make([]*eago.UrlRequest, 0, 10)
		for _, v := range crawler.SeedUrls {
			req := eago.NewUrlRequest(v, "GET", "byrbbs", "step2", 1)
			header := http.Header{}
			cookies := ""
			for _, v := range resp.Resp.Header["Set-Cookie"] {
				cookies += strings.Split(v, ";")[0] + ";"
			}
			cookies = strings.TrimRight(cookies, ";")
			log.Println("cookies***********", cookies)
			header.Set("Cookie", cookies)
			header.Set("X-Requested-With", "XMLHttpRequest")

			req.SetDepth(1).SetHeader(header)
			res = append(res, req)
		}
		return res
	})
	crawler.AddParser("step1", func(resp *eago.UrlResponse) (urls []*eago.UrlRequest) {
		log.Println("Url:", resp.Src.Url)
		unicode, _ := eago.GBKtoUTF(resp.Body)
		log.Println("Body:", unicode)
		log.Println("Content-Length:", len(resp.Body))
		type board struct {
			Tag string `json:"t, omitempty"`
			Id  string `json:"id, omitempty"`
		}
		boards := make([]board, 0)
		err := json.Unmarshal([]byte(unicode), &boards)
		if err != nil {
			log.Println(err)
			return nil
		}
		res := make([]*eago.UrlRequest, 0, 10)
		for _, v := range boards {
			log.Println(v.Tag)
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(v.Tag))
			if err != nil {
				continue
			}
			suburl := doc.Find("a").AttrOr("href", "")
			title := doc.Find("a").AttrOr("title", "")
			if v.Id != "" {
				url := "https://bbs.byr.cn/section/ajax_list.json?uid=" + userid + "&root=" + v.Id
				log.Println("url:", url)
				req := eago.NewUrlRequest(url, "GET", "byrbbs", "step1", 1)
				header := http.Header{}
				header.Set("X-Requested-With", "XMLHttpRequest")
				req.SetDepth(2).SetHeader(header)
				res = append(res, req)
			} else {
				url := "https://bbs.byr.cn" + suburl + "?_uid=" + userid
				req := eago.NewUrlRequest(url, "GET", "byrbbs", "step2", 1)
				header := http.Header{}
				header.Set("X-Requested-With", "XMLHttpRequest")
				req.SetDepth(2).SetHeader(header)
				res = append(res, req)
			}
			log.Println(suburl, title)

		}
		return res
	})
	crawler.AddParser("step2", func(resp *eago.UrlResponse) (urls []*eago.UrlRequest) {
		log.Println("Url:", resp.Src.Url)
		unicode, _ := eago.GBKtoUTF(resp.Body)
		log.Println("Body:", unicode)
		log.Println("Content-Length:", len(resp.Body))
		return nil
	})
	store := eago.NewDefaultStore(eago.GetRedisClient())
	crawler.SetStorage(store)
	return crawler
}

func main() {
	eago.LoadConfig()
	node := eago.GetNodeInstance()
	cluster := eago.GetClusterInstance()

	bbs := ByrBBSCrawler()
	node.AddCrawler(bbs)

	eago.NewRpcServer().Start()
	// Discover will Block the execution, until a master node
	// is found, or become master itself.
	cluster.Discover()
	// start the Http Server
	eago.NewHttpServer(node).Serve()
}
