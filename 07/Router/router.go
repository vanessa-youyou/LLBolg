package Router

import (
	"LlBlog/07/Controllers"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	router := gin.Default()
	userGroup := router.Group("/user")
	{
		userGroup.POST("/login", Controllers.UserLogin)
	}
	router.Run(":8080")
}

