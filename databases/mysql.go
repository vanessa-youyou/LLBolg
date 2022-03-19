package databases

import (
	"LlBlog/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func Init() {
	var err error
	DB, err = gorm.Open("mysql", "root:12345678@(127.0.0.1:3306)/lblog?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	// 应该放在这里
	DB.AutoMigrate(&models.UserInfo{})
}
