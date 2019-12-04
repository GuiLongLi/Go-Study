package main

import (
	"fmt"
)

func main() {
	checkpeople()
}

//Go 语言提供了另外一种数据类型即接口，它把所有的具有共性的方法定义在一起，任何其他类型只要实现了这些方法就是实现了这个接口。
/*
 定义接口
type interface_name interface {
	method_name1 [return_type]
	method_name2 [return_type]
	method_name3 [return_type]
	...
	method_namen [return_type]
}

定义结构体
type struct_name struct {
	variables
}

 实现接口方法
func (struct_name_variable struct_name) method_name1() [return_type] {
	 方法实现
}

func (struct_name_variable struct_name) method_namen() [return_type] {
	 方法实现
}
*/
//人物接口
type people interface {
	// param 是参数,类型是字符串
	// int,string 是返回值的类型 ，第一个返回值类型是 int ,第二个返回值类型是string
	checkname(param string)(int,string)
}
//人物属性
type attr struct {
	name string
	age int
	sex int
	height float32
}
//检查人物名字是否符合
func (attr attr) checkname(param string) (int,string){
	// len([]rune(string)) 可以判断中文或者中英文混合字符长度
	if(len([]rune(attr.name)) < 3){
		fmt.Printf("名字最少3个字\n")
	}else{
		fmt.Printf("名字符合\n")
	}
	//获取参数
	fmt.Printf("参数是%v\n", param)
	return 123,"string"
}

func checkpeople(){
	var people1  people  //声明接口变量
	people1 = attr {     //人物属性
		name:"张三a",
	}
	int1,string1 := people1.checkname("123123")  //调用接口方法 ,并获取返回值

	fmt.Printf("checkname 返回值是 %v 和 %v \n", int1,string1)
}

