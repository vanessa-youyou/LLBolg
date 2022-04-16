package main

import (
	"LlBlog/dao"
	"LlBlog/middleware"
	"LlBlog/router"

	"github.com/gin-gonic/gin"
)

func main() {
	rou := gin.Default()
	// Use 是用来做拦截的
	rou.Use(middleware.Auth())

	// 初始化路由
	router.InitRouter(rou)
	// 初始化数据库
	dao.Init()

	rou.Run(":8080")
}
