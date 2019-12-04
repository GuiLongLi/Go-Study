package main

import (
	"fmt"
	"net/http"
	"strings"
	"log"
)

func sayhelloName(w http.ResponseWriter,r *http.Request){
	r.ParseForm() //解析参数，默认是不会解析的
	fmt.Printf("r.Form -- %+v\n", r.Form) //这些是输出到服务器端的信息
	fmt.Printf("r.URL.Path -- %+v\n", r.URL.Path)
	fmt.Printf("r.URL.Scheme -- %+v\n", r.URL.Scheme)
	fmt.Printf("r.Form[url_long] -- %+v\n", r.Form["url_long"])
	for k,v := range r.Form {
		fmt.Printf("key -- %+v\n", k)
		fmt.Printf("val -- %+v\n", strings.Join(v,""))
	}
	fmt.Fprintf(w,"hello testweb1") //写入到w ，输出到客户端

}

func main(){
	http.HandleFunc("/",sayhelloName) //设置路由
	err := http.ListenAndServe("0.0.0.0:6665",nil) //设置端口
	if err != nil{
		log.Fatal("ListenAndServe: ",err)
	}
}