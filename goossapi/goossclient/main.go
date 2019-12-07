package main

import (
	"fmt"

	"github.com/spf13/viper"

	"goossclient/model"

)

func main() {
	viper.AddConfigPath("conf")
	if err := viper.ReadInConfig();err != nil{
		fmt.Println(err)
	}

	model.OpenWindow()
}
