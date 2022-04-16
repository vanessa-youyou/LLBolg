package handler

import (
	"LlBlog/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 基础函数
func BaseHandler(fun func(*gin.Context) (interface{}, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		rsp, erro := fun(c)
		if erro != nil {
			switch err := erro.(type) {
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
		} else {
			c.JSON(http.StatusOK, rsp)
		}
	}
}
