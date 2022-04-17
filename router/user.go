package router

import (
	"LlBlog/handler"
	"github.com/gin-gonic/gin"
)

func registerUser(router *gin.Engine) {
	userGroup := router.Group("/user")
	// 登录
	userGroup.POST("/login", handler.BaseHandler(handler.Login))
	// 注册
	userGroup.POST("/registered", handler.BaseHandler(handler.Register))
	// 修改（头像也可以在这里修改）
	userGroup.PUT("/updateInformation", handler.BaseHandler(handler.UserInformationUpdate))
	// 个人主页
	userGroup.GET("/selfPage", handler.BaseHandler(handler.ShowSelf))
	// 修改密码
	userGroup.PUT("/updatePassword", handler.BaseHandler(handler.PasswordUpdate))
	// 用户查找
	userGroup.POST("/searchUser", handler.BaseHandler(handler.SearchUSer))

}
