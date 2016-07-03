package eago

type UrlRequest struct {
	Url       string
	Method    string
	Node      string
	Parser    string
	Insite    bool
	Proxy     string
	CookieJar int
	Depth     int32
	Retry     int32
}

func NewUrlRequest(url, method, parser string, insite bool, depth, retry int32, cookie int) *UrlRequest {
	res := &UrlRequest{
		Url:       url,
		Method:    method,
		Parser:    parser,
		Insite:    insite,
		Depth:     depth,
		Retry:     retry,
		CookieJar: cookie,
	}
	return res
}

func (ur *UrlRequest) Incr() {
	ur.Retry++
}
