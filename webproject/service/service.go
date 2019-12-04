package service

import (
	"github.com/gin-gonic/gin"
	"webproject/model"
	. "webproject/handler"
	"webproject/pkg/errno"
	"fmt"
)

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
	SendResponseHtml(c,nil,html)
}

//新建用户
func AddUser(c *gin.Context){
	var r model.User
	if err := c.Bind(&r); err != nil{
		SendResponse(c,errno.ErrBind,nil)
		return
	}
	u := model.User{
		UserName: r.UserName,
		Password: r.Password,
	}
	//验证数据
	if err := u.Validate(); err != nil{
		SendResponse(c,errno.ErrValidation,nil)
		return
	}
	//插入数据
	if _,err := u.Create();err != nil{
		SendResponse(c,errno.ErrDatabase,nil)
		return
	}
	SendResponse(c,nil,u)
}

//查询用户
func SelectUser(c *gin.Context){
	name := c.Query("user_name")
	if name == ""{
		SendResponse(c,errno.ErrValidation,nil)
		return
	}
	var user model.User
	if err := user.SelectUserByName(name); err != nil{
		fmt.Println(err)
		SendResponse(c,errno.ErrUserNotFound,nil)
		return
	}
	//验证数据
	if err := user.Validate(); err != nil{
		SendResponse(c,errno.ErrUserNotFound,nil)
		return
	}
	SendResponse(c,nil,user)
}


//创建数据库
func Install(c *gin.Context){
	//检查是否安装数据库
	if err := model.CheckInstalled(); err != nil{
		fmt.Println(err)
		SendResponse(c,errno.INSTALLED,err)
		return
	}
	//安装数据库
	if err := model.CreateDatabase(); err != nil{
		fmt.Println(err)
		SendResponse(c,errno.INSTALLERROR,err)
		return
	}
	SendResponse(c,nil,"成功")
}