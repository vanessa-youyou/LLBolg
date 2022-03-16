package router

import (
	"LlBlog/controllers"

	"github.com/gin-gonic/gin"
)

func InitRouter(router *gin.Engine) {
	userGroup := router.Group("/user")
	{
		userGroup.POST("/login", controllers.UserLogin)
	}
	testGroup := router.Group("/test")
	{
		testGroup.GET("/succ", controllers.Success)
		testGroup.GET("/err", controllers.Err)
		testGroup.GET("/err2", controllers.Err2)
	}
}
