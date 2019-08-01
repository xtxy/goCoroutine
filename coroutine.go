package bvUtils

import (
	"time"
)

const (
	coroutine_init = iota
	coroutine_stop
)

type CoroutineCtx struct {
	waitChan chan int
	runChan  chan bool
	waitTime int64
	compare  func() bool
	state    int
}

func (self *CoroutineCtx) Start() bool {
	return <-self.runChan
}

func (self *CoroutineCtx) Wait(ms int) bool {
	if ms > 0 {
		self.waitTime = time.Now().Add(time.Duration(ms) * time.Millisecond).UnixNano()
	} else {
		self.waitTime = 0
	}
	self.compare = nil

	self.waitChan <- 1
	return <-self.runChan
}

func (self *CoroutineCtx) Until(compare func() bool) bool {
	self.compare = compare
	self.waitTime = 0

	self.waitChan <- 1
	return <-self.runChan
}

func (self *CoroutineCtx) Stop() {
	self.state = coroutine_stop
	self.waitChan <- 1
}

type CoroutineMgr struct {
	coroutines []*CoroutineCtx
}

func NewCoroutineMgr() *CoroutineMgr {
	coroutineMgr := new(CoroutineMgr)
	coroutineMgr.coroutines = make([]*CoroutineCtx, 0)

	return coroutineMgr
}

func (self *CoroutineMgr) NewCoroutine() *CoroutineCtx {
	coroutineCtx := new(CoroutineCtx)
	coroutineCtx.waitChan = make(chan int)
	coroutineCtx.runChan = make(chan bool)

	self.coroutines = append(self.coroutines, coroutineCtx)

	return coroutineCtx
}

func (self *CoroutineMgr) StopAll() {
	for _, v := range self.coroutines {
		if nil == v || v.state == coroutine_stop {
			continue
		}

		v.runChan <- false
		<-v.waitChan
	}

	self.coroutines = nil
}

func (self *CoroutineMgr) Update(duration int64) {
	needRemove := false
	for _, v := range self.coroutines {
		if nil == v || v.state == coroutine_stop {
			needRemove = true
			continue
		}

		if v.waitTime != 0 && time.Now().UnixNano() < v.waitTime {
			continue
		}

		if nil != v.compare && !v.compare() {
			continue
		}

		v.runChan <- true
		<-v.waitChan
	}

	if !needRemove {
		return
	}

	newArr := make([]*CoroutineCtx, 0)
	for _, v := range self.coroutines {
		if nil != v && v.state != coroutine_stop {
			newArr = append(newArr, v)
		}
	}

	self.coroutines = newArr
}
