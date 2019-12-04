package main

import (
	"fmt"
	"sync"
	"math/rand"
	"time"
)

/*
复制请求：
对于某些应用来说，尽可能快的接收响应是重中之重。
例如：程序正在处理用户的http 请求，或者检索一个数据块，
在这些情况下，你可以进行权衡：
可以将请求分发到多个处理程序（无论是goroutine ，进程，还是服务器），
其中一个将比其他处理程序返回更快，你可以立即返回结果。
缺点是为了维护多个实例的运行，你将不得不消耗更多的资源。
如果这种复制是在内存中进行的，消耗则没有那么大，但是如果多个处理程序需要多个进程，服务器甚至是数据中心，
那么可能会变得相对昂贵
*/

func main() {
	fmt.Println("testcprequest")
	testcprequest()
}

//复制请求
func testcprequest(){
	doWork := func(
		done <-chan interface{},
		id int,
		wg *sync.WaitGroup,
		result chan<- int,
		) {
		started := time.Now()
		defer wg.Done()

		//模拟随机负载
		simLoadTime := time.Duration(1+rand.Intn(5))*time.Second
		select {
		case <-done:
		case <-time.After(simLoadTime):
		}

		select {
		case <-done:
		case result <-id:
		}

		took := time.Since(started)
		//显示处理程序需要多长时间
		if took < simLoadTime{
			took = simLoadTime
		}
		fmt.Printf("%v took %v\n",id,took)
	}

	done := make(chan interface{})
	result := make(chan int)

	var wg sync.WaitGroup
	wg.Add(10)

	for i :=0;i<10;i++{
		go doWork(done,i,&wg,result)
	}

	firstReturned := <-result
	close(done)
	wg.Wait()

	fmt.Printf("received an answer from #%v\n", firstReturned)
	/*
	·for i :=0;i<10;i++{
	在这里，我们启动10个处理程序来处理请求

	·firstReturned := <-result
	在这里获得处理程序组的第一个返回值

	·close(done)
	在这里，我们取消其余的处理程序，以保证他们不会继续做多余的工作。

	*/

}