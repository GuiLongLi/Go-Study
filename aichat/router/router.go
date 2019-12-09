package router

import (
	"log"
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

	rou := gin.Default()
	rou.Use(TlsHandler())

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
	log.Println("sslhost:",sslhost)
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

		////判断get访问 ,
		//这里的只能手动访问https， 如果要跳转，就要用http监听多另外一个端口，监听后跳到https就可以了
		//log.Println("c.Request.Method",c.Request.Method)
		//if c.Request.Method == "GET"{
		//	/*
		//	fmt.Println(c.Request.Proto)
		//	// output:HTTP/1.1
		//	fmt.Println(c.Request.TLS)
		//	// output: <nil>
		//	fmt.Println(c.Request.Host)
		//	// output: localhost:9090
		//	fmt.Println(c.Request.RequestURI)
		//	// output: /index?id=1
		//	*/
		//
		//	//判断是否https
		//	log.Println("c.Request.TLS",c.Request.TLS)
		//	if c.Request.TLS == nil { //不是https
		//		newurls := strings.Join([]string{"https://", c.Request.Host, c.Request.RequestURI}, "")
		//		log.Println("newurls",newurls)
		//		c.Redirect(http.StatusMovedPermanently,newurls) //重定向到https
		//	}
		//}

		c.Next()
	}
}