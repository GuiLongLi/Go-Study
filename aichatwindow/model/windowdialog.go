package model

import (
	"fmt"
	"strconv"

	"github.com/tidwall/gjson"
	"github.com/spf13/viper"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"

	. "aichatwindow/handler"
)

/*-------------------window界面-------------------*/
func OpenWindow() {
	mw := &MyMainWindow{}

	if _, err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "AI对话",
		//Icon:     "test.ico", //图标文件路径
		MinSize:  Size{300, 400},
		Layout:   VBox{},
		Children: []Widget{
			GroupBox{
				Layout: HBox{},
				Children: []Widget{
					LineEdit{ //对话输入框
						AssignTo: &mw.searchBox,
						OnKeyDown: func(key walk.Key) { //键盘事件
							if key == walk.KeyReturn {//回车键
								mw.clicked()
							}
						},
					},
					PushButton{ //对话按钮
						Text:      "对话",
						OnClicked: mw.clicked,  //鼠标点击左键事件
					},
				},
			},
			ListBox{ //聊天记录框
				AssignTo: &mw.message,
			},
			ListBox{ //实时聊天框
				AssignTo: &mw.results,
				Row:      5,
			},
		},
	}.Run());
	err != nil {
		fmt.Printf("Run err: %+v\n",err)
	}
}

type MyMainWindow struct {
	*walk.MainWindow
	searchBox *walk.LineEdit
	message  *walk.ListBox
	results   *walk.ListBox
}

func (mw *MyMainWindow) clicked() {
	word := mw.searchBox.Text()
	message := mw.message.Model() //获取以前的聊天记录
	fmt.Println("message",message);
	model := []string{} //ai对话
	record := []string{} //聊天记录
	if message != nil {
		for _,v := range message.([]string){
			record = append(record,v)  //插入以前的聊天记录
		}
	}
	record = append(record,fmt.Sprintf("我：%v", word)) //先显示我的 信息
	for _, value := range search( word) {
		model = append(model, fmt.Sprintf("ai回复：%v", value))
		record = append(record, fmt.Sprintf("ai回复：%v", value))  //再显示ai回复信息
	}
	mw.results.SetModel(model) //实时聊天输出
	mw.message.SetModel(record) //聊天记录输出

}

/*-------------------图灵接口-------------------*/
//获取tuling接口回复
func TulingAi(info string) (string,error) {
	api := viper.GetString("common.tuling.api")

	//发送http请求图灵api  , body是http响应
	var body, resultErrs = HttpRequest(api,info,"POST")
	if resultErrs != nil {
		fmt.Printf("HttpRequest err: %+v\n",resultErrs)
	}

	return body, nil
}

//回复信息构造体
type tlReply struct {
	code int
	Text string `json:"text"`
}

//图灵搜索
func search( word string) (res []string) {
	res = []string{}
	var userId = "1";
	//图灵接口参数构造体
	var chattingInfo = BuildChatting(word,userId, viper.GetString("common.tuling.apikey"))
	fmt.Printf("chattingInfo: %+v\n",chattingInfo)
	// 参数构造体 转换成 字符串
	chatstr,err := ConvertJson(chattingInfo)
	if err != nil{
		fmt.Printf("ConvertJson err: %+v\n",err)
		return
	}

	//调用图灵接口
	body,err := TulingAi(chatstr)
	if err != nil{
		fmt.Printf("TulingAi err: %+v\n",err)
		return
	}
	fmt.Printf("body err: %+v\n",body)
	var results string
	// 使用gjson 获取返回结果的 resultType
	result := gjson.Get(body, "results.#.resultType")
	for key, name := range result.Array() {
		//如果 resultType 是 text格式
		if name.String() == "text"{
			//获取对应 key 的 values里的text ，就是图灵回复的文字
			getstring := "results."+strconv.Itoa(key)+".values.text"
			fmt.Printf("getstring: %+v\n",getstring)
			result_text := gjson.Get(body,getstring)
			results = result_text.String()
			res = append(res, results)
		}
	}
	return
}