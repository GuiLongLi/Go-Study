package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	"aichat/config"
	"aichat/router"
	"aichat/service"
)

func main() {
	if err := config.Init();err != nil{
		panic(err)
	}
	//设置gin模式
	gin.SetMode(viper.GetString("common.server.runmode"))

	//创建一个gin引擎
	g := gin.New()

	router.InitRouter(g)
	log.Printf("开始监听服务器地址: %s\n", viper.GetString("common.server.url"))

	// Listen and Server in https://127.0.0.1:8080
	//err := g.Run(viper.GetString("common.server.addr")) //使用http
	//err := g.RunTLS( //使用https
	//	viper.GetString("common.server.addr"),
	//	"/usr/local/orange/conf/cert/online/go.daily886.com.pem",
	//	"/usr/local/orange/conf/cert/online/go.daily886.com.key")
	//if err != nil {
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

