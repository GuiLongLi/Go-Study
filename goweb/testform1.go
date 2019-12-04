package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func sayhelloName(w http.ResponseWriter,r *http.Request)  {
	r.ParseForm() //解析url传递的参数，对于post则解析响应包的主体(request body)
	//注意：如果没有调用 ParseForm方法，下面无法获取表单的数据

	//下面这些信息是输出到服务器端的打印信息
	fmt.Println(r.Form)
	fmt.Println("path",r.URL.Path)
	fmt.Println("scheme",r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k,v := range r.Form{
		fmt.Println("key：",k)
		fmt.Println("val：",strings.Join(v,""))
	}

	//下面这个写入到w 的是输出到客户端的
	fmt.Fprintf(w,"hello testform1.go")
}

func login(w http.ResponseWriter,r *http.Request)  {
	fmt.Println("method：",r.Method) //获取请求的方法
	if r.Method == "GET"{
		t,_ := template.ParseFiles("testform1.gtpl")
		log.Println(t.Execute(w,nil))
	}else{
		r.ParseForm()       //解析url传递的参数，对于POST则解析响应包的主体（request body）
		//请求的是登录数据，name执行登录的逻辑判断
		fmt.Println("username：",r.Form["username"])
		fmt.Println("password：",r.Form["password"])

		//Request本身也提供了FormValue()函数来获取用户提交的参数。
		// 如r.Form["username"]也可写成r.FormValue("username")。
		// 调用r.FormValue时会自动调用r.ParseForm，所以不必提前调用。
		// r.FormValue只会返回同名参数中的第一个，若参数不存在则返回空字符串。
		fmt.Println()
		fmt.Println("FormValue")
		fmt.Println("username：",r.FormValue("username"))
		fmt.Println("password：",r.FormValue("password"))
		fmt.Println("verify_code：",r.FormValue("verify_code"))

		//request.Form是一个url.Values类型，里面存储的是对应的类似key=value的信息，下面展示了可以对form数据进行的一些操作:
		fmt.Println()
		fmt.Println("url.Values")
		v := url.Values{}
		v.Set("name", "Ava")
		v.Add("friend", "Jess")
		v.Add("friend", "Sarah")
		v.Add("friend", "Zoe")
		// v.Encode() == "name=Ava&friend=Jess&friend=Sarah&friend=Zoe"
		fmt.Println(v.Get("name"))
		fmt.Println(v.Get("friend"))
		fmt.Println(v["friend"])
	}

}

func main() {
	http.HandleFunc("/",sayhelloName) //设置路由
	http.HandleFunc("/login",login) //设置路由
	err := http.ListenAndServe(":6665",nil) //设置监听的端口
	if err != nil{
		log.Fatal("ListenAndServe：",err)
	}
}