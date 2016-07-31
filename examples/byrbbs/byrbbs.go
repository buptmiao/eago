package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/buptmiao/eago"
)

func ByrBBSCrawler() *eago.Crawler {
	userid := ""
	passwd := ""
	var cookies string
	crawler := eago.NewCrawler("byrbbs")
	crawler.SetTTL(2).AddSeedUrls(
		"https://bbs.byr.cn/section/ajax_list.json?uid="+userid+"&root=sec-0",
		"https://bbs.byr.cn/section/ajax_list.json?uid="+userid+"&root=sec-1",
		"https://bbs.byr.cn/section/ajax_list.json?uid="+userid+"&root=sec-2",
		"https://bbs.byr.cn/section/ajax_list.json?uid="+userid+"&root=sec-3",
		"https://bbs.byr.cn/section/ajax_list.json?uid="+userid+"&root=sec-4",
		"https://bbs.byr.cn/section/ajax_list.json?uid="+userid+"&root=sec-5",
		"https://bbs.byr.cn/section/ajax_list.json?uid="+userid+"&root=sec-6",
		"https://bbs.byr.cn/section/ajax_list.json?uid="+userid+"&root=sec-7",
		"https://bbs.byr.cn/section/ajax_list.json?uid="+userid+"&root=sec-8",
		"https://bbs.byr.cn/section/ajax_list.json?uid="+userid+"&root=sec-9",
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
		req := crawler.PostRequest("https://bbs.byr.cn/user/ajax_login.json", "login", 1)
		req.SetParams(params.Encode())
		req.SetHeader(header)
		return []*eago.UrlRequest{req}
	})
	crawler.AddParser("login", func(resp *eago.UrlResponse) (urls []*eago.UrlRequest) {
		res := make([]*eago.UrlRequest, 0, 10)
		for _, v := range resp.Resp.Header["Set-Cookie"] {
			cookies += strings.Split(v, ";")[0] + ";"
		}
		cookies = strings.TrimRight(cookies, ";")
		for _, v := range crawler.SeedUrls {
			req := crawler.Request(v, "board", 1)
			header := http.Header{}
			header.Set("Cookie", cookies)
			header.Set("X-Requested-With", "XMLHttpRequest")
			req.SetHeader(header)
			res = append(res, req)
		}
		return res
	})
	crawler.AddParser("board", func(resp *eago.UrlResponse) (urls []*eago.UrlRequest) {
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
		for _, v := range resp.Resp.Header["Set-Cookie"] {
			cookies += strings.Split(v, ";")[0] + ";"
		}
		cookies = strings.TrimRight(cookies, ";")
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
				req := crawler.Request(url, "board", 1)
				header := http.Header{}
				header.Set("Cookie", cookies)
				header.Set("X-Requested-With", "XMLHttpRequest")
				req.SetHeader(header)
				res = append(res, req)
			} else {
				url := "https://bbs.byr.cn" + suburl + "?_uid=" + userid
				req := crawler.Request(url, "artcles", 1)
				header := http.Header{}
				header.Set("Cookie", cookies)
				header.Set("X-Requested-With", "XMLHttpRequest")
				req.SetHeader(header)
				res = append(res, req)
			}
			log.Println(suburl, title)
		}
		return res
	})
	crawler.AddParser("artcles", func(resp *eago.UrlResponse) (urls []*eago.UrlRequest) {
		log.Println("Url:", resp.Src.Url)
		unicode, _ := eago.GBKtoUTF(resp.Body)
		log.Println("Body:", unicode)
		log.Println("Content-Length:", len(resp.Body))
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(unicode))
		type Article struct {
			Title    string
			Url      string
			Author   string
			UpTime   string
			Comments int
		}
		for _, v := range resp.Resp.Header["Set-Cookie"] {
			cookies += strings.Split(v, ";")[0] + ";"
		}
		cookies = strings.TrimRight(cookies, ";")
		articles := []*Article{}
		res := make([]*eago.UrlRequest, 0, 10)
		doc.Find("div.b-content table tbody tr").Each(func(i int, sel *goquery.Selection) {
			article := new(Article)
			article.Title = sel.Find("td.title_9 > a").Text()
			article.Url, _ = sel.Find("td.title_9 > a").Attr("href")
			article.Author = sel.Find("td.title_12 > a").Text()
			article.UpTime = sel.Find("td.title_10").Text()
			article.Comments, _ = strconv.Atoi(sel.Find("td.title_11").Text())
			articles = append(articles, article)
			req := crawler.Request("https://bbs.byr.cn"+article.Url, "text", 1)
			header := http.Header{}
			header.Set("Cookie", cookies)
			header.Set("X-Requested-With", "XMLHttpRequest")
			req.SetHeader(header)
			res = append(res, req)
		})
		pages := doc.Find("ul.pagination li ol li.page-normal > a")
		if count := pages.Last().Text(); count == ">>" {
			next_page_url, _ := pages.Last().Attr("href")
			next_page_url = "https://bbs.byr.cn" + next_page_url
			req := crawler.Request(next_page_url, "artcles", 1)
			header := http.Header{}
			header.Set("Cookie", cookies)
			header.Set("X-Requested-With", "XMLHttpRequest")
			req.SetHeader(header)
			res = append(res, req)
		}
		return res
	})
	crawler.AddParser("text", func(resp *eago.UrlResponse) (urls []*eago.UrlRequest) {
		log.Println("Url:", resp.Src.Url)
		unicode, _ := eago.GBKtoUTF(resp.Body)
		log.Println("Body:", unicode)
		log.Println("Content-Length:", len(resp.Body))
		doc, _ := goquery.NewDocumentFromReader(strings.NewReader(unicode))
		type Text struct {
			Url       string
			AuthorID  string
			Timestamp string
			Content   string
		}
		for _, v := range resp.Resp.Header["Set-Cookie"] {
			cookies += strings.Split(v, ";")[0] + ";"
		}
		cookies = strings.TrimRight(cookies, ";")
		texts := []*Text{}
		res := make([]*eago.UrlRequest, 0, 10)
		reg := regexp.MustCompile(`发信站:[^\(]+\(([^\)]+)\).+?站内([^※]+)`)
		doc.Find("div.b-content table.article").Each(func(i int, sel *goquery.Selection) {
			text := new(Text)
			text.Url = resp.Src.Url
			text.AuthorID = sel.Find("td.a-left a").Text()
			content := sel.Find("div.a-content-wrap").Text()
			match := reg.FindAllStringSubmatch(content, 1)
			if len(match) == 0 {
				log.Println("******************************************* Error Url:", resp.Src.Url, "\n", content)
				return
			}
			text.Timestamp = match[0][1]
			text.Content = match[0][2]
			texts = append(texts, text)
		})
		pages := doc.Find("div.t-pre ul.pagination li ol li.page-normal > a")
		if count := pages.Last().Text(); count == ">>" {
			next_page_url, _ := pages.Last().Attr("href")
			next_page_url = "https://bbs.byr.cn" + next_page_url
			req := crawler.Request(next_page_url, "text", 1)
			header := http.Header{}
			header.Set("Cookie", cookies)
			header.Set("X-Requested-With", "XMLHttpRequest")
			req.SetHeader(header)
			res = append(res, req)
		}
		return res
	})
	return crawler
}

// go run byrbbs.go
// launch browser
// access http://hostname:12002/?start
// then the crawler is running
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
