package main

import (
	"google.golang.org/grpc"

	"fmt"
	"context"
	"net"

	pb "grpc/protobuf"
)

//服务器端口
const port = ":6664"

//定义struct 来实现我们自定义的 helloworld.proto 对应的服务
type myServer struct {

}

func (m *myServer) SayHello(ctx context.Context,in *pb.HelloRequest) (*pb.HelloReply,error){
	return &pb.HelloReply{Message:"请求server端成功!"}, nil
}

/*
首先我们必须实现我们自定义的 rpc 服务，例如： rpc SayHello() -在此我们可以实现自己的逻辑
创建监听listener
创建grpc 服务
将我们的服务注册到grpc 的server 中
启动grpc 服务，将我们自定义的监听信息传递给grpc 服务器
*/

func main(){
	//创建server 端监听端口
	list,err := net.Listen("tcp",port)
	if err != nil{
		fmt.Println(err)
	}

	//创建grpc 的server
	server := grpc.NewServer()
	//注册我们自定义的helloworld 服务
	pb.RegisterGreeterServer(server,&myServer{})

	//启动grpc 服务
	fmt.Println("grpc 服务启动...")
	server.Serve(list)

}