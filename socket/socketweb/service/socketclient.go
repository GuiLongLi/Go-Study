package service

import (
	"log"
	"net"
	"strings"
)


func cConnHandler(c net.Conn,message string,timeout int) string {
	//缓存 conn 中的数据
	buf := make([]byte, 1024)
	response := ""

	log.Println("请输入客户端请求数据...")

	for {
		//发送消息
		response = messageHandle(c,message,buf)
		if response != "" {
			c.Close() //关闭连接
			return response
		}
	}
	timeout--
	return SocketConnect(message,timeout)
}

//消息处理
func messageHandle(c net.Conn,message string,buf []byte) string{
	//去除输入两端空格
	message = strings.TrimSpace(message)
	//客户端请求数据写入 conn，并传输
	c.Write([]byte(message))
	//服务器端返回的数据写入空buf
	cnt, err := c.Read(buf)

	if err != nil {
		log.Printf("客户端读取数据失败 %s\n", err)
		return "" //退出当前循环
	}

	//回显服务器端回传的信息
	log.Print("服务器端回复" + string(buf[0:cnt]))
	return string(buf[0:cnt]);
}

//timeout 超时次数
func SocketConnect(message string,timeout int) string {
	if(timeout == 0){ //0是没有次数
		return ""
	}
	conn, err := net.Dial("tcp", "47.75.74.233:6661")
	if err != nil {
		log.Println("客户端建立连接失败")
		return ""
	}

	return cConnHandler(conn,message,timeout)
}