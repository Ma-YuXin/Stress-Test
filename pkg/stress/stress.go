package stress

import (
	"context"
	"errors"
	"log"
	"stressTest/defs"
	"stressTest/util"
	"sync"
	"time"
)

var (
	funcs = map[string]func(*sync.WaitGroup, context.Context, string, int){
		"POST":   Post,
		"PATCH":  Patch,
		"DELETE": Delete,
		"PUT":    Put,
	}
	Resindex map[string][]string
)

func init() {
	Resindex = make(map[string][]string)
	for _, k := range defs.Reslist {
		Resindex[k] = util.GetResList(k)
	}
}

func RpsWithPercent(num int, actionRatio map[string]float64, resRatio map[string]map[string]float64, duration time.Duration) {
	err := verify(actionRatio, resRatio)
	if err != nil {
		log.Panicln(err)
	}
	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(duration))
	defer cancel()
	for action, v := range actionRatio {
		at := int(float64(num) * v)
		log.Println(action, at)
		for res, rat := range resRatio[action] {
			rt := int(float64(at) * rat)
			log.Println(action, res, rt)
			wg.Add(1)
			go funcs[action](wg, ctx, res, rt)
		}
	}
	wg.Wait()
}
func verify(actionRatio map[string]float64, resRatio map[string]map[string]float64) error {
	t1 := 0.0
	for action, v := range actionRatio {
		t1 += v
		t2 := 0.0
		for _, rat := range resRatio[action] {
			t2 += rat
		}
		if t2 > 1.0 {
			return errors.New("res ratio greater than 1")
		}
	}
	if t1 > 1.0 {
		return errors.New("action ratio greater than 1")
	}
	return nil
}
