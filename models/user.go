package models

// 用户信息表
type UserInfo struct {
	ID           uint   `json:"id"`                                   // id
	Name         string `json:"name" gorm:"not null"`                 // 昵称
	Gender       string `json:"gender" gorm:"not null"`               // 性别
	Password     string `json:"password" gorm:"not null"`             // 密码
	CustomerName string `json:"customer_name" gorm:"not null;unique"` // 用户名
	Introduction string `json:"introduction"`                         // 个人介绍
	Label        string `json:"label"`                                // 标签
	HeadPortrait string `json:"head_portrait" gorm:"head_portrait"`   // 头像
	CreateTime   int64  `json:"create_time"`                          // 创建时间
	UpdateTime   int64  `json:"update_time"`                          // 更新时间
}

func (u UserInfo) Clear() UserInfo {
	u.Password = ""
	return u
}
