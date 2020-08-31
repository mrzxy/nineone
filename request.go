package nineone

import (
	"flag"
	"github.com/franela/goreq"
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
