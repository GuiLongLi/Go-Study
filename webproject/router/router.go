package router

import (
	"net/http"
	"webproject/service"
	"webproject/router/middleware"
	"github.com/gin-gonic/gin"
)

//初始化路由
func InitRouter(g *gin.Engine){
	middlewares := []gin.HandlerFunc{}
	//中间件
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(middlewares...)

	//404处理
	g.NoRoute(func(c *gin.Context){
		c.String(http.StatusNotFound,"该路径不存在")
	})
	//健康检查中间件
	g.GET("/",service.Index)//主页
	g.GET("/install",service.Install)//创建数据库
	usergroup := g.Group("/user") //创建路由组
	{
		usergroup.POST("/addUser",service.AddUser)  //添加用户
		usergroup.GET("/selectUser",service.SelectUser)//查询用户
	}
}