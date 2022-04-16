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
