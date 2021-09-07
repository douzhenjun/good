package utils

import (
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

var client = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:           20,
		MaxConnsPerHost:        5,
		IdleConnTimeout:        time.Minute,
		TLSHandshakeTimeout:    5 * time.Second,
		ExpectContinueTimeout:  1 * time.Second,
		MaxResponseHeaderBytes: 5 * 1024,
	}}

func SimpleGet(url string) ([]byte, error) {
	response, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(response.Body)
	_ = response.Body.Close()
	return data, err
}
