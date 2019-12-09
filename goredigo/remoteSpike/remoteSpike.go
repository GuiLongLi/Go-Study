package remoteSpike

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"

	"goredigo/util"
)

//lua 脚本
const LuaScript = `
local ticket_key = KEYS[1]
local ticket_total_key = ARGV[1]
local ticket_sold_key = ARGV[2]
local ticket_total_nums = tonumber(redis.call('HGET', ticket_key, ticket_total_key))
local ticket_sold_nums = tonumber(redis.call('HGET', ticket_key, ticket_sold_key))
-- 查看是否还有余票，增加订单数量，返回结果值
if(ticket_total_nums >= ticket_sold_nums) then
	return redis.call('HINCRBY', ticket_key, ticket_sold_key, 1)
	-- Redis HINCRBY 命令用于为哈希表中的字段值加上指定增量值。增量可以正数也可以为负数，相当于对指定字段进行减法操作。
	-- 返回 哈希表中字段的值。
end
return 0
`

//远程订单存储键值
type RemoteSpikeKeys struct {
	SpikeOrderHashKey string //redis中秒杀订单hash结构key
	TotalInventoryKey string //hash结构中总订单库存key
	QuantityOfOrderKey string //hash结构中已有订单数量key
}

//初始化redis连接池
func NewPool() *redis.Pool{
	util.WriteLog("redis start","./stat.log")
	return &redis.Pool{
		// 最大空闲链接等待时间
		IdleTimeout: 5 * time.Second,
		// 最大空闲链接
		MaxIdle: 10000,
		// 最大数量连接
		MaxActive: 12000,
		Dial: func()(redis.Conn, error){
			//连接redis tpc协议 6639端口
			c,err := redis.Dial("tcp",":6379")
			if err != nil{
				panic(err.Error())
			}
			//密码登录
			if _, err := c.Do("AUTH", "xxx"); err != nil {
				c.Close()
				panic(err.Error())
			}
			util.WriteLog("redis connected ","./stat.log")
			return c,err
		},
	}
}

//远端统一扣库存
/*
远程连接redis，
通过lua脚本操作redis销量
*/
func (RS *RemoteSpikeKeys) RemoteDeductionStock(conn redis.Conn) bool{
	lua := redis.NewScript(1,LuaScript) //初始化lua脚本
	//lua.Do 调用脚本 ,
	//参数1 redis连接
	//参数2 订单的hash键
	//参数3 订单的库存键
	//参数4 订单的销量键
	result,err := redis.Int(lua.Do(conn, RS.SpikeOrderHashKey,RS.TotalInventoryKey, RS.QuantityOfOrderKey))
	//redis.Int 把lua.Do返回的结果转换为Int 类型
	logMsg := fmt.Sprintf("lua.Do result: %+v",result)
	util.WriteLog(logMsg,"./stat.log")
	if err != nil{
		return false
	}
	//判断 结果 是不是 不等于 0
	return result != 0
}