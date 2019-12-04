package main

import "fmt"

/*
声明数组
Go 语言数组声明需要指定元素类型及元素个数，语法格式如下：
var variable_name [SIZE] variable_type
以上为一维数组的定义方式。例如以下定义了数组 balance 长度为 10 类型为 float32：
var balance [10] float32
*/
func main() {
	initarray()
	visitarray()
	multidimenarray()

	//数组作为参数
	var arr = []int{10,20,30,50}
	var avg float32
	avg = getAverage(arr,len(arr))
	fmt.Printf("arr数组的平均数是%f\n", avg)

	a := 1.69           // 表示1.69
	b := 1.70           // 表示1.70
	c := a * b          // 结果应该表示 2.873
	fmt.Println(float64(c))// 输出的是2.8729999999999998

	//a 和 b转换为正整数 再计算
	//把 a * 100    b * 100
	//再重新计算
	a = a*100
	b = b*100
	c = a * b
	fmt.Println(c)
	fmt.Println(float64(c)/(100*100)) //把前面取整数时乘以的 100 *100 除去 ,最终得出结果


}

//初始化数组
func initarray(){
	//以下演示了数组初始化：

	var balance1 = [5]float32{1000.0, 2.0, 3.4, 7.0, 50.0}
	//初始化数组中 {} 中的元素个数不能大于 [] 中的数字。
	fmt.Println(balance1)

	//如果忽略 [] 中的数字不设置数组大小，Go 语言会根据元素的个数来设置数组的大小：

	//该实例与上面的实例是一样的，虽然没有设置数组的大小。
	var balance2 = [...]float32{1000.0, 2.0, 3.4, 7.0, 50.0}
	fmt.Println(balance2)
	balance2[4] = 100000.0
	fmt.Println(balance2)
	println()
}

//访问数组元素
func visitarray(){
	var n [10]int
	var i,j int

	//插入数组数据
	for i=0;i<10;i++{
		n[i] = 100+i
	}

	//取出数组数据
	for j=0;j<10 ;j++  {
		fmt.Printf("j的值是%v n[j]的值是%v\n", j,n[j])
	}
	println()
}

//多维数组
/*
Go 语言支持多维数组，以下为常用的多维数组声明方式：
var variable_name [SIZE1][SIZE2]...[SIZEN] variable_type
以下实例声明了三维的整型数组：
var threedim [5][10][4]int
*/
func multidimenarray(){
	var twodimenarray = [3][4]int{
		{0,1,3,5},
		{2,4,6,8},
		{10,11,13,15},
	}
	//注意：以上代码中倒数第二行的 } 必须要有逗号，因为最后一行的 } 不能单独一行，也可以写成这样：
	/*
	var twodimenarray = [3][4]int{
		{0,1,3,5},
		{2,4,6,8},
		{10,11,13,15}}
	*/
	fmt.Printf("twodimenarray 的值是%v\n", twodimenarray)

	twovalue1 := twodimenarray[2][3]
	fmt.Printf("twovalue1的值是%v\n", twovalue1)

	//输出数组各个元素
	var i,j int
	for i=0;i<len(twodimenarray) ;i++  {
		for j=0;j<len(twodimenarray[i]) ;j++  {
			fmt.Printf("i的值是%d j的值是%d twodimenarray[i][j]的值是%d\n", i,j,twodimenarray[i][j])
		}
	}

}

//----------------------------------------------
//向函数传递数组
//如果你想向函数传递数组参数，你需要在函数定义时，声明形参为数组，我们可以通过以下两种方式来声明：
/*
方式一
形参设定数组大小：
func myFunction1(param [10]int)

方式二
形参未设定数组大小：
func myFunction2(param []int)
*/
func getAverage(arr []int,size int) float32{
	var i,sum int
	var avg float32

	for i = 0;i < size ;i++  {
		sum += arr[i]
	}
	avg = float32(sum)/float32(size)
	return avg
}
