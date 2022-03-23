package databases

import (
	"LlBlog/models"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

var ip = os.Getenv("MYSQL_IP")
var password = os.Getenv("MYSQL_PASSWORD")

func Init() {

	if ip == "" {
		ip = "127.0.0.1"
	}
	if password == "" {
		password = "12345678"
	}

	var err error
	DB, err = gorm.Open(
		"mysql",
		fmt.Sprintf(
			"root:%s@(%s:3306)/lblog?charset=utf8mb4&parseTime=True&loc=Local",
			password, ip,
		),
	)
	if err != nil {
		panic(err)
	}

	// 应该放在这里
	DB.AutoMigrate(&models.UserInfo{})
}
