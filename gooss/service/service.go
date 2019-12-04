package service

import (
	"io"
	"log"
	"fmt"
	"time"
	"strconv"
	"net/http"
	"crypto/md5"

	"github.com/gin-gonic/gin"

	. "gooss/handler"
	"gooss/model"
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



//oss信息
func Gooss(c *gin.Context){
	var results string
	results = model.Aliossversion()
	log.Printf("results: %s\n",results)

	var list []string
	list,err := model.GetFilelist()
	if err != nil{
		log.Printf("model.Getlist err: %v\n",err)
	}
	log.Printf("文件列表: %v\n",list)

	SendResponse(c,nil,results)
}

func OssUpload(c *gin.Context){
	log.Printf("method: ",c.Request.Method) //请求方法
	if c.Request.Method == "GET"{
		curtime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h,strconv.FormatInt(curtime,10))
		token := fmt.Sprintf("%x",h.Sum(nil)) //生成token

		//设置cookie
		expiration := 60
		//SetCookie
		/*
		第一个参数为 cookie 名；
		第二个参数为 cookie 值；
		第三个参数为 cookie 有效时长/秒，当 cookie 存在的时间超过设定时间时，cookie 就会失效，它就不再是我们有效的 cookie；
		第四个参数为 cookie 所在的目录；
		第五个为所在域，表示我们的 cookie 作用范围；
		第六个表示是否只能通过 https 访问；
		第七个表示 cookie 是否可以通过 js代码进行操作。
		 */
		c.SetCookie("go-cookie","go-cookies", expiration, "/", "", false, true)

		c.HTML(http.StatusOK, "form.gtpl", gin.H{
			"token":token,
		})
	}else{
		//读取cookie
		cookie, _ := c.Cookie("go-cookie")
		SendResponse(c,nil,cookie)

		//防止多次重复提交表单
		//解决方案是在表单中添加一个带有唯一值的隐藏字段。
		// 在验证表单时，先检查带有该唯一值的表单是否已经递交过了。
		// 如果是，拒绝再次递交；如果不是，则处理表单进行逻辑处理。
		res1 := verifyToken(c)
		if !res1 {
			return
		}

		//上传文件
		header, err := c.FormFile("files")
		if err != nil {
			//ignore
			SendResponse(c,nil,"上传失败")
			return
		}
		localfile := "./uploads/"+header.Filename //本地文件路径
		// gin 简单做了封装,拷贝了文件流
		if err := c.SaveUploadedFile(header, localfile); err != nil {
			// ignore
			SendResponse(c,nil,"本地上传失败")
			return
		}
		SendResponse(c,nil,"本地上传成功")

		//上传到阿里云oss
		yunfiletmp := "uploads/"+header.Filename
		yunfile,err := model.UploadFile(localfile,yunfiletmp)
		if err != nil{
			// ignore
			SendResponse(c,nil,"阿里云上传失败")
			return
		}
		SendResponse(c,nil,"阿里云上传成功")
		SendResponse(c,nil,"阿里云路径"+yunfile)
	}
}

//防止多次重复提交表单
func verifyToken(c *gin.Context) bool{
	token := c.PostForm("token")

	SendResponse(c,nil,token)
	if token != ""{
		// 验证 token 的合法性
		if len(token) <10{
			SendResponse(c,nil,"token验证失败")
			return false
		}
	}else{
		//不存在token 报错
		SendResponse(c,nil,"token验证失败")
		return false
	}
	SendResponse(c,nil,"token验证通过")
	return true
}