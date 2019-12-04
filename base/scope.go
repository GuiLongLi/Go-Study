package main

import "fmt"

//初始化局部和全局变量
//不同类型的局部和全局变量默认值为：
//
//数据类型	初始化默认值
//int	0
//float32	0
//pointer	nil
func main() {
	localvar()
	overallvar()
	definitionvar()
}

//局部变量
//在函数体内声明的变量称之为局部变量，它们的作用域只在函数体内，参数和返回值变量也是局部变量。
func localvar(){
	/* 声明局部变量 */
	var a, b, c int

	/* 初始化参数 */
	a = 10
	b = 20
	c = a + b

	fmt.Printf ("localvar 结果： a = %d, b = %d and c = %d\n", a, b, c)
}

//全局变量
//在函数体外声明的变量称之为全局变量，全局变量可以在整个包甚至外部包（被导出后）使用。
/* 声明全局变量 */
var g int
func overallvar(){
	/* 声明局部变量 */
	var a, b int

	/* 初始化参数 */
	a = 10
	b = 20
	g = a + b
	fmt.Printf("overallvar 结果： a = %d, b = %d and g = %d\n", a, b, g)

	//Go 语言程序中全局变量与局部变量名称可以相同，但是函数内的局部变量会被优先考虑。实例如下：
	/* 声明局部变量 */
	var g int = 10
	fmt.Printf ("overallvar 结果： g = %d\n",  g)
}

//形式参数
//形式参数会作为函数的局部变量来使用。实例如下：
/* 声明全局变量 */
var a int = 20;
func definitionvar(){
	var a int = 10
	var b int = 20
	var c int = 0

	fmt.Printf("definitionvar 函数中 a = %d\n",  a);
	c = sum( a, b);
	fmt.Printf("definitionvar 函数中 c = %d\n",  c);
}

/* 函数定义-两数相加 */
func sum(a, b int) int {
	fmt.Printf("sum() 函数中 a = %d\n",  a);
	fmt.Printf("sum() 函数中 b = %d\n",  b);

	return a + b;
}