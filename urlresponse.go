package eago

import (
	"io/ioutil"
	"net/http"
)

type UrlResponse struct {
	Src  *UrlRequest
	Resp *http.Response
	Body string
}

func NewResponse(req *UrlRequest, resp *http.Response) *UrlResponse {
	res := &UrlResponse{
		Src:  req,
		Resp: resp,
	}
	if resp != nil {
		body, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			Error.Println("bad response", req.Url)
			return nil
		}
		// success, set body
		res.Body = string(body)
	} else {
		Error.Println("bad response nil", req.Url)
		return nil
	}
	return res
}
