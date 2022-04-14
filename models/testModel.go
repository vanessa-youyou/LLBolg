package models

type UserInfo struct {
	ID           uint   `json:"id"`
	Name         string `json:"name" gorm:"not null"`
	Gender       string `json:"gender" gorm:"not null"`
	Password     string `json:"password" gorm:"not null"`
	CustomerName string `json:"customer_name" gorm:"not null"`
	Introduction string `json:"introduction"`
	Label        string `json:"label"`
	HeadPortrait	string	`json:"head_portrait" gorm:"head_portrait"`
}

type ArticleInfo struct {
	ID          	uint	`json:"id"`
	Title			string	`json:"title" gorm:"column:title;not null"`
	Text 			string 	` json:"text" gorm:"column:text;not null"`
	AuthorID    	uint	`json:"author_id" gorm:"column:author_id;type:int unsigned;not null"`
	Original		int8	`json:"original" gorm:"column:original;not null"`
	PlacedTop		int8	`json:"placed_top" gorm:"column:placed_top;not null"`
	State			int	    `json:"state" gorm:"column:state;not null"`
	LikeNum			int		`json:"like_num" gorm:"like_num"`
	CommentsNum		int		`json:"comments_num" gorm:"commentsNum"`
}

// CommentInfo 评论的表
type CommentInfo struct {
	ID				uint	`json:"id" gorm:"column:id;not null"`
	UserID    		uint	`json:"user_id" gorm:"column:user_id;type:int unsigned;not null"`
	ArticleID		uint	`json:"article_id" gorm:"column:article_id;type:int unsigned;not null"`
	Content			string	`json:"content" gorm:"column:content"`
	LikeNum			int		`json:"like_num" gorm:"like_num"`
}

// GiveLike 文章点赞关系表
type GiveLike struct {
	ID				uint	`json:"id" gorm:"column:id"`
	UserID    		uint	`json:"user_id" gorm:"column:user_id;type:int unsigned;not null"`
	ArticleID		uint	`json:"article_id" gorm:"column:article_id;type:int unsigned;not null"`
}

// CommentLike 评论的点赞
type CommentLike struct {
	ID				uint	`json:"id" gorm:"column:id"`
	UserID    		uint	`json:"user_id" gorm:"column:user_id;type:int unsigned;not null"`
	CommentID		uint	`json:"comment_id" gorm:"comment_id;type:int unsigned;not null"`
	ArticleID		uint	`json:"article_id" gorm:"column:article_id;type:int unsigned;not null"`
}

// Search 查询
type Search struct {
	SearchWay	bool	`json:"search_way"`	// 0: 模糊 1：准确
	Check		string	`json:"check"`
	Content		string	`json:"content"`
}
// Page 分页
type Page struct {
	PageNow		int 	`json:"now_page"`
	PageNum		int		`json:"page_num"`
	PageSize	int		`json:"page_size"`
}

// Collection 收藏
type Collection struct {
	ArticleID		uint	`json:"article_id" gorm:"column:article_id;type:int unsigned;not null"`
	UserID    		uint	`json:"user_id" gorm:"column:user_id;type:int unsigned;not null"`
}

// Label 标签
type Label struct {
	ID			uint	`json:"id" gorm:"column:id"`
	Name		string	`json:"name" gorm:"column:name"`
}

// LabelRelation 标签-文章关系表
type LabelRelation struct {
	ID			uint	`json:"id" gorm:"id"`
	LabelId		uint	`json:"label_id" gorm:"label_id"`
	ArticleId	uint	`json:"article_id" gorm:"article_id"`
}
// LabelReceive 只用于接收
type LabelReceive struct {
	LabelId		[]uint	`json:"label_id"`
	ArticleId	uint 	`json:"article_id"`
}

func (u UserInfo) Clear() UserInfo {
	u.Password = ""
	return u
}
