package main

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	event1 := make(chan struct{}, 1)
	event2 := make(chan struct{}, 1)
	event3 := make(chan struct{}, 1)

	event1 <- struct{}{}
	wg.Add(3)
	start := time.Now().Unix()
	go Handle("event1", event1, event2)
	go Handle("event2", event2, event3)
	go Handle("event3", event3, event1)
	wg.Wait()

	end := time.Now().Unix()
	fmt.Println(end - start)
}

func Handle(event string, inputchan chan struct{}, outputchan chan struct{}) {

	for i := 0; i < 3; i++ {
		time.Sleep(1 * time.Second)
		goProcessId := GetGoProcessId()
		fmt.Println(goProcessId)
		select {
		case <-inputchan:
			fmt.Println(event)
			outputchan <- struct{}{}
		}
	}
	wg.Done()
}

func GetGoProcessId() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, err := strconv.ParseUint(string(b), 10, 64)
	if err != nil {
		panic(err)
	}
	return n
}

func lock(methodName string) bool {
	success := false
	custId := GetGoProcessId()

	var err error
	success, err = insertLock(methodName, fmt.Sprintf("%d", custId))

	if err != nil {
		return false
	}
	return success
}

func unLock(methodName string) bool {
	success := false
	custId := GetGoProcessId()

	var err error
	success, err = deleteLock(methodName, fmt.Sprintf("%d", custId))

	if err != nil {
		return false
	}
	return success
}

// 是否可以重入锁
func checkReentrantLock(methodName string) bool {
	return true
}

func insertLock(method, custId string) (bool, error) {
	return true, nil
}

func deleteLock(method, custId string) (bool, error) {
	return true, nil
}

// 测试案例
func Test() {
	methodName := "test"
	if !checkReentrantLock(methodName) {
		for !lock(methodName) {
			time.Sleep(time.Second)
		}
	}

	// TODO 业务

	unLock(methodName)
}
