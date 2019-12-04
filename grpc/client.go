package main

import (
	"google.golang.org/grpc"

	"fmt"
	"context"

	pb "grpc/protobuf"

)

//与服务器对应的端口
const address = "127.0.0.1:6664"

/*
创建grpc 连接器
创建grpc 客户端，并将连接器赋值给客户端
向grpc 服务器发起请求
获取grpc 服务器返回的结果
*/
func main() {
	//创建一个grpc 连接器
	conn,err := grpc.Dial(address,grpc.WithInsecure())
	if err != nil{
		fmt.Println(err)
	}
	//当请求完毕后记得关闭连接，否则大量连接会占用资源
	defer conn.Close()

	//创建grpc 客户端
	c := pb.NewGreeterClient(conn)

	name := "我是客户端，正在请求服务器！！"
	//客户端向grpc 服务器发起请求
	result,err := c.SayHello(context.Background(),&pb.HelloRequest{Name:name})
	fmt.Println(name)
	if err != nil{
		fmt.Println("请求失败！！")
		return
	}
	//获取服务器返回的结果
	fmt.Println(result.Message)
}