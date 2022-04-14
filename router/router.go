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
		userGroup.POST("/commentPick", controllers.CommentPick)	// redis
		userGroup.DELETE("/deleteComment", controllers.DeleteComment)

		// 前段页面展示部分
		userGroup.GET("/selfPage", controllers.ShowSelf)

		// 查找部分的接口
		userGroup.POST("/searchArticles", controllers.SearchArticles)	// 查找文章
		userGroup.POST("/searchUser", controllers.SearchUSer)		// 查找用户信息

		// 展示文章详情
		userGroup.POST("/articleDetails", controllers.ArticleDetails)
		// 收藏
		userGroup.POST("/collectionArticle", controllers.CollectionArticle)
		userGroup.DELETE("/cancelCollectionArticle", controllers.CancelCollectionArticle)

		// 标签
		userGroup.POST("/creatLabel", controllers.CreatLabel)	// 新建标签
		// 搜索标签
		userGroup.POST("/searchLabel", controllers.SearchLabel)
		// 为文章选择标签（传入的是 labels的列表）’
		userGroup.POST("/chooseLabels", controllers.ChooseLabels)

		// 首页
		userGroup.POST("/homePage", controllers.HomePage)
		}
	testGroup := router.Group("/test")
	{
		testGroup.GET("/succ", controllers.Success)
		testGroup.GET("/err", controllers.Err)
		testGroup.GET("/err2", controllers.Err2)
	}

}
