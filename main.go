package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"stressTest/defs"
	"stressTest/pkg/delete"
	"stressTest/pkg/patch"
	"stressTest/pkg/post"
	"stressTest/pkg/put"
	"stressTest/pkg/stress"
	"syscall"
	"time"
)

var (
	duration = time.Minute
)

func init() {
	rps := flag.Int("rps", 1, "the Base of rps")
	resnum := flag.Int("resnum", 10, "the num of res prepare for action such as : put patch and so on")
	period := flag.Duration("duration", time.Minute, "test duration time")
	// verbose := flag.Bool("verbose", true, "add more exclamations")
	flag.Parse()
	stress.SingleResNum = *resnum
	stress.RpsBase = *rps
	duration = *period
	log.Println("the base of rps is ", *rps)
	log.Println("the num of res is ", *resnum)
	log.Println("duration is ", *period)
}

// kubectl delete -n myx-test -l env=test no
// kubectl delete -n myx-test -l env=test pv
// kubectl delete -n myx-test -l env=test cm
// kubectl delete -n myx-test -l env=test ep
// kubectl delete -n myx-test -l env=test limits
// kubectl delete -n myx-test -l env=test pvc
// kubectl delete -n myx-test -l env=test po
// kubectl delete -n myx-test -l env=test podtemplate
// kubectl delete -n myx-test -l env=test rc
// kubectl delete -n myx-test -l env=test quota
// kubectl delete -n myx-test -l env=test secret
// kubectl delete -n myx-test -l env=test sa
// kubectl delete -n myx-test -l env=test svc
// kubectl delete -n myx-test -l env=test controllerrevision
// kubectl delete -n myx-test -l env=test ds
// kubectl delete -n myx-test -l env=test deploy
// kubectl delete -n myx-test -l env=test rs
// kubectl delete -n myx-test -l env=test sts
// kubectl delete -n myx-test -l env=test cj
// kubectl delete -n myx-test -l env=test job
// kubectl delete -n myx-test -l env=test ns

// put ds
func main() {
	// post.CreateRes(config.GetDefultNameSpace(), "sa", 50)
	// delete.ClearAll(config.GetDefultNameSpace(), config.GetDefaultLabelSelector(), config.GetDefultAuthor())
	rps()
	// postTest()
}

func rps() {
	// resRatio := map[string]map[string]int{
	// 	"POST":   {"pv": 4, "cm": 4, "ep": 4, "limits": 4, "pvc": 4, "podtemplate": 4, "rc": 4, "quota": 4, "secret": 4},
	// 	"PATCH":  {"pv": 4, "cm": 4, "ep": 4, "limits": 4, "pvc": 4, "podtemplate": 4, "rc": 4, "quota": 4, "secret": 4},
	// 	"PUT":    {"pv": 4, "cm": 4, "ep": 4, "limits": 4, "pvc": 4, "podtemplate": 4, "rc": 4, "quota": 4, "secret": 4},
	// 	"GET":    {"pv": 4, "cm": 4, "ep": 4, "limits": 4, "pvc": 4, "podtemplate": 4, "rc": 4, "quota": 4, "secret": 4},
	// 	"LIST":   {"pv": 4, "cm": 4, "ep": 4, "limits": 4, "pvc": 4, "podtemplate": 4, "rc": 4, "quota": 4, "secret": 4},
	// 	"DELETE": {"pv": 1, "cm": 1, "ep": 1, "limits": 1, "pvc": 1, "podtemplate": 1, "rc": 1, "quota": 1, "secret": 1},
	// }
	// "podtemplate":2,
	// "ns":2,
	// "no":2,
	// "pv":2,
	// "cm":2,
	// "ep":2,
	// "limits":2,
	// "pvc":2,
	// "po" :2,
	// "rc" :2,
	// "quota" :2,
	// "secret":2,
	// "sa":2,
	// "svc":2,
	// "ds" :2,
	// "deploy" :2,
	// "rs" :2,
	// "sts":2,
	// "cj" :2,
	// "job":2
	resRatio := map[string]map[string]int{
		"POST":   {"podtemplate": 2, "ns": 2, "pv": 2, "cm": 2, "ep": 2, "limits": 2, "pvc": 2, "rc": 2, "quota": 2, "secret": 2, "sa": 2, "svc": 2, "ds": 2, "deploy": 2, "rs": 2, "sts": 2, "cj": 2, "job": 2},
		"PATCH":  {"podtemplate": 2, "ns": 2, "pv": 2, "cm": 2, "ep": 2, "limits": 2, "pvc": 2, "rc": 2, "quota": 2, "secret": 2, "sa": 2, "svc": 2, "ds": 2, "deploy": 2, "rs": 2, "sts": 2, "cj": 2, "job": 2},
		"PUT":    {"podtemplate": 2, "ns": 2, "pv": 2, "cm": 2, "ep": 2, "limits": 2, "pvc": 2, "rc": 2, "quota": 2, "secret": 2, "sa": 2, "svc": 2, "ds": 2, "deploy": 2, "rs": 2, "sts": 2, "cj": 2},
		"GET":    {"podtemplate": 2, "ns": 2, "pv": 2, "cm": 2, "ep": 2, "limits": 2, "pvc": 2, "rc": 2, "quota": 2, "secret": 2, "sa": 2, "svc": 2, "ds": 2, "deploy": 2, "rs": 2, "sts": 2, "cj": 2, "job": 2},
		"LIST":   {"podtemplate": 2, "ns": 2, "pv": 2, "cm": 2, "ep": 2, "limits": 2, "pvc": 2, "rc": 2, "quota": 2, "secret": 2, "sa": 2, "svc": 2, "ds": 2, "deploy": 2, "rs": 2, "sts": 2, "cj": 2, "job": 2},
		"DELETE": {"pv": 1, "ns": 1, "limits": 1, "pvc": 1, "podtemplate": 1, "rc": 1, "secret": 1, "deploy": 1},
	}
	stress.RpsWithPercent(resRatio, duration)
}
func postRps() {

	resRatio := map[string]map[string]int{
		"POST": {"cm": 1},
	}
	stress.RpsWithPercent(resRatio, time.Second*60)
}
func postTest() {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGQUIT)
	done := make(chan bool, 1)
	go func() {
		sig := <-sigs
		done <- true
		fmt.Println()
		fmt.Println(sig)
		cancel()
		log.Println("received interrupt , run aftercare program ,and clearing")
		log.Println("start clear by interput conn : ")
	}()
	rl := []string{"sa", "svc",
		"ds", "deploy", "rs", "sts", "no", "pv",
		"cm", "ep", "limits", "pvc", "po", "podtemplate",
		"rc", "cj", "job"}
	concurrency_list := []int{3, 6, 9, 15, 30, 60}
	anno_num_list := []int{100, 200, 300, 400}
	for _, an := range anno_num_list {
		for _, cn := range concurrency_list {
			for _, v := range rl {
				select {
				case <-done:
					log.Println("pra stop")
					return
				default:
					s := post.NewStress(-1, cn, an, "myx-test", time.Minute)
					s.Res = v
					s.RpsPerConn = 100
					s.Run(ctx)
				}
			}
		}
	}
}
func clearSingle(res string) {
	if res == "ns" {
		delete.DeleteNameSpace(res, "myx-test", "env=test", "Bearer "+defs.Token)
	} else {
		delete.ClearPost(res, "myx-test", "env=test", "Bearer "+defs.Token)
	}
}
func patchTest(res string) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	concurrency_list := []int{3, 6, 9, 15, 30, 60, 90, 150, 300}
	anno_num_list := []int{0, 100, 200, 300, 400}
	for _, an := range anno_num_list {
		for _, cn := range concurrency_list {
			s := patch.NewStress(-1, cn, an, "myx-test", time.Minute)
			s.Res = res
			s.Run(ctx)
			time.Sleep(time.Second * 40)
		}
	}
}
func putTest(res string) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	concurrency_list := []int{3, 6, 9, 15, 30, 60, 90, 150, 300}
	anno_num_list := []int{0, 100, 200, 300, 400}
	for _, an := range anno_num_list {
		for _, cn := range concurrency_list {
			s := put.NewStress(-1, cn, an, "myx-test", time.Minute)
			s.Res = res
			s.Run(ctx)
			time.Sleep(time.Second * 40)
		}
	}
}
func deleteTest(res string) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()
	s := delete.NewStress(10, "myx-test", time.Second*100)
	s.Res = res
	s.Run(ctx)
}

func clear() {
	delete.ClearAll("myx-test", "env=test", "Bearer "+defs.Token)
}
