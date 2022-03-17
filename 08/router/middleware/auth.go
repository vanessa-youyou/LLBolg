package middleware

import (
	"LlBlog/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 登录态获取
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {

		// TODO 尝试获取session，如果有，获取当前登录用户信息
		// 并且写入上下文，这样后面的函数可以直接使用
		var userG services.LoginUser
		if err := c.ShouldBind(&userG); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{
				"error":err.Error(),
			})
		}
		c.Set("customer_name", userG.CustomerName)
		c.Set("is_login", true)
		// c.Set("account", xxxx)
		// c.Set("is_login", true)

		c.Next()
	}
}
