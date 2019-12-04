package model

import (
	"errors"
	"webproject/pkg/errno"
	"encoding/json"
	"log"
)

type User struct {
	UserName    string `json:"user_name"`
	Password    string `json:"password"`
}

//查询
func (user *User) SelectUserByName(name string) error {
	stmt,err := DB.Prepare("select user_name,password from user where user_name = ?")
	if err != nil{
		return err
	}
	defer stmt.Close()
	rows,err := stmt.Query(name)
	defer rows.Close()
	if err != nil{
		return err
	}
	//数据处理
	for rows.Next(){
		rows.Scan(&user.UserName,&user.Password)
	}
	if err := rows.Err();err != nil{
		return err
	}
	return nil
}

//验证字段
func (u *User) Validate() error{
	if u.UserName == "" || u.Password == ""{
		return errors.New(errno.ErrValidation.Message)
	}
	return nil
}

//创建用户
func (u *User) Create() (int64,error){
	id,err := Insert("insert into user(user_name,password) values(?,?)",u.UserName, &u.Password)
	if err != nil{
		return 0,err
	}
	return id,nil
}

//转换user 作为json字符串
func (user *User)UserToJson() string{
	jsonStr,err := json.Marshal(user)
	if err != nil{
		log.Println(err)
	}
	return string(jsonStr)
}

//转换json 作为user构造体
func (user *User)JsonToUser(jsonBlob string) error{
	err := json.Unmarshal([]byte(jsonBlob),&user)
	if err != nil{
		return err
	}
	return nil
}