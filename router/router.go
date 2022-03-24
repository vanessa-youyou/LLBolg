package router

import (
	"LlBlog/controllers"

	"github.com/gin-gonic/gin"
)

func InitRouter(router *gin.Engine) {
	userGroup := router.Group("/user")
	{
		userGroup.POST("/login", controllers.UserLogin)
		userGroup.POST("/registered", controllers.UserRegistered)
		userGroup.PUT("/updateInformation", controllers.UserInformationUpdate)

		// 文章
		userGroup.POST("/newArticles", controllers.CreatArticle)
		userGroup.PUT("/giveLike", controllers.GiveLike)
		userGroup.PUT("/updateArticle", controllers.ModifyArticle)
		userGroup.DELETE("/deleteArticles", controllers.DeleteArticle)
	}
	testGroup := router.Group("/test")
	{
		testGroup.GET("/succ", controllers.Success)
		testGroup.GET("/err", controllers.Err)
		testGroup.GET("/err2", controllers.Err2)
	}
}
