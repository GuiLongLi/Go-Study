package main

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type dbstruct struct {
	Db *sql.DB
}

//创建数据库
func (dbs *dbstruct) CreateDb(){
	//用户信息表
	stmt, err := dbs.Db.Prepare(`DROP TABLE IF EXISTS userinfo`)
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)

	stmt, err = dbs.Db.Prepare(`
CREATE TABLE userinfo (
	uid INT(10) NOT NULL AUTO_INCREMENT,
	username VARCHAR(64) NULL DEFAULT NULL,
	departname VARCHAR(64) NULL DEFAULT NULL,
	created DATE NULL DEFAULT NULL,
	PRIMARY KEY (uid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
`)
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)

	//用户详情表
	stmt, err = dbs.Db.Prepare(`DROP TABLE IF EXISTS userdetail`)
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)

	stmt, err = dbs.Db.Prepare(`
CREATE TABLE userdetail (
	uid INT(10) NOT NULL DEFAULT '0',
	intro TEXT NULL,
	profile TEXT NULL,
	PRIMARY KEY (uid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
`)
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	//--------------------------------
	//数据库连接
	db,err := sql.Open("mysql","root:Li880513@tcp(127.0.0.1:3306)/test?charset=utf8")
	checkErr(err)

	var dbconnect = dbstruct{
		Db: db,
	}

	//--------------------------------
	//创建数据库
	dbconnect.CreateDb()
	fmt.Println("创建数据库成功")

	//--------------------------------
	//插入数据
	stmt,err := db.Prepare("INSERT userinfo SET username=?,departname=?,created=?")
	checkErr(err)
	res,err := stmt.Exec("testinsert","研发部门","2019-09-09")
	checkErr(err)

	id,err := res.LastInsertId()
	checkErr(err)
	fmt.Println("testinsert 的id是：",id)

	//--------------------------------
	//更新数据
	stmt,err = db.Prepare("update userinfo set username=? where uid=?")
	checkErr(err)
	res,err = stmt.Exec("testupdate",id)
	checkErr(err)

	affectrow,err := res.RowsAffected()
	checkErr(err)

	fmt.Println("testupdate 影响的行数是：",affectrow)

	//--------------------------------
	//查询数据
	rows,err := db.Query("select * from userinfo")
	checkErr(err)
	for rows.Next(){
		var uid int
		var username string
		var department string
		var created string
		err = rows.Scan(&uid,&username,&department,&created)
		checkErr(err)
		fmt.Println("uid：",uid)
		fmt.Println("username：",username)
		fmt.Println("department：",department)
		fmt.Println("created：",created)
	}

	//--------------------------------
	//删除数据
	stmt,err = db.Prepare("delete from userinfo where uid=?")
	checkErr(err)

	res,err = stmt.Exec(id)
	checkErr(err)

	affectrow,err = res.RowsAffected()
	checkErr(err)

	fmt.Println("delete 影响的行数是：",affectrow)

	db.Close()
}