package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"httpredisqueue/service"
	"httpredisqueue/router/middleware"
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
	g.GET("/setRedis",service.SetRedis) //设置redis
	g.GET("/getRedis",service.GetRedis) //获取redis
	g.GET("/inqueue",service.Inqueue) //获取队列
	g.GET("/outqueue",service.Outqueue) //出队列
}