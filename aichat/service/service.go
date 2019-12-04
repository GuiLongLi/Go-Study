package service

import (
	"github.com/tidwall/gjson"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"strconv"
	"log"

	"aichat/pkg/errno"
	. "aichat/handler"
	"aichat/model"
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

//获取tuling接口回复
func TulingAi(info string) (string,error) {
	api := viper.GetString("common.tuling.api")

	//发送http请求图灵api  , body是http响应
	var body, resultErrs = HttpRequest(api,info,"POST")
	if resultErrs != nil {
		return "", errno.ApiServerError
	}

	return body, nil
}

//回复信息构造体
type tlReply struct {
	code int
	Text string `json:"text"`
}

//聊天函数
func AiChat(c *gin.Context){
	//获取聊天信息
	message := c.Query("message")
	if message == ""{
		SendResponse(c,errno.VALUEERROR,nil)
		return
	}
	var userId = "1"
	//图灵接口参数构造体
	var chattingInfo = model.BuildChatting(message,userId, viper.GetString("common.tuling.apikey"))
	log.Printf("chattingInfo: %+v\n",chattingInfo)
	// 参数构造体 转换成 字符串
	chatstr,err := model.ConvertJson(chattingInfo)
	if err != nil{
		SendResponse(c,errno.InternalServerError,nil)
		return
	}

	//调用图灵接口
	body,err := TulingAi(chatstr)
	if err != nil{
		SendResponse(c,errno.InternalServerError,nil)
		return
	}
	log.Printf("body: %+v\n",body)
	var results string
	// 使用gjson 获取返回结果的 resultType
	result := gjson.Get(body, "results.#.resultType")
	for key, name := range result.Array() {
		//如果 resultType 是 text格式
		if name.String() == "text"{
			//获取对应 key 的 values里的text ，就是图灵回复的文字
			getstring := "results."+strconv.Itoa(key)+".values.text"
			log.Printf("getstring: %+v\n",getstring)
			result_text := gjson.Get(body,getstring)
			results = result_text.String()
		}
	}

	SendResponse(c,nil,results)
}