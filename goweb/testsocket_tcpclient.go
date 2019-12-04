package main

import (
	"fmt"
	"net"
	"io/ioutil"
)

/*
Socket起源于Unix，而Unix基本哲学之一就是“一切皆文件”，都可以用“打开open –> 读写write/read –> 关闭close”模式来操作。

Socket基础知识
常用的Socket类型有两种：流式Socket（SOCK_STREAM）和数据报式Socket（SOCK_DGRAM）。流式是一种面向连接的Socket，针对于面向连接的TCP服务应用；数据报式Socket是一种无连接的Socket，对应于无连接的UDP服务应用。

通过上面的介绍我们知道Socket有两种：TCP Socket和UDP Socket，TCP和UDP是协议，而要确定一个进程的需要三元组，需要IP地址和端口。
*/
func main() {
	fmt.Println("testip")
	test_tcpclientip()

	fmt.Println()
	fmt.Println("test_tcpclient")
	test_tcpclient()
}

func test_tcpclientip()  {
	/*
	Go支持的IP类型
	在Go的net包中定义了很多类型、函数和方法用来网络编程，其中IP的定义如下：
	type IP []byte
	*/
	/*
	在net包中有很多函数来操作IP，但是其中比较有用的也就几个，其中ParseIP(s string) IP函数会把一个IPv4或者IPv6的地址转化成IP类型，请看下面的例子:
	*/
	ip := "10.10.87.1"
	addr := net.ParseIP(ip)
	if addr == nil{
		fmt.Println("无效的地址：",ip)
	}else{
		fmt.Println("地址：",addr.String())
	}
}


func test_tcpclient(){
	/*
	TCP Socket

	在Go语言的net包中有一个类型TCPConn，这个类型可以用来作为客户端和服务器端交互的通道，他有两个主要的函数：
	func (c *TCPConn) Write(b []byte) (int, error)
	func (c *TCPConn) Read(b []byte) (int, error)

	TCPConn可以用在客户端和服务器端来读写数据。
	还有我们需要知道一个TCPAddr类型，他表示一个TCP的地址信息，他的定义如下：
	type TCPAddr struct {
		IP IP
		Port int
		Zone string // IPv6 scoped addressing zone
	}

	在Go语言中通过net包中的ResolveTCPAddr获取一个TCPAddr
	func ResolveTCPAddr(net, addr string) (*TCPAddr, os.Error)
		net参数是"tcp4"、"tcp6"、"tcp"中的任意一个，分别表示TCP(IPv4-only), TCP(IPv6-only)或者TCP(IPv4, IPv6的任意一个)。
		addr表示域名或者IP地址，例如"www.google.com:80" 或者"127.0.0.1:22"。
	*/

	/*
	TCP client
	Go语言中通过net包中的DialTCP函数来建立一个TCP连接，并返回一个TCPConn类型的对象，当连接建立时服务器端也创建一个同类型的对象，此时客户端和服务器段通过各自拥有的TCPConn对象来进行数据交换。一般而言，客户端通过TCPConn对象将请求信息发送到服务器端，读取服务器端响应的信息。服务器端读取并解析来自客户端的请求，并返回应答信息，这个连接只有当任一端关闭了连接之后才失效，不然这连接可以一直在使用。建立连接的函数定义如下：

	func DialTCP(network string, laddr, raddr *TCPAddr) (*TCPConn, error)
		net参数是"tcp4"、"tcp6"、"tcp"中的任意一个，分别表示TCP(IPv4-only)、TCP(IPv6-only)或者TCP(IPv4,IPv6的任意一个)
		laddr表示本机地址，一般设置为nil
		raddr表示远程的服务地址
	*/

	/*
	接下来我们写一个简单的例子，模拟一个基于HTTP协议的客户端请求去连接一个Web服务端。我们要写一个简单的http请求头，格式类似如下：
		"HEAD / HTTP/1.0\r\n\r\n"

	从服务端接收到的响应信息格式可能如下：
		HTTP/1.0 200 OK
		ETag: "-9985996"
		Last-Modified: Thu, 25 Mar 2010 17:51:10 GMT
		Content-Length: 18074
		Connection: close
		Date: Sat, 28 Aug 2010 00:43:48 GMT
		Server: lighttpd/1.4.23
	*/

	//客户端发送 tcp 请求代码
	service := "127.0.0.1:12701"
	tcpAddr,err := net.ResolveTCPAddr("tcp4",service) //创建一个 tcp4 的解析地址
	checkTcpClientErr(err)
	conn, err := net.DialTCP("tcp",nil,tcpAddr) //通过解析地址建立 tcp 连接
	checkTcpClientErr(err)

	requestbody := "HEAD / HTTP/1.0\r\n\r\n"
	fmt.Println("client send：",requestbody)
	_,err = conn.Write([]byte(requestbody)) //发送请求
	checkTcpClientErr(err)
	result,err := ioutil.ReadAll(conn) //获取响应
	checkTcpClientErr(err)
	fmt.Println("server response：",string(result))

	fmt.Println()

	requestbody = "timestamp"
	fmt.Println("client send：",requestbody)
	_,err = conn.Write([]byte(requestbody)) //发送请求
	checkTcpClientErr(err)
	result,err = ioutil.ReadAll(conn) //获取响应
	checkTcpClientErr(err)
	fmt.Println("server response：",string(result))
}


func checkTcpClientErr(e error){
	if e != nil{
		fmt.Println(e)
	}
}