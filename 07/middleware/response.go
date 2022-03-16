package middleware

import (
	"LlBlog/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 登录态获取
func Response() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()
		// 存在异常的话
		if code, ok := c.Get("err"); ok {
			switch err := code.(type) {
			case errors.ErrorBase:
				c.JSON(http.StatusOK, gin.H{
					"code": err.Code,
					"msg":  err.Error(),
				})
			case error:
				// TODO 日志
				c.JSON(http.StatusOK, gin.H{
					"code": -1,
					"msg":  err.Error(),
				})
			}
		}
	}
}
