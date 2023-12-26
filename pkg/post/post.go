package post

import (
	"bytes"
	"context"
	"io"
	"log"
	"math"
	"net/http"
	"net/http/httputil"
	"strconv"
	"stressTest/client"
	"stressTest/defs"
	"stressTest/ioinfo"
	"stressTest/util"
	"strings"
	"sync"
	"time"
)

type stress struct {
	defs.Meta
	defs.Config
	RpsPerConn  int
	clientSet   []*http.Client
	ConnSend    []int
	ConnRecv    []int
	ConnSendNum []int
}

func NewStress(num, conn, antNum int, ns string, duration time.Duration) *stress {
	if num > 0 && num < conn {
		panic("num must greater or equal than conn")
	}
	if conn < 1 {
		panic("conn must greater than 0")
	}
	if antNum < 0 {
		panic("anntation number must greater or equal than zero ")
	}
	res := &stress{
		Meta: defs.Meta{
			Namespace: ns,
		},
		Config: defs.Config{
			Conn:          conn,
			Num:           num,
			Duration:      duration,
			Action:        "POST",
			Anntation:     antNum,
			LabelSelector: "env=test",
			Auth:          "Bearer " + defs.Token,
		},
	}
	return res
}
func CreateRes(ns, res string, num int) {
	s := NewStress(num, 1, 0, ns, time.Hour)
	s.Res = res
	s.run()
}
func (s *stress) Run() {
	start := time.Now()
	s.run()
	ioinfo.WriteInfo(start, s)
	time.Sleep(time.Minute * 1)
	start = time.Now()
	s.Clear()
	end := time.Now()
	s.Duration = end.Sub(start)
	ioinfo.WriteInfo(start, s)
}
func (s *stress) Info() (string, string, string, int, int, time.Duration, []int, []int, []int) {
	return s.Res, s.Namespace, s.Action, s.Conn, s.Anntation, s.Duration, s.ConnSend, s.ConnRecv, s.ConnSendNum
}

func (s *stress) initStress() {
	s.ConnSend = make([]int, s.Conn)
	s.ConnRecv = make([]int, s.Conn)
	s.ConnSendNum = make([]int, s.Conn)
	s.clientSet = client.ClientSet(s.Conn)
}
func (s *stress) run() {
	s.initStress()
	// defer client.PutClientSet(s.clientSet)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(s.Duration))
	defer cancel()
	wg := &sync.WaitGroup{}
	if s.Num < 0 {
		for i := 0; i < s.Conn; i++ {
			wg.Add(1)
			go s.startWithDeadline(ctx, wg, i, s.clientSet[i])
		}
	} else if s.Num > 0 {
		avg := s.Num / s.Conn
		remain := s.Num % s.Conn
		for i := 0; i < s.Conn; i++ {
			wg.Add(1)
			if i < remain {
				go s.start(ctx, wg, avg+1, i, s.clientSet[i])
			} else {
				go s.start(ctx, wg, avg, i, s.clientSet[i])
			}
		}
	}
	wg.Wait()
}

func (s *stress) startWithDeadline(ctx context.Context, wg *sync.WaitGroup, id int, httpClient *http.Client) {
	s.start(ctx, wg, math.MaxInt64, id, httpClient)
}
func (s *stress) start(ctx context.Context, wg *sync.WaitGroup, num, id int, httpClient *http.Client) {
	defer wg.Done()
	wgr := &sync.WaitGroup{}
	defer wgr.Wait()
	tick := time.NewTicker(time.Second / time.Duration(s.RpsPerConn))
	defer tick.Stop()
	for i := 0; i < num; i++ {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			s.ConnSendNum[id]++
			wgr.Add(1)
			go s.post(wgr, httpClient, i, id)
		}
	}
}

func (s *stress) post(wg *sync.WaitGroup, httpClient *http.Client, seq, id int) {
	defer wg.Done()
	data, url := util.GetPostDataAndUrl(s.Res, s.Namespace, s.Anntation, seq, id)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	log.Println("POST", req.URL.String())
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", s.Auth)
	req.Header.Set("Content-Type", "application/json")
	reques, err := httputil.DumpRequestOut(req, true)
	s.ConnSend[id] += len(reques)
	if err != nil {
		log.Println("DumpRequestOut err:", err)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println("err:", err)
		return
	}
	defer resp.Body.Close()
	log.Println("resp: ", resp.Status, resp.Request.Method, resp.Request.URL)
	respo, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Println("DumpResponse err:", err)
	}
	s.ConnRecv[id] += len(respo)
	io.Copy(io.Discard, resp.Body)
}

func (s *stress) Clear() {
	s.Action = "DELETE"
	connSendNum := make([]int, len(s.ConnSendNum))
	copy(connSendNum, s.ConnSendNum)
	s.initStress()
	defer client.PutClientSet(s.clientSet)
	s.runDel(connSendNum)
}

func (s *stress) runDel(connSendNum []int) {
	wg := &sync.WaitGroup{}
	f := func(id, num int, client *http.Client) {
		for i := 0; i < num; i++ {
			resName := "test-" + strings.ToLower(util.Res2kind(s.Res)) + "-" + strconv.Itoa(i) + "-" + strconv.Itoa(id)
			s.ConnSendNum[id]++
			s.delete(id, resName, client)
		}
		wg.Done()
	}
	for i, v := range connSendNum {
		wg.Add(1)
		go f(i, v, s.clientSet[i])
	}
	time.Sleep(time.Millisecond * 100)
	wg.Wait()
}

func (s *stress) delete(id int, resName string, client *http.Client) {
	_, _, request := util.GetBasic(s.Res, s.Namespace)
	req, err := http.NewRequest("DELETE", request+"/"+resName, nil)
	if err != nil {
		log.Fatal("new http request err", err)
	}
	req.Header.Set("Authorization", s.Auth)
	req.Header.Set("Content-Type", "application/json")
	log.Println("DELETE", req.URL.String())
	reqout, err := httputil.DumpRequestOut(req, true)
	if err != nil {
		log.Println("parse to request out err", err)
	}
	s.ConnSend[id] += len(reqout)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("do request has err", err)
		return
	}
	defer resp.Body.Close()
	log.Println("resp: ", resp.Status, resp.Request.Method, resp.Request.URL)
	repout, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Println("parse to reponse out err", err)
	}
	s.ConnRecv[id] += len(repout)
	io.Copy(io.Discard, resp.Body)
}
