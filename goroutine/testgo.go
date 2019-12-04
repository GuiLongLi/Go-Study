package main

import (
	"math"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	"bytes"
	"runtime"
)

//Go语言并发之道
func main() {
	fmt.Println("goroutineexample")
	goroutineexample()
	fmt.Println()
	fmt.Println("syncexample")
	syncexample()
	fmt.Println() //死锁
	fmt.Println("deathLock") //死锁
	//deathLock()
	fmt.Println()
	fmt.Println("liveLock") //活锁
	liveLock()
	fmt.Println()
	fmt.Println("hungry") //饥饿
	hungry()
	fmt.Println()
	fmt.Println("goroutineMemory")//内存测试
	goroutineMemory()
	fmt.Println()
	fmt.Println("mutexLock") //互斥锁
	mutexLock()
	fmt.Println()
	fmt.Println("mutexLock") //读写锁
	remutexLock()
	fmt.Println()
	fmt.Println("condSignal") //cond同步机制 Singal()
	condSignal()
	fmt.Println()
	fmt.Println("condBroadcast") //cond同步机制 Broadcast()
	condBroadcast()
	fmt.Println()
	fmt.Println("syncPool") //syncPool 池
	syncPool()
}

//并发例子
func goroutineexample(){
	var data int

	//使用go关键字 并发处理 data++
	go func () {
		data++
	}()

	if data == 0{
		fmt.Printf("the value is %v\n",data)
	}

	//可能出现的结果
	// 一. 不打印任何东西
	// data++ 先运行 -> if 再运行 -> fmt 最后运行

	// 三. 打印0
	// if 先运行 -> fmt 再运行 -> data++ 最后运行

	// 二. 打印1
	// if 先运行 -> data++ 再运行 -> fmt 最后运行
}

//内存访问同步
func syncexample(){
	//内存访问数据
	var memoryAccess sync.Mutex
	var value int

	go func(){
		//声明内存锁，独占内存访问
		memoryAccess.Lock()
		//有了锁，别的地方就不能访问value了
		value++
		//解除内存锁，释放内存访问
		memoryAccess.Unlock()
	}()

	//声明内存锁，独占内存访问
	memoryAccess.Lock()
	// 有了锁，别的地方就不能访问value了
	if value == 0 {
		fmt.Printf("the value is %v\n",value)
	}else{
		fmt.Printf("the value is %v\n",value)
	}
	//解除内存锁，释放内存访问
	memoryAccess.Unlock()

	//上面使用内存锁的访问 ，保护变量的内存访问同步
	//但会使我们的程序变慢
	//每次执行这些 Lock() 操作时，我们的程序就会暂停一段时间
}

//一个简单的死锁
func deathLock(){
	type value struct {
		mu sync.Mutex
		value int
	}

	var wg sync.WaitGroup
	printSum := func(v1,v2 *value){
		defer wg.Done()
		v1.mu.Lock()
		defer v1.mu.Unlock()

		time.Sleep(2*time.Second)
		v2.mu.Lock()
		defer v2.mu.Unlock()

		fmt.Printf("sum=%s\n",v1.value+v2.value)
	}

	var a,b value
	wg.Add(2)
	go printSum(&a,&b)
	go printSum(&b,&a)
	wg.Wait()

	//出现死锁了
	//fatal error: all goroutines are asleep - deadlock!
	/*
	本质上，我们创建了两个不能转动的齿轮
	第一次调用printSum锁定了a ,然后试图锁定b
	但在此期间
	第二次调用printSum锁定了b ,并试图锁定a
	这使得两个 goroutine 都无限地等待着
	------------------------
	coffman列举了死锁的条件：
	相互排斥
		并发进程同时拥有资源的独占权
	等待条件
		并发进程必须同时拥有一个资源，并等待额外的资源
	没有抢占
		并发进程拥有的资源只能被该进程释放，即可满足这个条件
	循环等待
		一个并发进程P1 必须等待一系列其他并发进程P2 ，
		其他并发进程P2 同时也在等待P1
	---------------------
	这里都满足了死锁条件：
	printSum函数需要a 和 b 的独占权
	printSum函数持有a 或 b 并且等待另外一个
	我们没有任何函数让我们的 goroutine被抢占
	第一次调用printSum正等待我们第二次调用 ,同时第二次调用也在等待第一次调用

	*/
}

//活锁
func liveLock(){
	cadence := sync.NewCond(&sync.Mutex{})
	go func(){
		for range time.Tick(1*time.Millisecond){
			cadence.Broadcast()
		}
	}()

	takeStep := func(){
		cadence.L.Lock()
		cadence.Wait()
		cadence.L.Unlock()
	}

	tryDir := func(dirName string,dir *int32, out *bytes.Buffer) bool{
		fmt.Fprintf(out, " %v",dirName)
		atomic.AddInt32(dir,1)
		takeStep()
		if atomic.LoadInt32(dir) == 1{
			fmt.Fprintf(out, ". Success!")
			return true
		}
		takeStep()
		atomic.AddInt32(dir,-1)
		return false
	}

	var left,right int32
	tryLeft := func(out *bytes.Buffer) bool {return tryDir("left",&left,out)}
	tryRight := func(out *bytes.Buffer) bool {return tryDir("right",&right,out)}

	walk := func(walking *sync.WaitGroup,name string){
		var out bytes.Buffer
		defer func() {fmt.Println(out.String())}()
		defer walking.Done()
		fmt.Fprintf(&out,"%v is trying to scoot:",name)
		for i := 0;i < 5 ;i++  {
			if tryLeft(&out) || tryRight(&out) {
				return
			}
		}
		fmt.Fprintf(&out,"\n%v tosses her hands up in exasperation!",name)
	}

	var peopleInHallway sync.WaitGroup
	peopleInHallway.Add(2)
	go walk(&peopleInHallway,"Alice")
	go walk(&peopleInHallway,"Barbara")
	peopleInHallway.Wait()

	//活锁的情形类似：
	/*
	一个走廊中，你走向她，两个人走在同一边
	她移动到另一边想让你通过，但你也移动到了另一边
	后来你再次移动回原来位置，她也移动到原来位置
	这种情况永远持续下去
	就成了活锁
	*/
}

//饥饿
func hungry(){
	var wg sync.WaitGroup
	var shareLock sync.Mutex
	const runtime = 1*time.Second

	//begin := time.Now() 获取当前时间
	//time.Since(begin) 计算程序运行时间
	// time.Second    表示1秒
	// time.Millisecond    表示1毫秒
	// time.Microsecond    表示1微妙
	// time.Nanosecond    表示1纳秒
	greedyWorker := func(){
		defer wg.Done()

		var count int
		for begin := time.Now();time.Since(begin) <= runtime;{
			shareLock.Lock()
			time.Sleep(3*time.Nanosecond)
			shareLock.Unlock()
			count++
		}
		fmt.Printf("Greedy worker was able to execute %v work loops\n",count)
	}

	politeWorker := func() {
		defer wg.Done()

		var count int
		for begin := time.Now();time.Since(begin) <= runtime;{
			shareLock.Lock()
			time.Sleep(1*time.Nanosecond)
			shareLock.Unlock()

			shareLock.Lock()
			time.Sleep(1*time.Nanosecond)
			shareLock.Unlock()

			shareLock.Lock()
			time.Sleep(1*time.Nanosecond)
			shareLock.Unlock()

			count++
		}
		fmt.Printf("Polite worker was able to execute %v work loops\n",count)
	}

	//添加两个并发
	wg.Add(2)
	go greedyWorker()
	go politeWorker()
	wg.Wait()
	//如果内存足够，将会输出以下
	/*----------------
	Polite  worker  was  able  to  exe cute  289777  work  loops.
	Greedy worker  was  able  to  execute  471287  work  loops.
	贪婪的worker会贪婪的抢占共享锁，以完成整个工作，
	而平和的worker则试图只在需要时锁定。
	两种worker都做同样多的模拟工作 sleeping时间都是 3ns
	但在同样的时间里，贪婪的worker工作量几乎是平和的worker工作量的两倍
	*/

	//如果内存不足，有可能会出现下面的状况
	/*-----------------
	Polite worker was able to execute 53352 work loops
	Greedy worker was able to execute 983 work loops
	这是因为同步访问内存是昂贵的，以至于当内存不足时，一个并发进程，有可能会饿死其他并发进程
	意思是内存不足时，一个并发进程占用了大部分内存，其他并发进程就有可能被阻止或者饿死
	*/

	// 我们还应该考虑来自于外部过程的饥饿。
	// 请记住，饥饿也应用于CPU、内存、文件句柄、数据库连接：任何必须共享的资源都是饥饿的候选者

}

//测算 goroutine运行前后的内存数量
func goroutineMemory(){
	memConsumed := func() uint64{
		runtime.GC()
		var s runtime.MemStats
		runtime.ReadMemStats(&s)
		return s.Sys
	}

	var c <- chan interface{}
	var wg sync.WaitGroup
	noop := func(){
		wg.Done()
		<-c
	}

	const numGoroutines = 1e4
	wg.Add(numGoroutines)
	before := memConsumed()
	for i := numGoroutines;i > 0;i--{
		go noop()
	}
	wg.Wait()
	after := memConsumed()
	fmt.Printf("%.3fkb\n",float64(after-before)/numGoroutines/1000)
}

//互斥锁
func mutexLock(){
	var count int
	var lock sync.Mutex

	increment := func(){
		lock.Lock()
		defer lock.Unlock()
		count++
		fmt.Printf("Incrementing: %d\n",count)
	}

	decrement := func(){
		lock.Lock()
		defer lock.Unlock()
		count--
		fmt.Printf("Decrement: %d\n",count)
	}

	//增量
	var arithmetic sync.WaitGroup
	for i := 0;i <= 5;i++  {
		arithmetic.Add(1)
		go func(){
			defer arithmetic.Done()
			increment()
		}()
	}

	//减量
	for i := 0;i <= 5 ;i++  {
		arithmetic.Add(1)
		go func() {
			defer arithmetic.Done()
			decrement()
		}()
	}

	arithmetic.Wait()
	fmt.Println("arithmetic complete")

}

//读写锁
func remutexLock(){
	producer := func(wg *sync.WaitGroup,lock sync.Locker){
		defer wg.Done()
		for i := 5;i > 0 ;i--  {
			lock.Lock()
			lock.Unlock()
			time.Sleep(1)
		}
	}

	observer := func(wg *sync.WaitGroup,lock sync.Locker){
		defer wg.Done()
		lock.Lock()
		defer lock.Unlock()
	}

	test := func(count int, mutex,rwMutex sync.Locker) time.Duration {
		var wg sync.WaitGroup
		wg.Add(count+1)
		beginTestTime := time.Now()
		go producer(&wg,mutex)
		for i := count;i > 0 ;i--  {
			go observer(&wg,rwMutex)
		}
		wg.Wait()
		return time.Since(beginTestTime)
	}

	var m sync.RWMutex
	fmt.Printf("Readers\tRWMutex\t\t Mutex\n")
	for i := 0;i < 10 ;i++  {
		count := int(math.Pow(2,float64(i)))
		fmt.Printf(
			"%d\t%v\t %v\n",
			count,
			test(count,&m,m.RLocker()),
			test(count,&m,&m),
			)
	}
}

//cond同步机制 Signal() ,Signal()会发现等待最长时间的goroutine并通知它
func condSignal(){
	//1.我们使用标准的sync.Mutex 作为锁
	cond := sync.NewCond(&sync.Mutex{})
	//2.创建一个长度为0的切片，最大容量是10
	queue := make([]interface{},0,10)

	removeFromQueue := func(delay time.Duration){
		time.Sleep(delay)
		//8.进入临界区，以便我们可以修改与条件相关的数据
		cond.L.Lock()
		//9.删除排在第一的queue元素
		queue = queue[1:]
		fmt.Println("removed from queue")
		//10.退出临界区
		cond.L.Unlock()
		//然一个正在登陆的goroutine知道发生了什么事情
		cond.Signal()
	}

	for i := 0 ;i < 10 ;i++  {
		//3.调用锁进入临界区
		cond.L.Lock()
		//4.检查循环queue 的长度，判断是否足够两个
		for len(queue) == 2 {
			//5.调用Wait,暂停main goroutine 直到一个信号的条件已发送
			cond.Wait()
		}
		fmt.Println("Adding to queue")
		//这里插入queue元素
		queue = append(queue,struct{}{})
		//6.创建一个新的goroutine ，它会在一秒后删除一个queue元素
		go removeFromQueue(500*time.Millisecond)
		//7.退出条件的临界区
		cond.L.Unlock()
	}

	//最终输出
	/*
	Adding to queue
	Adding to queue
	removed from queue
	Adding to queue
	removed from queue
	Adding to queue
	removed from queue
	Adding to queue
	removed from queue
	Adding to queue
	removed from queue
	Adding to queue
	removed from queue
	Adding to queue
	removed from queue
	Adding to queue
	removed from queue
	Adding to queue
	*/

	//逻辑顺序
	/*
	for循环10次
	循环时插入queue元素,并调用了 removeFromQueue
	当queue的长度等于2时
	暂停循环
	直到 removeFromqueue 停止了一秒后， 发出了一个 signal信号
	for继续循环
	直到循环结束
	*/
}

//cond同步机制 Broadcast() ,Broadcast()会向所有等待的goroutine并通知它们
func condBroadcast(){
	//声明一个Button 构造体 ，包含一个 Cond的点击事件
	type Button struct {
		Clicked *sync.Cond
	}
	button := Button{Clicked: sync.NewCond(&sync.Mutex{})}

	//事件处理  fn是回调函数
	subscribe := func(cond *sync.Cond,fn func()){
		//创建一个等待组
		var goroutining sync.WaitGroup
		goroutining.Add(1)
		go func() {
			goroutining.Done()
			cond.L.Lock()
			defer cond.L.Unlock()
			cond.Wait()  //暂停 go func ,并等待 button.Clicked.Broadcast() 发送通知
			fn() //调用回调函数
		}()
		//暂停 subscribe 并等待 go func 完成 goroutining.Done()
		goroutining.Wait()
	}

	/*
	WaitGroup顾名思义，就是用来等待一组操作完成的。
	WaitGroup内部实现了一个计数器，用来记录未完成的操作个数，它提供了三个方法，
	Add()用来添加计数。
	Done()用来在操作结束时调用，使计数减一。
	Wait()用来等待所有的操作结束，即计数变为0，
	该函数会在计数不为0时等待，在计数为0时立即返回。
	*/
	//创建一个等待组
	var clickRegistered sync.WaitGroup
	clickRegistered.Add(3)
	subscribe(button.Clicked,func(){
		fmt.Println("Maximizing window")
		clickRegistered.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("Displaying annoying dialog box")
		clickRegistered.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("Mouse clicked")
		clickRegistered.Done()
	})
	fmt.Println("Clicked.Broadcast case")
	//cond.Broadcast()向所有等待的 goroutine发送通知
	button.Clicked.Broadcast()

	fmt.Println("clickRegistered.Wait before")
	//暂停main goroutine，等待完成 计时器为0
	clickRegistered.Wait()
	fmt.Println("clickRegistered.Wait after")

	//显示结果
	/*
	Clicked.Broadcast case
	clickRegistered.Wait before
	Mouse clicked
	Maximizing window
	Displaying annoying dialog box
	clickRegistered.Wait after

	*/

	//逻辑顺序 ↓
	/*
	创建一个按钮，拥有点击属性
	创建 clickRegistered sync.WaitGroup 等待组
	clickRegistered.Add(3)添加了三个计时器
	调用三个 subscribe
	button.Clicked.Broadcast()给所有等待的goroutine  发送通知
	clickRegistered.Wait()暂停 main 的运行，等待 goroutine的运行

	subscribe内部的 cond.Wait() 接收到通知后，继续往下运行 fn() 回调函数
	回调调用 三个 subscribe的 fmt.Println ，并且clickRegistered.Done()减少 clickRegistered的计时器
	所有goroutine运行完毕，clickRegistered的计时器为0，
	clickRegistered.Wait() 继续往下运行fmt.Println("clickRegistered.Wait after")
	*/
}

//并发池
//Pool模式是一个创建和提供可供使用的固定数量实例或Pool实例的方法
func syncPool(){
	myPool := &sync.Pool{
		New: func() interface{}{
			fmt.Println("createing new instance")
			return struct {}{}
		},
	}

	myPool.Get()
	//调用Get 因为没有可用的实例，同时调用了 New

	instance := myPool.Get()
	//调用Get 因为没有可用的实例，同时调用了 New

	myPool.Put(instance)
	//归还一个 instance 实例到池中，池中获取一个可用的实例了

	myPool.Get()
	//调用Get 有一个可用的实例，不需要调用 New

	myPool.Get()
	//调用Get 因为唯一一个可用的实例被上面的Get调用了，所以继续调用了 New


	/*
	Pool的主接口是它的Get方法
	当调用时，Get将首先检查池中是否有可用的实例返回给调用者，
	如果没有，调用它的New方法来创建一个新实例。
	当完成时，调用者调用Put方法把工作的实例归还到池中，以供其他进程使用
	*/

	//当你使用Pool工作时，记住以下几点
	/*
	·当实例化sync.Pool ,使用New方法创建一个成员变量，在调用时是线程安全的。
	·当你收到一个来自Get的实例时，不要对所接收的对象的状态做出任何假设
	·当你用完了一个从Pool中取出来的对象时，一定要调用Put ,否则，Pool就无法复用这个实例了。通常情况
	*/
}