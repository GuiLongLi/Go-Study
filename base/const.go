package main

import "fmt"
import "unsafe"

/*
常量是一个简单值的标识符，在程序运行时，不会被修改的量。
常量中的数据类型只可以是布尔型、数字型（整数型、浮点型和复数）和字符串型。
常量的定义格式：
const identifier type = value
你可以省略类型说明符 [type]，因为编译器可以根据变量的值来推断其类型。

显式类型定义： const b string = "abc"
隐式类型定义： const b = "abc"
多个相同类型的声明可以简写为：

const c_name1, c_name2 = value1, value2
*/

const (
	c1 = "abc"
	c2 = len(c1)
	c3 = unsafe.Sizeof(c1)
	// unsafe.Sizeof 输出结果为：16
	//字符串类型在 go 里是个结构, 包含指向底层数组的指针和长度,这两部分每部分都是 8 个字节，所以字符串类型大小为 16 个字节。
)

func main(){
	const LENGTH int = 10
	const WIDTH int = 5
	var area int
	const a,b,c = 1,false,"str" //多重赋值

	area = LENGTH * WIDTH
	fmt.Printf("面积为: %d\n", area)
	println(a,b,c)

	println(c1,c2,c3)

	printiota()

	printbinary()
}

/*
iota
iota，特殊常量，可以认为是一个可以被编译器修改的常量。
iota 在 const关键字出现时将被重置为 0(const 内部的第一行之前)，const 中每新增一行常量声明将使 iota 计数一次(iota 可理解为 const 语句块中的行索引)。
iota 可以被用作枚举值：
*/

const (
	i1 = iota
	i2 = iota
	i3 = iota
)
//上面语句可以简写为下面的
//const (
//	i1 = iota
//	i2
//	i3
//)

func printiota(){
	println(i1,i2,i3)

	const (
		p1 = iota
		p2        //没有提供初始值，表示将使用上行的表达式。 即 p2 = iota
		p3        //没有提供初始值，表示将使用上行的表达式。 即 p3 = iota
		p4 = "p4"
		p5        //没有提供初始值，表示将使用上行的表达式。 即 p5 = "p4"
		p6 = 100
		p7        //没有提供初始值，表示将使用上行的表达式。 即 p7 = 100
		p8 = iota
		p9        //没有提供初始值，表示将使用上行的表达式。 即 p9 = iota
	)
	//在定义常量组时，如果不提供初始值，则表示将使用上行的表达式。

	println(p1,p2,p3,p4,p5,p6,p7,p8,p9)
}


//二进制偏移
const (
	b1 = 1<<iota
	b2 = 3<<iota
	b3     //没有提供初始值，表示将使用上行的表达式。 即 b3 = 3<<iota
	b4     //没有提供初始值，表示将使用上行的表达式。 即 b4 = 3<<iota
)
//iota 表示从 0 开始自动加 1，所以 1<<0, 3<<1（<< 表示二进制左移的意思）
func printbinary(){

	println(b1,b2,b3,b4)
	//以上实例运行结果为： 1  6  12  24
	/*
	二进制 即 逢二进一，从最右边一位开始，值等于 2 的就要往左偏移一位
		0000 0000 = 0
		0000 0001 = 1
		0000 0010 = 2
		0000 0011 = 3
		0000 0100 = 4
		0000 0101 = 5
		0000 0110 = 6
		...如此类推
	*/

	/*
	b1 = 1<<iota
	1 = 0000 0001
	iota = 0
	即 0000 0001 往左偏移 0 位 ,得出 0000 0001
	0000 0001 = 1
---------------------
	b2 = 3<<iota
	3 = 0000 0011
	iota = 1
	即 0000 0011 往左偏移 1 位 ,得出 0000 0110
	0000 0110 = 6
---------------------
	b3 = 3<<iota
	3 = 0000 0011
	iota = 2
	即 0000 0011 往左偏移 2 位 ,得出 0000 1100
	0000 1100 = 12
---------------------
	b4 = 3<<iota
	3 = 0000 0011
	iota = 3
	即 0000 0011 往左偏移 3 位 ,得出 0001 1000
	0001 1000 = 24

	*/


}
