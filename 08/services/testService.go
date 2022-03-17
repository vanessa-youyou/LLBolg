package services

import (
	"LlBlog/databases"
	"LlBlog/models"
	"crypto/md5"
	"encoding/hex"
)


type LoginUser struct {
	CustomerName       string   `json:"customer_name" form:"customer_name" `
	Password string `json:"password" form:"password" `
}


func (l LoginUser) LoginCheck() bool {
	// 进行查找 不应该用id作为账户 然后这里如果find到 应该返回id 没find到 则返回nil or 负数

	user, err := databases.GetUserByCustomerName(l.CustomerName)
	if err != nil {
		// TODO 日志
		return false
	}

	// 这种判断密码是否相同的逻辑应该放在
	// services里面，不能放在databases，因为不通用
	// TODO 是不是应该加密
	// md5 encryption
		d := []byte(l.Password)
		m := md5.New()
		m.Write(d)
	return user.Password == hex.EncodeToString(m.Sum(nil))
}

func RegisteredNameCheck(u models.UserInfo) bool {
	cn, err := databases.AccountRechecking(u.CustomerName)
	if err != nil {
		// TODO 日志
		return false
	}
	if !cn{
		// 没有相同的账号名 允许创建
		// TODO 日志
		return false
	}
	return true
}

func AddAccount(u models.UserInfo) bool{
	t, err := databases.AccountInsert(&u)
	if err != nil{
		return false
	}
	if !t {
		return false
	}
	return true
}