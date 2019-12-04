package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

var mylooger *log.Logger
func main() {
	Ldate()
	Ltime()
	Lmicroseconds()
	Llongfile()
	Lshortfile()
	LUTC()
	Ldefault()
	Lstander()
	LprintInfo()
	test_deferpanic()
	//test_fatalln()  //测试终止程序

	mylooger = LogInfo();

	//Prefix返回前缀，Flags返回Logger的输出选项属性值
	fmt.Printf("创建时前缀为:%s\n创建时输出项属性值为:%d\n",mylooger.Prefix(),mylooger.Flags())

	//SetFlags 重新设置输出选项
	mylooger.SetFlags(log.Ldate|log.Ltime|log.Llongfile)

	//重新设置输出前缀
	mylooger.SetPrefix("mylooger_")

	//获取修改后的前缀和输出项属性值
	fmt.Printf("修改后前缀为:%s\n修改后输出项属性值为:%d\n",mylooger.Prefix(),mylooger.Flags())

	//输出日志
	mylooger.Output(2,"使用Output进行日志输出")

	//格式化输出日志
	mylooger.Printf("我是%v方法在%d行内容为:%s","Printf",40,"其实我底层是以fmt.Printf的方式处理的，相当于Java里的Info级别")

	//开启这个注释，下面代码就不会继续走，并且程序停止
	//mylooger.Fatal("我是Fatal方法，我会停止程序，但不会抛出异常")

	//调用业务层代码
	serviceCode()

	mylooger.Printf("业务代码里的Panicln不会影响到我，因为他已经被处理干掉了，程序目前正常")

	//切换输出屏幕中
	mylooger.SetOutput(os.Stdout)
	mylooger.Printf("切换屏幕进行日志输出")
}

func Ldate(){
	log.SetFlags(log.Ldate)
	log.Println("这是Ldate格式\n")
}

func Ltime(){
	log.SetFlags(log.Ltime)
	log.Println("这是Ltime格式\n")
}

func Lmicroseconds(){
	log.SetFlags(log.Lmicroseconds)
	log.Println("这是Lmicroseconds格式\n")
}
func Llongfile(){
	log.SetFlags(log.Llongfile)
	log.Println("这是Llongfile格式\n")
}
func Lshortfile(){
	log.SetFlags(log.Lshortfile)
	log.Println("这是Lshortfile格式\n")
}
func LUTC(){
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC)
	log.Println("这是LUTC格式\n")
}
func Ldefault(){
	log.Println("这是默认格式\n")
}
//一个文件标准日志格式
func Lstander(){
	log.SetFlags(log.Llongfile | log.Ldate | log.Ltime)
	log.Println("这是日志标准格式\n")
}

//信息打印
func LprintInfo(){
	arr := []int {2,3}
	log.Print("Print array ",arr,"\n")
	log.Println("Println array",arr)
	log.Printf("Printf array with item [%d,%d]\n",arr[0],arr[1])
}

//抛出异常
func test_deferpanic(){
	//defer语句调用一个函数，这个函数执行会推迟，直到外围的函数返回，或者外围函数运行到最后，或者相应的goroutine panic
	defer func(){ //这个defer 函数 ，将会在 log.Panicln 后才会被调用
		log.Println("-- first --")
		//recover 函数会捕获 panic 抛出的异常信息
		if err := recover();err != nil{
			log.Println(err)
		}
	}()
	log.Panicln("test for defer panic") //抛出异常后，后面的操作将不会被执行
	defer func() {
		log.Println("-- second --")
	}()
}

//终止程序
//对于 log.Fatal 接口，会先将日志内容打印到标准输出，接着调用系统的 os.exit(1) 接口，退出程序并返回状态 1 。但是有一点需要注意，由于是直接调用系统接口退出，defer函数不会被调用，
func test_fatalln(){
	defer func() {
		//这个defer 函数 ，将会在 log.Fatalln 后才会被调用 ,由于 fatalln 直接终止程序的后续运行，所以 Println将不会运行
		log.Println("-- first --")
	}()
	log.Fatalln("test for defer fatalln")
}

func LogInfo() *log.Logger {
	//创建文件对象, 日志的格式为当前时间2006-01-02 15:04:05.log;据说是golang的诞生时间，固定写法
	timeString := time.Now().Format("2006-01-02")
	file := "./"+timeString+".log"

	logFile,err := os.OpenFile(file,os.O_RDWR|os.O_CREATE|os.O_APPEND,0766)
	if(err != nil){
		panic(err)
	}

	//创建一个Logger, 参数1：日志写入的文件, 参数2：每条日志的前缀；参数3：日志属性
	return log.New(logFile,"logpre_",log.Lshortfile)
}

func serviceCode()  {
	defer func() {
		if r := recover(); r != nil {
			//用以捕捉Panicln抛出的异常
			fmt.Printf("使用recover()捕获到的错误：%s\n", r)
		}
	}()
	// 模拟错误业务逻辑，使用抛出异常和捕捉的方式记录日志
	mylooger.Panicln("我是Panicln方法，我会抛异常信息的，相当于Error级别")
}

