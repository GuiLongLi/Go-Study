package main

import (
	"fmt"
	"net/http"
	"strings"    //strings包主要涉及字符串的基本操作
	"log"         //log包，跟踪日志的记录。
)

func main() {
	//设置访问路由
	http.HandleFunc("/",printfFunc)
	//设置访问的ip和端口
	err := http.ListenAndServe("10.10.87.243:2990",nil)
	if(err != nil){
		log.Fatal("ListenAndServe:",err)
	}
}

//创建输出函数
func printfFunc(httpwrite http.ResponseWriter,request *http.Request){
	//解析参数，默认是不会解析的
	request.ParseForm()

	//这些信息是输出到服务器端的打印信息
	fmt.Println(request.Form)
	fmt.Println(request.URL.Path)
	fmt.Println(request.URL.Scheme)
	fmt.Println(request.Form["url_long"])
	for key,val := range request.Form{
		fmt.Printf("key是%v val是%v\n",key,strings.Join(val,"") )
	}

	//这个写入到w的是输出到客户端的
	fmt.Fprintf(httpwrite,"hello world")

}
