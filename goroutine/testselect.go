package main

import (
	"fmt"
	"time"
)

/*
select语句是将channel绑定在一起的黏合剂，这就是我们如何在一个程序中组合channel以形式更大的抽象事物的方式。
*/

/*
声明select语句是一个具有并发性的Go语言程序中最重要的事情之一，这并不是夸大其词。
在一起系统中两个或多个组件的交集中，可以在本地、单个函数或类型以及全局范围内找到select语言绑定在一起的channel。
除了连接组件之外，在程序中的这些关键节点上，select语句可以帮助安全的将channel与诸如取消、超时、等待和默认值之类的概念结合在一起。
*/

func main() {
	fmt.Println("testSelect")
	testSelect()
}

func testSelect(){
	var c1 <-chan interface{}//读
	c1chan := func() <-chan interface{}{
		cmodule := make(chan interface{},1)
		cmodule<-"c1"
		return  cmodule
	}
	c1 = c1chan()
	select {
	case <-c1:
		fmt.Println("<-c1")
	}
	fmt.Println()

	//---------------------
	var c2 chan<- interface{} //写
	c2 = make(chan<- interface{},1)
	select {
	case c2<-struct {}{}:
		fmt.Println("c2<-struct {}{}")
	}
	close(c2)
	fmt.Println()

	//---------------------
	// 双向的 channel 在读取的情况下close后，变成可用的channel
	// 双向的 channel 在写入不具备下游消费能力的channel ，变成可用的channel
	// 可用的channel 才能被 select执行
	start := time.Now()
	c3 := make(chan interface{})
	go func() {
		time.Sleep(1*time.Second)
		close(c3)
	}()

	fmt.Println("blocking on read...")
	//读取的情况下close后，变成可用的channel
	select {
	case <-c3:
		fmt.Printf("unblocking %v later\n",time.Since(start))
	}
	//写入不具备下游消费能力的channel ，变成可用的channel
	c4 := make(chan interface{},1)
	select {
	case c4<-struct {}{}:
		fmt.Println("c4<-struct {}{}")
	}
	close(c4)
	fmt.Println()

	//---------------------
	//同时有两个可用的channel
	//select 将随机运行其中一个 case
	case1 := make(chan interface{});close(case1)
	case2 := make(chan interface{});close(case2)

	var c1Count,c2Count int
	for i := 0;i <= 100 ;i++  {
		select {
		case <-case1:
			c1Count++
		case <-case2:
			c2Count++
		}
	}
	fmt.Printf("c1Count: %d\nc2Count: %d\n",c1Count,c2Count)
	fmt.Println()

	//----------------------
	//select 超时
	var c5 <-chan int
	select {
	case <-c5: //这个case阻塞了，因为 c5是一个 nil channel
	case <-time.After(1*time.Second):
		fmt.Println("time out")
	}
	fmt.Println()

	//-----------------------
	//select default默认语句
	//当所有channel都被阻塞时，default也会被调用
	start1 := time.Now()
	var c6,c7 <-chan int
	select {
	case <-c6: //阻塞
	case <-c7://阻塞
	default:
		fmt.Printf("In default after %v\n",time.Since(start1))
	}
	fmt.Println()

	//------------------------
	//for-select循环
	done := make(chan interface{})
	go func() {
		time.Sleep(1*time.Second)
		close(done)
	}()
	workCounter := 0
	loop:
		for {
			select {
			case <-done: //case阻塞，不会被触发，直到close(done)
				break loop  //close(done)后，打断loop循环
			default:
			}
			//模拟工作行为
			workCounter++
			time.Sleep(200*time.Millisecond)
		}
	fmt.Printf("Achived %v cycles of work before signalled to stop\n",workCounter)
	fmt.Println()

	//-------------------------------
	//runtime.GOMAXPROCS(runtime.NumCPU())
	/*
	runtime.GOMAXPROCS(CPU的数量) 可以设置多核执行任务
	通常情况下它自动设置为主机上逻辑CPU的数量 runtime.NumCPU()

	特殊情况下，我们会去调整它：
	有一个项目，被竞争环境困扰。有时候测试失败，我们运行测试的主机有4个逻辑CPU
	因此在任何一个点上，我们都有4个goroutine同时执行
	通过增加GOMAXPROCS 已超过我们拥有的逻辑CPU数量，我们能够更频繁的触发竞争条件，从而更快的修复他们

	其他人可能通过实验发现，他们的程序在一定数量的工作队列和线程上运行的更换，但我更主张谨慎些。
	如果你通过调整这个方法来压缩性能，那么在每次提交之后，当你使用不同的硬件，以及使用不同版本的Go语言时，一定要这样做。
	调整这个值会使你的程序更接近它所运行的硬件，但以抽象和长期性能稳定为代价
	*/


}















