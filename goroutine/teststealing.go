package main

import (
	"os"
	"log"
	"time"
	"runtime/pprof"
	"fmt"
)
/*
工作窃取：
go语言遵循for-join模型进行并发。
在goroutine 开始的时候fork ,join点是两个或更多的goroutine 通过channel 或sync 包中的类型进行同步时。
工作窃取算法遵循一些基本原则，对于给定的线程：
·在fork 点，将任务添加到与线程关联的双端队列的尾部
·如果线程空闲，则选取一个随机的线程，从他关联的双端队列头部窃取工作
·如果在未准备好的join 点（即与其同步的goroutine 还没完成），则将工作从线程的双端队列尾部出栈
·如果线程的双端队列是空的，则：
	·暂停加入
	·从随机线程关联的双端队列中窃取工作
*/

func main() {
	fmt.Println("teststealing")
	teststealing()
	fmt.Println()
	fmt.Println("testpprof")
	testpprof()
}


//递归计算Fibonacci 数列
//窃取任务还是续体
func teststealing(){
	/*
	我们应该让什么样的任务进行排队和窃取：
	在fork-join 模式下，有两种选择：新任务和续体。
	*/
	var fib func(n int) <-chan int
	fib = func(n int) <-chan int {
		result := make(chan int)
		go func() {
			defer close(result)
			if n <= 2{
				result <- 1
				return
			}
			result <- <-fib(n-1) + <-fib(n-2)
		}()
		return result
	}

	fmt.Printf("fib(4) = %d\n",<-fib(4))
	/*
	·go func() {
	在go 语言中，goroutine 就是任务

	·return result
	在goroutine 之后的一切都被称为续体

	go语言的工作窃取算法是对续体进行入队和窃取
	*/

	/*
	假设我们机器有两个单核处理器，
	我们在每个处理器上生成一个系统线程：
	处理器1上生成T1 ,处理器2上生成T2

	·当续体在一个工作队列上排队时，我们将它列为 cont.of X
	·当一个续体被用于执行时，我们做一个隐式转换，将续体转换为fib 的下一个调用

	·我们从main goroutine 开始：
	T1 调用栈		T1 工作队列		T2 调用栈		T2 工作队列
	main

	·main goroutine 调用fib(4) 并将这个调用的续体追加到T1 工作队列尾部：
	T1 调用栈		T1 工作队列		T2 调用栈		T2 工作队列
	fib(4)			cont.of main

	·T2 是空闲的，所以他窃取了main 的续体：
	T1 调用栈		T1 工作队列		T2 调用栈		T2 工作队列
	fib(4)							cont.of main

	·fib(4) 调用了fib(3) ，并且立刻就开始运行，同时T1 将fib(4) 的续体入栈到他的队列尾部：
	T1 调用栈		T1 工作队列		T2 调用栈		T2 工作队列
	fib(3)			cont.of fib(4)	cont.of main

	·当T2 试图执行main 的续体时，他是一个等待join 的点，因此，他从T1 的队列中窃取了更多的工作，这一次，他得到了fib(4) 的续体：
	T1 调用栈		T1 工作队列		T2 调用栈				T2 工作队列
	fib(3)							cont.of main(等待join)
									cont.of fib(4)

	·接下来，T1 上的fib(3) 为fib(2) 启动了一个goroutine ，fib(3) 的续体被入栈到其工作队列的尾部：
	T1 调用栈		T1 工作队列		T2 调用栈		T2 工作队列
	fib(2)			cont.of fib(3)	cont.of main
									cont.of fib(4)

	·T2 执行了T1 未执行的fib(4) 续体，他又调用了一个新的fib(2)，fib(4) 的续体又继续进入队列
	T1 调用栈		T1 工作队列		T2 调用栈		T2 工作队列
	fib(2)			cont.of fib(3)	cont.of main	cont.of fib(4)
									fib(2)

	·接下来，T1 执行fib(2) 到了递归终结点，返回了1：
	T1 调用栈		T1 工作队列		T2 调用栈		T2 工作队列
	(returns 1)		cont.of fib(3)	cont.of main	cont.of fib(4)
									fib(2)

	·然后T2 也到达终结点并返回了1：
	T1 调用栈		T1 工作队列		T2 调用栈		T2 工作队列
	(returns 1)		cont.of fib(3)	cont.of main	cont.of fib(4)
									(returns 1)

	·然后，T1 从他自己的队列中取出了一个任务并开始执行，就是fib(1)。这里我们看一下T1 的调用链：fib(3) -> fib(2) ->fib(1)，这就是我们之前说的续体窃取算法的好处！
	T1 调用栈		T1 工作队列		T2 调用栈		T2 工作队列
	fib(1)							cont.of main	cont.of fib(4)
									(returns 1)

	·然后再fib(4) 的续体结束时，只有一个join 被实现：fib(2) ,对fib(3) 的调用仍然由T1 处理，因为没有工作窃取：
	T1 调用栈		T1 工作队列		T2 调用栈		T2 工作队列
	fib(1)							cont.of main
									fib(4)

	·T1 现在处于一个续体结束的阶段，fib(3) ，从fib(2) 和fib(1) join完成，T1 返回2：
	T1 调用栈		T1 工作队列		T2 调用栈		T2 工作队列
	(returns 2)						cont.of main
									(returns 2)

	·现在fib(4)、fib(3) 和fib(2) 的join都完成了，T2 能够执行其计算并返回结果(2+1=3)：
	T1 调用栈		T1 工作队列		T2 调用栈		T2 工作队列
									cont.of main
									(returns 3)

	·最后，main goroutine 已经join 完成了，他从fib(4) 收到返回值，然后打印出结果3：
	T1 调用栈		T1 工作队列		T2 调用栈		T2 工作队列
									main(prints 3)

	·当我们看过上面那些后，我们简要总结下续体在T1 上时是如何起作用的，如果我们看一下这个运行的统计数据(带有连续的窃取)，并与任务窃取相比，就会发现一个更清晰的优势：
	统计				续体窃取					任务窃取
	步骤				14							15
	队列最大长度		2							2
	延迟join			2(所有都在空闲的线程上)		3(所有都在忙碌的线程上)
	调用栈最大深度		2							3

	·这些统计数据似乎很接近，但我们可以推断如果在更大的程序中，我们就会发现续体窃取会带来显著的好处。
	*/

	/*
	go语言的调度器有三个主要的概念：
	G	goroutine
	M	OS线程(在源代码中也被称为机器)
	P	上下文(在源代码中也被称为处理器)

	在我们关于工作窃取的讨论中，M 等于T ，而P 等于工作队列(改变GOMAXPROCS 这个环境变量，可以改变分配数量)，G 是一个goroutine ，但是记住他只代表goroutine 的当前状态，最明显的是它的程序计数器(PC) ，G 相当于一个计算续体，使go语言可以实现续体窃取
	在go 语言的运行时中，首先启动M ，然后是P ，最后是调度运行G ：
	*/

}

/*
本书学习到此，已经接近尾声

竞争检测
在go1.1 中，为大多数命令增加了 -race 参数：
go test -race mypkg  		对此pkg 进行测试
go run -race mysrc.go 		编译程序并运行
go build -race mycmd		构建此命令
go install -race mypkg		安装此pkg
竞争检测是一个非常有用的工具，可用于自动检测代码中的竞争条件。
我们强烈建议将其整合为持续集合过程中的一部分，同样，由于竞争检测只能检测到已经产生的竞争，并且我们介绍过竞争条件有时难以触发的请，
因此，应该在持续集成中运行真实环境的场景以尝试触发竞争。


pprof
pprof是一个google 创造的工具，你可以在程序运行时显示当前的数据，或者使用这个工具来保存运行时的统计信息，这个程序的help 标签可以帮助你很细致的描述了该程序的使用方法，在这里我们将只讨论runtime/pprof 包，因为他涉及了并发
runtime/pprof 包非常简单，并具有预定义的配置文件，无需配置即可进行hook 和显示：
goroutine		堆栈跟踪当前所有goroutines
heap			所有堆分配的一个抽样
threadcreate	指向创建新系统线程的堆栈跟踪
block			导致同步原语阻塞的堆栈跟踪
mutex			争用互斥持有者的堆栈跟踪
*/

func testpprof(){
	log.SetFlags(log.Ltime|log.LUTC)
	log.SetOutput(os.Stdout)

	//每1 秒 log都会记录有多少个goroutine 在并发执行
	go func() {
		goroutines := pprof.Lookup("goroutine")
		for range time.Tick(1*time.Second){
			log.Printf("goroutine count:%d\n", goroutines.Count())
		}
	}()

	//创建一些永远不会退出的goroutines
	var blockForever chan struct{}
	for i:=0;i<10;i++{
		go func() {<-blockForever}()
		time.Sleep(500*time.Millisecond)
	}

	newProfIfNotDef := func (name string) *pprof.Profile{
		prof := pprof.Lookup(name)
		if prof == nil{
			prof = pprof.NewProfile(name)
		}
		return prof
	}

	prof := newProfIfNotDef("my_package_namespace")
	log.Println(prof)
}