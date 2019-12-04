package main

import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context)  {
	fmt.Fprintf(c.Writer,"Welcome index!")
}

func main() {
	router := gin.Default()
	router.GET("/",Login)

	router.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})

	router.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})

	router.POST("/your/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})

	router.Run("10.10.87.243:8080")
}