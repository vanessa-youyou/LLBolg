package controllers

import (
	"LlBlog/core"
	"LlBlog/databases"
	"LlBlog/errors"
	"LlBlog/models"
	"LlBlog/services"
	"LlBlog/utils"

	"github.com/gin-gonic/gin"
)

// UserLogin 登录检验
func UserLogin(c *gin.Context) {
	auth := c.MustGet("auth").(core.AuthAuthorization)
	// 数据库那边的操作(要接收的啊承诺书的结构体)
	var userG services.LoginUser

	// 接收数据
	err := c.ShouldBind(&userG)
	if err != nil {
		utils.Return(c, err)
		return
	}

	user, err := databases.GetUserByCustomerName(userG.CustomerName)
	if err != nil {
		// 其实这里跟上面哪个一样 不应该直接把系统错误显示给客户端，不过开发就随便啦
		utils.Return(c, err)
		return
	}

	// 登陆失败报错
	if !userG.LoginCheck(user) {
		utils.Return(c, errors.LoginFailed)
		return
	}

	// 登陆成功
	// 设置cookie
	auth.SetCookie(c, user.ID)
	utils.Return(c, gin.H{
		"user": user.Clear(),
	})
}

func UserRegistered(c *gin.Context) {
	// add
	var userN models.UserInfo
	// 接收数据
	err := c.ShouldBind(&userN)
	if err != nil {
		utils.Return(c, err)
		return
	}

	// 进行一个查找的动作 看看 账户名字有没有重复
	if !services.RegisteredNameCheck(userN) {
		utils.Return(c, errors.WrongAccountName)
		return
	}
	// 开始add账号进数据库
	if !services.AddAccount(userN) {
		utils.Return(c, errors.RegisteredFailed)
		return
	}

	// 注册成功
	utils.Return(c, gin.H{
		"message": "注册成功 这里应该跳转页面到 个人页面",
	})
}
