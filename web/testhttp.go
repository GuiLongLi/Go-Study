package main

import (
	"fmt"               //fmt 包提供了打印函数将数据以字符串形式输出
	"io"  	             //io 包提供了 io.Reader 和 io.Writer 接口，分别用于数据的输入和输出
	"os"                 //os包提供了操作系统的函数，设计向Unix风格，但是错误处理是go风格，当os包使用时，如果失败之后返回错误类型而不是错误数量．
	"net/http"          //net/http包提供HTTP客户端和服务器实现。Get，Head，Post和PostForm发出HTTP（或HTTPS）请求
)

//这里使用go 来搭建一个web服务器
func main(){
	testhttp()
}

//http测试
func testhttp(){
	//生成client 参数为默认
	client := &http.Client{}

	//生成要访问的url
	url := "http://baidu.com"

	//http请求
	request,error := http.NewRequest("GET",url,nil)

	if(error != nil){
		panic(error)
	}

	//处理返回结果
	response,_ := client.Do(request)

	//将结果定位到标准输出 也可以直接打印出来 或者定位到其他地方进行相应的处理
	stdout := os.Stdout
	_,error1 := io.Copy(stdout,response.Body)

	//返回的状态码
	status := response.StatusCode
	fmt.Printf("状态码是%v\n", status)
	fmt.Printf("错误信息是%v\n", error1)
	fmt.Printf("response是%v\n", response)
	fmt.Printf("request是%v\n", request)
}

