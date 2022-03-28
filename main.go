
package main

import (
	"LlBlog/databases"
	"LlBlog/middleware"
	"LlBlog/router"

	"github.com/gin-gonic/gin"
)

func main() {
	rou := gin.Default()
	// Use 是用来做拦截的
	rou.Use(middleware.Auth(), middleware.Response())

	// 初始化路由
	router.InitRouter(rou)
	// 初始化数据库
	databases.Init()
	databases.InitRedis()

	rou.Run(":8080")
}
