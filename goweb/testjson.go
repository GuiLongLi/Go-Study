package main

import (
	"encoding/json"
	"fmt"
)

type Server struct {
	ServerName string `json:"servername"`
	ServerIP string `json:"serverip"`
}

type Serverslice struct {
	Servers []Server `json:"servers"`
}

func main() {
	//已知数据类型的json解析
	fmt.Println("testdecodejson")
	testdecodejson()

	//未知数据类型的json解析
	fmt.Println()
	fmt.Println("testdecodejson2")
	testdecodejson2()

	//转换成json字符串
	fmt.Println()
	fmt.Println("testencodejson")
	testencodejson()

	//json构造体的其他用法
	fmt.Println()
	fmt.Println("structother")
	structother()

	//student
	fmt.Println()
	fmt.Println("teststudent")
	teststudent()
}

func testdecodejson(){
	var s Serverslice
	str := `{
	"servers":[
		{
		"serverName":"Shanghai_VPN",
		"serverIP":"127.0.0.1"
		},
		{
		"serverName":"Beijing_VPN",
		"serverIP":"127.0.0.2"
		}
]
}`
	json.Unmarshal([]byte(str),&s)
	fmt.Println(s)
}

func testdecodejson2(){
	b := []byte(`{
	"Name":"wednesday",
	"Age":6,
	"parents":["gomez","morticia"]
}`)
	var f interface{}
	err := json.Unmarshal(b,&f)
	if err != nil{
		fmt.Println(err)
		return
	}

	m := f.(map[string]interface{})

	for k,v := range m{
		switch vv := v.(type){
		case string:
			fmt.Println(k,"is string",vv)
		case int:
			fmt.Println(k,"is int",vv)
		case float64:
			fmt.Println(k,"is float64",vv)
		case []interface{}:
			fmt.Println(k,"is an array:")
			for i,u := range vv{
				fmt.Println(i,u)
			}
		default:
			fmt.Println(k,"is of a type i do not know how to handle")
		}
	}
}

func testencodejson(){
	/*
	Marshal函数只有在转换成功的时候才会返回数据，在转换的过程中我们需要注意几点：

	·JSON对象只支持string作为key，所以要编码一个map，那么必须是map[string]T这种类型(T是Go语言中任意的类型)

	·Channel, complex和function是不能被编码成JSON的

	·嵌套的数据是不能编码的，不然会让JSON编码进入死循环

	·指针在编码的时候会输出指针指向的内容，而空指针会输出null
	*/

	var s Serverslice
	s.Servers = append(s.Servers,Server{ServerName:"shanghai_vpn",ServerIP:"127.0.0.1"})
	s.Servers = append(s.Servers,Server{ServerName:"beijing_vpn",ServerIP:"127.0.0.2"})
	b,err := json.Marshal(s)
	if err != nil{
		fmt.Println("json err: ",err)
	}
	fmt.Println(string(b))
}

/*
·字段的tag 是 "-" , 那么这个字段不会输出到json

·tag 中带有自定义的名称，那么这个自定义名称会出现在json的字段名中，例如上面例子中的serverName

·tag 中如果带有 "omitempty" 选项，那么如果该字段值为空，就不会输出到json串中

·如果字段类型是bool,string,int,int64等，而tag中带有",string"选项，那么这个字段在输出到json的时候回把该字段对应的值转换成json字符串

*/
type ServerOther struct {
	//ID不会导出在json中
	ID int `json:"-"`

	//ServerName2 的值会进行二次json编码
	ServerName string `json:"serverName"`
	ServerName2 string `json:"serverName2,string"`

	//如果ServerIP 为空，则不会输出到json串中
	ServerIP string `json:"serverIP,omitempty"`
}

func structother(){
	s := ServerOther{
		ID:3,
		ServerName:`GO "1.13"`,
		ServerName2:`GO "1.13"`,
		ServerIP:``,
	}
	b,_ := json.Marshal(s)
	fmt.Printf("serverother: %s\n",b)
}


type Student struct {
	User []User
}

type User struct {
	Name string `json:"name"`
	Age int `json:"age"`
}

func teststudent(){
	user1 := `{
	"Name":"user1",
	"Age":123
}`
	var u1 User
	json.Unmarshal([]byte(user1),&u1)

	user2 := `{
	"Name":"user2",
	"Age":343
}`
	var u2 User
	json.Unmarshal([]byte(user2),&u2)

	fmt.Printf("user2 %+v\n",u2)

	u3 := make([]User,0)
	u3 = append(u3,u1,u2)
	fmt.Printf("u3 %+v\n",u3)

	u4 := Student{
		u3,
	}
	fmt.Printf("u4 %+v\n",u4)

	u5 := []User{u1,u2}
	fmt.Printf("u5 %+v\n",u5)
	u6 := Student{
		u5,
	}
	fmt.Printf("u6 %+v\n",u6)


}