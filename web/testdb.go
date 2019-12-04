package main

import (
	"fmt"
	"web/model"
)

func main() {
	initdb()
	testtransaction()
	testmysql()
}

//建立db的初始连接
func initdb(){
	var DB = model.Init()
	//建立链接
	err := DB.Ping()
	if(err != nil){
		panic("连接失败")
	}
	fmt.Printf("连接正常\n")
	DB.Close(); //关闭数据库连接池
	println();
}

//测试事务
func testtransaction(){
	var DB = model.Init()
	//开启事务
	var err error
	transaction_con,err := DB.Begin()
	if err != nil{
		panic(err)
	}
	result, err := transaction_con.Exec("insert into erp_admin(name,password) values(?,?)","test","123456")
	if err != nil{
		// 失败回滚
		transaction_con.Rollback()
		panic(err)
	}
	id, err := result.LastInsertId()//获取自增id
	//affected, err = result.RowsAffected() //获取update影响的行数
	if err != nil{
		// 失败回滚
		transaction_con.Rollback()
		panic(err)
	}
	fmt.Println(result)
	fmt.Println(id)

	result,err = transaction_con.Exec("insert into erp_admin_log(admin_id,action) values(?,?)",1,"test")
	if err != nil{
		// 失败回滚
		transaction_con.Rollback()
		panic(err)
	}
	id, err = result.LastInsertId() //获取自增id
	//affected, err = result.RowsAffected() //获取update影响的行数
	if err != nil{
		// 失败回滚
		transaction_con.Rollback()
		panic(err)
	}
	fmt.Println(result)
	fmt.Println(id)

	err = transaction_con.Commit()
	if(err != nil){
		// 失败回滚
		transaction_con.Rollback()
		panic(err)
	}
	fmt.Println("测试事务成功\n")

	DB.Close(); //关闭数据库连接池
	println();
}

func testmysql(){
	var DB = model.Init()
	//新建表
	/*
	DROP TABLE IF EXISTS `student`;
	CREATE TABLE `student` (
		`id` int(11) NOT NULL AUTO_INCREMENT,
		`name` varchar(50) CHARACTER SET latin1 NOT NULL DEFAULT '',
		`age` tinyint(4) DEFAULT '0',
		PRIMARY KEY (`id`)
	) ENGINE=MyISAM AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
	*/
	//创建表
	stmt, err := DB.Prepare(`DROP TABLE IF EXISTS student`)
	res, err := stmt.Exec()
	if(err != nil){
		panic(err)
	}
	stmt, err = DB.Prepare(`
	CREATE TABLE student (
		id int(11) NOT NULL AUTO_INCREMENT,
		name varchar(50) CHARACTER SET latin1 NOT NULL DEFAULT "",
		age tinyint(4) DEFAULT "0",
		PRIMARY KEY (id)
	) ENGINE=MyISAM AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
`)
	res, err = stmt.Exec()
	if(err != nil){
		panic(err)
	}

	//增加数据
	stmt, err = DB.Prepare(`INSERT student (name,age) values (?,?)`)
	res, err = stmt.Exec("wangwu", 26)
	id, err := res.LastInsertId()
	fmt.Println("自增id=", id)
	fmt.Printf("error %v\n", err)
	//修改数据
	stmt, err = DB.Prepare(`UPDATE student SET age=? WHERE id=?`)
	res, err = stmt.Exec(21, 5)
	num, err := res.RowsAffected() //影响行数
	fmt.Println(num)
	fmt.Printf("error %v\n", err)
	//删除数据
	stmt, err = DB.Prepare(`DELETE FROM student WHERE id=?`)
	res, err = stmt.Exec(5)
	num, err = res.RowsAffected()
	fmt.Println(num)
	fmt.Printf("error %v\n", err)
	//查询数据
	rows, err := DB.Query("SELECT * FROM student")

	//--------简单一行一行输出---start
	//    for rows.Next() { //满足条件依次下一层
	//        var id int
	//        var name string
	//        var age int
	//        rows.Columns()

	//        err = rows.Scan(&id, &name, &age)
	//        fmt.Println(id)
	//        fmt.Println(name)
	//        fmt.Println(age)
	//    }
	//--------简单一行一行输出---end

	//--------遍历放入map----start
	//构造scanArgs、values两个数组，scanArgs的每个值指向values相应值的地址
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))

	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		//将行数据保存到record字典
		err = rows.Scan(scanArgs...)
		record := make(map[string]string)
		for i, col := range values {
			if col != nil {
				record[columns[i]] = string(col.([]byte))
			}
		}
		fmt.Println(record)
	}
	//--------遍历放入map----end

	DB.Close(); //关闭数据库连接池
	println();
}