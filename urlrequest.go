package eago

type UrlRequest struct {
	url string
	method string
	node string
	parser string
	insite bool
	cookieJar int
	depth int32
	retry int32
}

func NewUrlRequest(url, method, parser string, insite bool, depth, retry int32, cookie int) *UrlRequest{
	res := &UrlRequest{
		url : url,
		method : method,
		parser : parser,
		insite : insite,
		depth : depth,
		retry : retry,
		cookieJar: cookie,
	}
	return res
}

func (ur *UrlRequest) Incr() {
	ur.retry++
}

