package main

import (
	"fmt"
	"strings"
	"strconv"
)


func main() {
	fmt.Println("teststring")
	teststring()

	fmt.Println()
	fmt.Println("stringexchange")
	stringexchange()
}

func teststring(){
	var result interface{}
	/*
	func Contains(s, substr string) bool
	字符串s中是否包含substr，返回bool值
	*/
	result = strings.Contains("hello world","hello")
	fmt.Println("判断字符串是否存在",result)
	result = strings.Contains("hello world","hellow")
	fmt.Println("判断字符串是否存在",result)

	/*
	func Join(a []string, sep string) string
	字符串链接，把slice a通过sep链接起来
	*/
	arr := []string{
		"hello",
		"world",
	}
	result = strings.Join(arr,",")
	fmt.Println("字符串链接",result)

	/*
	func Index(s, sep string) int
	在字符串s中查找sep所在的位置，返回位置值，找不到返回-1
	*/
	result = strings.Index("hello world","world")
	fmt.Println("字符串位置",result)
	result = strings.Index("hello world","world2")
	fmt.Println("字符串位置",result)

	/*
	func Repeat(s string, count int) string
	重复s字符串count次，最后返回重复的字符串
	*/
	result = strings.Repeat("hello ",2)
	fmt.Println("重复的字符串",result)

	/*
	func Replace(s, old, new string, n int) string
	在s字符串中，把old字符串替换为new字符串，n表示替换的次数，小于0表示全部替换
	*/
	result = strings.Replace("hello world","hello ","hello",-1)
	fmt.Println("字符串替换",result)

	/*
	func Split(s, sep string) []string
	把s字符串按照sep分割，返回slice
	*/
	result = strings.Split("hello world"," ")
	fmt.Println("字符串切割",result)

	/*
	func Trim(s string, cutset string) string
	在s字符串的头部和尾部去除cutset指定的字符串
	*/
	result = strings.Trim("hello world !!!"," !!!")
	fmt.Println("字符串Trim",result)

	/*
	func Fields(s string) []string
	去除s字符串的空格符，并且按照空格分割返回slice
	*/
	result = strings.Fields("a b  c   d    e")
	fmt.Println("字符串Fields切割",result)
}

func stringexchange(){
	/*
	字符串转换
	字符串转化的函数在strconv中，如下也只是列出一些常用的：
	*/
	/*
	Append 系列函数将整数等转换为字符串后，添加到现有的字节数组中。
	*/
	result := make([]byte,0,100)
	result = strconv.AppendInt(result,4567,10)
	result = strconv.AppendBool(result,false)
	result = strconv.AppendQuote(result,"abcdefg")
	result = strconv.AppendQuoteRune(result,'单')
	fmt.Println("字符串Append",string(result))

	/*
	Format 系列函数把其他类型的转换为字符串
	*/
	a := strconv.FormatBool(false)
	b := strconv.FormatFloat(123.33,'g',12,64)
	c := strconv.FormatInt(1234,10)
	d := strconv.FormatUint(12345,10)
	e := strconv.Itoa(1023)
	fmt.Println("字符串Format",a,b,c,d,e)

	/*
	Parse 系列函数把字符串转换为其他类型
	*/
	p1,err := strconv.ParseBool("false")
	checkError(err)
	p2,err := strconv.ParseFloat("123.23",64)
	checkError(err)
	p3,err := strconv.ParseInt("1234",10,64)
	checkError(err)
	p4,err := strconv.ParseUint("12345",10,64)
	checkError(err)
	p5,err := strconv.Atoi("1023")
	checkError(err)
	fmt.Println("字符串Parse",p1,p2,p3,p4,p5)
}

func checkError(e error){
	if e != nil{
		fmt.Println(e)
	}
}