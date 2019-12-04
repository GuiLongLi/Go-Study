package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
	"github.com/spf13/viper"

	"aichat/service"
	"aichat/router/middleware"
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

	g.Use(TlsHandler())

	//404处理
	g.NoRoute(func(c *gin.Context){
		c.String(http.StatusNotFound,"该路径不存在")
	})
	//健康检查中间件
	g.GET("/",service.Index)//主页
	g.GET("/chat",service.AiChat)//
}

//监听路由，自动跳转https
func TlsHandler() gin.HandlerFunc {
	sslhost := "go.daily886.com:"+viper.GetString("common.server.addr")
	return func(c *gin.Context) {
		secureMiddleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     sslhost,
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		// If there was an error, do not continue.
		if err != nil {
			return
		}

		c.Next()
	}
}