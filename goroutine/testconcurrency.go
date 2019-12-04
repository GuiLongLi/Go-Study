package main

import (
	"sync"
	"time"
	"bytes"
	"math/rand"
	"net/http"
	"fmt"
)


func main() {
	fmt.Println("testconstraint")
	testconstraint()
	fmt.Println()
	fmt.Println("testforselect")
	testforselect()
	fmt.Println()
	fmt.Println("avoidLeak")
	avoidLeak()
	fmt.Println()
	fmt.Println("testorchannel")
	testorchannel()
	fmt.Println()
	fmt.Println("testerror")
	testerror()
	fmt.Println()
	fmt.Println("testpipeline")
	testpipeline()
	fmt.Println()
	fmt.Println("testgenerator")
	testgenerator()
}


/*
约束
在编写并发代码的时候，有以下几种不同的保证操作安全的方法。我们已经介绍了其中两个：
·用于共享内存的同步原语（如 sync.Mutex ）
·通过通信共享内存来进行同步（如 channel ）

但是，在并发处理中还有其他几种情况也是隐式并发安全的：
·不会发生改变的数据
·受到保护的数据

约束是一种确保了信息只能从一个并发过程中获取到的简单且强大的方法。
有两种可能的约束：
·特定约束：
	特定约束是指通过公约实现约束时，无论是由语言社区、你所在的团队还是你的代码库设置
·词法约束：
	词法约束涉及使用词法作用域仅公开用于多个并发进程的正确数据和并发原语，这是的做错事是不可能的。
*/
func testconstraint(){
	//-----------------------
	//特定约束
	data1 := make([]int,4) //长度是4的数组

	//循环函数
	loopData := func(handleData chan<- int){
		defer close(handleData) //函数结束时，关闭 handleData
		for i := range data1 {
			handleData<-data1[i] //循环插入数据到 handleData
		}
	}

	handleData := make(chan int)
	go loopData(handleData)

	//range handleData 取出数据
	for num := range handleData{
		fmt.Println(num)
	}
	/*
	上面我们约束了 handleData的插入只能在 loopData函数里面，
	当随着代码被更多人触及，deadline缩短，就有可能会出错，
	并且约束可能被打破并导致问题发生。
	*/
	fmt.Println()

	//-----------------------
	//词法约束
	chanOwner := func() <-chan int {
		results := make(chan int,5)
		go func() {
			defer close(results)
			for i := 0;i <= 5 ;i++  {
				results<-i
			}
		}()
		return results
	}

	consumer := func(results <- chan int) {
		for result := range results{
			fmt.Printf("received: %d\n",result)
		}
		fmt.Println("done receiving")
	}

	results := chanOwner()
	consumer(results)
	fmt.Println()
	/*
	我们再 chanOwner函数的词法范围内实例化 channel .这将结果写入 channel 的处理的范围约束在他下面定义的闭包中。
	也就是说，chanOwner函数包含了这个 channel的写入处理，以防止其他 goroutine 写入他

	consumer 函数的词法范围内 收到了channel的读处理，约束了只能从中读取信息。
	这样将 main goroutine约束在 channel的只读视图中
	*/

	//--------------------------
	//使用buffer
	printData := func(wg *sync.WaitGroup,data []byte) {
		defer wg.Done()

		var buff bytes.Buffer
		for _,b := range data{
			fmt.Fprintf(&buff,"%c",b)
		}
		fmt.Println(buff.String())
	}

	var wg sync.WaitGroup
	wg.Add(2)
	data := []byte("golang")
	go printData(&wg,data[:3])
	go printData(&wg,data[3:])

	wg.Wait()
}

func testforselect(){
	/*
	for{ //要不就无限循环，要不就使用range语句循环
		select{
			//使用channel进行作业
		}
	}
	*/

	//-----------------
	//向channel发送迭代变量
	stringStream := make(chan string,1);
	done1 := make(chan interface{},1)
	done1<- struct {}{};
	defer close(done1);
	defer close(stringStream);
	for _,s := range []string{"a","b","c"}{
		select {
		case <-done1: //打断循环
			break;
		case stringStream<-s: //插入数据到 stringStream
			fmt.Println(<-stringStream)
		}
	}
	fmt.Println()

	//------------------
	//循环等待停止
	done2 := make(chan interface{},1)
	loopFor:
		for{
			select {
			case <-done2: //打断循环
				fmt.Println("<-done2")
				break loopFor;
			default:
				fmt.Println("select task")
			}
			//进行非抢占式任务
			fmt.Println("common task")
			close(done2)
		}
	fmt.Println()
}

//防止内存泄漏
func avoidLeak(){
	/*
	内存泄漏
	运行多个goroutine我们需要消耗资源，而且goroutine 在运行时，不会被垃圾回收
	所有无论 goroutine 所占用的内存有多么的少，当积累太多时整个操作系统越来越慢，甚至还会导致系统崩溃。
	*/
	/*
	处理内存泄漏
	我们需要确保每个 goroutine都被终止
	终止goroutine 有以下几种方式：
	·当它完成了它的工作
	·因为不可恢复的错误，它不能工作
	·当它被告知需要终止工作
	*/

	doWork := func(
		done <-chan interface{},
		strings <-chan string,
		) <-chan interface{}{
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited")
			defer close(terminated)
			for {
				select {
				case s := <-strings:
					//做一些有意思的操作
					fmt.Println(s)
				case <-done:
					return
				}
			}
		}()
		return terminated
	}

	done := make(chan interface{})
	terminated := doWork(done,nil)

	go func() {
		//取消本操作
		time.Sleep(500*time.Millisecond)
		fmt.Println("canceling doWork goroutine...")
		close(done)
	}()

	<-terminated
	fmt.Println("done.")
	fmt.Println()

	/*
	可以看到，尽管我们给 strings <-chan string 传递了个nil ,但我们的goroutine仍然成功退出
	在此我们加入了两个goroutine (一个是doWork里的go func() ,一个是 main goroutine) ,但没有造成死锁。
	这是因为我们在加入两个goroutine的外，还创建了第三个goroutine ( go func() close(done) ) ，
	他用来在doWork执行取消doWork中的goroutine .
	我们已经成功消除了我们的goroutine泄漏
	*/

	//---------------------------------
	//防止阻塞写入
	newRandStream := func(done1 <-chan interface{}) <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited.")
			defer close(randStream)
			for{
				select {
				case randStream<-rand.Int():
				case <-done1:
					return
				}
			}
		}()
		return randStream
	}

	done1 := make(chan interface{})
	randStream := newRandStream(done1)
	fmt.Println("3 random ints:")
	for i := 1;i <= 3 ;i++  {
		fmt.Printf("%d: %d\n",i,<-randStream)
		if i == 2{
			close(done1)
		}
	}

}

//or-channel
func testorchannel(){
	/*
	有时候我们可能会希望将一个或多个完成的channel 合并到一个完成的channel中，
	该channel 在任何组件channel 关闭时关闭。编写一个执行这种耦合的选择语句是完全可以接受的，
	*/
	/*
	or-channel模型是通过递归和goroutine 创建一个复合done channel
	*/
	var or func(channels ...<-chan interface{}) <-chan interface{}
	or = func(channels ...<-chan interface{}) <-chan interface{}{
		fmt.Printf("len(channels) : %v\n",len(channels))
		switch len(channels) {
		case 0:
			return nil
		case 1:
			return channels[0]
		}

		orDone := make(chan interface{})
		go func() {
			defer close(orDone)

			switch len(channels) {
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default:
				select{
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[3:],orDone)...):
				}
			}
		}()
		return orDone
	}

	/*
	·or = func(channels ...<-chan interface{}) <-chan interface{}{
	在这里，我们有我们的函数，或者，它采用可变的channel切片并返回单个channel

	·case 0:
			return nil
	由于这是一个递归函数，我们必须设置终止标准。
	首先，如果可变切片是空的，我们只返回一个空channel。这是由于不传递channel的观点所产生的，我们希望复合的channel做不任何事情。

	·case 1:
			return channels[0]
	我们的第二个终止标准是如果我们的变量切片只包含一个元素，我们只返回该元素。

	·go func()
	这是函数的主体，以及递归发生的地方。我们创建了一个goroutine,以便我们可以不接受阻塞地等待我们channel上的消息。

	·switch len(channels) {
			case 2:
	基于我们进行迭代的方式，每一次迭代调用都将至少有两个channel。在这里我们为需要两个channel的情况采用了约束goroutine数目的优化方法。

	·default:
				select{
	在这里，我们在循环到我们存放所有channel的slice的第三个索引的时候，我们创建了一个 or-channel 并从这个channel中选择了一个。这将形成一个由现有slice的剩余部分组成的树并且返回第一个信号量。为了使在建立这个树的goroutine退出的时候在树下的goroutine也可以跟着退出，我们将这个 orDone channel也传递到了调用中

	·这是一个相当简洁的函数，使你可以将任意数量的channel组合到单个channel中，只要任何组件channel关闭或写入，该channel就会关闭。
	*/

	//下面是一个使用该功能的例子，它将多个经过一段时间后关闭的channel ,合并到一个关闭的channel中：
	sig := func(after time.Duration) <-chan interface{}{
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}
	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
		)
	fmt.Printf("done after %v\n",time.Since(start))
	fmt.Println()
	/*
	·sig := func(after time.Duration) <-chan interface{}{
	此功能只是创建一个channel ,当后续时间中指定的时间结束时，将关闭该channel
	·start := time.Now()
	在这里，我们大致追踪来自or 函数的channel 何时开始阻塞
	在这里，我们打印读取发生的时间

	*/

	/*
	请注意，尽管在我们的 or 调用中放置了多个channel 或需要不同时间才能关闭，
	但我们在1秒后关闭的channel 会导致由 or 调用创建的整个channel 关闭。
	这是因为尽管它位于树或函数构建的树种，它将始终关闭，因此依赖于其关闭的channel 也将关闭。
	*/
}

//错误处理
func testerror(){
	type Result struct {
		Response *http.Response
		Error error
	}
	checkStatus := func(
		done <-chan interface{},
		urls ...string,
		) <-chan Result {
		results := make(chan Result)
		go func() {
			defer close(results)
			for _,url := range urls{
				var result Result
				resp,err := http.Get(url)
				result = Result{Error:err,Response:resp}
				select {
				case <-done:
					return
				case results <-result:
				}
			}
		}()
		return results
	}
	done := make(chan interface{})
	defer close(done)
	urls := []string{"https://www.baidu.com","https://badhost"}
	for result := range checkStatus(done,urls...){
		if(result.Error != nil){
			fmt.Printf("error:%v\n",result.Error)
			continue
		}
		fmt.Printf("response: %v\n",result.Response.Status)
	}
	fmt.Println()
	/*
	·type Result struct {
	在这里，我们创建了一个包含 *http.Response 和从我们的goroutine中循环迭代中可能出现的错误的类型

	·checkStatus := func(
	该函数返回了一个可读取的channel ，以检索循环迭代的结果.

	·result = Result{Error:err,Response:resp}
	在这里，我们创建了一个Result实例，并设置错误和响应字段

	·case results <-result:
	这是我们将结果写入我们的channel 的地方

	·if(result.Error != nil){
	在这里，在我们的main goroutine 中，我们能够智能的处理由 checkStatus 启动的goroutine 中出现的错误
	*/

	/*
	错误处理的关键是我们如何将潜在的结果和潜在的错误结合起来。
	这表示从goroutine checkStatus 创建的完整可能的结果集，并且允许我们的主要常规关于发生错误时做什么的决定。
	从更广泛的角度来说，我们已经成功的将错误处理的担忧从我们的生产者 goroutine 中分离出来。
	这是可取的，因为生产goroutine 的goroutine（在这种情况下是我们的main goroutine ）具有更多关于正在运行的程序的上下文，并且可以做出关于如何处理错误的更明智的决定。
	*/

	//下面，我们可以尝试修改一下程序，以便在出现三个或更多错误时停止尝试检查状态：
	done2 := make(chan interface{})
	defer close(done2)

	errCount := 0
	urls2 := []string{"https://www.google.com","a","b","c","d"}
	for result := range checkStatus(done2,urls2...){
		if result.Error != nil {
			fmt.Printf("error: %v\n",result.Error)
			errCount++
			if errCount >= 3{
				fmt.Println("Too many errors,breaking!\n")
				break
			}
			continue
		}
		fmt.Printf("Response: %v\n",result.Response.Status)
	}
}

//管道
func testpipeline(){
	/*
	pipeline 是一系列将数据输入，执行操作并将结果数据传回的系统。
	我们将这些操作称为pipeline 的一个 stage

	通过使用pipeline，你可以分离每个 stage 的关注点
	你可以独立的修改每个stage ,
	也可以混合搭配stage 的组合方法，而无需修改stage ，你可以将每个stage 同时处理到上游或下游stage ，并且可以扇出或限制部分 pipeline
	*/

	//一个stage 只是将数据输入，对其进行转换并将数据返回的例子：
	multiply := func(values []int,multiplier int) []int {
		multipliedValues := make([]int,len(values))
		for i,v := range values{
			multipliedValues[i] = v*multiplier
		}
		return multipliedValues
	}

	//另一个stage
	add := func(values []int,additive int) []int {
		addedValues := make([]int,len(values))
		for i,v := range values{
			addedValues[i] = v + additive
		}
		return addedValues
	}

	//我们尝试将两个stage 合并：
	ints := []int{1,2,3,4}
	for _,v := range add(multiply(ints,2),1){
		fmt.Printf("add value: %v\n",v)
	}
	fmt.Println()

	//------------------------------------------------
	//pipeline与channel 结合
	generator := func(done <-chan interface{},integers ...int) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for _,i := range integers{
				select {
				case <- done:
					return;
				case intStream<-i:
				}
			}
		}()
		return intStream
	}

	multiter := func(
		done <-chan interface{},
		intStream <-chan int,
		multiplier int,
		) <-chan int {
		multiStream := make(chan int)
		go func() {
			defer close(multiStream)
			for i := range intStream{
				select {
				case <-done:
					return
				case multiStream<-i*multiplier:
				}
			}
		}()
		return multiStream
	}
	adder := func(
		done <-chan interface{},
		intStream <-chan int,
		additive int,
		)<-chan int {
		addStream := make(chan int)
		go func() {
			defer close(addStream)
			for i := range intStream{
				select {
				case <-done:
					return
				case addStream<-i+additive:
				}
			}
		}()
		return addStream
	}
	done := make(chan interface{})
	defer close(done)

	intStream := generator(done,1,2,3,4)
	pipeline := adder(done,multiter(done,intStream,2),1)
	for v := range pipeline {
		fmt.Printf("pipeline v: %v\n",v)
	}
	fmt.Println()
	/*
	generator函数将一组离散值转换为一个channel 上的数据流。
	适当的说，这种类型的功能称为生成器。
	在使用流水线时，你会经常看到这一点，因为在流水线开始时，你总是会有一些需要转换为channel 的数据。

	无论pipeline stage 所处的是什么（在等待传入的channel 的状态还是等待发送完成）
	关闭done channel 都将会迫使pipeline stage 被终止.
	*/

}

//生成器
func testgenerator(){
	//生成器是将一组离散值转换为channel 上的数据流的任何函数：
	repeat := func(
		done <-chan interface{},
		values ...interface{},
		) <-chan interface{}{
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for{
				for _,v := range values{
					select {
					case <-done:
						return
					case valueStream<-v:
					}
				}
			}
		}()
		return valueStream
	}

	take := func(
		done <-chan interface{},
		valueStream <-chan interface{},
		num int,
		) <-chan interface{}{
			takeStream := make(chan interface{})
			go func() {
				defer close(takeStream)
				for i:=0;i<num;i++{
					select {
					case <- done:
						return
					case takeStream<-<-valueStream:
					}
				}
			}()
			return takeStream
	}

	done := make(chan interface{})
	defer close(done)

	for num := range take(done,repeat(done,1),10){
		fmt.Printf("take num: %v\n",num)
	}
	fmt.Println()

	//重复函数生成器
	repeatFn := func(
		done <-chan interface{},
		fn func() interface{},
		)<-chan interface {}{
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for{
				select {
				case <-done:
					return
				case valueStream<-fn():
				}
			}
		}()
		return valueStream
	}

	done2 := make(chan interface{})
	defer close(done2)

	randFn := func() interface{}{ return rand.Int()}

	for num := range take(done2,repeatFn(done2,randFn),10){
		fmt.Printf("repeatFn num: %v\n",num)
	}
	fmt.Println()

	//转换字符串的pipeline
	toString := func(
		done <-chan interface{},
		valueStream <-chan interface{},
	)<-chan string{
		stringStream := make(chan string)
		go func() {
			defer close(stringStream)
			for v := range valueStream{
				select {
				case <-done:
					return
				case stringStream<-v.(string):
				}
			}
		}()
		return stringStream
	}

	done3 := make(chan interface{})
	defer close(done3)

	var message string
	for token := range toString(done,take(done,repeat(done,"I","am."),5)){
		message += token+" "
	}
	fmt.Printf("message: %s...\n",message)

}