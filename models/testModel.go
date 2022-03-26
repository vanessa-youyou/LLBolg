package models

type UserInfo struct {
	ID           uint   `json:"id"`
	Name         string `json:"name" gorm:"not null"`
	Gender       string `json:"gender" gorm:"not null"`
	Password     string `json:"password" gorm:"not null"`
	CustomerName string `json:"customer_name" gorm:"not null"`
	Introduction string `json:"introduction"`
	Label        string `json:"label"`
}
type ArticleInfo struct {
	ID          	uint	`json:"id"`
	Title			string	`json:"title" gorm:"column:title;not null"`
	Text 			string 	` json:"text" gorm:"column:text;not null"`
	AuthorID    	uint	`json:"author_id" gorm:"column:author_id;type:int unsigned;not null"`
	Original		int8	`json:"original" gorm:"column:original;not null"`
	PlacedTop		int8	`json:"placed_top" gorm:"column:placed_top;not null"`
	State			string	`json:"state" gorm:"column:state;not null"`

}

type CommentInfo struct {
	ID				uint	`json:"id" gorm:"column:id;not null"`
	UserID    		uint	`json:"user_id" gorm:"column:user_id;type:int unsigned;not null"`
	ArticleID		uint	`json:"article_id" gorm:"column:article_id;type:int unsigned;not null"`
	Content			string	`json:"content" gorm:"column:content"`
	Like			int		`json:"like" gorm:"column:like"`
}

// GiveLike 点赞关系表
type GiveLike struct {
	ID				uint	`json:"id" gorm:"column:id"`
	UserID    		uint	`json:"user_id" gorm:"column:user_id;type:int unsigned;not null"`
	ArticleID		uint	`json:"article_id" gorm:"column:article_id;type:int unsigned;not null"`
}

func (u UserInfo) Clear() UserInfo {
	u.Password = ""
	return u
}


