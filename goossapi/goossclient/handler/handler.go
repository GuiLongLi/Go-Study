package handler

import (
	"errors"
	"bytes"
	"net/http"
	"io/ioutil"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

//http请求
func HttpRequest(api string,contenttype string,buff *bytes.Buffer,method string) (string, error) {
	//jsonStr := []byte(json)
	//req, err := http.NewRequest(method, api, bytes.NewBuffer(jsonStr))
	req, err := http.NewRequest(method, api, buff)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", contenttype) //使用json格式传参

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if !(resp.StatusCode == 200) {
		return "",  errors.New("api服务器出错")
	}
	return string(body), nil
}