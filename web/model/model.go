package model

import (
	"fmt"
	"strconv"
	"web/config"
	"database/sql"
	_ "github.com/go-sql-driver/mysql" //这个引用是必不可少的，因为需要调用driver.go文件里的init方法来提供一个数据库驱动程序
)
/*
导入MySQL数据库驱动
import (
   "database/sql"
   _ "github.com/go-sql-driver/mysql"
)
通常来说, 不应该直接使用驱动所提供的方法, 而是应该使用 sql.DB, 因此在导入 mysql 驱动时, 这里使用了匿名导入的方式(在包路径前添加 _), 当导入了一个数据库驱动后, 此驱动会自行初始化并注册自己到Golang的database/sql上下文中, 因此我们就可以通过 database/sql 包提供的方法访问数据库了.
*/
var DB *sql.DB  //声明全局变量

func Init() *sql.DB{
	//加载配置文件
	//这行代码的作用就是初始化一个sql.DB对象
	config.Init()
	//获取数据库配置
	var dbconfig = config.GetDatabaseInfo();

	var err error
	var constring string
	max_idle_conns := dbconfig["max_idle_conns"]
	root := dbconfig["username"]
	password := dbconfig["password"]
	host := dbconfig["host"]
	port := dbconfig["port"]
	dbname := dbconfig["dbname"]
	fmt.Printf("dbconfig的值是%v\n", dbconfig)
	// constring 它的配置规则:
	// 账号:密码@tcp(IP:端口号)/数据库名?parseTime=true&charset=utf8&loc=Local
	constring = fmt.Sprintf("%s%s%s%s%s%s%d%s%s%s",root,":",password,"@tcp(",host,":",port,")/",dbname,"?parseTime=true&charset=utf8&loc=Local")
	fmt.Printf("constring的值是%v\n", constring)

	//打开mysql连接
	DB,err = sql.Open("mysql",constring)
	if(err != nil){
		panic(err)
	}
	//设置最大超时时间 DB.SetMaxIdleConns(int)
	//  max_idle_conns类型是 type interface {}
	//fmt.Sprintf("%d",max_idle_conns) 把 max_idle_conns 转换成字符串类型
	//strconv.Atoi(string) 把 string转换成 int类型
	max_idle_conns_int, _ := strconv.Atoi(fmt.Sprintf("%d",max_idle_conns))
	DB.SetMaxIdleConns(max_idle_conns_int)
	//建立链接
	err = DB.Ping()
	if nil != err{
		panic(err)
	}else{
		fmt.Println("Mysql Startup Normal!")
	}
	return DB
}

