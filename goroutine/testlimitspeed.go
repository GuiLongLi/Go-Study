package main

import (
	"os"
	"sync"
	"log"
	"fmt"
	"sort"
	"time"
	"context"

	"golang.org/x/time/rate"
)

/*
大多数的限速都是基于令牌桶算法的：
如果要访问资源，你必须拥有资源的访问令牌，没有令牌的请求会被拒绝。
假设这些令牌存储在一个等待被检索使用的桶中。
桶的深度为d ，表示一个桶可以容纳d 个访问令牌。（如：d 为5，则可以存放5 个令牌）
每当你需要访问资料时，都会在桶中删除一个令牌。
如果你的存储桶包含5 个令牌，前5 次访问没有问题，操作正常进行。
但是到了第6 次尝试时，就没有访问令牌可用，你的请求必须排队等待，直到令牌可用，或者被拒绝操作。
*/

func main() {
	fmt.Println("simpleLimit")
	simpleLimit()
	fmt.Println()
	fmt.Println("multiLimit")
	multiLimit()
}

type APIConnection struct {
	rateLimiter *rate.Limiter
}

func (a *APIConnection) ReadFile(ctx context.Context)error{
	if err := a.rateLimiter.Wait(ctx); err != nil{
		return err
	}
	//我们登陆限速器有足够的令牌来完成我们的请求
	return nil
}
func (a *APIConnection) ResulveAddress(ctx context.Context) error{
	if err := a.rateLimiter.Wait(ctx); err != nil{
		return err
	}
	//我们登陆限速器有足够的令牌来完成我们的请求
	return nil
}
func Open() *APIConnection{
	return &APIConnection{
		rateLimiter:rate.NewLimiter(rate.Limit(1),1),
	}
	//在这里，我们将所有api连接的速率限制设置为每秒一次
}

//单一速率限制
func simpleLimit(){
	defer log.Println("done.")
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime|log.LUTC)

	apiConnection := Open()
	var wg sync.WaitGroup
	wg.Add(20)

	for i:=0;i<10;i++{
		go func() {
			defer wg.Done()
			err := apiConnection.ReadFile(context.Background())
			if err != nil{
				log.Printf("cannot readfile: %v\n",err)
			}
			log.Printf("readfile\n")
		}()
	}

	for i:=0;i<10;i++{
		go func() {
			defer wg.Done()
			err := apiConnection.ResulveAddress(context.Background())
			if err != nil {
				log.Printf("cannot resolveaddress: %v\n",err)
			}
			log.Printf("resolveaddress")
		}()
	}

	wg.Wait()

	/*
	下面我们使用 golang.org/x/time/rate 包中的令牌桶限速器实现限速

	·Limit定义了某些事件的最大频率
	·Limit表示为每秒事件数
	·zero Limit不允许发生任何事件
	type Limit float64

	·NewLimiter 返回一个新的Limit
	·他允许事件速率为r ,并允许最大数为b 的token
	·r是我们之前说的速率 ,b 是桶深度
	func NewLimiter(r Limit,b int)*Limiter

	·每一个都将事件之间的最小事件间隔转换为一个Limit
	func Every(interval time.Duration) Limit

	·获取间隔时间，如果想针对每次操作的间隔时间进行测量，而不是请求的间隔长度，我们可以这样写：
	rate.Limit(events/timePeriod.Seconds())

	·获取间隔时间，如果不想每次都手动输入，可以使用Every 函数，将返回rate.Inf
	func Per(eventCount int,duration time.Duration) rate.Limit{
		return rate.Every(duration/time.Duration(eventCount))
	}

	·在创建rate.Limiter之后，我们将使用他来阻塞我们的请求，直到获得访问令牌。
	·Wait是WaitN(ctx,1)的缩写
	func (lim *Limiter) Wait(ctx context.Context)

	·WaitN 会执行直到有n 个事件发生
	·如果n 超过Limiter 的突发大小，context被取消，或者预期等待时间超过
	·context的deadline ,他会返回一个错误
	func (lim *Limiter) WaitN(ctx context.Context,n int) (err error)
	*/
}

//速率限制组
/*
在生产中，我们可能会想要建立多层次的限制：
用细粒度的控制来限制每秒的请求，
用粗粒度的控制来限制每分钟、每小时或每天的请求

为此，我们把不同粒度的限速器独立，然后把他们组合成一个限速器组来管理
*/
type RateLimiter interface {
	Wait(context.Context) error
	Limit() rate.Limit
}
type multiLimiterStruct struct {
	limiters []RateLimiter
}
func MultiLimiter(limiters ...RateLimiter) *multiLimiterStruct{
	byLimit := func(i,j int) bool {
		return limiters[i].Limit() < limiters[j].Limit()
	}
	sort.Slice(limiters,byLimit)
	return &multiLimiterStruct{limiters:limiters}
}
func (l *multiLimiterStruct) Wait(ctx context.Context) error{
	for _,l := range l.limiters{
		if err := l.Wait(ctx);err != nil{
			return err
		}
	}
	return nil
}

func(l *multiLimiterStruct) Limit() rate.Limit{
	return l.limiters[0].Limit()
}
/*
·type RateLimiter interface {
在这里我们定义了一个RateLimiter 接口，使MultiLimiter 可以递归的定义其他的MultiLimiter 实例

·sort.Slice(limiters,byLimit)
这里我们实现了一个优化，并根据每个RateLimiter 的Limit() 进行排序

·return l.limiters[0].Limit()
因为我们在multiLimiter 实例化时对子RateLimiter 实例进行了排序，
所以我们可以直接返回限制最多的限制器，这将是切片slice 中的第一个元素
*/
type MultiAPIConnection struct {
	networkLimit,
	diskLimit,
	apiLimit RateLimiter
}

func Per(eventCount int,duration time.Duration) rate.Limit{

	return rate.Every(duration/time.Duration(eventCount))
}

func MultiOpen() *MultiAPIConnection{
	//在这里，我们定义每秒的限制，避免突发请求
	//限制每秒钟生成2个令牌，就是每0.5秒生成一个令牌 ,最多存储有2个令牌
	secondLimit := rate.NewLimiter(Per(2,time.Second),2)

	//在这里，我们定义每分钟的限制，为用户提供初始池，每秒的限制将确保我们的系统不会被突发请求而超载
	//限制每分钟生成10个令牌，就是每6秒生成一个令牌 ，最多存储有10个令牌
	minuteLimit := rate.NewLimiter(Per(10,time.Minute),10)

	return &MultiAPIConnection{
		//然后我们组合这两个限制，并将其设置为 MultiAPIConnection 的 apiLimit 限速器。
		apiLimit:MultiLimiter(secondLimit,minuteLimit),

		//这里我们为硬盘读取设置了一个限速器，每秒只能读取一次 ，最多存储有1个令牌
		diskLimit:rate.NewLimiter(rate.Limit(1),1),

		//对于网络，我们设置了每秒3 个请求的限速 ，最多存储有3个令牌
		networkLimit:rate.NewLimiter(Per(3,time.Second),3),
	}
}
func (a *MultiAPIConnection) MultiReadFile(ctx context.Context) error{
	//当我们读取文件时，我们融合api 限速器和硬盘限速器的限制。
	err := MultiLimiter(a.apiLimit,a.diskLimit).Wait(ctx)
	if err != nil{
		return err
	}
	//假设我们在这里执行一些逻辑
	return nil
}
func (a *MultiAPIConnection) MultiResolveAddress(ctx context.Context) error{
	//当我们需要访问网络时，我们融合api 限速器和网络限速器的限制。
	err := MultiLimiter(a.apiLimit,a.networkLimit).Wait(ctx)
	if err != nil{
		return err
	}
	//假设我们在这里执行一些逻辑
	return nil
}


func multiLimit(){
	defer log.Println("multiLimit done.")
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime|log.LUTC)

	multiapiConnection := MultiOpen()
	var wg sync.WaitGroup
	wg.Add(20)

	for i:=0;i<10;i++{
		go func() {
			defer wg.Done()
			err := multiapiConnection.MultiReadFile(context.Background())
			if err != nil{
				log.Printf("cannot MultiReadFile: %v\n",err)
			}
			log.Printf("MultiReadFile\n")
		}()
	}

	for i:=0;i<10;i++{
		go func() {
			defer wg.Done()
			err := multiapiConnection.MultiResolveAddress(context.Background())
			if err != nil {
				log.Printf("cannot MultiResolveAddress: %v\n",err)
			}
			log.Printf("MultiResolveAddress")
		}()
	}

	wg.Wait()

	/*
	我们可以看到，客户端每秒发出两个请求，直到第11个请求，我们开始每6秒发出一次请求。
	这是因为我们耗尽了分钟级限速器的可用令牌，所以限制了请求速度。

	为什么第11个请求只延时了2秒，而不是像之后的请求那样隔了6秒，这可能有点违反直觉。
	原因是这样的，虽然我们将api 请求限制为10次/分钟，但是在一分钟内令牌增加是一个递增的过程。
	当我们第11个请求到达时，我们的分钟级线速度去已经积累了一个令牌

	*/
}
