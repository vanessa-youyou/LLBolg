package services

import (
	"LlBlog/databases"
	"LlBlog/models"
	"crypto/md5"
	"encoding/hex"
)

type LoginUser struct {
	CustomerName string `json:"customer_name" form:"customer_name" `
	Password     string `json:"password" form:"password" `
}

func (l LoginUser) LoginCheck(user *models.UserInfo) bool {
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
	if !cn {
		// 没有相同的账号名 允许创建
		// TODO 日志
		return false
	}
	return true
}

func AddAccount(u models.UserInfo) bool {
	t, err := databases.AccountInsert(&u)
	if err != nil {
		return false
	}
	if !t {
		return false
	}
	return true
}

// UpdateUserInformation 进行更新
func UpdateUserInformation(u models.UserInfo) bool {
	t,err :=databases.UserInformationUpdate(&u)
	if err != nil {
		return false
	}
	return t
}
