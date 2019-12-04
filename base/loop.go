package main

import "fmt"

func main(){
	forloop()
	forforloop()
	breakloop()
	continueloop()
	gotoloop()
}

/*
语法
Go语言的For循环有3中形式，只有其中的一种使用分号。

和 C 语言的 for 一样：
for init; condition; post { }

和 C 的 while 一样：
for condition { }

和 C 的 for(;;) 一样：
for { }

init： 一般为赋值表达式，给控制变量赋初值；
condition： 关系表达式或逻辑表达式，循环控制条件；
post： 一般为赋值表达式，给控制变量增量或减量。
for语句执行过程如下：

①先对表达式1赋初值；
②判别赋值表达式 init 是否满足给定条件，若其值为真，满足循环条件，则执行循环体内语句，然后执行 post，进入第二次循环，再判别 condition；否则判断 condition 的值为假，不满足条件，就终止for循环，执行循环体外语句。
for 循环的 range 格式可以对 slice、map、数组、字符串等进行迭代循环。格式如下：

for key, value := range oldMap {
    newMap[key] = value
}
*/
func forloop(){
	var a int = 15
	var c int

	//声明数组
	//Go 语言数组声明需要指定元素类型及元素个数，语法格式如下：
	//var variable_name [SIZE] variable_type
	// [SIZE] 填写为 [...] 则数组长度自动根据元素个数填充
	//以上为一维数组的定义方式。例如以下定义了数组 balance 长度为 10 类型为 float32
	//var balance [10] float32{1.1, 2.2, 3.3}
	//初始化数组中 {} 中的元素个数不能大于 [] 中的数字。

	//以下定义了数组 numbers 长度为 6 类型为 int 数组内元素 有 1,3,5,7
	numbers := [...]int{1,3,5,7}

	for b:=0;b<10 ;b++  {
		fmt.Printf("b的值是%v\n", b)
	}

	println()
	for c < a  {
		c++
		fmt.Printf("c的值是%v\n",c )
	}
	println()

	for key,value := range numbers  {
		fmt.Printf("numbers数组的key是%v, value是%v\n",key,value )
	}
	println()
}

//循环嵌套
func forforloop(){
	var i,j int
	var max = 100

	//输出 100以内的素数
	// 1 即不是素数又不是合数，所以初始值是2
	for i = 2;i < max ;i++  {
		//j=2 , i/j 就是 i 的 1/2
		//所有合数可能的最大因数就是 本身的 1/2
		//如果 j 是 i 的因数，就退出循环
		for j = 2;j <= (i/j);j++  {
			if(i%j==0){
				break;
			}
			//如果循环到最后都没有因数，j 就会递增到 i 的 1/2 的正整数 +1
			// 例如 i = 3 的 1/2 就是 1.5 ，正整数就是 1 ，再加1  最后得出 j = 2
		}

		// j > (i/j) 的是素数
		// j <= (i/j) 的是合数
		if(j > (i/j)){
			fmt.Printf("%d是素数\n", i)
		}
	}
	println()
}

//break  经常用于中断当前 for 循环或跳出 switch 语句
func breakloop(){
	var max int = 10
	var a int

	for a = 1;a <= max ;a++  {
		if(a > 5){ //大于5就打断循环
			break
		}
		fmt.Printf("breakloop a的值是%v\n", a)
	}
	println()
}

//continue  跳过当前循环执行下一次循环语句。
func continueloop(){
	var max int = 5
	var a int
	for a = 1;a <= max;a++  {
		if(a == 3){ //跳出当前循环，不执行下面的打印操作
			continue
		}
		fmt.Printf("continueloop a的值是%v\n", a)
	}
	println()
}

//goto 语句可以无条件地转移到过程中指定的行。
//goto 语句通常与条件语句配合使用。可用来实现条件转移， 构成循环，跳出循环体等功能。
//但是，在结构化程序设计中一般不主张使用 goto 语句， 以免造成程序流程的混乱，使理解和调试程序都产生困难。
func gotoloop(){
	var max int = 10
	var a int = 1

	LOOP:for a < max {
		if a == 5 {
			a++
			goto LOOP
		}
		fmt.Printf("gotoloop a的值是%v\n", a)
		a++
	}
	println()
}