package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"httpredisqueue/router"
	"httpredisqueue/config"
	"httpredisqueue/model"
	"httpredisqueue/service"
)

func main() {
	//初始化配置
	if err := config.Init();err != nil{
		panic(err)
	}

	//初始化redis
	redisCli,err := model.InitRedis();
	if err != nil{
		panic(err)
	}
	service.SetRedisCli(redisCli)

	//设置gin模式
	gin.SetMode(viper.GetString("common.server.runmode"))

	//创建一个gin引擎
	g := gin.New()

	router.InitRouter(g)
	log.Printf("开始监听服务器地址: %s\n", viper.GetString("common.server.url"))
	//不使用热重启
	//if err := g.Run(viper.GetString("common.server.addr"));err != nil {
	//	log.Fatal("监听错误:", err)
	//}

	//使用热重启
	// kill -USR2 pid 重启
	// kill -INT pid 关闭
	add := viper.GetString("common.server.addr")
	srv := &http.Server{
		Addr:    add,
		Handler: g,
	}
	log.Printf( "srv.Addr  %v  \n", srv.Addr)
	service.Listenserver(srv)
}
