package client

import (
	"crypto/tls"
	"net/http"
	"time"
)

func ClientSet(num int) []*http.Client {
	httpClientSet := make([]*http.Client, num)
	for i := 0; i < num; i++ {
		httpClientSet[i] = &http.Client{
			Timeout: time.Minute * 20,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		}
	}
	return httpClientSet
}
func ClientSetShort(num int) []*http.Client {
	httpClientSet := make([]*http.Client, num)
	for i := 0; i < num; i++ {
		httpClientSet[i] = &http.Client{
			Timeout: time.Minute * 20,
			Transport: &http.Transport{
				TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
				DisableKeepAlives: true,
			},
		}
	}
	return httpClientSet
}
