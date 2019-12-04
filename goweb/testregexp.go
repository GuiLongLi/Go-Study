package main

import (
	"fmt"
	"regexp"
	"net/http"
	"io/ioutil"
	"strings"
)

/*
通过正则判断是否匹配
regexp包中含有三个函数用来判断是否匹配，如果匹配返回true，否则返回false

func Match(pattern string, b []byte) (matched bool, error error)
func MatchReader(pattern string, r io.RuneReader) (matched bool, error error)
func MatchString(pattern string, s string) (matched bool, error error)

面的三个函数实现了同一个功能，就是判断pattern是否和输入源匹配，匹配的话就返回true，如果解析正则出错则返回error。三个函数的输入源分别是byte slice、RuneReader和string。
*/


func main() {
	ip := "127.0.0.1"
	res := isip(ip)
	if !res{
		fmt.Println(ip," 不是ip")
	}else{
		fmt.Println(ip," 是ip")
	}

	ip = "127.0.0.1.0"
	res = isip(ip)
	if !res{
		fmt.Println(ip," 不是ip")
	}else{
		fmt.Println(ip," 是ip")
	}

	fmt.Println()

	vars := "1134"
	res = isnumber(vars)
	if !res {
		fmt.Println(vars," 不是数字")
	}else{
		fmt.Println(vars," 是数字")
	}

	vars = "13.2"
	res = isnumber(vars)
	if !res {
		fmt.Println(vars," 不是数字")
	}else{
		fmt.Println(vars," 是数字")
	}

	vars = ".234"
	res = isnumber(vars)
	if !res {
		fmt.Println(vars," 不是数字")
	}else{
		fmt.Println(vars," 是数字")
	}

	vars = "123abv"
	res = isnumber(vars)
	if !res {
		fmt.Println(vars," 不是数字")
	}else{
		fmt.Println(vars," 是数字")
	}

	//正则替换
	fmt.Println()
	fmt.Println("spiders")
	spiders()

	//正则查找
	fmt.Println()
	fmt.Println("regexpfind")
	regexpfind()

	//expend
	fmt.Println()
	fmt.Println("testexpend")
	testexpend()
}


//验证是不是ip地址
func isip(ip string) bool {
	regexp1 := `^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`
	if m,_ := regexp.MatchString(regexp1,ip);!m{
		return false
	}
	return true
}

//验证数字
func isnumber(num string) bool{
	regexp1 := `^[0-9]+$`
	//判断是否整数
	if m,_ := regexp.MatchString(regexp1,num);!m{
		//判断是否浮点数
		regexp2 := `^[0-9]*\.[0-9]+$`
		if m,_ := regexp.MatchString(regexp2,num);!m {
			return false
		}else{
			return true
		}
	}
	return true
}

//正则替换字符串
func spiders(){
	resp,err := http.Get("http://www.baidu.com")
	if err != nil{
		fmt.Println("http get error.")
	}
	defer resp.Body.Close()
	body,err := ioutil.ReadAll(resp.Body)
	if err != nil{
		fmt.Println("http read error")
		return
	}

	src := string(body)
	//fmt.Println();
	//fmt.Println("src1",src);

	//将html标签全部替换成小写
	re,_ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src,strings.ToLower)

	//去除style
	re,_ = regexp.Compile(`<style([\s\S]*?)</style>`)
	src = re.ReplaceAllString(src,"")

	//去除script
	re,_ = regexp.Compile(`<script([\s\S]*?)</script>`)
	src = re.ReplaceAllString(src,"")

	//去除所有尖括号内的html代码，并换成换行符
	re,_ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src,"\n")

	//去除连续的换行符
	re,_ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src,"\n")

	fmt.Println(strings.TrimSpace(src))
}

func regexpfind(){
	a := "i am learning go language"

	//正则匹配 a-z 字母 2-4 个
	re,_ := regexp.Compile("[a-z]{2,4}")

	//查找符合正则的第一个
	one := re.Find([]byte(a))
	fmt.Println("find:",string(one))

	//查找符合正则的所有slice ，n小于0表示返回全部符合的字符串,不然就是返回指定的长度
	//regexp.Compile(regexp).FindAll([]byte(string),n)
	all := re.FindAll([]byte(a),-1)
	fmt.Println("FindAll",all)

	//查找符合条件的index位置，开始位置和结束位置
	//regexp.Compile(regexp).FindIndex([]byte(string))
	index := re.FindIndex([]byte(a))
	fmt.Println("FindIndex",index)

	//查找符合条件的所有index 的位置, n同上
	//regexp.Compile(regexp).FindAllIndex([]byte(string),n)
	allindex := re.FindAllIndex([]byte(a),-1)
	fmt.Println("FindAllIndex",allindex)

	//正则匹配 am所有lang所有
	re2,_ := regexp.Compile("am(.*)lang(.*)")

	//查找Submatch, 返回数组,第一个元素是匹配的全部元素，第二个元素是第一个()里面的，第三个是第二个()里面的
	//下面的输出第一个元素是"am learning Go language"
	//第二个元素是"learning Go",注意保护空格的输出
	//第三个元素是"uage"
	submatch := re2.FindSubmatch([]byte(a))
	fmt.Println("FindSubmatch",submatch)
	for _,v := range submatch{
		fmt.Println(string(v))
	}

	//定义和上面的FindIndex一样
	submatchindex := re2.FindSubmatchIndex([]byte(a))
	fmt.Println(submatchindex)

	//FindAllSubmatch,查找所有符合条件的子匹配
	submatchall := re2.FindAllSubmatch([]byte(a),-1)
	fmt.Println(submatchall)

	//FindAllSubmatchIndex,查找所有字匹配的index
	submatchallindex := re2.FindAllSubmatchIndex([]byte(a),-1)
	fmt.Println(submatchallindex)
}

func testexpend(){
	src := []byte(`
		call hello alice
		hello bob
		call hello eve
	`)
	pat := regexp.MustCompile(`(?m)(call)\s+(?P<cmd>\w+)\s+(?P<arg>.+)\s*$`)
	res := []byte{}
	for _, s := range pat.FindAllSubmatchIndex(src, -1) {
		res = pat.Expand(res, []byte("$cmd('$arg')\n"), src, s)
	}
	fmt.Println(string(res))
}