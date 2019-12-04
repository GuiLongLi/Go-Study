package main

import (
	"log"

	"github.com/spf13/viper"

	"socketserver/config"
	"socketserver/service"
)

func main() {
	if err := config.Init();err != nil{
		panic(err)
	}

	log.Printf("开始监听服务器端口: %s\n", viper.GetString("common.server.addr"))

	service.Listenserver()

}

