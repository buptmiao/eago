package eago

type UrlRequest struct {
	Url    string
	Method string
	// Params include some key-value pairs, URL_Encode
	Params string
	Proxy  string

	Node      string
	CookieJar int

	Crawler string
	Parser  string

	Depth int32
	Retry int32
}

func NewUrlRequest(url, method, crawler, parser, proxy, params string, depth, retry int32, cookie int) *UrlRequest {
	res := &UrlRequest{
		Url:       url,
		Method:    method,
		Crawler:   crawler,
		Parser:    parser,
		Depth:     depth,
		Retry:     retry,
		Proxy:     proxy,
		Params:    params,
		CookieJar: cookie,
	}
	return res
}

func (u *UrlRequest) Incr() {
	u.Retry++
}
