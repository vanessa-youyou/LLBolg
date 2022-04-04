package services

import (
	"LlBlog/databases"
	"LlBlog/models"
	"context"
	"crypto/md5"
	"encoding/hex"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"mime/multipart"
)

type LoginUser struct {
	CustomerName string `json:"customer_name" form:"customer_name" `
	Password     string `json:"password" form:"password" `
}

type UpdatePassword struct {
	OriginalPassword string	`json:"original_password" form:"original_password"`
	NewPassword string `json:"new_password" form:"new_password"`
	AgainPassword string `json:"again_password" form:"again_password"`
}

// LoginCheck 登录时候的密码校验
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

// RegisteredNameCheck 注册的检查
func RegisteredNameCheck(u models.UserInfo) bool {
	cn, err := databases.AccountRechecking(u)
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

// AddAccount 增加用户
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

// UpdateUserInformation 用户信息进行更新
func UpdateUserInformation(u *models.UserInfo) bool {
	t,err :=databases.UserInformationUpdate(u)
	if err != nil {
		return false
	}
	return t
}

// SearchUser 查找用户
func SearchUser(search models.Search) ([]models.UserInfo, bool) {

	t,err,users := databases.SearchUser(search)
	if err != nil || !t{
		return nil, false
	}
	return users, true
}

// UploadFile 上传
func UploadFile(file multipart.File, fileSize int64) (string, error) {
	var Bucket = databases.Bucket
	var AccessKey = databases.AccessKey
	var ImgUrl = databases.QiniuServer
	var SecretKey =databases.SecretKey
	putPolicy := storage.PutPolicy{
		Scope: Bucket,
	}
	mac := qbox.NewMac(AccessKey,SecretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{
		Zone: &storage.ZoneHuanan,
		UseHTTPS: false,
	}

	putExtra := storage.PutExtra{}

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	err := formUploader.PutWithoutKey(context.Background(),&ret,upToken,file,fileSize,&putExtra)

	if err != nil{
		return "",err
	}
	url := ImgUrl + ret.Key
	return url, nil
}

// UpdateHeadPortrait 上传/更新头像
func UpdateHeadPortrait(url string, user *models.UserInfo) bool {
	// 对此人进行update 头像
	t,err  := databases.UpdateHeadPortrait(url, user)
	if err != nil {
		return false
	}
	return t
}

// CheckPassword 旧密码的校验：
func (u UpdatePassword) CheckPassword(user *models.UserInfo) bool {
	d := []byte(u.OriginalPassword)
	m := md5.New()
	m.Write(d)
	u.OriginalPassword = hex.EncodeToString(m.Sum(nil))
	// 取出密码 （根据id取出信息）
	var err error
	user,err = databases.GetUserById(user.ID)
	if err != nil{
		return false
	}

	if u.OriginalPassword != user.Password{
		return false
	}
	// 2：验证两个新密码是否相等
	if u.NewPassword != u.AgainPassword{
		return false
	}else{
		// 存储新的密码
		d := []byte(u.NewPassword)
		m := md5.New()
		m.Write(d)
		u.NewPassword = hex.EncodeToString(m.Sum(nil))
		user.Password = u.NewPassword
		t := UpdateUserInformation(user)
		if !t{
			return false
		}
		return true
	}
}