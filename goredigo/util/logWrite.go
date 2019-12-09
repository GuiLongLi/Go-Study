package util

import (
	"os"
	"strings"
)

func WriteLog(msg string,logPath string){
	//创建并已读写形式打开文件并追加插入数据，文件权限644
	fd,_ := os.OpenFile(logPath,os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
	//函数的最后，把文件资源关闭
	defer fd.Close()
	//数组合并成字符串
	content := strings.Join([]string{msg,"\r\n"},"")
	//字符串转换byte类型
	buf := []byte(content)
	//把buf 写入到文件中
	fd.Write(buf)
}
