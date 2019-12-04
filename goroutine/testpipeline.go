package main

import (
	"fmt"
	"time"
	"sync"
	"runtime"
	"math/rand"
)

func main() {
	fmt.Println("pipe1")
	pipe1()
	fmt.Println()
	fmt.Println("ordonechannel")
	ordonechannel()
	fmt.Println()
	fmt.Println("teechannel")
	teechannel()
	fmt.Println()
	fmt.Println("bridgechannel")
	bridgechannel()
}

func pipe1(){
	repeatFn := func(
		done <-chan interface{},
		fn func() interface{},
		)<-chan interface{}{
		numStream := make(chan interface{})
		go func() {
			defer close(numStream)
			for{
				select {
				case <-done:
					return
				case numStream<-fn():
				}
			}
		}()
		return numStream
	}
	toInt := func(
		done <-chan interface{},
		valueStream <-chan interface{},
		) <-chan int{
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for v := range valueStream{
				select {
				case <- done:
					return
				case intStream<-v.(int):
				}
			}
		}()
		return intStream
	}

	//查找素数
	primeFinder := func(
		done <-chan interface{},
		randStream <-chan int,
		) <-chan interface{}{
		primeStream := make(chan interface{})
		go func() {
			defer close(primeStream)
			for integer := range randStream {
				if integer <= 1{
					continue
				}
				prime := true
				for divisor := integer -1;divisor > 1;divisor--{
					if integer%divisor == 0{
						prime = false
						break
					}
				}
				if prime {
					select {
					case <-done:
						return
					case primeStream<-integer:
					}
				}
			}
		}()
		return primeStream
	}

	take := func(
		done <-chan interface{},
		primeStream <-chan interface{},
		num int,
		) <-chan interface{}{
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i := num;i > 0;i--{
				select {
				case <-done:
					return
				case takeStream<-<-primeStream:
				}
			}
		}()
		return takeStream
	}

	rand.Seed(time.Now().UnixNano()) //利用当前时间的UNIX时间戳初始化rand包
	rand := func() interface {}{return rand.Intn(5000000)}

	done := make(chan interface{})
	defer close(done)

	start1 := time.Now()

	randIntStream := toInt(done,repeatFn(done,rand))
	fmt.Println("Primes:")
	for prime := range take(done,primeFinder(done,randIntStream),10){
		fmt.Printf("prime1 %d\n",prime)
	}
	fmt.Printf("search1 took: %v\n", time.Since(start1))
	fmt.Println()

	//-----------------------------
	//优化stage ,使用多个CPU
	/*
	因为我们有多个cpu 核心，我们可以启动 primeFinder 这个stage的许多副本
	现在我们将有 cpu 数量的 goroutine 从随机数发生器中拉出并试图确定数字是否为素数。
	*/
	//扇入意味着将多个数据流复用或合并成一个流：
	fanIn := func(
		done <-chan interface{},
		channels ...<-chan interface{},
		) <-chan interface{}{
		var wg sync.WaitGroup
		multiStream := make(chan interface{})
		multi := func(c <- chan interface{}) {
			defer wg.Done()
			for i := range c {
				select {
				case <- done:
					return
				case multiStream<-i:
				}
			}
		}
		// 从所有channel里取值
		wg.Add(len(channels))
		for _,c := range channels{
			go multi(c)
		}
		// 等待所有的读操作结束
		go func() {
			wg.Wait()
			close(multiStream)
		}()
		return multiStream
	}

	done2 := make(chan interface{})
	defer close(done2)

	start2 := time.Now()

	numFinders := runtime.NumCPU()
	fmt.Printf("Spinning up %d prime2 finders\n",numFinders)
	finders := make([]<-chan interface{},numFinders)
	for i := 0;i < numFinders ;i++  {
		finders[i] = primeFinder(done,randIntStream)
	}
	for prime := range take(done,fanIn(done,finders...),10){
		fmt.Printf("prime2 %d\n",prime)
	}
	fmt.Printf("search2 took: %v\n",time.Since(start2))


}

// or-done-channel
func ordonechannel(){
	/*
	有时候，我们在处理来自系统各个分散部分的channel时，并不知道我们的 goroutine是否被取消。
	我们需要用channel 的select 语句来包装我们的读操作，并从中选择正确的channel
	*/
	total := 3 //共有3个 valueStream
	channels := func(done <-chan interface{}) <-chan interface{}{
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for i := 0;i <= total;i++{
				select {
				case <- done:
					return
				case valueStream<-i:
				}
			}
		}()
		return valueStream
	}
	orDone := func(
		done <-chan interface{},
		channel <-chan interface{},
		)<-chan interface{}{
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			i := 0
			for{
				//当到达 i = 2 channel 时，我们提前取出 channel 的内容 ，channel 就会关闭了
				//这样后面的select case 取出 channel 时， ok 就会接收到 false
				if i == 2 {
					<-channel
				}
				i++
				select {
				case <-done:
					return
				case v,ok := <-channel:
					fmt.Printf("v %v -- ok %v\n",v,ok)
					//如果是已关闭的 channel ，我们就不用进行操作了
					if ok == false{
						return
					}
					select {
					case <-done:
					case valueStream<-v:
					}
				}
			}
		}()
		return valueStream
	}

	done := make(chan interface{})
	defer close(done)

	for val := range orDone(done,channels(done)){
		fmt.Printf("orDone %v\n",val)
	}
}

//tee-channel
func teechannel(){
	/*
	tee-channel 可以将一个传递的读channel ,返回两个单独的channel ,以获取两个相同的值：
	*/
	orDone := func(
		done <-chan interface{},
		intStream <-chan interface{},
	) <-chan interface {}{
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for {
				select {
				case <-done:
					return
				case val,ok := <-intStream:
					if ok == false {
						return
					}
					valueStream<-val
				}
			}
		}()
		return valueStream
	}

	tee := func(
		done <- chan interface{},
		in <-chan interface{},
		)(_,_ <-chan interface{}){
		out1 := make(chan interface{})
		out2 := make(chan interface{})
		go func() {
			defer close(out1)
			defer close(out2)
			for val := range orDone(done,in) {
				var out1, out2 = out1,out2
				for i := 0;i < 2 ;i++  {
					select {
					case <- done:
					case out1<-val:
						out1 = nil
					case out2<-val:
						out2 = nil
					}
				}
			}
		}()
		return out1,out2
	}

	repeat := func(
		done <-chan interface{},
		nums ...int,
		) <-chan interface{}{
		intStream := make(chan interface{})
		go func() {
			defer close(intStream)
			for{
				for num := range nums{
					select {
					case <-done:
						return
					case intStream<-num:
					}
				}
			}
		}()
		return intStream
	}

	take := func(
		done <-chan interface{},
		intStream <-chan interface{},
		num int,
		) <-chan interface{}{
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i := 0;i < num;i++{
				select {
				case <-done:
					return
				case takeStream<-<-intStream:
				}
			}
		}()
		return takeStream
	}

	done := make(chan interface{})
	defer close(done)

	out1,out2 := tee(done,take(done,repeat(done,1,2),4))

	for val1 := range out1{
		fmt.Printf("out1: %v, out2: %v\n",val1,<-out2)
	}
	/*
	·var out1, out2 = out1,out2
	我们将要使用的是out1 和 out2 的本地版本，所有我们会隐藏这些变量。

	·for i := 0;i < 2 ;i++  {
	我们将使用一条select 语句，以便不阻塞的写入out1 和 out2
	为确保两个都被写入 ,我们将执行select 语句的两次迭代：每个出站一个 channel

	·out1 = nil out2 = nil
	一旦我们写入了channel ,我们将其影副本设置为nil ,以便进一步阻塞写入，而另一个channel 可以继续

	注意写入out1 和out2 是紧密耦合的。
	知道out1 和out2 都被写入，迭代才能继续。
	*/
}

//桥接channel 模式
func bridgechannel(){
	/*
	我们将一个充满channel 的channel 拆解为一个简单的channel ，称为桥接channel
	*/
	orDone := func(
		done <-chan interface{},
		valueStream <-chan interface{},
		)<-chan interface{}{
		doneStream := make(chan interface{})
		go func() {
			defer close(doneStream)
			for{
				select {
				case <-done:
					return
				case val,ok := <-valueStream:
					if ok == false{
						return
					}
					doneStream<-val
				}
			}
		}()
		return doneStream
	}
	//桥接channel
	bridge := func(
		done <-chan interface{},
		chanStream <-chan <-chan interface{},
		)<-chan interface{}{
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for{
				var stream <-chan interface{}
				select{
				case maybeStream,ok := <-chanStream:
					if ok == false{
						return
					}
					stream = maybeStream
				case <-done:
					return
				}
				for val := range orDone(done,stream){
					select{
					case <-done:
					case valueStream <-val:
					}
				}
			}
		}()
		return valueStream
	}
	/*
	·valueStream := make(chan interface{})
	这是将返回的bridge 中的所有值的channel

	·for{
	这个循环负责从chanStream 中提取channel 并将其提供给嵌套循环来使用

	·for val := range orDone(done,stream){
	该循环负责读取已经给出的channel 中的值，并将这些值重复到valueStream中。
	当我们当前正在循环的流关闭时，我们从执行从此channel 读取的循环中跳出,并继续循环的下一次迭代，
	选择要读取的channel ，这为我们提供了一个不间断的结果值的流。
	*/
	generator := func(max int)<-chan <-chan interface{}{
		chanStream := make(chan (<-chan interface{}))
		go func() {
			defer close(chanStream)
			for i:=0;i<max;i++{
				stream := make(chan interface{},1)
				stream<-i
				close(stream)
				chanStream<-stream
			}
		}()
		return chanStream
	}
	done := make(chan interface{})
	defer close(done)
	for v := range bridge(done,generator(10))  {
		fmt.Printf("bridge %v\n",v)
	}

}



