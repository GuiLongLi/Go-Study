package main

import "fmt"

func main() {
	rangesum()
	rangeother()
}

//range 求和
func rangesum(){
	arr := []int{1,2,3,4,5}
	sum := 0

	for _,val := range arr  {
		sum += val
	}
	fmt.Printf("rangesum 的值是%v\n", sum)
}

//range其他用法
func rangeother(){
	arr := []int{1,2,3,4,5}
	//在数组上使用range将传入index和值两个变量。上面那个例子我们不需要使用该元素的序号，所以我们使用空白符"_"省略了。有时侯我们确实需要知道它的索引。
	for i, num := range arr {
		if num == 3 {
			fmt.Println("index:", i)
		}
	}
	//range也可以用在map的键值对上。
	kvs := map[string]string{"a": "apple", "b": "banana"}
	for k, v := range kvs {
		fmt.Printf("%s -> %s\n", k, v)
	}
	//range也可以用来枚举Unicode字符串。第一个参数是字符的索引，第二个是字符（Unicode的值）本身。
	for i, c := range "go" {
		fmt.Println(i, c)
	}
}