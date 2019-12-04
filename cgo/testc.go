package main

import (
	"fmt"
)

/*
在很多场景下，在Go的程序中需要调用c函数或者是用c编写的库（底层驱动，算法等，不想用Go语言再去造一遍轮子，复用现有的c库）。
那么该如何调用呢？Go可是更好的C语言啊，当然提供了和c语言交互的功能，称为Cgo！
Cgo封装了#cgo伪c文法，参数CFLAGS用来传入编译选项，LDFLAGS来传入链接选项。这个用来调用非c标准的第三方c库。

1）先从最简单的写起吧，Go代码直接调用c函数，下面的示例中在代码注释块调用了标准的c库，并写了一个c函数(本例只是简单打印了一句话，在该注释块中可以写任意合法的c代码)，在Go代码部分直接调用该c函数hi()
*/



/*
//------------------------------
//直接调用c代码
#include <stdio.h>

void hi(){
	printf("hello world!\n");
}

//--------------------------------
//调用非标准的c的第三方动态库
#cgo CFLAGS: -I./
//CFLAGS 中的 -I 大写i 参数表示 .h 头文件所在的路径
#cgo LDFLAGS: -L./ -lhiso
//LDFLAGS 中的 -L 大写 表示.so文件所在的路径  -l(小写的L)表示指定该路径下的库的名称,
//比如要使用libhi.so ，则只需用-lhi (省略了libhi.so 中的lib 和 .so 字符)
#include "hiso.h"
//非标准c头文件，所有用引号
//编译成动态库.so
//gcc -c -fPIC -o hiso.o hiso.c
//gcc -shared -o libhiso.so hiso.o

//--------------------------------
//调用非标准的c的第三方静态库
#cgo CFLAGS: -I./
//CFLAGS 中的 -I 大写i 参数表示 .h 头文件所在的路径
#cgo LDFLAGS: -L./ -lhia
//LDFLAGS 中的 -L 大写 表示.so文件所在的路径  -l(小写的L)表示指定该路径下的库的名称,
//比如要使用libhi.so ，则只需用-lhi (省略了libhi.so 中的lib 和 .so 字符)
#include "hia.h"
//非标准c头文件，所有用引号
//编译生成静态库.a
//gcc -c hia.c
//ar rv libhia.a hia.o

*/
import "C"
//这里可看做封装的伪包C ,这条语句要紧挨着上面的注释块，不可在它俩之间间隔空行！

//直接调用c代码
func testc(){
	C.hi()
	fmt.Println("hi,go C")
}

//调用非标准的c的第三方动态库
// go build -ldflags="-r ./" testc.go
// -ldflags 是 .so 动态库的路径
func testso(){
	C.hiso()
	fmt.Println("hi,go c-so")
}

//调用非标准的c的第三方静态库
func testa(){
	C.hia()
	fmt.Println("hi,go c-a")
}

func main() {
	fmt.Println("testc")
	testc()
	fmt.Println()
	fmt.Println("testso")
	testso()
	fmt.Println()
	fmt.Println("testa")
	testa()
}
