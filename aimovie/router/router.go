package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"aimovie/service"
	"aimovie/router/middleware"
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
	g.GET("/movie",service.SearchMovie)//获取电影下载链接
}