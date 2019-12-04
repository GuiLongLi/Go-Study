package main

import (
	"html/template"
	"os"
	"fmt"
	"strings"
)


func main() {
	fmt.Fprintf(os.Stdout,"%v\n","testtemplate")
	testtemplate()

	fmt.Fprintf(os.Stdout,"\n")
	fmt.Fprintf(os.Stdout,"\n")
	fmt.Fprintf(os.Stdout,"%v\n","testforeach")
	testforeach()

	fmt.Fprintf(os.Stdout,"\n")
	fmt.Fprintf(os.Stdout,"%v\n","testcondition")
	testcondition()

	fmt.Fprintf(os.Stdout,"\n")
	fmt.Fprintf(os.Stdout,"%v\n","testtemplatevar")
	testtemplatevar()

	fmt.Fprintf(os.Stdout,"\n")
	fmt.Fprintf(os.Stdout,"%v\n","testtemplatefunction")
	testtemplatefunction()

	fmt.Fprintf(os.Stdout,"\n")
	fmt.Fprintf(os.Stdout,"%v\n","testmust")
	testmust();
}


type Person struct {
	UserName string
	email string //未导出的字段，首字母是小写的
}
func testtemplate(){
	t := template.New("fieldname example")
	t,_ = t.Parse("hello {{.UserName}}! email: {{.email}}")
	p := Person{
		UserName:"template",
	}
	t.Execute(os.Stdout,p)
}

/*
输出嵌套字段内容
上面我们例子展示了如何针对一个对象的字段输出，那么如果字段里面还有对象，如何来循环的输出这些内容呢？我们可以使用{{with …}}…{{end}}和{{range …}}{{end}}来进行数据的输出。

{{range}} 这个和Go语法里面的range类似，循环操作数据
{{with}}操作是指当前对象的值，类似上下文的概念
*/
type Friend struct {
	Fname string
}
type PersonForeach struct {
	UserName string
	Emails []string
	Friends []*Friend
}
func testforeach(){
	f1 := Friend{Fname:"friend1"}
	f2 := Friend{Fname:"friend2"}
	t := template.New("testforeach")
	t,_ = t.Parse(`hello {{.UserName}}!
{{range .Emails}}
an email {{.}}
{{end}}
{{with .Friends}}
{{range .}}
my friend name is {{.Fname}}
{{end}}
{{end}}
`)
	p := PersonForeach{
		UserName:"hello world",
		Emails:[]string{"hello@163.com","world@qq.com"},
		Friends:[]*Friend{&f1,&f2},
	}
	t.Execute(os.Stdout,p)
}

/*
在Go模板里面如果需要进行条件判断，那么我们可以使用和Go语言的if-else语法类似的方式来处理，如果pipeline为空，那么if就认为是false，下面的例子展示了如何使用if-else语法：
注意：if里面无法使用条件判断，例如.Mail=="astaxie@gmail.com"，这样的判断是不正确的，if里面只能是bool值
*/
func testcondition(){
	tE := template.New("empty template")
	tE = template.Must(tE.Parse("空 pipeline if demo: {{if ``}} 不会输出. {{end}} \n"))
	tE.Execute(os.Stdout,nil)

	tW := template.New("with template")
	tW = template.Must(tW.Parse("不为空的 pipeline if demo: {{if `anything`}} 我有内容，我会输出. {{end}}\n"))
	tW.Execute(os.Stdout,nil)

	tElse := template.New("else template")
	tElse = template.Must(tElse.Parse("if-else demo: {{if `anything`}} if部分 {{else}} else部分.{{end}}\n"))
	tElse.Execute(os.Stdout,nil)
}

/*
模板变量
有时候，我们在模板使用过程中需要定义一些局部变量，我们可以在一些操作中申明局部变量，例如with``range``if过程中申明局部变量，这个变量的作用域是{{end}}之前，Go语言通过申明的局部变量格式如下所示：
*/
func testtemplatevar(){
	tVar := template.New("var template")
	//{{"output" | printf "%q"}}
	//竖线 |  左边的结果output 作为printf 函数最后一个参数。（等同于：printf("%q", "output")。）
	tVar = template.Must(tVar.Parse(`
{{with $x := "output" | printf "%s"}}{{$x}}{{end}}
{{with $x := "output"}}{{printf "%v" $x}}{{end}}
{{with $x := "output"}}{{$x | printf "%q"}}{{end}}
`))
	tVar.Execute(os.Stdout,nil)
}

/*
模板函数
下面我们将使用自定义函数把 @ 转换为 at
每个模板函数都有一个唯一值的名字，然后与一个 go 函数关联
type FuncMap map[string]interface{}
*/
func EmailDealWith(args ...interface{}) string {
	ok := false
	var s string
	if len(args) == 1 {
		s, ok = args[0].(string)
	}
	if !ok {
		s = fmt.Sprint(args...)
	}
	// find the @ symbol
	substrs := strings.Split(s, "@")
	if len(substrs) != 2 {
		return s
	}
	// replace the @ by " at "
	return (substrs[0] + " at " + substrs[1])
}

func testtemplatefunction(){
	t := template.New("fieldname example")
	//把 go 函数 EmailDealWith 绑定到 emailDeal 模板函数上
	t = t.Funcs(template.FuncMap{"emailDeal":EmailDealWith})
	t,_ = t.Parse(`exchange email {{.|emailDeal}}`)
	email := "admin@baidu.com"
	t.Execute(os.Stdout,email)
}

/*
Must操作
Must函数，可以检测模板是否正确，如大括号是否匹配，
如果不符合，则返回错误
*/
func testmust(){
	t := template.New("must")
	template.Must(t.Parse("some static text {{.Name}}"))
	fmt.Println("must ok")

}