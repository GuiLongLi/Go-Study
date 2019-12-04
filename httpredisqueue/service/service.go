package service

import (
	"fmt"
	"log"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"

	"httpredisqueue/pkg/errno"
	"httpredisqueue/model"
	. "httpredisqueue/handler"
)
var cli *redis.Client

//列表数据的结构体
type ResultList struct {
	Total int `json:"total"`
	List []string `json:"list"`
}

//设置redis客户端构造体
func SetRedisCli(client *redis.Client){
	cli = client
}

//首页
func Index(c *gin.Context){
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>hello world</title>
</head>
<body>
    hello world
</body>
</html>
`
	SendResponseHtml(c,nil,html)
}

func SetRedis(c *gin.Context){
	key := c.Query("key")
	value := c.Query("value")
	if key == "" || value == ""{
		SendResponse(c,errno.VALUEERROR,nil)
		return
	}
	_,err := cli.Set(key,value,0).Result();
	if(err != nil){
		SendResponse(c,errno.ApiServerError,err)
	}
	results := "成功"
	SendResponse(c,nil,results)
}

func GetRedis(c *gin.Context){
	key := c.Query("key")
	if key == ""{
		SendResponse(c,errno.VALUEERROR,nil)
		return
	}
	log.Printf("key  %+v\n",key)
	results,err := cli.Get(key).Result();
	if(err != nil){
		SendResponse(c,errno.ApiServerError,err)
	}
	log.Printf("results  %+v\n",results)
	SendResponse(c,nil,results)
}

func Inqueue(c *gin.Context){
	value := c.Query("value")
	if value == ""{
		SendResponse(c,errno.VALUEERROR,nil)
		return
	}
	//取出redis缓存队列
	redisstr,err := cli.Get("queuelist").Result()
	valueGet := []byte(redisstr) //队列数据转换为 byte
	if err != nil {
		fmt.Println(err)
	}
	var redislist []string
	//json解析队列数据
	errShal := json.Unmarshal(valueGet, &redislist)
	if errShal != nil {
		fmt.Println(err)
	}
	fmt.Printf("redislist %+v\n",redislist)

	//使用队列模型 插入队列
	queue := model.NormalQueue{}
	queue.List = redislist  //队列数据赋值
	list := queue.Innormalqueue(value) //插入队列
	total := queue.Length()			//获取队列长度

	//数据转换成json  保存到redis里面
	queuelist, _ := json.Marshal(list)
	_,err = cli.Set("queuelist",queuelist,0).Result()  //保存到redis里面
	if(err != nil){
		SendResponse(c,errno.ApiServerError,err)
	}
	//fmt.Printf("str %v",str)

	//返回数据
	var res interface{}
	res = ResultList{
		Total:total,
		List:list,
	}
	fmt.Printf("res %+v\n",res)
	SendResponse(c,nil,res)
}

func Outqueue(c *gin.Context){
	//取出redis缓存队列
	redisstr,err := cli.Get("queuelist").Result()
	valueGet := []byte(redisstr) //队列数据转换为 byte
	if err != nil {
		fmt.Println(err)
	}
	var redislist []string
	//json解析队列数据
	errShal := json.Unmarshal(valueGet, &redislist)
	if errShal != nil {
		fmt.Println(err)
	}
	fmt.Printf("redislist %+v\n",redislist)

	//使用队列模型 出队列
	queue := model.NormalQueue{}
	queue.List = redislist  //队列数据赋值
	list := queue.Outnormalqueue() //出队列
	total := queue.Length()			//获取队列长度//返回数据

	//数据转换成json  保存到redis里面
	queuelist, _ := json.Marshal(list)
	_,err = cli.Set("queuelist",queuelist,0).Result()  //保存到redis里面
	if(err != nil){
		SendResponse(c,errno.ApiServerError,err)
	}

	var res interface{}
	res = ResultList{
		Total:total,
		List:list,
	}
	fmt.Printf("res %+v\n",res)
	SendResponse(c,nil,res)
}