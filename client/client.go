package client

import (
	"crypto/tls"
	"net/http"
	"sync"
	"time"
)

var (
	shortConnPool = sync.Pool{
		New: func() any {
			return &http.Client{
				Timeout: time.Minute * 20,
				Transport: &http.Transport{
					TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
					DisableKeepAlives: true,
				},
			}
		},
	}
	longConnPool = sync.Pool{
		New: func() any {
			return &http.Client{
				Timeout: time.Minute * 20,
				Transport: &http.Transport{
					TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
				},
			}
		},
	}
)

func ClientSet(num int) []*http.Client {
	httpClientSet := make([]*http.Client, num)
	for i := 0; i < num; i++ {
		httpClientSet[i] = GetClient()
	}
	return httpClientSet
}
func PutClientSet(set []*http.Client){
	for i := 0; i < len(set); i++ {
		PutClient(set[i]) 
	}
}
func ClientSetShort(num int) []*http.Client {
	httpClientSet := make([]*http.Client, num)
	for i := 0; i < num; i++ {
		httpClientSet[i] = GetShortClient()
	}
	return httpClientSet
}
func PutClientSetShort(set []*http.Client){
	for i := 0; i < len(set); i++ {
		PutShortClient(set[i]) 
	}
}
func GetShortClient() *http.Client {
	return shortConnPool.Get().(*http.Client)
}
func PutShortClient(client *http.Client) {
	shortConnPool.Put(client)
}
func GetClient() *http.Client {
	return longConnPool.Get().(*http.Client)
}
func PutClient(client *http.Client) {
	longConnPool.Put(client)
}
