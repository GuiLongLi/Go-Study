package main
import "fmt"

//声明变量的一般形式是使用 var 关键字：
func main() {
	//var identifier type
	//	可以一次声明多个变量：
	var a string = "run go"
	fmt.Println(a)

	//var identifier1, identifier2 type
	var b,c int = 1,2
	fmt.Println(b,c)


	//变量声明
	//第一种，指定变量类型，如果没有初始化，则变量默认为零值。
	//var v_name v_type
	//v_name = value
	var no_value1 string
	fmt.Println(no_value1)//输出 空

	var no_value2 int
	fmt.Println(no_value2)//输出 0

	var no_value3 bool
	fmt.Println(no_value3)//输出 false

	var no_value4 map[string] int
	fmt.Println(no_value4)

	var no_value5 []int
	fmt.Println(no_value5)


	//以下几种类型为 nil：
	var no_value6 *int
	fmt.Println(no_value6)

	var no_value7 chan int
	fmt.Println(no_value7)

	var no_value8 func(string) int
	fmt.Println(no_value8)

	var no_value9 error // error 是接口
	fmt.Println(no_value9)

	//第二种，根据值自行判定变量类型。
	var auto_value = true
	fmt.Println(auto_value)

	//第三种，省略 var, 注意 := 左侧如果没有声明新的变量，就产生编译错误，格式：
	//no_var := value;
	//var f string = "Runoob" 简写为 f := "Runoob"：
	no_var1, no_var2, no_var3 := 1, 2, 3;
	fmt.Printf("%v %v\n",no_var1, no_var2, no_var3)

	//类型相同多个变量, 非全局变量
	//var vname1, vname2, vname3 type
	//	vname1, vname2, vname3 = v1, v2, v3
	var no_var4, no_var5, no_var6 = no_var1, no_var2, no_var3 // 和 python 很像,不需要显示声明类型，自动推断
	fmt.Printf("%v %v %v\n",no_var4,no_var5,no_var6)

	no_var7, no_var8, no_var9 := no_var1, no_var2, no_var3 // 出现在 := 左侧的变量不应该是已经被声明过的，否则会导致编译错误
	fmt.Printf("%v %v %v\n",no_var7,no_var8,no_var9)

	// 这种因式分解关键字的写法一般用于声明全局变量
	//var (
	//	vname1 v_type1
	//	vname2 v_type2
	//)
	var(
		no_var10 int
		no_var11 string
	)
	fmt.Printf("%v '%v'\n",no_var10, no_var11)

	//函数与变量声明使用
	more()


	// _ 空白标识符在函数返回值时的使用：
	_,number,string := numbers() // _ 空白标识符的值 被丢弃 //只获取函数返回值的后两个
	println(number,string)
}

var x, y int
var (  // 这种因式分解关键字的写法一般用于声明全局变量
	a int
	b bool
)

var c, d int = 1, 2
var e, f = 123, "hello"

//这种不带声明格式的只能在函数体中出现
//g, h := 123, "hello"
func more(){
	g, h := 123, "hello"
	println(x, y, a, b, c, d, e, f, g, h)
}


//空白标识符 _ 也被用于抛弃值，
// _, b = 5, 7
// 如值 5 上面的语句中被抛弃。
//_ 实际上是一个只写变量，你不能得到它的值。这样做是因为 Go 语言中你必须使用所有被声明的变量，但有时你并不需要使用从一个函数得到的所有返回值。
//空白标识符在函数返回值时的使用：
func numbers()(int,int,string){
	a,b,c := 1,2,"abc"
	return a,b,c
}