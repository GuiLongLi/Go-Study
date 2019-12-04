package main

import (
	"log"
	"fmt"
	"strings"
	"time"
	"regexp"

	"github.com/gocolly/colly"
	"github.com/antchfx/htmlquery"
)

func main() {
	//使用dom方式抓取
	testcollydom()
	log.Println("2秒后使用xpath抓取...........")
	//停止2秒后，使用xpath抓取
	time.Sleep(2*time.Second)
	fmt.Println()
	testcollyxpath()
}

/*
Collector对象接受多种回调方法，有不同的作用，按调用顺序我列出来：

OnRequest。请求前
OnError。请求过程中发生错误
OnResponse。收到响应后
OnHTML。如果收到的响应内容是HTML调用它。
OnXML。如果收到的响应内容是XML 调用它。写爬虫基本用不到，所以上面我没有使用它。
OnScraped。在OnXML/OnHTML回调完成后调用。不过官网写的是Called after OnXML callbacks，实际上对于OnHTML也有效，大家可以注意一下。
*/

func testcollydom(){
	//创建新的采集器
	c := colly.NewCollector(
			//这次在colly.NewCollector里面加了一项colly.Async(true)，表示抓取时异步的
			colly.Async(true),
			//模拟浏览器
			colly.UserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"),
		)
	//限制采集规则
	//在Colly里面非常方便控制并发度，只抓取符合某个(些)规则的URLS，有一句c.Limit(&colly.LimitRule{DomainGlob: "*.douban.*", Parallelism: 5})，表示限制只抓取域名是douban(域名后缀和二级域名不限制)的地址，当然还支持正则匹配某些符合的 URLS，具体的可以看官方文档。
	c.Limit(&colly.LimitRule{DomainGlob: "*.douban.*",Parallelism:5})
	/*
另外Limit方法中也限制了并发是5。为什么要控制并发度呢？因为抓取的瓶颈往往来自对方网站的抓取频率的限制，如果在一段时间内达到某个抓取频率很容易被封，所以我们要控制抓取的频率。另外为了不给对方网站带来额外的压力和资源消耗，也应该控制你的抓取机制。
	*/

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})

	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	c.OnHTML(".hd", func(e *colly.HTMLElement) {
		log.Println(strings.Split(e.ChildAttr("a", "href"), "/")[4],
			strings.TrimSpace(e.DOM.Find("span.title").Eq(0).Text()))
	})

	c.OnHTML(".paginator a", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	c.Visit("https://movie.douban.com/top250?start=0&filter=")
	c.Wait()

}


func testcollyxpath(){
	//创建新的采集器
	c := colly.NewCollector(
		//这次在colly.NewCollector里面加了一项colly.Async(true)，表示抓取时异步的
		colly.Async(true),
		//模拟浏览器
		colly.UserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"),
		//最大深度2
		colly.MaxDepth(2),
	)
	//限制采集规格
	c.Limit(&colly.LimitRule{DomainGlob: "*.douban.*",Parallelism:5})
	//请求前
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	//出现错误
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong:", err)
	})
    //收到响应后
	c.OnResponse(func(r *colly.Response) {
		doc, err := htmlquery.Parse(strings.NewReader(string(r.Body)))
		if err != nil {
			log.Fatal(err)
		}
		nodes := htmlquery.Find(doc, `//ol[@class="grid_view"]/li//div[@class="hd"]`)
		for _, node := range nodes {
			url := htmlquery.FindOne(node, "./a/@href")
			title := htmlquery.FindOne(node, `.//span[@class="title"]/text()`)
			log.Println(strings.Split(htmlquery.InnerText(url), "/")[4],
				htmlquery.InnerText(title))
		}
	})
	//因为最大深度设置2，
	//当前第一级 html里的 每个a标签都会回调访问
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")

		// 查找行首以 ?start=0&filter= 的字符串（非贪婪模式）
		reg := regexp.MustCompile(`(?U)^\?start=(\d+)&filter=`)
		regMatch := reg.FindAllString(link, -1)
		//如果找的到的话
		if(len(regMatch) > 0){

			link = "https://movie.douban.com/top250"+regMatch[0]
			//访问该链接
			e.Request.Visit(link)
		}

		// Visit link found on page
	})

	//结束
	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})
	//采集开始
	c.Visit("https://movie.douban.com/top250?start=0&filter=")
	c.Wait()

}
