package main

import (
	"webproject/config"
	"webproject/model"
	"github.com/gin-gonic/gin"
	"log"
	"github.com/spf13/viper"
	"webproject/router"
)

func main() {
	if err := config.Init();err != nil{
		panic(err)
	}
	if err := model.Init();err != nil{
		panic(err)
	}
	//g := gin.default

	//设置gin模式
	gin.SetMode(viper.GetString("common.runmode"))

	//创建一个gin引擎
	g := gin.New()

	router.InitRouter(g)
	log.Printf("开始监听服务器地址: %s\n", viper.GetString("common.url"))
	//log.Println(http.ListenAndServe(viper.GetString("common.addr"), g).Error())
	if err := g.Run(viper.GetString("common.addr"));err != nil {
		log.Fatal("监听错误:", err)
	}

}
