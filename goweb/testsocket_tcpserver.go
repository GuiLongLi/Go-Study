package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("test_tcpserver")
	//test_tcpserver(); //单任务形式监听端口

	fmt.Println()
	fmt.Println("test_tcpservergoruntine")
	//test_tcpservergoruntine(); //支持多并发,使用goroutine机制，


	fmt.Println()
	fmt.Println("test_tcpservergoruntinemore")
	test_tcpservergoruntinemore(); //支持多并发,使用goroutine机制，不同的请求来获取不同的时间格式，而且需要一个长连接
}


func test_tcpserver(){
	/*
	TCP server
	上面我们编写了一个TCP的客户端程序，也可以通过net包来创建一个服务器端程序，在服务器端我们需要绑定服务到指定的非激活端口，并监听此端口，当有客户端请求到达的时候可以接收到来自客户端连接的请求。net包中有相应功能的函数，函数定义如下：

	func ListenTCP(network string, laddr *TCPAddr) (*TCPListener, error)
	func (l *TCPListener) Accept() (Conn, error)
	*/

	/*
	参数说明同DialTCP的参数一样。下面我们实现一个简单的时间同步服务，监听7777端口
	*/
	service := ":12701"
	tcpAddr,err := net.ResolveTCPAddr("tcp4",service) //解析地址
	checkTcpServerErr(err)
	listener,err := net.ListenTCP("tcp",tcpAddr) //监听端口
	checkTcpServerErr(err)
	for{ //无限循环，监听客户端请求
		conn,err := listener.Accept() //接收数据
		if err != nil{
			continue;
		}
		daytime := time.Now().String()
		conn.Write([]byte(daytime)) //响应二进制的日期
		conn.Close() //关闭请求连接
	}

}

func test_tcpservergoruntine(){
	service := ":12701"
	tcpAddr,err := net.ResolveTCPAddr("tcp4",service)
	checkTcpServerErr(err)
	listener,err := net.ListenTCP("tcp",tcpAddr)
	checkTcpServerErr(err)
	for{
		conn,err := listener.Accept()
		if err != nil{
			continue
		}
		go handlerClient(conn) //使用goruntine 多并发处理客户端请求
	}
}

func handlerClient(conn net.Conn){
	defer conn.Close()
	daytime := time.Now().String()
	conn.Write([]byte(daytime))

}

func test_tcpservergoruntinemore(){
	service := ":12701";
	tcpAddr,err := net.ResolveTCPAddr("tcp4",service)
	checkTcpServerErr(err)
	listener,err := net.ListenTCP("tcp",tcpAddr)
	checkTcpServerErr(err)
	for{
		conn,err := listener.Accept()
		if err != nil{
			continue
		}
		go handlerClientMore(conn)
	}
}

func handlerClientMore(conn net.Conn){
	conn.SetReadDeadline(time.Now().Add(2*time.Minute)) //设置2分钟后超时
	request := make([]byte,128) //设置最大的请求长度为128b，为了防止洪水攻击
	defer conn.Close() //最后执行自动关闭连接
	for{
		read_len,err := conn.Read(request)

		if err != nil{
			fmt.Println(err)
			break
		}

		if read_len == 0{
			break //客户端已经关闭连接
		}else if strings.TrimSpace(string(request[:read_len])) == "timestamp"{
			daytime := strconv.FormatInt(time.Now().Unix(),10) //格式化时间戳为数字
			conn.Write([]byte(daytime)) //响应时间戳给客户端
		}else{
			daytime := time.Now().String() //日期形式
			conn.Write([]byte(daytime))
		}
		request = make([]byte,128) //清除最后读取的请求内容
	}
}

func checkTcpServerErr(e error){
	if e != nil{
		fmt.Fprintf(os.Stderr,"fata error: %s",e.Error()) //把错误信息输出到屏幕
		os.Exit(1)
	}
}