package nineone

import (
	"flag"
	"github.com/franela/goreq"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

var proxy = flag.String("proxy", "", "")

func BuildGetReq(uri string, method string, values url.Values) (*goreq.Request, error) {

	req := goreq.Request{
		Method:      method, //"POST",
		Uri:         uri,
		Accept:      "application/json",
		ContentType: "application/json",
		Timeout:     30 * time.Second, //30s
		Proxy:       *proxy,
		QueryString: values,
	}

	return &req, nil
}


var s5Proxy = flag.String("socket5", "", "")

func RequestBySocket5(uri string) ([]byte, error){
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(*s5Proxy)
	}

	httpTransport := &http.Transport{
		Proxy: proxy,
	}

	httpClient := &http.Client{
		Transport: httpTransport,
	}

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	resp, err2 := httpClient.Do(req)
	if err2 != nil {
		return nil, err2
	}
	defer resp.Body.Close()

	body, err3 := ioutil.ReadAll(resp.Body)
	if err3 != nil {
		return nil, err3
	}

	return body, nil
}

