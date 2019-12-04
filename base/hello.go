package main

import "fmt"

var Group1 = "Group1"

func main() {
	var word = "hello world"
	fmt.Println(word)
	fmt.Println(Group1)
}
//在系统中运行  go run hello.go  即可看到输出结果

//包声明 ，第一行代码 package main 定义了包名。你必须在源文件中，非注释的第一行指明这个文件属于哪个包，例如 package main 。 package main 表示一个可独立执行的程序，每个go应用程序都包含一个名为main 的包.

//引入包 ， import "fmt" 告诉go编译器需要使用 fmt 包(的函数或其他元素)，fmt包实现了格式化io (输入/输出)的函数.

//函数，func main()是程序开始执行的函数. main函数是每一个可执行程序所必须包含的，一般来说都是在启动后第一个执行的函数(如果存在 init() 函数，则会优先执行该函数)

// fmt.Println(....) 可以将字符串输出到控制台，并在最后自动增加换行符 \n.
// 使用 fmt.Print("hello world\n") 可以得到相同的结果

//标识符（常量、变量、类型、函数名、结构字段等等）以一个大写字母开头，如：Group1 ,那么使用这种形式的标识符的对象就可以被外部包的代码所使用（客户端程序需要先导入这个包）,这被称为导出(像面向对象语言中的public ),标识符如果以小写字母开头 如：word，则对包外是不可见的，但是他们在整个包的内部是可见并且可用的（就像面向对象语言中的protected）