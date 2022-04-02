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
		userGroup.PUT("/updatePassword", controllers.PasswordUpdate)
		userGroup.POST("/upload", controllers.Upload)

		// 文章
		userGroup.POST("/newArticles", controllers.CreatArticle)
		userGroup.POST("/PickArticle",controllers.PickArticle)	// redis
		userGroup.PUT("/updateArticle", controllers.ModifyArticle)
		userGroup.DELETE("/deleteArticles", controllers.DeleteArticle)
		// 评论
		userGroup.POST("/newComment", controllers.WriteComment)
		userGroup.POST("commentPick", controllers.CommentPick)	// redis
		userGroup.DELETE("/deleteComment", controllers.DeleteComment)

		// 前段页面展示部分
		userGroup.GET("/selfPage", controllers.ShowSelf)

		// 查找部分的接口
		userGroup.POST("/searchArticles", controllers.SearchArticles)	// 查找文章
		userGroup.POST("/searchUser", controllers.SearchUSer)		// 查找用户信息
		}
	testGroup := router.Group("/test")
	{
		testGroup.GET("/succ", controllers.Success)
		testGroup.GET("/err", controllers.Err)
		testGroup.GET("/err2", controllers.Err2)
	}

}
