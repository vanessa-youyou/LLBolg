package utils

import (
	"LlBlog/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Return(c *gin.Context, v interface{}) {
	switch v.(type) {
	// 异常消息则写入上下文
	case error, errors.ErrorBase:
		c.Set("err", v)
	default:
		// 否则输出
		c.JSON(http.StatusOK, v)
	}
}
