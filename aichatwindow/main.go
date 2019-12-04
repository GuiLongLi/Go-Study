package main

import (
	"fmt"

	"aichatwindow/config"
	"aichatwindow/model"
)

func main() {
	if err := config.Init();err != nil{
		fmt.Println(err)
	}

	model.OpenWindow()
}
