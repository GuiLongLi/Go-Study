package router

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gooss/service"
	"gooss/router/middleware"
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

	g.LoadHTMLGlob("template/**/*") //加载模板路径

	//健康检查中间件
	g.GET("/",service.Index)//主页
	g.GET("/gooss",service.Gooss)//oss信息
	g.GET("/ossupload",service.OssUpload)//上传oss
	g.POST("/ossupload",service.OssUpload)//上传oss
}
