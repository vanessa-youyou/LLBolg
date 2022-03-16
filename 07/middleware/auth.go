package middleware

import "github.com/gin-gonic/gin"

// 登录态获取
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {

		// TODO 尝试获取session，如果有，获取当前登录用户信息
		// 并且写入上下文，这样后面的函数可以直接使用

		// c.Set("account", xxxx)
		// c.Set("is_login", true)

		c.Next()
	}
}
