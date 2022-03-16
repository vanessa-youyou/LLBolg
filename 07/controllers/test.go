package controllers

import (
	"LlBlog/errors"
	"LlBlog/utils"

	errors2 "errors"

	"github.com/gin-gonic/gin"
)

// 返回成功
func Success(c *gin.Context) {
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
