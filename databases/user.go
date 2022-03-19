package databases

import (
	"LlBlog/models"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/jinzhu/gorm"
)

// databases应该关注于单纯的通用逻辑
// 例如通过id获取记录这种
func GetUserById(id uint) (*models.UserInfo, error) {
	var u models.UserInfo
	if err := DB.Where("id = ?", id).Find(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// search user by CustomerName
func GetUserByCustomerName(cn string) (*models.UserInfo, error){
	var u models.UserInfo
	if err := DB.Where("customer_name = ?", cn).Find(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// 账号查重
func AccountRechecking(cn string) (bool, error){
	var s string
	err := DB.Where("customer_name = ?", cn).Find(&s).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 没有相同的账号名 允许创建
		return true, nil
	}
	return false, err

}

// 添加账号
func AccountInsert(u *models.UserInfo)  (bool, error){
	d := []byte(u.Password)
	m := md5.New()
	m.Write(d)
	u.Password = hex.EncodeToString(m.Sum(nil))
	err := DB.Create(&u).Error
	if err != nil{
		return false, err
	}
	return true, nil
}