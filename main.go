package main

import (
	"log"
	"stressTest/defs"
	"stressTest/pkg/delete"
	"stressTest/pkg/patch"
	"stressTest/pkg/post"
	"stressTest/pkg/put"
	"stressTest/pkg/stress"
	"time"
)

// put ds
func main() {
	rl := []string{ "no", "pv",
		"cm", "ep", "limits", "pvc", "po", "podtemplate",
		"rc", "quota", "secret", "sa", "svc",
		"controllerrevision", "ds", "deploy", "rs", "sts",
		"cj", "job"}
	concurrency_list := []int{3, 6, 9, 15, 30, 60}
	anno_num_list := []int{100, 200, 300, 400}
	for _, v := range rl {
		for _, an := range anno_num_list {
			for _, cn := range concurrency_list {
				log.Println("start conn : ", cn, ", annotation : ", an)
				s := post.NewStress(-1, cn, an, "myx-test", time.Minute)
				s.Res = v
				s.RpsPerConn = 20
				s.Run()
				time.Sleep(time.Minute * 1)
			}
		}
	}
}
// crictl pull registry.k8s.io/kube-apiserver:v1.29.0
// crictl pull registry.k8s.io/kube-controller-manager:v1.29.0
// crictl pull registry.k8s.io/kube-scheduler:v1.29.0
// crictl pull registry.k8s.io/kube-proxy:v1.29.0
// crictl pull registry.k8s.io/coredns/coredns:v1.11.1
// crictl pull registry.k8s.io/pause:3.9
// crictl pull registry.k8s.io/etcd:3.5.10-0

func rps() {
	actionRatio := map[string]float64{
		"POST":   0.25,
		"PATCH":  0.25,
		"DELETE": 0.25,
		"PUT":    0.25,
	}
	// "ns", "no", "pv",
	// 	"cm", "ep", "limits", "pvc", "po", "podtemplate",
	// 	"rc", "quota", "secret", "sa", "svc",
	// 	"controllerrevision", "ds", "deploy", "rs", "sts",
	// 	"cj", "job"
	resRatio := map[string]map[string]float64{
		"POST":   {"pv": 0.1, "cm": 0.1, "ep": 0.1, "limits": 0.1, "pvc": 0.1, "podtemplate": 0.1, "rc": 0.1, "quota": 0.1, "secret": 0.1},
		"PATCH":  {"pv": 0.1, "cm": 0.1, "ep": 0.1, "limits": 0.1, "pvc": 0.1, "podtemplate": 0.1, "rc": 0.1, "quota": 0.1, "secret": 0.1},
		"DELETE": {"pv": 0.1, "cm": 0.1, "ep": 0.1, "limits": 0.1, "pvc": 0.1, "podtemplate": 0.1, "rc": 0.1, "quota": 0.1, "secret": 0.1},
		"PUT":    {"pv": 0.1, "cm": 0.1, "ep": 0.1, "limits": 0.1, "pvc": 0.1, "podtemplate": 0.1, "rc": 0.1, "quota": 0.1, "secret": 0.1},
	}
	stress.RpsWithPercent(100, actionRatio, resRatio, time.Second)
}
func postRps() {
	actionRatio := map[string]float64{
		"POST": 1.0,
	}
	resRatio := map[string]map[string]float64{
		"POST": {"cm": 1.0},
	}
	stress.RpsWithPercent(100, actionRatio, resRatio, time.Second*30)
}
func postTest(res string) {
	concurrency_list := []int{3, 6, 9, 15, 30, 60, 90, 150, 300}
	anno_num_list := []int{0, 100, 200, 300, 400}
	for _, an := range anno_num_list {
		for _, cn := range concurrency_list {
			log.Println("start conn : ", cn, ", annotation : ", an)
			s := post.NewStress(-1, cn, an, "myx-test", time.Minute)
			s.Res = res
			s.RpsPerConn = 40
			s.Run()
			time.Sleep(time.Minute * 5)
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
	concurrency_list := []int{3, 6, 9, 15, 30, 60, 90, 150, 300}
	anno_num_list := []int{0, 100, 200, 300, 400}
	for _, an := range anno_num_list {
		for _, cn := range concurrency_list {
			s := patch.NewStress(-1, cn, an, "myx-test", time.Minute)
			s.Res = res
			s.Run()
			time.Sleep(time.Second * 40)
		}
	}
}
func putTest(res string) {
	concurrency_list := []int{3, 6, 9, 15, 30, 60, 90, 150, 300}
	anno_num_list := []int{0, 100, 200, 300, 400}
	for _, an := range anno_num_list {
		for _, cn := range concurrency_list {
			s := put.NewStress(-1, cn, an, "myx-test", time.Minute)
			s.Res = res
			s.Run()
			time.Sleep(time.Second * 40)
		}
	}
}
func deleteTest(res string) {
	s := delete.NewStress(10, "myx-test", time.Second*100)
	s.Res = res
	s.Run()
}

func clear() {
	delete.ClearAll("myx-test", "env=test", "Bearer "+defs.Token)
}
