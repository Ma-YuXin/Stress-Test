package main

import (
	"stressTest/defs"
	"stressTest/pkg/delete"
	"stressTest/pkg/patch"
	"stressTest/pkg/post"
	"stressTest/pkg/put"
	"time"
)

// patch : ns
func main() {
	rl := []string{
		"cm", "ep", "limits", "pvc", "po", "podtemplate",
		"rc", "quota", "secret", "sa", "svc",
		"controllerrevision", "ds", "deploy", "rs", "sts",
		"cj", "job"}
	for _, v := range rl {
		// v := "pv"
		post.CreateRes("myx-test", v, 1000)
		time.Sleep(time.Minute)
		patchTest(v)
		time.Sleep(time.Minute)
		clearSingle(v)
		time.Sleep(time.Minute * 15)
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
			time.Sleep(time.Second * 70)
		}
	}
}
func putTest(res string) {
	s := put.NewStress(100, 2, "myx-test", time.Second*10)
	s.Res = res
	s.Run()
}
func deleteTest(res string) {
	s := delete.NewStress(10, "myx-test", time.Second*100)
	s.Res = res
	s.Run()
}

func clearSingle(res string) {
	delete.ClearPost(res, "myx-test", "env=test", "Bearer "+defs.Token)
}
func clear() {
	delete.ClearAll("myx-test", "env=test", "Bearer "+defs.Token)
}
func postTest() {
	// fmt.Println(os.Getpid())
	// time.Sleep(time.Second * 5)
	s := post.NewStress(-1, 1, 0, "myx-test", time.Second*10)
	// for _, res := range defs.Reslist {
	// 	s.Res = res
	// 	anno_num_list := []int{0, 10, 20, 30, 40}
	// 	concurrency_list_no_queue := []int{3, 6, 9, 15, 30, 60, 90, 150, 300, 600, 900}
	// 	for _, v := range anno_num_list {
	// 		s.Anntation = v
	// 		for _, cn := range concurrency_list_no_queue {
	// 			s.Conn = cn
	// 			s.Run()
	// 		}
	// 	}
	// }
	// for _, res := range defs.Reslist {
	// 	s.Res = res
	// 	s.Anntation = 1s
	// 	s.Conn = 1
	// 	s.Run()
	// }
	// successful:
	// job cj sts rs deploy ds rc podtemplate pvc ns no cm ep limits po quota secret sa svc controllerrevision pv
	s.Res = "ns"
	s.Run()
	// DeleteTerminated()
}
