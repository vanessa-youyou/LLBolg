package databases

import (
	"LlBlog/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func Init() {
	var err error
	DB, err = gorm.Open("mysql", "root:balabalamiaomiaomiao@(120.78.167.231:3306)/lblog?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}

	// 应该放在这里
	DB.AutoMigrate(&models.UserInfo{})
}
