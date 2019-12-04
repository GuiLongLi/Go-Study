package main

import "fmt"

func main() {
	getmemory()
	getpointer()
	emptypointer()
	pointerarray()
	pointerpointer()

	var a,b int = 10,100
	fmt.Printf("a的值是%d b的值是%d\n", a,b)
	result := pointerfunction(&a,&b)
	fmt.Printf("pointerfunction后 a的值是%d b的值是%d result的值是%d\n", a,b,result)


	/* 变量值交换 */
	a = 100
	b = 200
	a, b = b, a

	fmt.Printf("交换后 a 的值 : %d\n", a )
	fmt.Printf("交换后 b 的值 : %d\n", b )
}

//获取变量内存指针
func getmemory(){
	var a int = 10
	fmt.Printf("getmemory a的内存指针地址是%v\n",&a )
	a = 20
	//变量的值变化，但是内存地址不变
	fmt.Printf("getmemory a的内存指针地址是%v\n",&a )
}
//-------------------------------------------------
//什么是指针
//一个指针变量指向了一个值的内存地址。
//类似于变量和常量，在使用指针前你需要声明指针。指针声明格式如下：
/*
var var_name *var-type
var-type 为指针类型，var_name 为指针变量名，* 号用于指定变量是作为一个指针。以下是有效的指针声明：

var ip *int        指向整型
var fp *float32    指向浮点型
*/
func getpointer(){
	var a int = 20      // 声明实际变量
	var address *int    // 声明指针变量

	address = &a    //获取变量的指针地址(内存地址)
	fmt.Printf("getpointer a的内存地址是%x\n", &a)

	fmt.Printf("getpointer a的指针地址是%x\n", address)

	fmt.Printf("getpointer a的值是%d\n", *address)
}

//空指针
func emptypointer(){
	var ptr *int //声明指针变量
	fmt.Printf("ptr 的值是%v\n",ptr)
}

//指针数组
const MAX int = 3
func pointerarray(){
	a := []int{10,100,1000}
	var i int
	var ptr [MAX]*int

	for i = 0;i < MAX ;i++  {
		ptr[i] = &a[i]  //把数组的元素指针地址赋值给指针数组
	}

	for i = 0;i < MAX ;i++  {
		fmt.Printf("a[%d]的值%v 是指针地址是%v\n",i,a[i], ptr[i])
	}

}

//指针的指针
/*
如果一个指针变量存放的又是另一个指针变量的地址，则称这个指针变量为指向指针的指针变量。
var ptr **int;
*/
func pointerpointer(){
	var a int
	var ptr *int
	var ptrptr **int

	a = 3

	ptr = &a //指向 a 的指针地址

	ptrptr = &ptr //指向 ptr 的指针地址

	fmt.Printf("a的值是%v\n",a )
	fmt.Printf("ptr的值是%v 指针地址是%v\n",*ptr,ptr )
	fmt.Printf("ptrptr的值是%v 指针地址是%v\n",*ptrptr,ptrptr )
	fmt.Printf("ptrptr **两次指向的值是%v \n",**ptrptr )

}

//指针作为函数参数
func pointerfunction(a *int,b *int) int {
	var result int
	result = *a
	*a = *b
	*b = result
	return result
}