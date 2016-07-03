package eago

import (
	"io/ioutil"
	"net/http"
)

type UrlResponse struct {
	src    *UrlRequest
	resp   *http.Response
	parser string

	body string
}

func NewResponse(req *UrlRequest, resp *http.Response) *UrlResponse {
	res := &UrlResponse{
		src:    req,
		resp:   resp,
		parser: req.Parser,
	}
	if resp != nil {
		body, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			Error.Println("bad response", req.Url)
			return nil
		}
		// success, set body
		res.body = string(body)
	} else {
		Error.Println("bad response nil", req.Url)
		return nil
	}
	return res
}
