package main

import (
	"fmt"
	"crypto/md5"
	"html/template"
	texttemplate "text/template"
	"io"
	"os"
	"log"
	"time"
	"net/http"
	"strconv"
	"regexp"
	"unicode"
)



//判断是否在切片中
func In_slice(val string, slice []string) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}
//判断两个切片的差异
//返回差异的数组
func Slice_diff(slice1, slice2 []string) (diffslice []string) {
	for _, v := range slice1 {
		if !In_slice(v, slice2) {
			diffslice = append(diffslice, v)
		}
	}
	return
}


func sayhelloName(w http.ResponseWriter,r *http.Request)  {
	//下面这个写入到w 的是输出到客户端的
	fmt.Fprintf(w,"hello testform2.go")
}

//验证表单输入
func info(w http.ResponseWriter,r *http.Request)  {
	fmt.Println("method：",r.Method) //获取请求的方法
	if r.Method == "GET"{
		curtime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h,strconv.FormatInt(curtime,10))
		token := fmt.Sprintf("%x",h.Sum(nil)) //生成token

		//设置cookie
		expiration := time.Now()
		expiration = expiration.AddDate(1, 0, 0)
		cookie := http.Cookie{Name: "username", Value: "go-cookie", Expires: expiration}
		http.SetCookie(w, &cookie)

		t,_ := template.ParseFiles("testform2.gtpl")
		t.Execute(w,token) //把token 写入到模板中
	}else{
		//读取cookie
		cookie, _ := r.Cookie("username")
		fmt.Println("cookie：",cookie)

		//防止多次重复提交表单
		//解决方案是在表单中添加一个带有唯一值的隐藏字段。
		// 在验证表单时，先检查带有该唯一值的表单是否已经递交过了。
		// 如果是，拒绝再次递交；如果不是，则处理表单进行逻辑处理。
		res1 := verifyToken(w,r)
		if !res1 {
			return
		}

		//表单验证
		res2 := formValidate(w,r)
		if !res2 {
			return
		}

		//预防跨站脚本
		res3 := xss(w,r)
		if !res3 {
			return
		}

		//获取上传文件
		res4 := uploadHandler(w,r)
		if !res4{
			return
		}

	}

}
//防止多次重复提交表单
func verifyToken(w http.ResponseWriter,r *http.Request) bool{
	//r.ParseForm() //普通表单

	r.ParseMultipartForm(32<<20) //加入了 enctype="multipart/form-data" 的表单
	token := r.Form.Get("token")
	fmt.Println("token:",token)
	if token != ""{
		// 验证 token 的合法性
		if len(token) <10{
			fmt.Fprintf(w,"token验证失败")
			return false
		}
		fmt.Fprintf(w,"token验证通过")
	}else{
		//不存在token 报错
		fmt.Fprintf(w,"token验证失败")
		return false
	}
	fmt.Fprintf(w,"\n")
	return true
}

//表单验证
func formValidate(w http.ResponseWriter,r *http.Request) bool{
	//r.ParseForm()       //解析url传递的参数，对于POST则解析响应包的主体（request body）

	r.ParseMultipartForm(32<<20) //加入了 enctype="multipart/form-data" 的表单

	//请求的是登录数据，name执行登录的逻辑判断
	fmt.Println("username：",r.Form["username"])
	fmt.Println("password：",r.Form["password"])
	fmt.Println("age：",r.Form["age"])
	fmt.Println("email：",r.Form["email"])
	fmt.Println("mobile：",r.Form["mobile"])
	fmt.Println("sex：",r.Form["sex"])
	fmt.Println("interest：",r.Form["interest"])
	fmt.Println("usercard：",r.Form["usercard"])

	//验证表单数据--------------------------------------------------------
	//必填字段
	if len(r.Form["username"][0]) == 0{
		//返回到客户端
		fmt.Fprintf(w,"用户名不能为空")
		return false
	}
	//unicode判断中文
	for _, val := range r.Form.Get("username"){
		if unicode.Is(unicode.Scripts["Han"], val) {
			//有中文
			fmt.Fprintf(w,"unicode姓名不能有中文\n")
		}
	}
	//正则判断中文
	chiReg := regexp.MustCompile("[\u4e00-\u9fa5]+")
	for _, val := range r.Form.Get("username") {
		if chiReg.MatchString(string(val)){
			//有中文
			fmt.Fprintf(w,"regexp姓名不能有中文\n")
			return false
		}
	}

	//判断英文
	match2,_ := regexp.MatchString("^[a-zA-Z]+$",r.Form.Get("username"))
	if !match2{
		fmt.Fprintf(w,"姓名只能是英文字母")
		return false
	}


	//使用转化判断数字
	getint,err := strconv.Atoi(r.Form.Get("age"))
	if err != nil{
		//数字转化出错了，那么可能就不是数字了
		fmt.Fprintf(w,"strconv年龄只能是数字\n")
	}
	//使用正则判断数字 ,应该尽量避免使用正则表达式，因为使用正则表达式的速度会比较慢
	m,_ := regexp.MatchString("^[0-9]+$",r.Form.Get("age"));
	if !m {
		fmt.Fprintf(w,"regexp年龄只能是数字\n")
		return false
	}

	//判断数字的大小范围
	if getint > 200{
		fmt.Fprintf(w,"年龄不能大于200")
		return false
	}


	//判断电子邮件
	match3,_ := regexp.MatchString(`^([\w\.\_]{2,})@(\w{1,})\.([a-z]{2,4})$`,r.Form.Get("email"))
	if !match3{
		fmt.Fprintf(w,"电子邮箱格式不正确")
		return false
	}

	//判断手机号码
	match4,_ := regexp.MatchString(`^1[0-9][0-9]\d{8}$`,r.Form.Get("mobile"))
	if !match4{
		fmt.Fprintf(w,"手机号码不正确")
		return false
	}

	//判断身份证号码
	//15位 或18位
	usercard := r.Form.Get("usercard")
	match5,_ := regexp.MatchString(`^(\d{15})$`,usercard)
	match6,_ := regexp.MatchString(`^(\d{17})([0-9]|X)$`,usercard)
	fmt.Println(usercard)
	fmt.Println(len(usercard))
	fmt.Println(!match6)
	if len(usercard) <= 15 && !match5{
		fmt.Fprintf(w,"身份证不正确")
		return false
	}else if len(usercard) > 15 && !match6{
		fmt.Fprintf(w,"身份证不正确")
		return false
	}

	//判断下拉菜单的值是否符合预设值
	sex_slice := []string{"1","2"}
	sex_value := r.Form.Get("sex")
	is_set := 0 //是否符合预设值 ，1 是，0否
	for _,val := range sex_slice{
		if sex_value == val{
			is_set = 1
		}
	}
	if is_set == 0{
		fmt.Fprintf(w,"性别不正确")
		return false
	}

	//判断复选框
	interest_slice := []string{"football","basketball","tennis"}
	interest := r.Form["interest"]
	array := Slice_diff(interest,interest_slice)
	fmt.Println("interest：",array)
	if array != nil{
		//有差异
		fmt.Fprintf(w,"兴趣选择不正确")
		return false
	}
	return true
}

//预防跨站脚本
func xss(w http.ResponseWriter,r *http.Request) bool{
	//template.HTMLEscape(w io.Writer, b []byte) //把b进行转义之后写到w
	//template.HTMLEscapeString(s string) string //转义s之后返回结果字符串
	//template.HTMLEscaper(args ...interface{}) string //支持多个参数一起转义，返回结果字符串

	//输出到服务器端
	fmt.Println("xss:", template.HTMLEscapeString(r.Form.Get("xss")))

	//转义后输出到客户端
	template.HTMLEscape(w, []byte(r.Form.Get("xss")))
	fmt.Fprintf(w, "\n") //换行

	//转义后输出到客户端
	t,err := template.New("foo").Parse(`{{define "T"}}hello,{{.}}!{{end}}`)
	err = t.ExecuteTemplate(w,"T","<script>alert('you have been pwned')</script>")
	if err != nil{
		fmt.Println(err)
	}
	fmt.Fprintf(w, "\n") //换行

	//使用 text/template 完全输出到客户端
	t1,err1 := texttemplate.New("foo").Parse(`{{define "T"}}hello,{{.}}!{{end}}`)
	err1 = t1.ExecuteTemplate(w,"T","<script>alert('you have been pwned')</script>")
	if err1 != nil{
		fmt.Println(err1)
	}
	fmt.Fprintf(w, "\n") //换行

	//使用 html/template template.HTML 完全输出到客户端
	t2, err2 := template.New("foo").Parse(`{{define "T"}}hello,{{.}}!{{end}}`)
	err2 = t2.ExecuteTemplate(w,"T",template.HTML("<script>alert('you have been HTML')</script>"))
	if err2 != nil{
		fmt.Println(err2)
	}
	return true
}

//获取上传文件
func uploadHandler(w http.ResponseWriter,r *http.Request) bool{
	//通过上面的代码可以看到，处理文件上传我们需要调用r.ParseMultipartForm，
	// 里面的参数表示maxMemory，调用ParseMultipartForm之后，
	// 上传的文件存储在maxMemory大小的内存里面，
	// 如果文件大小超过了maxMemory，那么剩下的部分将存储在系统的临时文件中。
	// 我们可以通过r.FormFile获取上面的文件句柄，然后实例中使用了io.Copy来存储文件。
	r.ParseMultipartForm(32<<20)
	file,handler,err := r.FormFile("userheader")
	//文件的handler 是multipart.FileHeader ，里面存储了如下结构信息
	/*
	type FileHeader struct{
		Filename string
		Header textproto.MIMEHeader
		//contains filtered or unexported fields
	}
	*/
	if err != nil{
		fmt.Println(err)
		return false
	}
	defer file.Close()
	fmt.Fprintf(w,"%v",handler.Header)
	f,err := os.OpenFile("./"+handler.Filename,os.O_WRONLY|os.O_CREATE,0777)
	if err != nil{
		fmt.Println(err)
		return false
	}
	defer f.Close()
	io.Copy(f,file)

	return true
}

func main() {
	http.HandleFunc("/",sayhelloName) //设置路由
	http.HandleFunc("/info",info) //设置路由
	err := http.ListenAndServe(":6665",nil) //设置监听的端口
	if err != nil{
		log.Fatal("ListenAndServe：",err)
	}
}