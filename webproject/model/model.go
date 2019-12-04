package model

import (
	"fmt"
	"log"
)
var code int64
//  Insert 插入操作
func Insert(sql string,args...interface{})(int64,error){
	stmt,err := DB.Prepare(sql)
	defer stmt.Close()
	code,err = CheckErr(err,"sql语句设置失败",1)
	if(code != 1){return code,err}

	result,err := stmt.Exec(args...)
	code,err = CheckErr(err,"参数添加失败",1)
	if(code != 1){return code,err}

	id,err := result.LastInsertId()
	code,err = CheckErr(err,"插入失败",1)
	if(code != 1){return code,err}

	fmt.Printf("插入成功,id为%v\n", id)
	return id,err
}

// Delete删除
func Delete(sql string,args...interface{})(int64,error){
	stmt,err := DB.Prepare(sql)
	defer stmt.Close()
	code,err = CheckErr(err,"sql语句设置失败",1)
	if(code != 1){return code,err}

	result,err := stmt.Exec(args...)
	code,err = CheckErr(err,"参数添加失败",1)
	if(code != 1){return code,err}

	num,err := result.RowsAffected()
	code,err = CheckErr(err,"删除失败",1)
	if(code != 1){return code,err}

	fmt.Printf("删除成功，删除行数为%d\n",num)
	return num,err
}

// Update 更新
func Update(sql string,args...interface{})(int64,error){
	stmt, err := DB.Prepare(sql)
	defer stmt.Close()
	code,err = CheckErr(err, "SQL语句设置失败",1)
	if(code != 1){return code,err}

	result, err := stmt.Exec(args...)
	code,err = CheckErr(err, "参数添加失败",1)
	if(code != 1){return code,err}

	num, err := result.RowsAffected()
	code,err = CheckErr(err,"修改失败",1)
	if(code != 1){return code,err}

	fmt.Printf("修改成功，修改行数为%d\n",num)
	return num,err
}

// 检查错误
func CheckErr(err error,msg string,code int)(int64,error){
	if err != nil{
		log.Panicln(msg,err,code)
		return 0,err
	}
	return 1,nil
}