package stress

import (
	"bytes"
	"context"
	"log"
	"math/rand"
	"net/http"
	"stressTest/client"
	"stressTest/config"
	"stressTest/util"
	"sync"
	"time"
)

func Post(wg *sync.WaitGroup, ctx context.Context, res string, num int) {
	defer wg.Done()
	tick := time.NewTicker(time.Second / time.Duration(num))
	defer tick.Stop()
	wgr := &sync.WaitGroup{}
	defer wgr.Wait()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			wgr.Add(1)
			go post(wgr, res)
		}
	}
}

func Patch(wg *sync.WaitGroup, ctx context.Context, res string, num int) {
	defer wg.Done()
	tick := time.NewTicker(time.Second / time.Duration(num))
	defer tick.Stop()
	wgr := &sync.WaitGroup{}
	defer wgr.Wait()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			reslist := Resindex[res]
			if len(reslist) == 0 {
				return
			}
			wgr.Add(1)
			go patch(wgr, res, reslist[rand.Intn(len(reslist)>>1)])
		}
	}
}

func Delete(wg *sync.WaitGroup, ctx context.Context, res string, num int) {
	defer wg.Done()
	reslist := Resindex[res]
	pos := len(reslist) >> 1
	tick := time.NewTicker(time.Second / time.Duration(num))
	defer tick.Stop()
	wgr := &sync.WaitGroup{}
	defer wgr.Wait()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			if pos < len(reslist) {
				wgr.Add(1)
				go delete(wgr, res, reslist[pos])
			}
			pos++
		}
	}
}

func Put(wg *sync.WaitGroup, ctx context.Context, res string, num int) {
	defer wg.Done()
	tick := time.NewTicker(time.Second / time.Duration(num))
	defer tick.Stop()
	wgr := &sync.WaitGroup{}
	defer wgr.Wait()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			reslist := Resindex[res]
			wgr.Add(1)
			go put(wgr, res, reslist[rand.Intn(len(reslist)>>1)])
		}
	}
}

func post(wg *sync.WaitGroup, res string) {
	defer wg.Done()
	httpClient := client.GetShortClient()
	defer client.PutShortClient(httpClient)
	data, request := util.GetPostDataAndUrl(res, config.GetDefultNameSpace(), rand.Intn(10), rand.Int(), rand.Int())
	req, err := http.NewRequest("POST", request, bytes.NewBuffer(data))
	log.Println("POST", req.URL.String())
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", config.GetDefultAuthor())
	req.Header.Set("Content-Type", "application/json")
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println("err:", err)
	}
	defer resp.Body.Close()
	log.Println("resp: ", resp.Status, resp.Request.Method, resp.Request.URL)
}

func delete(wg *sync.WaitGroup, res, resName string) {
	defer wg.Done()
	cli := client.GetShortClient()
	defer client.PutShortClient(cli)
	_, _, request := util.GetBasic(res, config.GetDefultNameSpace())
	req, err := http.NewRequest("DELETE", request+"/"+resName, nil)
	log.Println("DELETE", req.URL.String())
	if err != nil {
		log.Fatal("new http request err", err)
	}
	req.Header.Set("Authorization", config.GetDefultAuthor())
	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		log.Panicln("do request has err", err)
	}
	defer resp.Body.Close()
	log.Println("resp: ", resp.Status, resp.Request.Method, resp.Request.URL)
}

func patch(wg *sync.WaitGroup, res, resName string) {
	defer wg.Done()
	cli := client.GetShortClient()
	defer client.PutShortClient(cli)
	_, _, request := util.GetBasic(res, config.GetDefultNameSpace())
	anno := util.GetPatchAnnotations(rand.Intn(10))
	req, err := http.NewRequest("PATCH", request+"/"+resName, bytes.NewBufferString(anno))
	log.Println("PATCH", req.URL.String())
	if err != nil {
		log.Fatal("new http request err", err)
	}
	req.Header.Set("Authorization", config.GetDefultAuthor())
	req.Header.Set("Content-Type", "application/strategic-merge-patch+json")
	resp, err := cli.Do(req)
	if err != nil {
		log.Panicln("do request has err", err)
	}
	defer resp.Body.Close()
	log.Println("resp: ", resp.Status, resp.Request.Method, resp.Request.URL)
}

func put(wg *sync.WaitGroup, res, resName string) {
	defer wg.Done()
	cli := client.GetShortClient()
	defer client.PutShortClient(cli)
	data, request := util.GetPutDataAndUrl(res, config.GetDefultNameSpace(), resName, rand.Intn(10))
	req, err := http.NewRequest("PUT", request+"/"+resName, bytes.NewBuffer(data))
	log.Println("PUT", req.URL.String())
	if err != nil {
		log.Fatal("new http request err", err)
	}
	req.Header.Set("Authorization", config.GetDefultAuthor())
	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		log.Println("do request has err", err)
	}
	defer resp.Body.Close()
	log.Println("resp: ", resp.Status, resp.Request.Method, resp.Request.URL)
}
