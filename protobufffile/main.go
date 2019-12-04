package main

import (
	"fmt"
	"os"
	"io"

	"github.com/golang/protobuf/proto"

	stProto "protobufffile/protobuf"
)

func main() {
	//写入数据
	protobuffWrite()
	//读取数据
	protobuffRead()
}

func protobuffWrite(){

	//初始化protobuf数据格式
	msg := &stProto.HelloWorld{
		Id:     *proto.Int32(17),
		Name:   *proto.String("hello world"),
		Sin:    *proto.Int32(18),

	}

	filename := "./protobufffile.txt"
	fmt.Printf("使用protobuf创建文件 %s\n",filename)
	fObj,_ := os.Create(filename)  //创建文件
	defer fObj.Close() //关闭文件 ,defer 会在程序最后运行
	buffer,_ := proto.Marshal(msg)  //序列化数据
	fObj.Write(buffer) //写入文件
}

func protobuffRead(){
	filename := "protobufffile.txt"
	file,fileErr := os.Open(filename)  //打开文件
	checkError(fileErr)

	defer file.Close()//关闭文件 ,defer 会在程序最后运行
	fs,fsErr := file.Stat()
	checkError(fsErr)
	buffer := make([]byte,fs.Size()) //创建 byte切片
	//把file文件内容读取到buffer
	_,readErr := io.ReadFull(file,buffer)
	checkError(readErr)

	//初始化pb结构体对象并将buffer中的文件内容读取到pb结构体中
	msg := &stProto.HelloWorld{}
	pbErr := proto.Unmarshal(buffer, msg) //反序列化数据
	checkError(pbErr)
	fmt.Printf("读取文件:%s \r\nname:%s\nid:%d\nopt:%d\n",filename,msg.GetName(),msg.GetId(),msg.GetSin())
}

//检查错误
func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
}