package main

import (
	"crypto/tls"
	"net/http"
	"net/url"
)

func newClient(proxy string) *http.Client {
	client := &http.Client{}
	// using proxy
	if proxy != "" {
		p, err := url.Parse(proxy)
		if err != nil {
			logger.Fatal("failed to parse proxy string: %v", err)
		}

		client.Transport = &http.Transport{
			Proxy:           http.ProxyURL(p),
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	return client
}
