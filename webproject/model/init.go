package model

import (
	"database/sql"
	"log"
	"fmt"
	"os"
	"errors"
	"strings"
	"io/ioutil"
	_ "github.com/go-sql-driver/mysql" //这个引用是必不可少的，因为需要调用driver.go文件里的init方法来提供一个数据库驱动程序
	"github.com/spf13/viper"
)

var DB *sql.DB     //全局变量，这样可以在别处调用

func Init() error{
	var err error
	var connectstring string
	max_idle_conns := viper.GetInt("common.database.max_idle_conns")
	username := viper.GetString("common.database.username")
	password := viper.GetString("common.database.password")
	host := viper.GetString("common.database.host")
	port := viper.GetInt("common.database.port")
	dbname := viper.GetString("common.database.dbname")
	// 账号:密码@tcp(IP:端口号)/数据库名?parseTime=true&charset=utf8&loc=Local
	connectstring = fmt.Sprintf("%s%s%s%s%s%s%d%s%s%s",username,":",password,"@tcp(",host,":",port,")/",dbname,"?parseTime=true&charset=utf8&loc=Local")

	//这行代码的作用就是初始化一个sql.DB对象
	DB,err = sql.Open("mysql",connectstring)
	if err != nil{
		return err
	}

	//设置最大超时时间
	DB.SetMaxIdleConns(max_idle_conns)

	//建立连接
	err = DB.Ping()
	if err != nil{
		return err
	}else{
		log.Println("数据库连接成功")
	}
	return nil
}

//检查是否安装数据库
func CheckInstalled() error{
	//先判断是否已经生成 install.lock文件，表示已经安装过了
	exists,_ := PathExists("install.lock")
	if exists {
		//fmt.Println("已经安装过数据库")
		return errors.New("已经安装过数据库")
	}
	return nil

}

//安装数据库
func CreateDatabase() error{
	//读取 database.sql 全部内容
	database, err := ioutil.ReadFile("database.sql")
	if err != nil {
		//fmt.Println("ioutil ReadFile database.sql error: ", err)
		return err
	}
	// 读取的数据转换成字符串类型
	databasestring := string(database)
	//去除换行符
	databasestring = strings.Replace(databasestring, "\n", "", -1)
	databasestring = strings.Replace(databasestring, "\r", "", -1)
	//strings.Split(string , "切割符") 把字符串转换成 数组
	databasearr := strings.Split(string(databasestring),";")
	//fmt.Printf("%v\n", reflect.TypeOf(databasearr))
	//fmt.Printf("%v\n", databasearr)

	//循环sql 数组，使用事务提交到数据库
	if(len(databasearr) > 0){
		transaction_con,err := DB.Begin()
		if err != nil{
			return err
		}
		for i := 0;i < len(databasearr);i++{
			if databasearr[i] == ""{
				continue
			}
			_, err = transaction_con.Exec(databasearr[i])
			if err != nil{
				// 失败回滚
				transaction_con.Rollback()
				return err
			}
		}
		//提交数据库
		err = transaction_con.Commit()
		if(err != nil){
			// 失败回滚
			transaction_con.Rollback()
			return err
		}
		//写入install.lock文件
		lock := []byte("1\n")
		err = ioutil.WriteFile("install.lock", lock, 0777)
		if err != nil{
			return err
		}
	}

	//
	//stmt,err := DB.Prepare(databasestring)
	//if err != nil{
	//	return err
	//}
	//defer stmt.Close()
	//_,err = stmt.Exec()
	//if err != nil{
	//	return err
	//}
	return nil
}

//检查文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}