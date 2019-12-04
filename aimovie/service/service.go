package service

import (
	"sync"

	"github.com/gin-gonic/gin"

	. "aimovie/handler"
	"aimovie/model"
	"aimovie/pkg/errno"
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



//获取电影下载链接
func SearchMovie(c *gin.Context){
	//获取聊天信息
	movie := c.Query("movie")
	if movie == ""{
		SendResponse(c,errno.VALUEERROR,nil)
		return
	}
	results:= make(chan string)
	go DownloadMovie(results,movie)

	html := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>电影资源</title>
</head>
<body>
`
	for {
		msg, ok := <-results //retrive result from channel
		if !ok {
			return
		}
		html = html+msg

		SendResponseHtml(c,nil,html)
	}

}

func DownloadMovie(results chan<- string,movie string) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		results <- getResourceFromLbldy(movie) //使用龙部落电影资源
	}()
	wg.Wait()
	close(results)
}

func getResourceFromLbldy(movie string) (string){
	//获取电影id
	movieId,_ := model.SearchLbldy(movie)
	//获取下载链接
	movieLink,_ := model.DownloadLbldy(movieId)
	return movieLink
}