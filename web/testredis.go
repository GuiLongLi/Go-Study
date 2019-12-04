package main

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

var client *redis.Client

func main() {
	err := redisInits()
	if err != nil {
		fmt.Println("链接失败: ",err)
	}

	//redis set 10秒生存时间
	_, err = client.SetNX("test", "value", 10*time.Second).Result()
	if(err != nil){
		fmt.Println("设置失败: ",err)
	}
	//redis get
	value := client.Get("test");
	if err != nil {
		fmt.Println("获取失败: ",err)
	}
	fmt.Println("test", value)

	// 自定义redis命令
	cmd1 := client.Do("get", "test")
	client.Process(cmd1)
	test, err := cmd1.Result()
	if err != nil {
		fmt.Println("test failed.", err)
	} else {
		fmt.Println("test is", test)
	}

	//使用管道 pipeline 获取数据
	_, err = client.SetNX("test1", "value1", 10*time.Second).Result()
	_, err = client.SetNX("test2", "value2", 10*time.Second).Result()
	_, err = client.SetNX("test3", "value3", 10*time.Second).Result()

	/*
Redis的pipeline功能的原理是 Client通过一次性将多条redis命令发往Redis Server，减少了每条命令分别传输的IO开销。同时减少了系统调用的次数，因此提升了整体的吞吐能力。

我们在主-从模式的Redis中，pipeline功能应该用的很多，但是Cluster模式下，估计还没有几个人用过。
我们知道 redis cluster 默认分配了 16384 个slot，当我们set一个key 时，会用CRC16算法来取模得到所属的slot，然后将这个key 分到哈希槽区间的节点上，具体算法就是：CRC16(key) % 16384。如果我们使用pipeline功能，一个批次中包含的多条命令，每条命令涉及的key可能属于不同的slot

go-redis 为了解决这个问题, 分为3步
源码可以阅读 defaultProcessPipeline
1) 将计算command 所属的slot, 根据slot选择合适的Cluster Node
2）将同一个Cluster Node 的所有command，放在一个批次中发送（并发操作）
3）接收结果
	*/
	pipe := client.Pipeline()
	pipe.Get("test1")
	pipe.Get("test2")
	pipe.Get("test3")
	cmders, err := pipe.Exec()
	if err != nil {
		fmt.Println("err", err)
	}
	for _, cmder := range cmders {
		cmd := cmder.(*redis.StringCmd)
		str, err := cmd.Result()
		if err != nil {
			fmt.Println("err", err)
		}
		fmt.Println("str", str)
	}
}

//初始化链接
func redisInits() error{
	client = redis.NewClient(&redis.Options{
		Addr:     "10.10.87.242:6001",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	if(err != nil){
		return err
	}
	fmt.Println(pong)
	// Output: PONG <nil>
	return nil
}
