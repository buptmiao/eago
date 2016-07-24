package eago

import (
	"net/http"
	"strings"
)

type UrlRequest struct {
	Url    string
	Method string
	// Params include some key-value pairs, URL_Encode
	Headers http.Header
	Params  string
	Proxy   string

	Node      string
	CookieJar int

	Crawler string
	Parser  string

	Depth int32
	Retry int32
}

func NewUrlRequest(url, method, crawler, parser string, cookie int) *UrlRequest {
	res := &UrlRequest{
		Url:       url,
		Method:    method,
		Headers:   make(http.Header),
		Crawler:   crawler,
		Parser:    parser,
		Proxy:     "",
		Params:    "",
		CookieJar: cookie,
	}
	return res
}

func (u *UrlRequest) SetHeader(header http.Header) *UrlRequest {
	u.Headers = header
	return u
}

func (u *UrlRequest) SetDepth(depth int32) *UrlRequest {
	u.Depth = depth
	return u
}

func (u *UrlRequest) SetRetry(retry int32) *UrlRequest {
	u.Retry = retry
	return u
}

func (u *UrlRequest) Incr() *UrlRequest {
	u.Retry++
	return u
}

func (u *UrlRequest) SetProxy(proxy string) *UrlRequest {
	u.Proxy = proxy
	return u
}

func (u *UrlRequest) SetParams(params string) *UrlRequest {
	u.Params = params
	return u
}

func (u *UrlRequest) ToRequest() (*http.Request, error) {
	res, err := http.NewRequest(u.Method, u.Url, strings.NewReader(u.Params))
	if err != nil {
		Error.Println("create request failed, ", err, u.Url)
		return nil, err
	}
	res.Header = u.Headers
	return res, nil
}
