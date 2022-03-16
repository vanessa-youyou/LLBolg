package Databases

import (
	"LlBlog/07/Models"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)
var DB *gorm.DB

func init(){
	var err error
	DB, err = gorm.Open("mysql", "root:12345678@(127.0.0.1:3306)/lblog?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil{
		panic(err)
	}
	if DB.Error != nil{
		fmt.Printf("database error %v", DB.Error)
	}
}

func UserLoginCheck(id uint, psw string) bool{
	// 自动迁移 感觉放这里不对..
	DB.AutoMigrate(&Models.UserInfo{})
	var u Models.UserInfo
	if err := DB.First("id = ? AND password = ?", id, psw).Find(&u).Error; err != nil{
		return false
	}
	return true
}
