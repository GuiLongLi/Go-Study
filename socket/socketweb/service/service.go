package service

import (
	"log"
	"strings"
	"github.com/gin-gonic/gin"

	. "socketweb/handler"
	"socketweb/pkg/errno"
)

/*
1、使用接收单个参数各种方法：

c.Param()
c.Query
c.DefaultQuery
c.PostForm
c.DefaultPostForm
c.QueryMap
c.PostFormMap
c.FormFile
c.MultipartForm

2、使用各种绑定方法

c.Bind
c.BindJSON
c.BindXML
c.BindQuery
c.BindYAML
c.ShouldBind
c.ShouldBindJSON
c.ShouldBindXML
c.ShouldBindQuery
c.ShouldBindYAML
*/

//首页
func Index(c *gin.Context){
	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>hello world</title>
</head>
<body>
    hello world
</body>
</html>
`
	SendResponseString(c,nil,html)
}

//消息html
func SocketWeb(c *gin.Context){
	param := gin.H{
		"title": "web socket",
	}
	SendResponseHtml(c,"socketweb.html",&param)
}

// 发送消息
// Binding from JSON
type FormMessage struct {
	Message string `form:"message" json:"message" binding:"required"`
	Timestamp string `form:"time" json:"time" binding:"required"`
}
type ResponseMessage struct {
	Message string `json:"message"`
}
func SendMessage(c *gin.Context){
	var form FormMessage
	// 你可以显式声明来绑定多媒体表单：
	// c.BindWith(&form, binding.Form)
	// 或者使用自动推断:
	var result interface{}
	result = ResponseMessage{
		Message: "参数错误",
	}
	var err error
	//获取 content type 类型
	contentType := c.Request.Header.Get("Content-Type")
	semicol := strings.Index(contentType, ";")
	if semicol <= 0 {
		semicol = len(contentType)
	}
	contentType = contentType[0:semicol]
	//绑定参数
	switch contentType {
	case "application/json":
		err = c.ShouldBindJSON(&form)
	case "application/x-www-form-urlencoded":
		err = c.ShouldBind(&form)
	}
	log.Println(form)
	if err != nil || form.Message == "" || form.Timestamp == "" {
		result = ResponseMessage{
			Message: "参数不能为空",
		}
		SendResponseJSON(c,errno.VALUEERROR,result)
		return
	}

	response := SocketConnect(form.Message,2)
	//替换无效字符
	response = strings.Replace(response, "服务器端回复", "", -1)
	response = strings.Replace(response, "->", "", -1)
	//去除输入两端空格
	response = strings.TrimSpace(response)

	if response == ""{
		result = ResponseMessage{
			Message: "无法连接服务器",
		}
		SendResponseJSON(c,errno.VALUEERROR,result)
	}
	log.Println(response)
	result = ResponseMessage{
		Message: response,
	}
	SendResponseJSON(c,nil,result)
}