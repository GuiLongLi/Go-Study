package handler

import (
	"net/http"
	"webproject/pkg/errno"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

//返回json 格式
func SendResponse(c *gin.Context,err error,data interface{}){
	code,message := errno.DecodeErr(err)

	//总是返回http状态ok
	c.JSON(http.StatusOK,Response{
		Code: code,
		Message:message,
		Data: data,
	})

}

//返回html 格式
func SendResponseHtml(c *gin.Context,err error,data string){
	c.Header("Content-Type", "text/html; charset=utf-8")
	//总是返回http状态ok
	c.String(http.StatusOK,data)
}