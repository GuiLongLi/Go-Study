package main

import (
	"os"
	"fmt"
	"bytes"
	"reflect"
	"encoding/json"
	"github.com/json-iterator/go"
)

func main() {
	testjsonencode()
	testjsondecode()
	testunknowjsondecode()
	testjsontag()
	userrequester()
	testjsoniter()
	teststructjson()
	testjsonnum()

	jsonstr := []byte(`{"jsonrpc":"2.0","result":[{"host":"10297"}]}`)
	fmt.Println("usestruct")
	usestruct(jsonstr)
	fmt.Println()
	fmt.Println("notusestruct")
	notusestruct(jsonstr)

	jsonendecode()

}


type Animal struct {
	Name string `json:"name"`
	Weight string `json:"weight"`
}

//构造体 转 json字符串
func testjsonencode(){
	var animals []Animal
	animals = append(animals,Animal{Name:"Elephant",Weight:"3 ton"})
	animals = append(animals,Animal{Name:"Whale",Weight:"10 ton"})

	jsonstr,err := json.Marshal(animals)
	if(err != nil){
		fmt.Println("error:%v\n",err)
	}
	fmt.Println(string(jsonstr))
	for _,v := range animals {
		fmt.Println("v : %+v\n",v.Name);
	}
}

//字符串 转 json
func testjsondecode(){
	var jsonstr = []byte(`[
		{"Name":"李四","Weight":"45kg"},
		{"Name":"张三","Weight":"88kg"}
	]`)

	var people []Animal

	err := json.Unmarshal(jsonstr,&people)
	if(err != nil){
		fmt.Println("error:%v\n",err)
	}
	for key,val := range people{
		fmt.Printf("people的key是%v val是%+v\n", key,val)

		getType := reflect.TypeOf(val)
		fmt.Println("get Type is :", getType.Name())

		getValue := reflect.ValueOf(val)
		fmt.Println("get all Fields is:", getValue)

		// 获取方法字段
		// 1. 先获取interface的reflect.Type，然后通过NumField进行遍历
		// 2. 再通过reflect.Type的Field获取其Field
		// 3. 最后通过Field的Interface()得到对应的value
		for i := 0; i < getType.NumField(); i++ {
			field := getType.Field(i)
			value := getValue.Field(i).Interface()
			fmt.Printf("%s: %v = %v\n", field.Name, field.Type, value)
		}
		println()
	}
}

//未知的json 解析
func testunknowjsondecode(){
	//在解析 JSON 的时候，任意动态的内容都可以解析成 interface{}。
	var unknow interface{}
	jsonstr := []byte(`{"Name":"hello","Sex":1,"parents":["Whale","Elaphant"]}`)
	json.Unmarshal(jsonstr,&unknow)
	for key,val := range unknow.(map[string]interface{}){

		switch va := val.(type) { //判断 unknow下级的类型
		case string:
			fmt.Println(key," is string ",va)
		case int:
			fmt.Println(key, "is int ", va)
		case float64:
			fmt.Println(key, "is float64 ", va)
		case []interface{}:
			fmt.Println(key, "is array:")
			for i, j := range va {
				fmt.Println(i, j)
			}
		}
	}
	println()
}
//-----------------------------------------------------------------
/*
json tag 有很多值可以取，同时有着不同的含义，比如：

-：不要解析这个字段，表示该字段不会输出到 JSON

omitempty 当字段为空（默认值）时，不要解析这个字段。比如 false、0、nil、长度为 0 的 array，map，slice，string，就不会输出到JSON 串中

FieldName，当解析 json 的时候，使用这个名字

,string当字段类型是 bool, string, int, int64 等，而 tag 中带有该选项时，那么该字段在输出到 JSON 时，会把该字段对应的值转换成 JSON 字符串.

----------------------------
示例：

// 解析的时候忽略该字段。默认情况下会解析这个字段，因为它是大写字母开头的
Field int   `json:"-"`

// 解析（encode/decode） 的时候，使用 `other_name`，而不是 `Field`
Field int   `json:"other_name"`

// 解析的时候使用 `other_name`，如果struct 中这个值为空，就忽略它
Field int   `json:"other_name,omitempty"`


// 解析的时候会将接受到的字符串类型转为int类型
Field int   `json:"other_name,string"`
*/
type People struct {
	Name string `json:"-"`
	Age int `json:"age"`  //注意 取别名不能使用中文
	Sex int `json:"sex,omitempty"`
	Weight int `json:"weight,string"`
}
func testjsontag(){
	var people1 []People  //二维数组
	var people2 People    //一维数组
	var jsonstr []byte
	var err error
	jsonstr = []byte(`[
		{"Name":"张三","Age":111,"Sex":2,"Weight":"46"}
	]`)
	err = json.Unmarshal(jsonstr,&people1)
	fmt.Printf("err是%v\n", err)
	fmt.Printf("people1是%+v\n", people1) //忽略了 Name属性

	// 获取tag中的内容
	t := reflect.TypeOf(people1)
	field := t.Elem().Field(0)
	fmt.Println(field.Tag)
	fmt.Println(field.Tag.Get("json"))

	jsonstr = []byte(`{
		"Name":"李四",
		"Age":666,
		"Sex":0,
		"Weight":"33"
	}`)
	err = json.Unmarshal(jsonstr,&people2)
	fmt.Printf("err是%v\n", err)
	fmt.Printf("people2是%+v\n", people2)

	println()
}

//自定义解析方法
/*
// Marshaler 接口定义了怎么把某个类型 encode 成 JSON 数据
type Marshaler interface {
        MarshalJSON() ([]byte, error)
}
// Unmarshaler 接口定义了怎么把 JSON 数据 decode 成特定的类型数据。如果后续还要使用 JSON 数据，必须把数据拷贝一份
type Unmarshaler interface {
        UnmarshalJSON([]byte) error
}
*/
type UserRequest struct {
	Name string
	Mail Mail
	Phone Phone
}
type Mail struct {
	Value string
}
type Phone struct {
	Value string
}
func (mailer *Mail) MarshalJSON() (data []byte, err error){
	if(mailer != nil){
		data = []byte(mailer.Value)
	}
	return
}
func (mailer *Mail) UnmarshalJSON(data []byte) error{
	//判断 data 中 ，是否含有 @
	if(!bytes.Contains(data,[]byte("@"))){
		return fmt.Errorf("mail format error")
	}
	mailer.Value = string(data)
	fmt.Printf("current mail format\n")
	return nil
}
func (phone *Phone) MarshalJSON() (data []byte,err error){
	if(phone != nil){
		data = []byte(phone.Value)
	}
	return
}
func (phone *Phone) UnmarshalJSON(data []byte) error{
	//判断手机号码是否11位
	if(len(data) != 11){
		return fmt.Errorf("phone format error")
	}
	phone.Value = string(data)
	fmt.Printf("current phone format\n")
	return nil
}
func userrequester(){
	user := UserRequest{}
	user.Name = "Tellphone"

	var err error
	err = user.Mail.UnmarshalJSON([]byte("callmephone.com"))
	if(err != nil){
		fmt.Printf("%v\n", err)
	}
	err = user.Phone.UnmarshalJSON([]byte("1351234567"))
	if(err != nil){
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("%s的邮箱是%s 电话是%s\n", user.Name,user.Mail,user.Phone)

	err = user.Mail.UnmarshalJSON([]byte("callme@phone.com"))
	if(err != nil){
		fmt.Printf("%v\n", err)
	}
	err = user.Phone.UnmarshalJSON([]byte("13512345678"))
	if(err != nil){
		fmt.Printf("%v\n", err)
	}
	fmt.Printf("%s的邮箱是%s 电话是%s\n", user.Name,user.Mail,user.Phone)

	println()
}
//Json的编码器和解码器
/*
json包提供了解码器和编码器类型，以支持读取和写入json数据流的常见操作。在该包中使用NewDecoder和NewEncoder函数包装io。

func NewDecoder(r io.Reader) *Decoder
func NewEncoder(w io.Writer) *Encoder
*/
func jsonendecode(){
	//{"Name": "Platypus", "Order": "Monotremata"}
	fmt.Printf("请输入json字符串：")
	encoder := json.NewEncoder(os.Stdout)
	decoder := json.NewDecoder(os.Stdin)
	for{
		var v map[string]interface{}
		if err := decoder.Decode(&v);err != nil{
			fmt.Println(err) //解析出错，打印错误信息
			jsonendecode() //重新调用自身
			return
		}
		for k := range v{
			if k != "Name"{
				delete(v,k)
			}
		}
		if err := encoder.Encode(&v);err != nil{
			fmt.Println(err)
		}
	}

}

//推荐的 json 解析库
//jsoniter（json-iterator）是一款快且灵活的 JSON 解析器，同时提供 Java 和 Go 两个版本。从 dsljson 和 jsonparser 借鉴了大量代码。
/*
基本用法如下：
jsoniter.Marshal(&data)
jsoniter.Unmarshal(input, &data)
*/
var testString = `{"Name": "Platypus", "Order": "Monotremata"}`
func testjsoniter(){
	var animal interface{}
	var err error
	var jsonBlob = []byte(testString)

	err = jsoniter.Unmarshal(jsonBlob,&animal)
	if(err != nil){
		fmt.Printf("error %v\n",err)
	}
	fmt.Printf("animal的值是%v\n", animal)
}
//复合结构的解析
type Car struct {
	Name string
	Engine Engine
	Tire Tire
}
type Engine struct {
	Value string
}
type Tire struct {
	Value string
}
func teststructjson(){
	var jsonstr = []byte(`{"Name":"奔驰","Engine":{"Value":"自由梦"},"Tire":{"Value":"米其林"}}`)
	car := Car{}
	var engine Engine
	var tire Tire
	jsoniter.Unmarshal(jsonstr,&struct {
		*Car
		*Engine
		*Tire
	}{&car,&engine,&tire})
	fmt.Printf("%+v\n",car)
	fmt.Printf("小明从小红的%v上卸下了%v发动机用来改造他的车，连%v轮胎都不放过，真实孤终生啊！\n",car.Name,car.Engine.Value,car.Tire.Value)
}
//Unmarshal 精度问题
//golang使用json.Unmarshal的时候，有些数字类型的数据会默认转为float64，而一些数据因其比较大，导致输出的时候会出现数据与原数据不等的现象，解决办法是，将此数据类型变为json.Number
type Numb struct {
	Nid jsoniter.Number `json:"nid"`
}
func testjsonnum(){
	var jsonstr = `{"nid":114420234065740369922}`
	var number Numb

	jsoniter.Unmarshal([]byte(jsonstr),&number)

	fmt.Printf("使用Number转换前 %+v\n", number)
	fmt.Printf("使用Number转换后 %+v\n", number.Nid.String())
}



//自定义json 解析
type ResultStruct struct {
	Jsonrpc string `json:"jsonrpc"`
	Result []HostStruct `json:"result"`
}

type HostStruct struct {
	Host string `json:"host"`
}


func usestruct(jsonstr []byte){
	var res ResultStruct

	err := json.Unmarshal(jsonstr,&res)
	if(err != nil){
		fmt.Printf("error:%v\n",err)
	}
	fmt.Printf("res: %+v\n", res)
}

func notusestruct(jsonstr []byte){
	var unknow interface{}
	json.Unmarshal(jsonstr,&unknow)
	reprinln(unknow)
}

func reprinln(unknow interface{}){
	for key,val := range unknow.(map[string]interface{}){
		switch va := val.(type) { //判断 unknow下级的类型
		case string:
			fmt.Println(key," is string ",va)
		case int:
			fmt.Println(key, "is int ", va)
		case float64:
			fmt.Println(key, "is float64 ", va)
		case []interface{}:
			fmt.Println(key, "is array:")
			for i, val := range va {
				fmt.Println(i, val)
				reprinln(val)
			}
		}
	}
}