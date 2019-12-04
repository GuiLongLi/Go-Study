package main

import "fmt"

//条件语句需要开发者通过指定一个或多个条件，并通过测试条件是否为 true 来决定是否执行指定语句，并在条件为 false 的情况在执行另外的语句。
//
//下图展示了程序语言中条件语句的结构：

func main(){
	ifcondition()
	ifelsecondition()
	ififcondition()
	switchcondition()
	selectcondition()
	passwords()
}

// if
func ifcondition(){
	var a int = 10
	if (a <= 10){
		fmt.Printf("a的值小于等于10\n")
	}
	fmt.Printf("a的值是%v\n", a)
}

//if else
func ifelsecondition(){
	var b int = 100
	if(b <= 10){
		fmt.Printf("b的值小于等于10\n")
	}else{
		fmt.Printf("b的值大于10\n")
	}
	fmt.Printf("b的值是%v\n", b)
}

//if 内嵌套 if
func ififcondition(){
	var c int = 100
	var d int = 200
	if(c == 100){
		if(d == 200){
			fmt.Printf("c的值等于100，同时d的值等于200\n" )
		}
	}
	fmt.Printf("c的值是%v ,d的值是%v\n",c,d )
}

//switch
func switchcondition(){
	var grade string
	var mark int = 90
	switch mark {
	case 90:grade = "A"
	case 80:grade = "B"
	case 70,60,50:grade = "C"
	default:grade = "D"
	}

	switch {
	case grade == "A":
		fmt.Printf("优秀\n" )
	case grade == "B",grade == "C":
		fmt.Printf("良好\n")
	case grade == "D":
		fmt.Printf("及格\n")
	case grade == "F":
		fmt.Printf("差\n")
	default:
		fmt.Printf("不及格\n")
	}
	fmt.Printf("你的分数是%v,评级是%v\n", mark,grade)

	//switch 语句还可以被用于 type-switch 来判断某个 interface 变量中实际存储的变量类型。
	//Type Switch 语法格式如下：
	var x interface{}

	switch i := x.(type) {
	case nil:
		fmt.Printf(" x 的类型 :%T\n",i)
	case int:
		fmt.Printf("x 是 int 型\n")
	case float64:
		fmt.Printf("x 是 float64 型\n")
	case func(int) float64:
		fmt.Printf("x 是 func(int) 型\n")
	case bool, string:
		fmt.Printf("x 是 bool 或 string 型\n")
	default:
		fmt.Printf("未知型\n")
	}

	//fallthrough
	//使用 fallthrough 会强制执行后面的 case 语句，fallthrough 不会判断下一条 case 的表达式结果是否为 true。
	switch {
	case false:
		fmt.Println("1、case 条件语句为 false")
		fallthrough
	case true:
		fmt.Println("2、case 条件语句为 true")
		fallthrough
	case false:
		fmt.Println("3、case 条件语句为 false")
		fallthrough
	case true:
		fmt.Println("4、case 条件语句为 true") //没有fallthrough ，将会打断后面的case操作
	case false:
		fmt.Println("5、case 条件语句为 false")
		fallthrough
	default:
		fmt.Println("6、默认 case")
	}
}

//select 是 Go 中的一个控制结构，类似于用于通信的 switch 语句。每个 case 必须是一个通信操作，要么是发送要么是接收。
//select 随机执行一个可运行的 case。如果没有 case 可运行，它将阻塞，直到有 case 可运行。一个默认的子句应该总是可运行的。
//select
/*
Go 编程语言中 select 语句的语法如下：

select {
    case communication clause  :
       statement(s);
    case communication clause  :
       statement(s);
     你可以定义任意数量的 case
	default :  可选
	statement(s);
}
以下描述了 select 语句的语法：

每个 case 都必须是一个通信
所有 channel 表达式都会被求值
所有被发送的表达式都会被求值
如果任意某个通信可以进行，它就执行，其他被忽略。
如果有多个 case 都可以运行，Select 会随机公平地选出一个执行。其他不会执行。
否则：
如果有 default 子句，则执行该语句。
如果没有 default 子句，select 将阻塞，直到某个通信可以运行；Go 不会重新对 channel 或值进行求值。
*/
func selectcondition(){
	var c1, c2, c3 chan int
	var i1, i2 int
	select {
	case i1 = <-c1:
		fmt.Printf("received ", i1, " from c1\n")
	case c2 <- i2:
		fmt.Printf("sent ", i2, " to c2\n")
	case i3, ok := (<-c3):  // same as: i3, ok := <-c3
		if ok {
			fmt.Printf("received ", i3, " from c3\n")
		} else {
			fmt.Printf("c3 is closed\n")
		}
	default:
		fmt.Printf("no communication\n")
	}
}


//输入密码
func passwords(){
	var a int
	var b int
	fmt.Printf("请输入密码：   \n")
	fmt.Scan(&a)
	if a == 123 {
		fmt.Printf("请再次输入密码：\n")
		fmt.Scan(&b)
		if b == 123 {
			fmt.Printf("密码正确，门锁已打开\n")
		}else{
			fmt.Printf("非法入侵，已自动报警\n")
		}
	}else{
		fmt.Printf("非法入侵，已自动报警\n")
	}
}