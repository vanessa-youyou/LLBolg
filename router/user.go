package router

import (
	"github.com/gin-gonic/gin"
)

func registerUser(router *gin.Engine) {
	userGroup := router.Group("/user")
	// 登录
	userGroup.POST("/login", handler.BaseHandler(handler.Login))
	// 注册
	userGroup.POST("/registered", handler.BaseHandler(handler.Registered))
}
