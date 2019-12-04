package main

import (
	"log"
	"time"
	"strings"
	"strconv"
	"net/http"

	//htmlquery 包
	"golang.org/x/net/html"
	"github.com/antchfx/htmlquery"
)

//测试爬虫
func main() {
	htmlquerypath()

}

//使用htmlquery 包
func htmlquerypath(){
	start := time.Now()
	ch := make(chan bool)
	for i := 0; i < 10; i++ {
		go htmlparseUrls("https://movie.douban.com/top250?start="+strconv.Itoa(25*i), ch)
	}

	for i := 0; i < 10; i++ {
		<-ch
	}

	elapsed := time.Since(start)
	log.Printf("Took %s", elapsed)
}

func htmlfetch(url string) *html.Node {
	log.Println("Fetch Url", url)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Http get err:", err)
	}
	if resp.StatusCode != 200 {
		log.Fatal("Http status code:", resp.StatusCode)
	}
	defer resp.Body.Close()
	doc, err := htmlquery.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return doc
}

func htmlparseUrls(url string, ch chan bool) {
	doc := htmlfetch(url)
	pic := htmlquery.Find(doc, `//ol[@class="grid_view"]/li//div[@class="pic"]`)
	nodes := htmlquery.Find(doc, `//ol[@class="grid_view"]/li//div[@class="hd"]`)
	for key, node := range nodes {
		num := htmlquery.FindOne(pic[key], `./em[@class=""]/text()`)
		url := htmlquery.FindOne(node, "./a/@href")
		title := htmlquery.FindOne(node, `.//span[@class="title"]/text()`)
		log.Println(htmlquery.InnerText(num),
			strings.Split(htmlquery.InnerText(url), "/")[4],
			htmlquery.InnerText(title))
	}
	time.Sleep(2 * time.Second)
	ch <- true
}