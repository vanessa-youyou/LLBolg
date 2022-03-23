package middleware

import (
	"LlBlog/core"

	"github.com/gin-gonic/gin"
)

// 登录态获取
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {

		auth := core.AuthAuthorization{}
		auth.LoadAuthenticationInfo(c)
		// 设置基础信息
		c.Set("auth", auth)

		c.Next()
	}
}
