package model

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OssStruct struct {
	Url string `json:"url"`
}

//sdk 版本
func Aliossversion()(version string){
	return oss.Version
}

//初始化oss服务
func Initserver()(client *oss.Client,err error){
	// Endpoint以杭州为例，其它Region请按实际情况填写。
	endpoint := viper.GetString("common.aliyunoss.endpoint")
	// 阿里云主账号AccessKey拥有所有API的访问权限，风险很高。强烈建议您创建并使用RAM账号进行API访问或日常运维，请登录 https://ram.console.aliyun.com 创建RAM账号。
	accessKeyId := viper.GetString("common.aliyunoss.accessid")
	accessKeySecret := viper.GetString("common.aliyunoss.accesskey")
	// 创建OSSClient实例。
	client, err = oss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil{
		return
	}
	return
}

//获取文件列表
func GetFilelist()(list []string,err error){
	list = make([]string,100)
	client,err := Initserver()
	// 获取存储空间。
	bucketName := viper.GetString("common.aliyunoss.bucket")
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return list,err
	}
	// 列举文件。
	marker := ""
	for {
		lsRes, err := bucket.ListObjects(oss.Marker(marker))
		if err != nil {
			return list,err
		}
		// 打印列举文件，默认情况下一次返回100条记录。
		for _, object := range lsRes.Objects {
			log.Printf("object.Key:%v\n",object.Key)
			list = append(list,object.Key)
		}
		if lsRes.IsTruncated {
			marker = lsRes.NextMarker
		} else {
			break
		}
	}
	return list,err
}

//上传文件
func UploadFile(localfile string,uploadfile string)(resultfile string,err error){
	resultfile = ""
	// 创建OSSClient实例。
	client,err := Initserver()

	bucketName := viper.GetString("common.aliyunoss.bucket")
	// <yourObjectName>上传文件到OSS时需要指定包含文件后缀在内的完整路径，例如abc/efg/123.jpg。
	uploaddir := viper.GetString("common.aliyunoss.uploaddir")
	uploadfile = strings.Trim(uploadfile,"/")
	objectName := fmt.Sprintf("%s/%s",uploaddir,uploadfile) //完整的oss路径
	// <yourLocalFileName>由本地文件路径加文件名包括后缀组成，例如/users/local/myfile.txt。
	localFileName := localfile
	// 获取存储空间。
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return
	}
	// 上传文件。
	err = bucket.PutObjectFromFile(objectName, localFileName)
	if err != nil {
		return
	}
	resultfile = objectName
	return
}