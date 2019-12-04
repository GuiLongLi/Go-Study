package main

import (
	"math/rand"
	"time"
	"fmt"
)

/*
并发进程取消的原因：
·超时：超时是隐式取消

·用户干预：为了获取良好的用户体验，通常建议维持一个长连接，然后以轮询间隔将状态报告发送给用户，或允许用户查看他们认为合适的状态。
当用户使用并发程序时，有时需要允许用户取消他们已经开始的操作。

·父进程取消：对于这个问题，如果任何一种并发操作的父进程停止，那么子进程也将被取消。

·复制请求：我们可能希望将数据发送到多个并发进程，以尝试从其中一个进程获得更快的响应。
当第一个响应回来的时候，我们将会取消其余的进程。
*/

func main() {
	fmt.Println("testheart")
	testheart()
	fmt.Println()
	fmt.Println("testheart2")
	testheart2()
}

//心跳
func testheart(){
	doWork := func(
		done <-chan interface{},
		pulseInterval time.Duration,
		) (<-chan interface{},<-chan time.Time) {
		heartbeat := make(chan interface{})
		results := make(chan time.Time)
		go func() {
			defer close(heartbeat)
			defer close(results)

			pulse := time.Tick(pulseInterval)
			workGen := time.Tick(2*pulseInterval)

			sendPulse := func() {
				select {
				case heartbeat <- struct {}{}:
				default:
				}
			}
			sendResult := func(r time.Time) {
				for {
					select {
					case <-done:
						return
					case <-pulse:
						sendPulse()
					case results<-r:
						return
					}
				}
			}
			for {
				select {
				case <- done:
					return
				case <-pulse:
					sendPulse()
				case r := <-workGen:
					sendResult(r)
				}
			}
		}()
		return heartbeat,results
	}
	/*
	·heartbeat := make(chan interface{})
	我们建立了一个发送心跳的channel ，我们把这个返回给 doWork。

	·pulse := time.Tick(pulseInterval)
	我们设定心跳的间隔时间为我们接到的pulseInterval ，每隔一个pulseInterval 的时长都会有一些东西读取这个channel

	·workGen := time.Tick(2*pulseInterval)
	这是另一个用来模拟滴答声的channel ，我们选择的持续时间大于pulseInterval ，这样我们就能看到从goroutine中发出的一些心跳。

	·case heartbeat <- struct {}{}:
	  default:
	注意，我们在这里加入了一个默认语句，我们必须时刻警惕这样一个事实：
	可能会没有人接收我们的心跳，从goroutine 发出的信息是重要的，但心跳却不一定重要。

	·case <-pulse:
	  sendPulse()
	就像done channel 一样，当你执行发送或接收时，你也需要包含一个发送心跳的分支。

	*/

	done := make(chan interface{})
	time.AfterFunc(10*time.Second, func() {close(done)})

	const timeout = 2*time.Second
	heartbeat,results := doWork(done,timeout/2)
	for {
		select {
		case _,ok := <-heartbeat:
			if ok == false{
				return
			}
			fmt.Println("pulse")
		case r,ok := <-results:
			if ok == false{
				return
			}
			fmt.Printf("results %v\n", r.Second())
		case <-time.After(timeout):
			return
		}
	}

	/*
	·time.AfterFunc(10*time.Second, func() {close(done)})
	我们声明了一个标准的done channel ，并在10秒后关闭，这给我们的goroutine 做一些工作的时间

	·const timeout = 2*time.Second
	这里我们设置了超时时间，我们使用此方法将心跳间隔与超时时间联系起来

	·heartbeat,results := doWork(done,timeout/2)
	我们在这里timeout/2 ，这使得我们的心跳有额外的响应时间，以便我们的超时有一定缓冲时间。

	·case _,ok := <-heartbeat:
	在这里，我们处理心跳，当没有消息时，我们至少知道每过timeout/2 的时间会从心跳channel 发出一条消息。
	如果我们什么都没有收到，我们便知道是goroutine 本身出了问题。

	·case r,ok := <-results:
	在这里，我们处理results channel

	·case <-time.After(timeout):
	如果我们没有收到心跳或其他消息，就会超时

	*/
}

func testheart2(){
	doWork := func(
		done <-chan interface{},
		) (<-chan interface{},<-chan int){
		heartbeatStream := make(chan interface{},1)
		workStream := make(chan int)
		go func() {
			defer close(heartbeatStream)
			defer close(workStream)

			for i:=0;i < 10;i++{
				select {
				case heartbeatStream<- struct {}{}:
				default:
				}

				select {
				case <-done:
					return
				case workStream<-rand.Intn(10):
				}
			}
		}()
		return heartbeatStream,workStream
	}

	done:=make(chan interface{})
	defer close(done)

	heartbeat,results := doWork(done)
	for{
		select {
		case _,ok := <-heartbeat:
			if ok {
				fmt.Println("pulse")
			} else {
				return
			}
		case r,ok := <-results:
			if ok {
				fmt.Printf("results %v\n", r)
			}else{
				return
			}
		}
	}

	/*
	·heartbeatStream := make(chan interface{},1)
	在这里，我们创建一个缓冲区大小为1的 heartbeat channel ，这确保了即使没有及时接受发送的消息，至少也会发送一个心跳。

	·for i:=0;i < 10;i++{
	  select {
	在这里，我们为心跳设置了一个单独的select块，我们希望将发送results 和心跳分开，
	因为如果接受者没有准备好接收结果，作为替代，他将接收到一个心跳，而代表当前结果的值将会丢失。
	由于我们有默认逻辑，所以这里也没有包含对done channel 的处理。

	·default:
	在这里，为了防止没人接收我们的心跳，我们增加了默认逻辑，
	因为我们的heartbeat channel 创建时有一个缓冲区那么大，所以如果有人正在监听，
	但是没有及时收到第一个心跳，接收者仍然可以接收到心跳。


	*/
}
