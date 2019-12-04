package main

import "fmt"

func main() {
	initslice()

	//获取切片长度和容量
	var slice1 = make([]int,3,5)
	var len1,cap1 int
	len1,cap1 = getslicelencap(slice1)
	fmt.Printf("slice1的值是%v 的长度是%v 容量是%v\n",slice1, len1,cap1)

	//空切片 一个切片在未初始化之前默认为 nil，长度为 0，
	var slice2 []int
	len2,cap2 := getslicelencap(slice2)
	fmt.Printf("slice2的值是%v 的长度是%v 容量是%v\n",slice2, len2,cap2)

	//切片截取
	sliceintercept()

	//切片插入和复制
	appendcopyslice()
}

//定义切片
/*
你可以声明一个未指定大小的数组来定义切片：切片不需要说明数组长度。
var identifier []type

或使用make()函数来创建切片:
var slice1 []type = make([]type, length)

也可以简写为
slice1 := make([]type, length)

也可以指定容量，其中capacity为可选参数。
make([]type, length, capacity)

这里 length 是数组的长度并且也是切片的初始长度。
*/
func initslice(){
	//直接初始化切片 s = []切片类型 {1,2,3}
	// 初始化值依次是1,2,3.其capacity=length=3
	s := []int{1,2,3}
	fmt.Printf("s 的值是%v\n", s)

	//初始化切片s1 ,引用数组arr
	var arr = [...]int {1,3,5,7,9}
	s1 := arr[:]
	fmt.Printf("s1 的值是%v\n", s1)

	//通过内置函数make()初始化切片s5, make([]type, length, capacity)
	s5 := make([]int,4,5)
	fmt.Printf("s5 的值是%v\n", s5)
}

//len() 和 cap() 函数
//切片是可索引的，并且可以由 len() 方法获取长度。
//切片提供了计算容量的方法 cap() 可以测量切片最长可以达到多少。
func getslicelencap(slice []int) (int,int){
	return len(slice),cap(slice)
}

//切片截取
//可以通过设置下限及上限来设置截取切片 [lower-bound:upper-bound]，
func sliceintercept(){
	var arr = [...]int {1,3,5,7,9}
	fmt.Printf("arr 的值是%v\n", arr)
	//s := arr[startIndex:endIndex] -------------------------------------------
	//将arr中从下标startIndex到endIndex-1 下的元素创建为一个新的切片
	s2 := arr[0:1]
	fmt.Printf("s2 的值是%v\n", s2)

	//缺省endIndex时将表示一直到arr的最后一个元素
	s3 := arr[2:]
	fmt.Printf("s3 的值是%v\n", s3)

	//缺省startIndex时将表示从arr的第一个元素开始
	s4 := arr[:3]
	fmt.Printf("s4 的值是%v\n", s4)

}

//切片插入和复制
/*
append() 和 copy() 函数
如果想增加切片的容量，我们必须创建一个新的更大的切片并把原分片的内容都拷贝过来。
*/
func appendcopyslice(){
	var slice1 []int
	var len1,cap1 int
	var slice2 []int
	var len2,cap2 int

	len1,cap1 = getslicelencap(slice1)
	fmt.Printf("slice1的值是%v 的长度是%v 容量是%v\n",slice1, len1,cap1)

	//插入空切片
	slice1 = append(slice1,0)
	len1,cap1 = getslicelencap(slice1)
	fmt.Printf("slice1的值是%v 的长度是%v 容量是%v\n",slice1, len1,cap1)

	//插入一个元素
	slice1 = append(slice1,1)
	len1,cap1 = getslicelencap(slice1)
	fmt.Printf("slice1的值是%v 的长度是%v 容量是%v\n",slice1, len1,cap1)

	//插入多个元素
	slice1 = append(slice1,2,3,4)
	len1,cap1 = getslicelencap(slice1)
	fmt.Printf("slice1的值是%v 的长度是%v 容量是%v\n",slice1, len1,cap1)

	//创建一个两倍容量的切片
	slice2 = make([]int,len(slice1),(cap(slice1)*2))
	len2,cap2 = getslicelencap(slice2)
	fmt.Printf("slice2的值是%v 的长度是%v 容量是%v\n",slice2, len2,cap2)

	//复制 slice1 到 slice2
	copy(slice2,slice1)
	len2,cap2 = getslicelencap(slice2)
	fmt.Printf("slice2的值是%v 的长度是%v 容量是%v\n",slice2, len2,cap2)

}