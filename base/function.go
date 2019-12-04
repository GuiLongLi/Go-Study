package main

import "fmt"

func main() {
	max := maxfunction(4,8)
	fmt.Printf("maxfunction 4和8的最大值是%v\n", max)

	str1,str2 := swapfunction("baidu","google")
	fmt.Printf("swapfunction 第一是%v,第二是%v\n", str1,str2)

	//函数值传递
	var int1,int2 int = 5,10
	intresult := valueswapfunction(int1,int2)
	fmt.Printf("valueswapfunction int1的值是%v int2的值是%v 传递结果是%v\n",int1,int2, intresult)

	//引用传递
	var int3,int4 int = 3,11
	quoteswapfunction(&int3,&int4)
	fmt.Printf("quoteswapfunction int3的值是%v int4的值是%v\n", int3,int4)

	//函数作为实参
	callbackresult1 := callbackfunction1(6,callbackfunction2)
	fmt.Printf("callbackfunction1 的返回值是%v\n", callbackresult1)
	callbackresult2 := callbackfunction1(1, func(int1 int) int {
		return int1*100
	})
	fmt.Printf("callbackfunction1 的返回值是%v\n", callbackresult2)

	//函数闭包
	setNum := getSequence()
	fmt.Printf("函数闭包setNum的值%v\n", setNum())
	fmt.Printf("函数闭包setNum的值%v\n", setNum())
	fmt.Printf("函数闭包setNum的值%v\n", setNum())
	//重新获取闭包函数
	setNum1 := getSequence()
	fmt.Printf("函数闭包setNum1的值%v\n", setNum1())
	fmt.Printf("函数闭包setNum1的值%v\n", setNum1())
	fmt.Printf("函数闭包setNum1的值%v\n", setNum1())
	//闭包传参
	paramNum1 := paramcolsure(5,1)
	fmt.Printf("闭包传参paramNum1的值是")
	fmt.Println(paramNum1(1,1,1))
	fmt.Printf("闭包传参paramNum1的值是")
	fmt.Println(paramNum1(2,3,4))

	//函数方法
	//计算圆面积
	var c1 Circle
	c1.radius = 10
	fmt.Printf("圆的面积是%v\n", c1.getArea())
}

//求最大值
func maxfunction(num1,num2 int) int {
	var result int

	if(num1 > num2){
		result = num1
	}else{
		result = num2
	}
	return result
}

//交换位置
func swapfunction(str1, str2 string) (string,string) {
	return str2,str1
}

//函数值传递是指在调用函数时将实际参数复制一份传递到函数中，这样在函数中如果对参数进行修改，将不会影响到实际参数。
//默认情况下，Go 语言使用的是值传递，即在调用过程中不会影响到实际参数。
func valueswapfunction(int1 ,int2 int) int {
	var result int
	result = int1
	int1 = int2
	int2 = result

	return result
}

//引用传递是指在调用函数时将实际参数的地址传递到函数中，那么在函数中对参数所进行的修改，将影响到实际参数。
//引用传递指针参数传递到函数内，以下是交换函数 使用了引用传递：
func quoteswapfunction(int1 *int,int2 *int) {
	var result int
	result = *int1
	*int1 = *int2
	*int2 = result
}

//函数作为实参 把函数作为参数传递给另一个函数
// 声明一个函数类型
type callback func(int) int
func callbackfunction1(int1 int,fun callback) int {
	return fun(int1)
}
func callbackfunction2(int2 int) int {
	int2 = int2*10
	return int2
}

//-------------------------------------------------------------
//函数闭包
func getSequence() func() int{
	i := 0
	return func () int {
		i += 1
		return i
	}
}
//闭包带参数
func paramcolsure(int1 int,int2 int) func(int3 int,int4 int,int5 int) (int,int){
	int1 = int1*100
	return func(int3 int,int4 int,int5 int) (int,int){
		int1++
		return int1,(int2+int3+int4+int5)
	}
}
//--------------------------------------------------------------
//函数方法
/*
Go 语言中同时有函数和方法。一个方法就是一个包含了接受者的函数，接受者可以是命名类型或者结构体类型的一个值或者是一个指针。所有给定类型的方法属于该类型的方法集。语法格式如下：
func (variable_name variable_data_type) function_name() [return_type]{
	函数体
}
*/
//定义结构体
type Circle struct {
	radius float64 //结构体的属性 radius ,类型是float64
}
//计算圆的面积
func (c Circle) getArea() float64{
	//S = πr²
	return 3.14 * c.radius * c.radius
}
