package model

import (
	"io"
	"log"
	"fmt"
	"strings"
	"regexp"

	"github.com/spf13/viper"
	"github.com/PuerkitoBio/goquery"

	"aimovie/pkg/errno"
	. "aimovie/handler"
)

//龙部落电影获取电影资源
type Media struct {
	Name string
	Size string
	Link string
}

//获取电影id
func SearchLbldy(movie string) (string,error){
	api := viper.GetString("common.lbldy.search")
	api = fmt.Sprintf(api+"%s",movie)
	log.Println("api ", api)
	//请求接口，获取电影信息
	result, err := HttpGet(api)
	if err != nil{
		return "", errno.ModelError
	}
	//正则获取 post-id  <div> class="postlist" id="post-64115">
	re, err := regexp.Compile("<div class=\"postlist\" id=\"post-(.*?)\">")
	if err != nil{
		return "", errno.ModelError
	}
	firstId := re.FindSubmatch([]byte(result)) //find first match case
	//log.Println("firstId ", string(firstId[1]))
	if len(firstId) == 0 {
		return "", errno.ModelError
	}
	return string(firstId[1]),nil

}

//获取电影下载链接
func DownloadLbldy(movieId string) (string,error){
	var ms []Media
	api := viper.GetString("common.lbldy.download")
	//请求接口，获取电影信息
	api = fmt.Sprintf(api,movieId)
	//log.Println("api ", api)
	result, err := HttpGetBody(api)
	//log.Println("err :", err)
	if err != nil{
		return "", errno.ModelError
	}
	defer result.Body.Close()
	doc, err := goquery.NewDocumentFromReader(io.Reader(result.Body))
	if err != nil {
		return "", errno.ModelError
	}
	//正则匹配寻找 a标签链接
	doc.Find("p").Each(func(i int, selection *goquery.Selection) {
		name := selection.Find("a").Text()
		link, _ := selection.Find("a").Attr("href")
		if strings.HasPrefix(link, "ed2k") || strings.HasPrefix(link, "magnet") || strings.HasPrefix(link, "thunder") {
			m := Media{
				Name: name,
				Link: link,
			}
			ms = append(ms, m)
		}
	})
	message := ConvertMsg(ms)
	return message,nil
}

//数组转换字符串
func ConvertMsg(ms []Media) string{
	ret := "<h2>龙部落电影资源列表</h2>"
	for i, m := range ms {
		ret += fmt.Sprintf("<p>*%s*</p><p>```<a href=%s target=_blank >%s</a>```</p>", m.Name, m.Link,m.Link)
		//when results are too large, we split it.
		if i%4 == 0 && i < len(ms)-1 && i > 0 {
			ret += fmt.Sprintf("<p>*切割部分 %d*</p>", i/4+1)
		}
	}
	return ret
}