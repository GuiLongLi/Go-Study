package main

import (
	"fmt"
	"os"
	"io/ioutil"
)

//7.5 文件操作
func main() {
	fmt.Println("testdir")
	testdir()

	fmt.Println()
	fmt.Println("testfile")
	testfile()
}

func testdir(){
	/*
	目录操作
		文件操作的大多数函数都是在os包里面，下面列举了几个目录操作的：

		func Mkdir(name string, perm FileMode) error

		创建名称为name的目录，权限设置是perm，例如0777

		func MkdirAll(path string, perm FileMode) error

		根据path创建多级子目录，例如astaxie/test1/test2。

		func Remove(name string) error

		删除名称为name的目录，当目录下有文件或者其他目录时会出错

		func RemoveAll(path string) error

		根据path删除多级子目录，如果path是单个名称，那么该目录下的子目录全部删除。
	*/
	os.Mkdir("test",0777)
	os.MkdirAll("test/test1/test11",0777)
	nowdir := "./test"
	files, _ := ioutil.ReadDir(nowdir)
	fmt.Println("当前目录："+nowdir)
	for _, f := range files {
		fmt.Printf(" ",f.Name())
	}
	fmt.Println()

	err := os.Remove("test")
	if err != nil {
		fmt.Println(err)
	}
	err = os.RemoveAll("test")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("testdir done")
}

func testfile(){
	/*
	文件操作
		建立与打开文件
		新建文件可以通过如下两个方法

		func Create(name string) (file *File, err Error)

		根据提供的文件名创建新的文件，返回一个文件对象，默认权限是0666的文件，返回的文件对象是可读写的。

		func NewFile(fd uintptr, name string) *File

		根据文件描述符创建相应的文件，返回一个文件对象

		通过如下两个方法来打开文件：

		func Open(name string) (file *File, err Error)

		该方法打开一个名称为name的文件，但是是只读方式，内部实现其实调用了OpenFile。

		func OpenFile(name string, flag int, perm uint32) (file *File, err Error)

		打开名称为name的文件，flag是打开的方式，只读、读写等，perm是权限

		写文件
		写文件函数：

		func (file *File) Write(b []byte) (n int, err Error)

		写入byte类型的信息到文件

		func (file *File) WriteAt(b []byte, off int64) (n int, err Error)

		在指定位置开始写入byte类型的信息

		func (file *File) WriteString(s string) (ret int, err Error)

		写入string信息到文件
	*/
	userFile := "testfile.txt"
	fileresource,err := os.Create(userFile)
	if err != nil {
		fmt.Println(userFile,err)
		return
	}
	defer fileresource.Close()
	for i := 0;i < 10;i++{
		fileresource.WriteString("Just a string!\r\n")
		fileresource.Write([]byte("just a byte!\r\n"))
	}

	fmt.Println("testfile write done")

	/*
	读文件
		读文件函数：

		func (file *File) Read(b []byte) (n int, err Error)

		读取数据到b中

		func (file *File) ReadAt(b []byte, off int64) (n int, err Error)

		从off开始读取数据到b中
	*/

	fileresource,err = os.Open(userFile)
	if err != nil {
		fmt.Println(userFile,err)
		return
	}
	defer fileresource.Close()
	filebuffer := make([]byte,1024) //创建一个文件缓存，大小1024byte

	for{
		n,_ := fileresource.Read(filebuffer)
		if 0 == n{
			break
		}
		//os.Stdout.Write(filebuffer[:n]) //直接输出到屏幕
	}
	fmt.Println(userFile+"的内容是：",string(filebuffer)) //转换成字符串输出
	fmt.Println("testfile read done")

	/*
	删除文件
		Go语言里面删除文件和删除文件夹是同一个函数

		func Remove(name string) Error

		调用该函数就可以删除文件名为name的文件
	*/
	err = os.Remove(userFile);
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println("testfile remove done")

}