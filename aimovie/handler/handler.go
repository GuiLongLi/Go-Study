package handler

import (
	"bytes"
	"net/http"
	"io/ioutil"

	"github.com/gin-gonic/gin"

	"aimovie/pkg/errno"
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

//http请求 post
func HttpPost(api string,json string) (string, error) {
	jsonStr := []byte(json)
	req, err := http.NewRequest("POST", api, bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", errno.ApiServerError
	}
	req.Header.Set("Content-Type", "application/json") //使用json格式传参

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", errno.ApiServerError
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if !(resp.StatusCode == 200) {
		return "",  errno.ApiServerError
	}
	return string(body), nil
}

//http请求 get
func HttpGet(api string) (string,error){
	resp, err := http.Get(api)
	if err != nil {
		return "", errno.ApiServerError
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errno.ApiServerError
	}
	if !(resp.StatusCode == 200) {
		return "",  errno.ApiServerError
	}
	return string(body), nil
}

//http请求 get 没有解析body
func HttpGetBody(api string) (*http.Response,error){
	resp, err := http.Get(api)
	return resp,err
}