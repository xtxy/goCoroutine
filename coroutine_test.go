package bvUtils

import (
	"fmt"
	"testing"
	"time"
)

func Test_main(t *testing.T) {
	mgr := NewCoroutineMgr()
	lastTime := time.Now().UnixNano()

	dataMap := map[string]int{"a": 1, "b": 2}

	ctx1 := mgr.NewCoroutine()
	go job1(ctx1, dataMap)

	ctx2 := mgr.NewCoroutine()
	go job2(ctx2, dataMap)

	for i := 0; i < 100; i++ {
		time.Sleep(100 * time.Millisecond)

		timeNow := time.Now().UnixNano()
		mgr.Update(timeNow - lastTime)
		lastTime = timeNow
	}
}

func job1(ctx *CoroutineCtx, dataMap map[string]int) {
	ctx.Start()
	defer ctx.Stop()

	fmt.Println("job1, a:", dataMap["a"])

	if !ctx.Until(func() bool {
		return dataMap["a"] == 5
	}) {
		return
	}

	fmt.Println("job1, a:", dataMap["a"])

	fmt.Println("job1, change b to 6")
	dataMap["b"] = 6

	fmt.Println("job1, wait 600 ms")
	if !ctx.Wait(600) {
		return
	}

	fmt.Println("job1, b:", dataMap["b"])
}

func job2(ctx *CoroutineCtx, dataMap map[string]int) {
	ctx.Start()
	defer ctx.Stop()

	fmt.Println("job2, change a to 5")
	dataMap["a"] = 5
	fmt.Println("job2, b:", dataMap["b"])

	if !ctx.Until(func() bool {
		return dataMap["b"] == 6
	}) {
		return
	}

	fmt.Println("job2, b:", dataMap["b"])

	fmt.Println("job2, change b to 7")
	dataMap["b"] = 7
}
