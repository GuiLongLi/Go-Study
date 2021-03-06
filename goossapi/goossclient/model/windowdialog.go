package model

import (
	"io"
	"os"
	"os/exec"
	"fmt"
	"time"
	"bytes"
	"strings"
	"strconv"
	"crypto/md5"
	"path/filepath"
	"mime/multipart"

	"github.com/tidwall/gjson"
	"github.com/spf13/viper"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"

	. "goossclient/handler"
)

/*-------------------window界面-------------------*/
/*
# 不隐藏cmd 窗口
[root@localhost aichatwindow]# CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build

# 隐藏cmd 窗口
[root@localhost aichatwindow]# CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui"
*/
type MyMainWindow struct {
	*walk.MainWindow
	model *MessageModel
	selectedfile *walk.LineEdit
	message  *walk.ListBox
	wv *walk.WebView
	//imageview  *walk.ImageView
}

func OpenWindow() {
	mw := &MyMainWindow{model: NewMessageModel()}

	if _, err := (MainWindow{
		AssignTo: &mw.MainWindow,
		Title:    "阿里云文件上传",
		//Icon:     "test.ico", //图标文件路径
		MinSize:  Size{300, 400},
		Layout:   VBox{},
		Children: []Widget{
			GroupBox{
				MaxSize: Size{500, 500},
				Layout: HBox{},
				Children: []Widget{
					PushButton{ //选择文件
						Text:      "打开文件",
						OnClicked: mw.selectFile, //点击事件
					},
					Label{Text: "选中的文件 "},
					LineEdit{
						AssignTo: &mw.selectedfile, //选中的文件
					},
					PushButton{ //上传
						Text:      "上传",
						OnClicked: mw.uploadFile,  //上传
					},
				},
			},
			ListBox{ //记录框
				AssignTo: &mw.message,
				OnCurrentIndexChanged: mw.lb_CurrentIndexChanged, //单击
				OnItemActivated:       mw.lb_ItemActivated, //双击
			},
			Composite{
				Layout: Grid{Columns: 2, Spacing: 10},
				Children: []Widget{
					WebView{
						//MinSize:  Size{1000, 0},
						AssignTo: &mw.wv,
					},
				},
			},
		},
	}.Run());
	err != nil {
		fmt.Printf("Run err: %+v\n",err)
	}
}


func (mw *MyMainWindow) selectFile() {
	allowType := viper.GetStringSlice("common.server.allowtype")
	fmt.Printf("allowType: %+v\n",allowType)

	dlg := new(walk.FileDialog)
	dlg.Title = "选择文件"
	//dlg.Filter = "可上传jpg (*.jpg)|*.jpg|可上传png (*.png)|*.png|可上传gif (*.gif)|*.gif|所有文件 (*.*)|*.*"
	//判断可允许上传的文件
	filter := []string{}
	filterstring := ""
	for _,v := range allowType{
		if(v != "*"){
			filterstring = "可上传"+v+" (*."+v+")|*."+v
			filter = append(filter,filterstring)
		}else{
			filterstring = "所有文件"+v+" (*."+v+")|*."+v
			filter = append(filter,filterstring)
		}
	}
	dlg.Filter = strings.Join(filter,"|") //切片转换字符串
	fmt.Printf("dlg.Filter: %+v\n",dlg.Filter)

	record := getMessage(mw) //获取记录

	if ok, err := dlg.ShowOpen(mw); err != nil {
		mw.selectedfile.SetText("") //通过重定向变量设置TextEdit的Text
		_=writeMessage(mw,"Error : File Open",record) //写入记录
		return
	} else if !ok {
		mw.selectedfile.SetText("") //通过重定向变量设置TextEdit的Text
		_=writeMessage(mw,"cancel",record) //写入记录
		return
	}
	s := fmt.Sprintf("Select : %s", dlg.FilePath)
	_=writeMessage(mw,s,record) //写入记录
	mw.selectedfile.SetText(dlg.FilePath) //通过重定向变量设置TextEdit的Text
}

func (mw *MyMainWindow) uploadFile(){
	filename := mw.selectedfile.Text() //上传的文件
	record := getMessage(mw) //获取记录

	fmt.Printf("filename: %+v\n",filename)
	if len(filename) == 0{
		fmt.Println("select a file")
		_=writeMessage(mw,"请选择文件",record)
		return
	}


	//-----------------------------------------------
	//模拟表单上传文件到服务器接口
	bodyBuf := &bytes.Buffer{} //使用二进制缓冲
	bodyBuf.Reset() //重置缓冲
	bodyWriter := multipart.NewWriter(bodyBuf) //创建表单

	//加入token
	curtime := time.Now().Unix()
	h := md5.New()
	io.WriteString(h,strconv.FormatInt(curtime,10))
	token := fmt.Sprintf("%x",h.Sum(nil)) //生成token
	//模拟表单数据
	_ = bodyWriter.WriteField("token", token)
	_ = bodyWriter.WriteField("username", "awef")

	//关键的一步操作 ，模拟上传文件
	//获取上传文件名称
	filenamesplit := strings.Split(filename,"\\")  //使用 \ 分隔符，将字符串切片
	filenamesplit = filenamesplit[len(filenamesplit)-1:] //获取切片最后一个
	newname := strings.Join(filenamesplit,"")  //切片转换为字符串
	fileWriter, err := bodyWriter.CreateFormFile("files", newname) //表单创建一个文件参数
	if err != nil {
		fmt.Println("error writing to buffer")
		_=writeMessage(mw,"error writing to buffer",record)
		return
	}

	//打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		_=writeMessage(mw,"error opening file",record)
		return
	}
	defer fh.Close()

	//iocopy ，复制源文件到 模拟表单的文件里
	_, err = io.Copy(fileWriter, fh) //fileWriter目标文件 ，fh源文件 ，复制源文件到目标文件
	if err != nil {
		_=writeMessage(mw,"error io.Copy",record)
		return
	}

	contentType := bodyWriter.FormDataContentType() //获取表单类型
	fmt.Printf("contentType: %+v\n",contentType)
	bodyWriter.Close()

	//resp, err := http.Post(targetUrl, contentType, bodyBuf)
	resp, err := UploadRequest(contentType,bodyBuf)
	if err != nil{
		fmt.Printf("UploadApi err: %+v\n",err)
		_=writeMessage(mw,filename+" 上传失败",record)
		return
	}
	fmt.Printf("body : %+v\n",resp)
	//var results string
	// 使用gjson 获取返回结果的
	coderesult := gjson.Get(resp, "code")
	code := coderesult.String() //断言转换成字符串
	if code != "1"{
		_=writeMessage(mw,filename+" 上传失败",record)
		return
	}
	record=writeMessage(mw,filename+" 上传成功",record)

	result := gjson.Get(resp, "data.0.0.url")
	fmt.Printf("result: %+v\n",result)
	resultstring := result.String()//断言转换成字符串
	_=writeMessage(mw,resultstring,record)

	return
}

//消息记录改变事件
func (mw *MyMainWindow) lb_CurrentIndexChanged() {
	fmt.Printf("mw.message.CurrentIndex(): ",mw.message.CurrentIndex())
	fmt.Println()
	return
}
//消息记录点击事件
func (mw *MyMainWindow) lb_ItemActivated() {
	fmt.Println("mw.message.CurrentIndex(): ",mw.message.CurrentIndex())
	fmt.Println();
	fmt.Printf("mw.model.items: %+v ",mw.model.items)
	fmt.Println();

	index := mw.message.CurrentIndex()
	imagename := mw.model.items[index].Name //获取当前选中名称
	image := mw.model.items[index].Value //获取当前选中值
	is_http := strings.Index(image,"http") //查找字符串位置
	is_ie := strings.Index(imagename,"默认浏览器") //查找字符串位置
	if is_http != -1{ // -1 是找不到
		//walk.MsgBox(mw, "Value", value, walk.MsgBoxIconInformation) //提示框
		fmt.Printf("image : %+v ",image)
		fmt.Println();

		if is_ie != -1{
			openImageExplorer(image) ////使用默认浏览器打开图片
		}else{
			openImageWebview(mw,image) ////使用webview打开图片
		}
	}
	return
}

//使用默认浏览器打开图片
func openImageExplorer(image string){
	cmd := exec.Command("explorer", image)
	err := cmd.Start()
	if err != nil {
		fmt.Println(err)
	}
	return
}

//创建html,使用webview打开图片
func openImageWebview(mw *MyMainWindow,image string){
	go func() {
		h := md5.New()
		io.WriteString(h,image)
		htmlname := fmt.Sprintf("%x",h.Sum(nil)) //生成html名称
		userFile := ""+htmlname+".html"
		fileall := getCurrentDirectory() + "/"+userFile //html的绝对路径

		fout, err := os.Create(fileall) //创建html文件
		defer fout.Close()
		if err != nil {
			fmt.Println(err)
			return
		}

		html := `<!DOCTYPE html><html><head><meta http-equiv="Content-Type" content="text/html; charset=utf-8" /><title>图片展示</title></head><body><a href="%s" target="_self">%s</a><br><img src="%s" alt=""></body></html>`
		html = fmt.Sprintf(html,image,image,image)  //插入图片
		fout.WriteString(html) //把字符串写入html

		mw.wv.SetURL("file:///" + getCurrentDirectory() + "/"+userFile)
	}()
}
//获取当前文件路径
func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println(err.Error())
	}
	return strings.Replace(dir, "\\", "/", -1)
}

//消息记录每个模型
type MessageItem struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
//消息记录模型
type MessageModel struct {
	walk.ListModelBase
	items []MessageItem
}
//新建消息模型
func NewMessageModel() *MessageModel {
	m := &MessageModel{items: []MessageItem{}}
	return m
}
//插入到消息模型中
func appendMessageModel(mw *MyMainWindow,message string) {
	item := MessageItem{
		Name:  message,
		Value: message,
	}
	mw.model.items = append(mw.model.items,item)

}


//获取记录
func getMessage(mw *MyMainWindow) []string{
	message := mw.message.Model() //获取以前的记录
	fmt.Println("message",message)
	record := []string{} //记录
	if message != nil {
		for _,v := range message.([]string){
			record = append(record,v)  //插入以前的记录
		}
	}
	return record
}

//写入记录
func writeMessage(mw *MyMainWindow,message string,record []string)[]string{
	is_http := strings.Index(message,"http") //查找字符串位置
	message_record := message
	if is_http != -1{ // -1 是找不到
		message_record = "双击查看图片："+message

		//插入默认浏览器打开记录
		item := MessageItem{
			Name:  "双击用默认浏览器打开图片"+message,
			Value: message,
		}
		mw.model.items = append(mw.model.items,item)
		record = append(record,"双击用默认浏览器打开图片"+message)  //插入记录
	}
	record = append(record,message_record)  //插入记录
	appendMessageModel(mw,message) //插入模型
	mw.message.SetModel(record) //记录输出

	return record
}

/*-------------------上传接口-------------------*/

type OssStruct struct {
	Url string `json:"url"`
}
//上传请求
func UploadRequest(contenttype string,buff *bytes.Buffer) (string,error) {
	server := viper.GetString("common.server.server")
	port := viper.GetString("common.server.port")
	serverapi := viper.GetString("common.server.uploadapi")

	api := fmt.Sprintf("%s:%s/%s",server,port,serverapi)

	//发送http请求图灵api  , body是http响应
	var body, resultErrs = HttpRequest(api,contenttype,buff,"POST")
	if resultErrs != nil {
		fmt.Printf("HttpRequest err: %+v\n",resultErrs)
	}

	return body, nil
}
