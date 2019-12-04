package main

import "fmt"

func main() {
	start := 3
	//阶乘 1*2*3*4*5* ... n
	result := factorial(uint64(start))
	fmt.Printf("%v的阶乘是%v\n",start,result )

	//斐波那契数列
	var i int
	for i = 0; i < 10; i++ {
		fmt.Printf("%d\t", fibonacci(i))
	}
	println()
}

//递归，就是在运行的过程中调用自己。
//阶乘
//以下实例通过 Go 语言的递归函数实例阶乘：
func factorial(n uint64) (result uint64){
	if(n > 0){
		result = n * factorial(n-1)
		return result
	}
	return 1
}

//斐波那契数列
//以下实例通过 Go 语言的递归函数实现斐波那契数列：
func fibonacci(n int) int {
	if n < 2 {
		return n
	}
	return fibonacci(n-2) + fibonacci(n-1)
}