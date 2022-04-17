package dto

import (
	"LlBlog/models"
)

// CreateArticleReq 新文章
type CreateArticleReq struct {
	Title       string `json:"title"`                         // 文章标题
	Text        string ` json:"text"`                         // 文章详情
	AuthorID    uint   `json:"author_id"`					  // 作者id
	Original    int8   `json:"original"`                      // 原创与否
	PlacedTop   int8   `json:"placed_top"`                    // 置顶与否
	State       int    `json:"state"`                         // 状态
}

type CreateArticleRsp struct {
	Article 	models.ArticleInfo		`json:"article"`
}

// UpdateArticleReq 修改文章
type UpdateArticleReq struct{
	ArticleId	uint	`json:"article_id"`
	Title       string `json:"title"`                         // 文章标题
	Text        string ` json:"text"`
	Original    int8   `json:"original"`                      // 原创与否
	PlacedTop   int8   `json:"placed_top"`                    // 置顶与否
	State       int    `json:"state"`
}

type UpdateArticleRsp struct {
	Message 	string		`json:"message"`
}

// DeleteArticleReq 删除文章
type DeleteArticleReq struct {
	ArticleId 	uint	`json:"article_id"`
}

type DeleteArticleRsp struct {
	Message 	string		`json:"message"`
}

// CreateCommentReq 新建评论
type CreateCommentReq struct {
	UserId 		uint	`json:"user_id"`
	ArticleId	uint	`json:"article_id"`
	Context		string	`json:"context"`
}

type CreateCommentRsp struct {
	Message 	string		`json:"message"`
}

// DeleteCommentReq 删除评论
type DeleteCommentReq struct {
	CommentId	uint	`json:"comment_id"`
	ArticleId	uint	`json:"article_id"`
}

type DeleteCommentRsp struct {
	Message 	string		`json:"message"`
}

// PickArticleReq 点赞
type PickArticleReq struct {
	ArticleId	uint	`json:"article_id"`
}
type PickArticleRsp struct {
	Message 	string		`json:"message"`
}

// CancelPickArticleReq 取消点赞
type CancelPickArticleReq struct {
	ArticleId	uint	`json:"article_id"`
}
type CancelPickArticleRsp struct {
	Message 	string		`json:"message"`
}

// PickCommentReq 点赞评论
type PickCommentReq struct {
	CommentId 	uint	`json:"comment_id"`
	ArticleId	uint	`json:"article_id"`
}
type PickCommentRsp struct {
	Message 	string		`json:"message"`
}

// CancelPickCommentReq 取消点赞
type CancelPickCommentReq struct {
	CommentId 	uint	`json:"comment_id"`
	ArticleId	uint	`json:"article_id"`
}
type CancelPickCommentRsp struct {
	Message 	string		`json:"message"`
}

// SearchArticleReq 查询文章
type SearchArticleReq struct {
	SearchWay	bool	`json:"search_way"`	// 0: 模糊 1：准确
	Check		string	`json:"check"`		// title：搜标题 text:搜内容
	Content		string	`json:"content"`	// 搜索框的内容
}

type SearchArticleRsq struct {
	Articles	[]models.ArticleInfo	`json:"articles"`
	Message 	string					`json:"message"`
}

// ArticleDetailsReq 文章详情页面
type ArticleDetailsReq struct {
	ArticleId 	uint	`json:"article_id"`
}

type ArticleDetailsRsp struct {
	Article 	models.ArticleInfo		`json:"article"`
	Comment     []models.CommentInfo	`json:"comment"`
}

// CollectionArticleRep 文章收藏
type CollectionArticleRep struct {
	ArticleId 	uint	`json:"article_id"`
}

type CollectionArticleRsq struct {
	Message 	string		`json:"message"`
}

// CancelCollectionArticleRep 取消收藏
type CancelCollectionArticleRep struct {
	ArticleId	uint	`json:"article_id"`
}

type CancelCollectionArticleRsq struct {
	Message 	string		`json:"message"`
}

// CreateLabelRep 新建标签
type CreateLabelRep struct {
	Name		string	`json:"name"`
}

type CreateLabelRsq struct {
	Message 	string		`json:"message"`
}

// SearchLabelReq 搜索标签
type SearchLabelReq struct {
	Name		string	`json:"name"`
}

type SearchLabelRsp struct {
	Labels 		[]models.Label	`json:"labels"`
	Message 	string			`json:"message"`
}

// LabelReceiveReq 联系（文章 标签）
type LabelReceiveReq struct {
	LabelId		[]uint	`json:"label_id"`
	ArticleId	uint 	`json:"article_id"`
}

type LabelReceiveRsp struct {
	Message 	string			`json:"message"`
}

// ShowPageReq 首页
type ShowPageReq struct {
	NowPage int `json:"now_page"`
}

type ShowPageRsp struct {
	Page			Page						`json:"page"`
	Articles		[]models.ArticleInfo		`json:"articles"`
	SelfPart		models.UserInfo				`json:"self_part"`
}
