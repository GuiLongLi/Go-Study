package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)


func main() {
	//读取xml
	fmt.Println("xmlRead")
	xmlRead()

	//编写xml
	fmt.Println()
	fmt.Println("xmlWrite")
	xmlWrite()
}


type Recurlyservers struct {
	XMLName xml.Name `xml:"servers"`
	Version string `xml:"version,attr"`
	Svs []serverRead `xml:"server"`
	Description string `xml:",innerxml"`
}

type serverRead struct {
	XMLName xml.Name `xml:"server"`
	ServerName string `xml:"serverName"`
	ServerIP string `xml:"serverIP"`
}
//读取xml
func xmlRead(){
	file,err := os.Open("textprocess/testxml.xml") //for read access
	if err != nil{
		fmt.Printf("error:%v",err)
		return
	}
	defer file.Close()
	data,err := ioutil.ReadAll(file)
	if err != nil{
		fmt.Printf("error:%v",err)
		return
	}
	v := Recurlyservers{}
	err = xml.Unmarshal(data,&v)
	if err != nil{
		fmt.Printf("error:%v",err)
		return
	}
	fmt.Printf("%+v\n",v)
}


type Servers struct {
	XMLName xml.Name `xml:"servers"`
	Version string `xml:"version,attr"`
	Svs []serverWrite `xml:"server"`
}
type serverWrite struct {
	ServerName string `xml:"serverName"`
	ServerIP string `xml:"serverIP"`
}
//编写xml
func xmlWrite(){
	v := &Servers{Version:"1"}
	v.Svs = append(v.Svs,serverWrite{"Shanghai_VPN","127.0.0.1"})
	v.Svs = append(v.Svs,serverWrite{"Beijing_VPN","127.0.0.2"})
	output,err := xml.MarshalIndent(v," ","  ")
	if err != nil{
		fmt.Printf("error:%v\n",err)
	}
	os.Stdout.Write([]byte(xml.Header))
	os.Stdout.Write(output)
}