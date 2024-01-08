package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

func newClient(proxy string, timeout int) *http.Client {
	client := &http.Client{}
	// using proxy
	if proxy != "" {
		p, err := url.Parse(proxy)
		if err != nil {
			logger.Fatal("failed to parse proxy string: %v", err)
		}

		client.Transport = &http.Transport{
			Proxy:                 http.ProxyURL(p),
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
			ResponseHeaderTimeout: time.Second * time.Duration(timeout),
		}
	}
	return client
}

func RespContent(req *http.Request) ([]byte, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %v", err)
	}

	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
