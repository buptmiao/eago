package eago

import (
	"log"
	"net/http/httputil"

	"net/http"

	"gopkg.in/iconv.v1"
)

func GBKtoUTF(input string) (string, error) {
	cd, err := iconv.Open("utf-8", "gbk")
	defer cd.Close()
	if err != nil {
		log.Println("error", err)
		return "", err
	}
	out := make([]byte, len(input))
	a, _, _ := cd.Conv([]byte(input), out)

	return string(a), nil
}

// this is to test the http requests
func DumpHttp(req *UrlRequest) {
	req.Node = GetNodeInstance().GetName()
	request, err := req.ToRequest()
	bytes, _ := httputil.DumpRequest(request, true)
	Debug.Println("\n******** Request ********\n", string(bytes))
	if err != nil {
		Error.Println("create request failed, ", err, req.Url)
		return
	}
	response, err := http.DefaultClient.Do(request)
	bytes, _ = httputil.DumpResponse(response, true)
	Debug.Println("\n******** Response ********\n", string(bytes))
}
