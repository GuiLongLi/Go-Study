package main

import (
	"fmt"
)

func main() {
	initmap()
	mapdelete()
}

/*
定义 Map
可以使用内建函数 make 也可以使用 map 关键字来定义 Map:

声明变量，默认 map 是 nil
var map_variable map[key_data_type]value_data_type

使用 make 函数
map_variable := make(map[key_data_type]value_data_type)
如果不初始化 map，那么就会创建一个 nil map。nil map 不能用来存放键值对
*/
func initmap(){
	var map1 map[string]string
	map1 = make(map[string]string)  //map 一定要使用 make初始化
	map2 := make(map[string]string)

	map1["name"] = "张三"
	map1["age"] = "18"
	map1["sex"] = "1"

	map2["name"] = "李四"
	map2["age"] = "33"
	map2["sex"] = "2"

	for key,val := range map1{
		fmt.Printf("map1的%v是%v\n", key,val)
	}

	for key,val := range map2{
		fmt.Printf("map2的%v是%v\n", key,val)
	}

	//获取某个属性
	value,flag := map1["name"]
	fmt.Printf("map1的name属性是否存在%v 值是%v\n", flag,value)

}

//delete() 函数
//delete() 函数用于删除集合的元素, 参数为 map 和其对应的 key
func mapdelete(){
	map1 := map[string]string{
		"people1":"张三",
		"people2":"李四",
		"people3":"王五",
		"people4":"陈六",
	}
	fmt.Printf("map1的值是%v\n", map1)

	//删除某个元素
	delete(map1,"people1")

	fmt.Printf("map1删除后的值是%v\n", map1)

}
