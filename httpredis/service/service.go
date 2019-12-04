package service

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"

	"httpredis/pkg/errno"
	. "httpredis/handler"
)
var cli *redis.Client

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