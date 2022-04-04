package controllers

import (
	"LlBlog/core"
	"LlBlog/errors"
	"LlBlog/utils"

	errors2 "errors"

	"github.com/gin-gonic/gin"
)

// Success 返回成功
func Success(c *gin.Context) {
	auth := c.MustGet("auth").(core.AuthAuthorization)

	if auth.IsLogin() {
		utils.Return(c, gin.H{
			"user": auth.User.Clear(),
		})
		return
	}

	utils.Return(c, gin.H{
		"message": "登陆成功 这里应该跳转页面到 个人页面",
	})
}

// 自定义失败
func Err(c *gin.Context) {
	utils.Return(c, errors.LoginFailed)
}

// 公共error
func Err2(c *gin.Context) {
	utils.Return(c, errors2.New("gggggg"))
}
