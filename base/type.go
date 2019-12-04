package main

import (
	"fmt"
	"reflect"
)

/*
类型转换用于将一种数据类型的变量转换为另外一种类型的变量。Go 语言类型转换基本格式如下：
type_name(expression)
type_name 为类型，expression 为表达式。
*/
func main() {
	var sum int = 17
	var count int = 5
	var mean float32

	fmt.Printf("sum转换前的类型是%v 转换后的类型是%v\n",reflect.TypeOf(sum), reflect.TypeOf(float32(sum)))
	fmt.Printf("count转换前的类型是%v 转换后的类型是%v\n",reflect.TypeOf(count), reflect.TypeOf(float32(count)))

	mean = float32(sum)/float32(count)
	fmt.Printf("mean 的值为: %f\n",mean)
}
