package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"time"
	"crypto/md5"
	"strconv"
)

func postFile(filename string, targetUrl string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//加入token
	curtime := time.Now().Unix()
	h := md5.New()
	io.WriteString(h,strconv.FormatInt(curtime,10))
	token := fmt.Sprintf("%x",h.Sum(nil)) //生成token
	//模拟表单数据
	_ = bodyWriter.WriteField("token", token)
	_ = bodyWriter.WriteField("username", "awef")
	_ = bodyWriter.WriteField("password", "password")
	_ = bodyWriter.WriteField("age", "3")
	_ = bodyWriter.WriteField("email", "123@qq.com")
	_ = bodyWriter.WriteField("mobile", "12312345678")
	_ = bodyWriter.WriteField("sex", "1")
	_ = bodyWriter.WriteField("interest", "football")
	_ = bodyWriter.WriteField("usercard", "123456789123456123")

	//关键的一步操作 ，模拟上传文件
	fileWriter, err := bodyWriter.CreateFormFile("userheader", filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return err
	}

	//打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return err
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	resp, err := http.Post(targetUrl, contentType, bodyBuf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)
	fmt.Println(string(resp_body))
	return nil

	//上面的例子详细展示了客户端如何向服务器上传一个文件的例子，客户端通过multipart.Write把文件的文本流写入一个缓存中，然后调用http的Post方法把缓存传到服务器。

	//如果你还有其他普通字段例如username之类的需要同时写入，那么可以调用multipart的WriteField方法写很多其他类似的字段。
}

// sample usage
func main() {
	target_url := "http://127.0.0.1:6665/info"
	filename := "./testformupload.go"
	postFile(filename, target_url)
}