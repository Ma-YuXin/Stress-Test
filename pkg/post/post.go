package post

import (
	"bytes"
	"context"
	"io"
	"log"
	"math"
	"net/http"
	"net/http/httputil"
	"stressTest/client"
	"stressTest/defs"
	"stressTest/ioinfo"
	"stressTest/pkg/delete"
	"stressTest/util"
	"sync"
	"time"
)

type stress struct {
	defs.Meta
	defs.Config
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
	time.Sleep(time.Second * 3)
	s.clearIfAllowed()
}
func (s *stress) Info() (string, string, string, int, int, time.Duration, []int, []int, []int) {
	return s.Res, s.Namespace, s.Action, s.Conn, s.Anntation, s.Duration, s.ConnSend, s.ConnRecv, s.ConnSendNum
}

func (s *stress) initStress() {
	s.ConnSend = make([]int, s.Conn)
	s.ConnRecv = make([]int, s.Conn)
	s.ConnSendNum = make([]int, s.Conn)
}

func (s *stress) clearIfAllowed() {
	if s.Res == "ns" {
		delete.DeleteNameSpace(s.Res, s.Namespace, s.LabelSelector, s.Auth)
		// time.Sleep(time.Second * 5)
		// s.deleteTerminatedNamespace()
	} else {
		delete.ClearPost(s.Res, s.Namespace, s.LabelSelector, s.Auth)
	}
}
func (s *stress) run() {
	s.initStress()
	clientSet := client.ClientSet(s.Conn)
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(s.Duration))
	defer cancel()
	wg := &sync.WaitGroup{}
	if s.Num < 0 {
		for i := 0; i < s.Conn; i++ {
			wg.Add(1)
			go s.startWithDeadline(ctx, wg, i, clientSet[i])
		}
	} else if s.Num > 0 {
		avg := s.Num / s.Conn
		remain := s.Num % s.Conn
		for i := 0; i < s.Conn; i++ {
			wg.Add(1)
			if i < remain {
				go s.start(ctx, wg, avg+1, i, clientSet[i])
			} else {
				go s.start(ctx, wg, avg, i, clientSet[i])
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
	for i := 0; i < num; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			s.ConnSendNum[id]++
			s.post(httpClient, i, id)
		}
	}
}

func (s *stress) post(httpClient *http.Client, seq, id int) {
	// data := []byte(``)
	data, url := util.GetPostDataAndUrl(s.Res, s.Namespace, s.Anntation, seq, id)
	log.Println("data:", string(data))
	log.Println("url :", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", s.Auth)
	req.Header.Set("Content-Type", "application/json")
	reques, err := httputil.DumpRequestOut(req, true)
	s.ConnSend[id] += len(reques)
	// fmt.Println("request:", string(reques))
	if err != nil {
		log.Println("DumpRequestOut err:", err)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println("err:", err)
	}
	defer resp.Body.Close()
	respo, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Println("DumpResponse err:", err)
	}
	s.ConnRecv[id] += len(respo)
	// fmt.Println("response:", len(respo))
	// fmt.Println()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("err:", err)
	}
	log.Println("body:", string(body))
}
