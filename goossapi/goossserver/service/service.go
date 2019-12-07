package service

import (
	"log"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"

	. "goossserver/handler"
	"goossserver/model"
)


func OssUpload(c *gin.Context){
	//读取cookie
	cookie, _ := c.Cookie("go-cookie")
	log.Printf("go-cookie: ",cookie)

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
		SendFailure(c,"上传失败")
		return
	}
	localfile := "./uploads/"+header.Filename //本地文件路径
	// gin 简单做了封装,拷贝了文件流
	if err := c.SaveUploadedFile(header, localfile); err != nil {
		log.Printf("SaveUploadedFile err: ",err)
		// ignore
		SendFailure(c,"本地上传失败")
		return
	}
	log.Printf("本地上传成功")

	//上传到阿里云oss
	dateyear := time.Now().Format("2006") //获取当前年
	datemonth := time.Now().Format("01")//获取当前月
	dateday := time.Now().Format("02")//获取当前日
	yunfiletmp := fmt.Sprintf("uploads/%v/%v/%v/%v",dateyear,datemonth,dateday,header.Filename)
	yunfile,err := model.UploadFile(localfile,yunfiletmp)
	if err != nil{
		log.Printf("UploadFile err: ",err)
		// ignore
		SendFailure(c,"阿里云上传失败")
		return
	}
	log.Printf("阿里云路径: ",yunfile)
	domain := viper.GetString("common.aliyunoss.domain")
	domain = fmt.Sprintf("%s/%s",domain,yunfile)
	oss := model.OssStruct{
		Url: domain,
	}
	log.Printf("阿里云上传成功: %s",oss)
	SendSuccess(c,"阿里云路径",oss)
}

//防止多次重复提交表单
func verifyToken(c *gin.Context) bool{
	token := c.PostForm("token")

	log.Printf("token: %s",token)
	if token != ""{
		// 验证 token 的合法性
		if len(token) <10{
			SendFailure(c,"token验证失败")
			return false
		}
	}else{
		//不存在token 报错
		SendFailure(c,"token验证失败")
		return false
	}
	log.Printf("token验证通过")
	return true
}