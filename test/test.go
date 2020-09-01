package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	HttpProxy  = "http://127.0.0.1:1087"
	SocksProxy = "socks5://127.0.0.1:1087"
)

func main() {
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(HttpProxy)
	}

	httpTransport := &http.Transport{
		Proxy: proxy,
	}

	httpClient := &http.Client{
		Transport: httpTransport,
	}

	req, err := http.NewRequest("GET", "https://api.ip.sb/ip", nil)
	if err != nil {
		// handle error
	}

	resp, err := httpClient.Do(req)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))

	proxy = func(_ *http.Request) (*url.URL, error) {
		return url.Parse(SocksProxy)
	}

	httpTransport = &http.Transport{
		Proxy: proxy,
	}

	httpClient = &http.Client{
		Transport: httpTransport,
	}

	req, err = http.NewRequest("GET", "https://down.zhaiclub.com/385770.mp4?act=url&key=TJw92fDnLsChYvkX", nil)
	if err != nil {
		// handle error
	}

	resp, err = httpClient.Do(req)
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}
