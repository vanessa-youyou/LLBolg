package dao

import (
	"LlBlog/models"
	"fmt"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

var ip = os.Getenv("MYSQL_IP")
var password = os.Getenv("MYSQL_PASSWORD")

func mysqlInit() {

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

	DB.Callback().Create().Before("gorm:create").Register("update_time_in_create", updateTimeForCreateCallback)
	DB.Callback().Update().Before("gorm:update").Register("update_time_in_update", updateTimeForUpdateCallback)

	// 应该放在这里
	DB.AutoMigrate( &models.UserInfo{},
					&models.ArticleInfo{},
					&models.CommentInfo{},
					&models.GiveLike{},
					&models.CommentLike{},
					&models.Collection{})
					//&models.Label{},
					//&models.LabelRelation{})
}

// 更新创建时间
func updateTimeForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		if scope.HasColumn("CreateTime") {
			scope.SetColumn("CreateTime", time.Now().Unix())
		}
		updateTimeForUpdateCallback(scope)
	}
}

// 更新更新时间
func updateTimeForUpdateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		// 更新时间
		if scope.HasColumn("UpdateTime") {
			scope.SetColumn("UpdateTime", time.Now().Unix())
		}
	}
}
