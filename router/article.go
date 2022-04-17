package router

import (
	"LlBlog/handler"
	"github.com/gin-gonic/gin"
)

func registerArticle(router *gin.Engine)  {
	articleGroup := router.Group("/article")
	// 写文章
	articleGroup.POST("/newArticles", handler.BaseHandler(handler.CreatArticle))
	// 修改文章
	articleGroup.PUT("/updateArticle", handler.BaseHandler(handler.ModifyArticle))
	// 删除文章
	articleGroup.DELETE("/deleteArticles", handler.BaseHandler(handler.DeleteArticle))
	// 新建评论
	articleGroup.POST("/newComment", handler.BaseHandler(handler.CreateComment))
	// 删除评论
	articleGroup.DELETE("/deleteComment", handler.BaseHandler(handler.DeleteComment))
	// 点赞文章
	articleGroup.POST("/pickArticle", handler.BaseHandler(handler.PickArticle))
	// 取消点赞文章
	articleGroup.DELETE("/cancelPickArticle", handler.BaseHandler(handler.CancelPickArticle))
	// 点赞评论
	articleGroup.POST("/pickComment", handler.BaseHandler(handler.CommentPick))
	// 取消点赞评论
	articleGroup.DELETE("/cancelCommentPick", handler.BaseHandler(handler.CancelCommentPick))
	// 查询文章
	articleGroup.POST("/searchArticles", handler.BaseHandler(handler.SearchArticles))
	// 文章详情
	articleGroup.POST("/articleDetails", handler.BaseHandler(handler.ArticleDetails))
	// 收藏文章
	articleGroup.POST("/collectionArticle", handler.BaseHandler(handler.CollectionArticle))
	// 取消收藏
	articleGroup.DELETE("/cancelCollectionArticle", handler.BaseHandler(handler.CancelCollectionArticle))
	// 添加标签
	articleGroup.POST("/creatLabel", handler.BaseHandler(handler.CreatLabel))
	// 搜索标签
	articleGroup.POST("/searchLabel", handler.BaseHandler(handler.SearchLabel))
	// 建立联系（标签和文章）
	articleGroup.POST("/chooseLabels", handler.BaseHandler(handler.ChooseLabels))
	// 首页
	articleGroup.POST("/homePage", handler.BaseHandler(handler.HomePage))
}
