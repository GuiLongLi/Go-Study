package main

import (
	"log"
	"time"
	"os"
)

/*
治愈异常的goroutine：
在长期运行的后台程序中，经常会有一些长时间运行的goroutine。
这些goroutine 经常处于阻塞状态，等待数据以某种方式到达，然后唤醒它们，进行一些处理，再返回一些数据。
如果没有外部干预，一个goroutine 很容易进入一个不正常的状态，并且无法恢复。
我们需要建立一个机制来监控你的goroutine 是否处于健康的状态是很有用的，
当他们变得异常时，就可以尽快重启。

我们将这种重启goroutine的过程称为治愈。

为了治愈goroutine ，我们需要使用心跳模式来检查我们正在监控的goroutine 是否活跃。
心跳的类型取决于你想要监控的内容，但是如果你的goroutine 有可能会产生活锁，确保心跳包含某些信息，表明该goroutine 在正常的工作而不仅仅是活着。
我们把监控goroutine 的健康这段逻辑称为管理员，它监视一个管理区的goroutine，
如果有goroutine 变得不健康，管理员将负责重新启动这个管理区的goroutine 。

*/

func main() {
	log.Println("testhealing")
	testhealing()
	log.Println()
	log.Println("testhealing2")
	testhealink2()
}


func testhealing()  {
	var or func(channels ...<-chan interface{}) <-chan interface{}
	or = func(channels ...<-chan interface{}) <-chan interface{}{
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


	type startGoroutineFn func(
		done <-chan interface{},
		pulseInterval time.Duration,
		)(heartbeat <-chan interface{})
	
	newSteward := func(
		timeout time.Duration,
		startGoroutine startGoroutineFn,
		) startGoroutineFn {
			return func(
				done <-chan interface{},
				pulseInterval time.Duration,
				) (<-chan interface{}) {
				heartbeat := make(chan interface{})
				go func() {
					defer close(heartbeat)

					var wardDone chan interface{}
					var wardHeartbeat <-chan interface{}
					startWard := func() {
						wardDone = make(chan interface{})
						wardHeartbeat = startGoroutine(or(wardDone,done),timeout/2)
					}
					startWard()
					pulse := time.Tick(pulseInterval)

					monitorLoop:
						for{
							timeoutSignal := time.After(timeout)
							for {
								select {
								case <-pulse:
									select {
									case heartbeat<- struct {}{}:
									default:
									}
								case <-wardHeartbeat:
									continue monitorLoop
								case <-timeoutSignal:
									log.Println("steward:ward unhealthy;restarting")
									close(wardDone)
									startWard()
									continue monitorLoop
								case <-done:
									return
								}
							}
						}
				}()
				return heartbeat
			}
		
	}

	/*
	· pulseInterval time.Duration,
	   )(heartbeat <-chan interface{})
	在这里，我们定义一个可以监控和重启的goroutine 的信号，
	我们看到了熟悉的channel ，以及来自心跳模式的pulseInterval 和heartbeat

	·timeout time.Duration,
	  startGoroutine startGoroutineFn,
	  ) startGoroutineFn {
	在这，我们看到一个管理员监控goroutine 需要timeout 变量，
	还有一个函数startGoroutine 来启动他监控的goroutine
	管理员本身会返回一个startGoroutineFn ，表示管理员本身也是可监控的

	·startWard := func() {
	在这里，我们定义了一个闭包，他实现了一个统一的startWard 方法来启动我们正在监视的goroutine

	·wardDone = make(chan interface{})
	这是我们创建的一个新的channel ,如果我们需要发出一个停止的信号，就会通过他传入goroutine 中

	·wardHeartbeat = startGoroutine(or(wardDone,done),timeout/2)
	在这，我们启动将要监控的goroutine ，如果管理员被停止了，或者管理员想要停止goroutine
	我们希望这些信息都能传递给管理区里的goroutine ，所有我们把两个done channel 用逻辑 or 包装一下。
	我们设定心跳间隔时间是超时时间的一半

	·timeoutSignal := time.After(timeout)
	  for {
	这里我们的内部循环，他确保管理员可以发出自己的心跳

	·case <-wardHeartbeat:
	在这，如果我们收到goroutine 的心跳，将继续执行监控的循环

	·case <-timeoutSignal:
	这里如果我们在暂停期间没有收到管理区里goroutine 的心跳，
	我们会要求管理区里的goroutine 停下来，并启动一个新的goroutine ，然后我们继续监控

	*/

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime|log.LUTC)

	doWork := func(
		done <-chan interface{},
		_ time.Duration,
		) <-chan interface{}{
			log.Println("ward:hello, i,m irresponsible!")
			go func() {
				<-done
				log.Println("ward:i am halting.")
			}()
			return nil
	}
		doWorkWithSteward := newSteward(1*time.Second,doWork)

		done := make(chan interface{})
		time.AfterFunc(5*time.Second, func() {
			log.Println("main:halting steward and ward.")
			close(done)
		})

		for range doWorkWithSteward(done,2*time.Second){

		}
		log.Println("Done")

		/*
		·go func() {
		  <-done
		在这里，我们看到这个goroutine 没有做任何事情，只是等待被取消，他也没有发出任何心跳

		·doWorkWithSteward := newSteward(4*time.Second,doWork)
		这里创建了一个函数，为goroutine doWork 创建一个管理员，我们设置doWork 的超时时间为4秒

		·time.AfterFunc(9*time.Second, func() {
		在这我们设置9秒后停止管理员和goroutine ，这样我们的测试就会结束。

		·for range doWorkWithSteward(done,4*time.Second){
		最后，我们启动管理员并在其心跳范围内防止我们的测试停止

		*/

}

func testhealink2(){
	var or func(channels ...<-chan interface{}) <-chan interface{}
	or = func(channels ...<-chan interface{}) <-chan interface{}{
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
	take := func(
		done <-chan interface{},
		intStream <-chan interface{},
		num int,
		) <-chan interface {}{
			takeStream := make(chan interface{})
			go func() {
				defer close(takeStream)
				for i :=0;i<num;i++{
					select {
					case <-done:
						return
					case takeStream<-<-intStream:
					}
				}
			}()
			return takeStream
	}


	type startGoroutineFn func(
		done <-chan interface{},
		pulseInterval time.Duration,
	)(heartbeat <-chan interface{})

	newSteward := func(
		timeout time.Duration,
		startGoroutine startGoroutineFn,
	) startGoroutineFn {
		return func(
			done <-chan interface{},
			pulseInterval time.Duration,
		) (<-chan interface{}) {
			heartbeat := make(chan interface{})
			go func() {
				defer close(heartbeat)

				var wardDone chan interface{}
				var wardHeartbeat <-chan interface{}
				startWard := func() {
					wardDone = make(chan interface{})
					wardHeartbeat = startGoroutine(or(wardDone,done),timeout/2)
				}
				startWard()
				pulse := time.Tick(pulseInterval)

			monitorLoop:
				for{
					timeoutSignal := time.After(timeout)
					for {
						select {
						case <-pulse:
							select {
							case heartbeat<- struct {}{}:
							default:
							}
						case <-wardHeartbeat:
							continue monitorLoop
						case <-timeoutSignal:
							log.Println("steward:ward unhealthy;restarting")
							close(wardDone)
							startWard()
							continue monitorLoop
						case <-done:
							return
						}
					}
				}
			}()
			return heartbeat
		}

	}

	doWorkFn := func(
		done <-chan interface{},
		intList ...int,
		) (startGoroutineFn,<-chan interface{}){
			intChanStream := make(chan (<-chan interface{}))
			intStream := bridge(done,intChanStream)

			doWork := func(
				done <-chan interface{},
				pulseInterval time.Duration,
				) <-chan interface {}{
					intStream := make(chan interface{})
					heartbeat := make(chan interface{})
					go func() {
						defer close(intStream)
						select {
						case intChanStream <- intStream:
						case <-done:
							return
						}
						pulse := time.Tick(pulseInterval)

						for{
						valueLoop:
							for _,intVal := range intList{
								if intVal < 0{
									log.Printf("negative value: %v\n",intVal)
									return
								}
								for{
									select {
									case <- pulse:
										select {
										case heartbeat<- struct {}{}:
										default:
										}
									case intStream<-intVal:
										continue valueLoop
									case <-done:
										return
									}
								}
							}
						}
					}()
					return heartbeat
			}
		return doWork,intStream
	}

	/*
	·) (startGoroutineFn,<-chan interface{}){
	在这里，我们填入一些我们管理区所需的参数，并返回我们管理区用来通信的channel

	·intChanStream := make(chan (<-chan interface{}))
	这里我们创建了作为桥接模式一部分的channel

	·) <-chan interface {}{
	我们创建一个将被管理员启动和监控的闭包

	·intStream := make(chan interface{})
	这是我们实例化channel 的地方，我们将利用这些channel 与管理区中的goroutine 通信

	·case intChanStream <- intStream:
	这里我们把我们即将用来通信的channel 通知给bridge

	·if intVal < 0{
	  log.Printf("negative value: %v\n",intVal)
	当我们处理到负数时，在这里打印出一个错误信息，然后从goroutine 中返回
	*/

	log.SetFlags(log.Ltime|log.LUTC)
	log.SetOutput(os.Stdout)

	done := make(chan interface{})
	defer close(done)

	doWork,intStream := doWorkFn(done,1,2,-1,3,4,5)
	doWorkWithSteward := newSteward(1*time.Millisecond,doWork)
	doWorkWithSteward(done,1*time.Hour)

	for intVal := range take(done,intStream,6){
		log.Printf("Received:%v\n",intVal)
	}

	/*
	·doWork,intStream := doWorkFn(done,1,2,-1,3,4,5)
	这里我们创建管理区的函数，允许他结束我们的可变整数切片，并返回一个用来返回的流

	·doWorkWithSteward := newSteward(1*time.Millisecond,doWork)
	我们创建一个管理员，用来监听doWork 闭包，
	因为我们希望能尽快知道失败的信息，所以我们将监听时间间隔设置为1 毫秒

	·doWorkWithSteward(done,1*time.Hour)
	我们通知管理员启动管理区并开始监控

	·for intVal := range take(done,intStream,6){
	最后，使用我们开发的一段管道代码，并从intStream 中取出前6 个值
	*/

}
