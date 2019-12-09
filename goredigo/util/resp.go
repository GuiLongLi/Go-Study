package util

import (
	"encoding/json"
	"log"
	"net/http"
)

//返回结构体
type ResponseData struct {
	Code int `json:"code"`
	Msg string `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

//返回json格式的数据给客户端
func RespJson(writer http.ResponseWriter,code int,msg string,data interface{}){
	Resp(writer,code,msg,data)
}
func Resp(writer http.ResponseWriter,code int,msg string,data interface{}){
	//设置header 为JSON ，默认是text/html ,所以特别指出返回数据类型是 application/json
	writer.Header().Set("Content-Type","application/json")
	writer.WriteHeader(http.StatusOK)
	rep := ResponseData{
		Code:code,
		Msg:msg,
		Data:data,
	}
	//将结构体转换为json字符串
	ret,err := json.Marshal(rep)
	if err != nil{
		log.Panicln(err.Error())
	}

	//返回json ok
	writer.Write(ret)
}
