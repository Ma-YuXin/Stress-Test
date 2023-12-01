package client

import (
	"crypto/tls"
	"net/http"
)

func ClientSet(num int) []*http.Client {
	httpClientSet := make([]*http.Client, num)
	for i := 0; i < num; i++ {
		httpClientSet[i] = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}
	return httpClientSet
}
