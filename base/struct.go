package main

import "fmt"

func main() {
	getpeople()

	//结构体作为函数参数
	var people4 = people{
		"abc",
		30,
		1,
		165.5,
	}
	structparam(people4)

	//结构体指针
	var people5 = people{
		"aaaa",
		88,
		0,
		177.5,
	}
	structpointer(&people5)

	//通过指针修改结构体
	structchange(&people5)
	fmt.Printf("people5的值是%v\n", people5)
}

//定义结构体
//结构体定义需要使用 type 和 struct 语句。struct 语句定义一个新的数据类型，结构体有中有一个或多个成员。type 语句设定了结构体的名称。
/*
type struct_variable_type struct {
   member definition;
   member definition;
   ...
   member definition;
}
*/
type people struct{
	name string
	age int
	sex int
	height float32
}
func getpeople(){
	var people1 = people{
		"张三",
		18,
		1,
		168.5,
	}
	fmt.Printf("people1的值是%v\n", people1)

	fmt.Printf("people2的值是%v\n", people{"李四",20,0,199})

	var people3 = people{
		name:"王五",
		sex:2,
		age:99,
		height:155.5,
	}
	fmt.Printf("people3的值是%v\n", people3)

	fmt.Printf("people4的值是%v\n", people{name:"陈六"})  //忽略部分属性，忽略的属性为 0 或 空

	people3.age = 55
	fmt.Printf("people3的重新修改值是%v\n", people3) //重新修改属性
	/* 打印 people3 信息 */
	fmt.Printf( "people3 name : %s\n", people3.name)
	fmt.Printf( "people3 sex : %d\n", people3.sex)
	fmt.Printf( "people3 age : %d\n", people3.age)
	fmt.Printf( "people3 height : %f\n", people3.height)
	println()
}

//结构体作为函数参数
func structparam(people4 people){
	fmt.Printf( "people4 name : %s\n", people4.name)
	fmt.Printf( "people4 sex : %d\n", people4.sex)
	fmt.Printf( "people4 age : %d\n", people4.age)
	fmt.Printf( "people4 height : %f\n", people4.height)
	println()
}

//结构体指针
func structpointer(people5 *people){
	fmt.Printf( "people5 name : %s\n", people5.name)
	fmt.Printf( "people5 sex : %d\n", people5.sex)
	fmt.Printf( "people5 age : %d\n", people5.age)
	fmt.Printf( "people5 height : %f\n", people5.height)
	println()
}

//通过指针修改结构体属性
func structchange(people5 *people){
	people5.name = "who are your"
}