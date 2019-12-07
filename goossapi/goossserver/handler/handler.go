package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

//返回json 格式
/*
param c *gin.Context		 gin上下文 必传
param code					 int 代号 必传
param message string		 提示语 必传
param data ...interface{} 	 返回数据 选填
*/
func SendResponse(c *gin.Context,code int,message string,data ...interface{}){
	if data == nil{
		data = make([]interface{},0)
	}
	//总是返回http状态ok
	c.JSON(http.StatusOK,Response{
		Code: code,
		Message:message,
		Data: data,
	})
}

//返回成功
func SendSuccess(c *gin.Context,message string,data ...interface{}){
	SendResponse(c,1,message,data)
}

//返回提示
func SendTips(c *gin.Context,message string,data ...interface{}){
	SendResponse(c,0,message,data)
}

//返回失败
func SendFailure(c *gin.Context,message string,data ...interface{}){
	SendResponse(c,-1,message,data)
}