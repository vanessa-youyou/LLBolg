package dto

import "LlBlog/models"

/* 登录 */
type LoginReq struct {
	CustomerName string `json:"customer_name"` // 用户名
	Password     string `json:"password"`      // 密码
}

type LoginRsp struct {
	User models.UserInfo `json:"user"` // 用户基础信息
}

/* 注册 */
type RegisterReq struct {
	CustomerName string `json:"customer_name"` // 用户名
	Password     string `json:"password"`      // 密码
}

type RegisterRsp struct {
	Id uint `json:"id"` // id
}

// UpdateReq 修改
type UpdateReq struct {
	Name         	string `json:"name"`
	Gender       	string `json:"gender"`
	Introduction 	string `json:"introduction"`
	HeadPortrait	string	`json:"head_portrait"`
}

type UpdateRsp struct {
	User models.UserInfo `json:"user"` // 用户基础信息
}

// ShowSelfPageReq 个人主页
type ShowSelfPageReq struct {
	NowPage int `json:"now_page"`
}

type Page struct {
	NowPage		int		`json:"now_page"`
	AllPage		int		`json:"all_page"`
	PageSize	int		`json:"page_size"`
}

type ShowSelfPageRsp struct {
	Page			Page						`json:"page"`
	Articles		[]models.ArticleInfo		`json:"articles"`
	SelfPart		models.UserInfo				`json:"self_part"`
}

// UpdatePasswordReq 修改密码
type UpdatePasswordReq struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type UpdatePasswordRsp struct {
	User	models.UserInfo		`json:"user"`
}

// SearchUserReq 查找用户 by name
type SearchUserReq struct {
	SearchWay	bool	`json:"search_way"`	// 0: 模糊 1：准确
	Content		string	`json:"content"`
}

type SearchUserRsp struct {
	Users 	[]models.UserInfo	`json:"users"`
}
