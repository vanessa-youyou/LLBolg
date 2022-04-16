package models

// 文章信息表
type ArticleInfo struct {
	ID          uint   `json:"id"`                                                           // id
	Title       string `json:"title" gorm:"column:title;not null"`                           // 文章标题
	Text        string ` json:"text" gorm:"column:text;not null"`                            // 文章详情
	AuthorID    uint   `json:"author_id" gorm:"column:author_id;type:int unsigned;not null"` // 作者id
	Original    int8   `json:"original" gorm:"column:original;not null"`                     // 原创与否
	PlacedTop   int8   `json:"placed_top" gorm:"column:placed_top;not null"`                 // 置顶与否
	State       int    `json:"state" gorm:"column:state;not null"`                           // 状态
	LikeNum     int    `json:"like_num" gorm:"like_num"`                                     // 点赞
	CommentsNum int    `json:"comments_num" gorm:"commentsNum"`                              // 评论数量
}

// CommentInfo 评论的表
type CommentInfo struct {
	ID        uint   `json:"id" gorm:"column:id;not null"`                                   // id
	UserID    uint   `json:"user_id" gorm:"column:user_id;type:int unsigned;not null"`       // 评论用户id
	ArticleID uint   `json:"article_id" gorm:"column:article_id;type:int unsigned;not null"` // 文章id
	Content   string `json:"content" gorm:"column:content"`                                  // 评论正文
	LikeNum   int    `json:"like_num" gorm:"like_num"`                                       // 点赞数量
}

// GiveLike 文章点赞关系表
type GiveLike struct {
	ID        uint `json:"id" gorm:"column:id"`
	UserID    uint `json:"user_id" gorm:"column:user_id;type:int unsigned;not null"`
	ArticleID uint `json:"article_id" gorm:"column:article_id;type:int unsigned;not null"`
}

// CommentLike 评论的点赞
type CommentLike struct {
	ID        uint `json:"id" gorm:"column:id"`
	UserID    uint `json:"user_id" gorm:"column:user_id;type:int unsigned;not null"`
	CommentID uint `json:"comment_id" gorm:"comment_id;type:int unsigned;not null"`
	ArticleID uint `json:"article_id" gorm:"column:article_id;type:int unsigned;not null"`
}
