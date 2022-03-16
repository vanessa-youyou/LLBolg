package Controllers

import (
	"LlBlog/07/Services"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"

)

// UserLogin 登录检验
func UserLogin(c *gin.Context)  {
	// 数据库那边的操作(要接收的啊承诺书的结构体)
	var userG Services.LoginUser
	err := c.ShouldBind(&userG)
	if err != nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"error":err.Error(),
		})
	}else{
		fmt.Printf("%#v\n",userG)
		// 进行一个查找的动作

		if userG.LoginCheck() {
			// 如果true
			c.JSON(http.StatusBadRequest,gin.H{
				"message":"登陆成功 这里应该跳转页面到 个人页面",
			})
		}else{
			c.JSON(http.StatusBadRequest,gin.H{
				"message":"The account or password is incorrect",
			})
		}
	}
}