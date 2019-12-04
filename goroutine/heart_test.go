package main

import (
	"time"
	"testing"
)

func DoWork(
	done <-chan interface{},
	pulseInterval time.Duration,
	nums ...int,
	)(<-chan interface{},<-chan int){
	heartbeat := make(chan interface{},1)
	intStream := make(chan int)

	go func() {
		defer close(heartbeat)
		defer close(intStream)

		time.Sleep(2*time.Second)

		pulse := time.Tick(pulseInterval)
		numLoop:
			for _,n := range nums {
				for {
					select {
					case <- done:
						return
					case <-pulse:
						select {
						case heartbeat<- struct {}{}:
						default:
						}
					case intStream <- n:
						continue numLoop
					}
				}
			}
	}()
	return heartbeat,intStream
}

func TestDoWork_GeneratesAllNumbers(t *testing.T){
	done := make(chan interface{})
	defer close(done)

	intSlice := []int{0,1,2,3,5}
	const timeout = 2*time.Second
	heartbeat,results := DoWork(done,timeout/2,intSlice...)

	<-heartbeat
	i:=0
	for {
		select {
		case r,ok := <-results:
			if ok == false{
				return
			}else if expected := intSlice[i];r != expected{
				t.Errorf(
					"index %v: expected %v,but received %v,",
					i,
					expected,
					r,
					)
			}
			i++
		case <-heartbeat:
		case <-time.After(timeout):
			t.Fatal("test timed out")
		}
	}
	/*
	·numLoop:
	使用一个标签来标识内部循环

	·for _,n := range nums {
	  for {
	我们需要两个循环，一个外部循环遍历数列，另一个内部循环持续执行，直到intStream 中的数字成功发送。

	·continue numLoop
	跳回numLoop 标签继续执行外部循环

	·<-heartbeat
	等待第一次心跳到达，来确认goroutine 已经进入了循环。

	·case <-heartbeat:
	接收心跳，以防止超时。

	*/
}