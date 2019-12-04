package main

import (
	"fmt"
	"sync"
)

//----------------------------channel状态表--------------------------
/*
操作		Channel状态			结果
-----------------------------------------
Read		nil					阻塞
    		打开且非空			输出
    		打开且空			阻塞
    		关闭				<默认值>,false
    		只写				编译错误
------------------------------------------
Write		nil					阻塞
    		打开但填满			阻塞
    		打开且不满			写入
    		关闭				panic
    		只读				编译错误
------------------------------------------
Close		nil					panic
    		打开且非空			关闭Channel;读取成功，直到通道耗尽，然后读取生产者的默认值
    		打开且空			关闭Channel;读取生产者的默认值
    		关闭				panic
    		只读				编译错误


*/

func main() {
	fmt.Println("channelStream") //channel 通道
	declareChannel()
	fmt.Println()
	fmt.Println("testChan") //channel 通道
	testChan()
	fmt.Println()
	fmt.Println("bufferChan") //缓冲通道
	bufferChan()
	fmt.Println()
	fmt.Println("chanBelongto") //channel所有者
	chanBelongto()
}

//声明通道
func declareChannel(){
	//声明一个双向通道
	var dataStream chan interface{}
	dataStream = make(chan interface{})
	fmt.Printf("dataStream: %v\n",dataStream)

	//声明读通道
	var readStream <-chan interface{}
	readStream = make(<-chan interface{})
	fmt.Printf("readStream: %v\n",readStream)

	//声明写通道
	var writeStream chan<- interface{}
	writeStream = make(chan<- interface{})
	fmt.Printf("writeStream: %v\n",writeStream)

	//隐式声明
	var receiveChan <-chan interface{}
	var sendChan chan<- interface{}
	bothChan := make(chan interface{})

	receiveChan = bothChan
	sendChan = bothChan
	fmt.Printf("receiveChan: %v\n",receiveChan)
	fmt.Printf("sendChan: %v\n",sendChan)

}

//实例
func testChan(){
	//实例1
	stringChan := make(chan string)
	go func(){
		stringChan <- "Hello Channel"
	}()
	fmt.Println(<-stringChan)

	/*
	go语言中的 channel 是阻塞的
	意味着只有 channel 内的数据被消费后，新的数据才可以插入
	而任何试图从空的 channel 读取数据的 goroutine ，将会等待至少一条数据被写入 channel 后才能读到
	*/
	/*
	在上面例子中 fmt.Println会从stringChan 这个 channel 中消费一条数据，所以他会等待 channel 中有数据才开始消费
	同样，匿名的goroutine 试图往 stringChan 中写入一条数据，所以在写入数据之前，goroutine不会退出
	因此，main goroutine 和匿名的 goroutine 都被阻塞住
	*/

	stringStream := make(chan string)
	go func() {
		stringStream <- "hello channles"
	}()
	salutation,ok := <-stringStream
	fmt.Printf("ok:(%v) salutation:%v\n",ok,salutation)

	//第二个返回值是读取操作的一种方式，
	//用于表示该channel 上有新数据写入，
	//或者是由 closed channel 生成的默认值

	//下面是一个close channel的实例
	valueStream := make(chan interface{})
	close(valueStream)

	//从 close channel中读取数据
	intStream := make(chan int)
	close(intStream)
	integer,ok := <- intStream
	fmt.Printf("ok:(%v) integer:%v\n",ok,integer)

	//从 range channel中读取数据
	intChan := make(chan int)
	go func() {
		defer close(intChan)
		for i := 0;i < 5 ;i++  {
			intChan <- i
		}
	}()

	for integer := range intChan{
		fmt.Printf("integer: %v\n",integer)
	}

	//从WaitGroup中 close channel
	begin := make(chan interface{})
	var wg sync.WaitGroup
	for i := 0;i < 5 ;i++  {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-begin
			fmt.Printf("%v has begun\n",i)
		}(i)
	}
	fmt.Println("Unblocking goroutines...")
	close(begin)
	wg.Wait()


}


//buffer channel缓冲通道
func bufferChan(){
	var bufferChan chan interface{}
	bufferChan = make(chan interface{}, 4)
	bufferChan<-"A"
	bufferChan<-"B"
	bufferChan<-"C"
	bufferChan<-"D"
	//bufferChan<-"E" //提示 fatal error: all goroutines are asleep - deadlock!
	fmt.Println("bufferChan push all")

	/*
	我们创建了一个具有四个槽的缓冲通道
	第一次写入 ，数据会被放进第一个槽中
	第二次写入 ，数据会被放进第二个槽中
	第三次写入 ，数据会被放进第三个槽中
	第四次写入 ，数据会被放进第四个槽中
	缓冲通道满了，
	当第五次写入时，数据将会被阻塞，直到有一些goroutine 执行读取了缓冲通道的数据
	*/

	var bufferChan1 chan interface{}
	bufferChan1 = make(chan interface{},3)
	var wg1 sync.WaitGroup
	wg1.Add(1)
	go func() {
		wg1.Done()
		fmt.Printf("bufferChan1 read: %v\n",<-bufferChan1)
	}()
	bufferChan1<-"write 1"
	wg1.Wait()
	/*
	如果一个缓存通道是空的，并且有一个下游接收，那么缓冲通道将被忽略，并且该值将直接从发送方传递到接收方。
	*/

	intChan := make(chan int,4)
	go func() {
		defer close(intChan)
		defer fmt.Println("Producer done")
		for i := 0;i < 5 ;i++  {
			fmt.Printf("Sending:%v\n",i)
			intChan<-i
		}
	}()

	for integer := range intChan{
		fmt.Printf("Received:%v\n",integer)
	}
}

//正确配置channel
func chanBelongto(){
	/*
	分配channel的所有权
	把所有权定义为实例化、写入和关闭channel的goroutine.
	*/

	/*
	单向的channel声明是一种工具
	它将允许我们区分channel的所有者和channel的使用者
	channel的所有者对channel有一个写访问视图(chan或chan<-)
	而channel的使用者对channel有一个只读视图(<-chan)
	*/

	/*
	channel的所有者应具备如下：
	·实例化channel
	·执行写操作，或将所有权传递给另一个goroutine
	·关闭channel
	·压缩前面三件事，并通过一个只读channel 将他们暴露出来
	*/

	/*
	channel的消费者只需要担心两件事：
	知道channel是何时关闭的
	正确的处理阻塞
	*/

	chanOwner := func() <-chan int{
		resultStream := make(chan int,5)
		go func(){
			defer close(resultStream)
			for i := 0;i <= 5 ;i++  {
				resultStream <- i
			}
		}()
		return resultStream
	}

	resultOutside := chanOwner()
	for result := range resultOutside{
		fmt.Printf("Received: %d\n",result)
	}
	fmt.Println("Done receiving")

	/*
	resultStream := make(chan int,5) 实例化一个缓存channel.因为知道将产生6个结果，我们创建一个5个缓存的channel,这样 goroutine就能尽快完成
	go func()启动一个匿名的 goroutine ，他在resultStream上执行写操作。注意，我们已经在外围函数chanOwner 中封装了goroutine的创建
	defer close(resultStream) 确保一旦执行完成了goroutine, resultStream就会关闭。作为channel的所有者，这是我们必须做的
	return resultStream 在这里我们返回channel.由于返回值被声明为一个只读channel, 因此resultStream将隐式的转换为只读消费者
	for result := range resultOutside 遍历 resultOutside。作为消费者，我们只关心阻塞和channel的关闭
	*/

}