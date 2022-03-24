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
	ID          uint	`json:"id"`
	Title		string	`json:"title" gorm:"title;not null"`
	Text 		string 	` json:"text" gorm:"text;not null"`
	UserInfoID    uint	`json:"author_id" gorm:"user_info_id;type:int unsigned;not null"`
	Original	int8	`json:"original" gorm:"original;not null"`
	PlacedTop	int8	`json:"placed_top" gorm:"placed_top;not null"`
	Release	int8	`json:"release" gorm:"release;not null"`
	Praise	int 	`json:"praise" gorm:"praise"`
}

func (u UserInfo) Clear() UserInfo {
	u.Password = ""
	return u
}


