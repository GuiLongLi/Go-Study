package main

//原文地址：https://github.com/GuoZhaoran/spikeSystem

import (
	"fmt"
	"strconv"
	"net/http"

	"github.com/garyburd/redigo/redis"

	localSpike2 "goredigo/localSpike"
	remoteSpike2 "goredigo/remoteSpike"
	"goredigo/util"
)

/*
---------------------------------------------------------
在启动服务之前，我们需要初始化redis的初始库存信息:
hmset ticket_hash_key "ticket_total_nums" 10000 "ticket_sold_nums" 0
*/

var (
	localSpike localSpike2.LocalSpike
	remoteSpike remoteSpike2.RemoteSpikeKeys
	redisPool *redis.Pool
	done chan int
)

//初始化要使用的结构体和redis连接池
func init(){
	//init函数会在main之前执行
	//初始化本地库存和销量
	localSpike = localSpike2.LocalSpike{
		LocalInStock: 150,
		LocalSalesVolume:0,
	}
	//初始化远程redis 键
	remoteSpike = remoteSpike2.RemoteSpikeKeys{
		SpikeOrderHashKey: "ticket_hash_key",
		TotalInventoryKey: "ticket_total_nums",
		QuantityOfOrderKey: "ticket_sold_nums",
	}
	//新建redis连接池
	redisPool = remoteSpike2.NewPool()
	//创建chan done
	done = make(chan int,1)
	done <- 1 //插入1 到 done
}

func main() {
	//新建路由
	/*
	---------------------------------------------------------
	在启动服务之前，我们需要初始化redis的初始库存信息:
	hmset ticket_hash_key "ticket_total_nums" 10000 "ticket_sold_nums" 0
	*/

	http.HandleFunc("/buy/ticket",handleReq)
	http.ListenAndServe(":6666",nil) //启动并监听服务器端口 6666
}

//处理请求函数，根据请求将响应结果信息写入日志
func handleReq(w http.ResponseWriter,r *http.Request){
	redisConn := redisPool.Get() //创建redis连接，返回连接
	LogMsg := ""
	<- done  //取出done ，形成阻塞
	/*
	在上面例子中 <- done 会从 done 这个 channel 中消费一条数据，
	所以他会等待 channel 中有数据才开始消费，在写入数据之前，goroutine不会退出
	因此，main goroutine 和匿名的 goroutine 都被阻塞住
	*/
	/*
	通过
	LogMsg = fmt.Sprintf("<- done的值是: %+v",<- done)
	util.WriteLog(LogMsg,"./stat.log")
	我们可以看出 done 形成过程
	----------------------------------
	done <- 1   插入 1
	<- done  取出 done 数据，也就是 1
	done <- -1 再插入 -1
	<- done  再取出数据 ，也就变成了 -1
	后面循环 插入 -1 和 取出 -1
	----------------------------------
	在此过程 <- done 取出数据 会形成阻塞，直到 done <- -1 插入数据后，阻塞才会被解除，这样就形成原子过程
	----------------------------------
	同原理， done <- -1 如果已经被插入数据，那么他必须 <- done 取出数据后，再能重新插入数据
	*/

	//全局读写锁
	if localSpike.LocalDeductionStock() && remoteSpike.RemoteDeductionStock(redisConn){
		util.RespJson(w,1,"抢票成功",nil)
		LogMsg = LogMsg+"result:1,localSales:"+strconv.FormatInt(localSpike.LocalSalesVolume,10)
	}else{
		util.RespJson(w,-1,"已售罄",nil)
		LogMsg = LogMsg+"result:0,localSales:"+strconv.FormatInt(localSpike.LocalSalesVolume,10)
	}
	//插入到done,解除阻塞
	done <- -1
	//将抢票状态写入到log中
	util.WriteLog(LogMsg,"./stat.log")
}