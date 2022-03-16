package services

import "LlBlog/databases"

type LoginUser struct {
	ID       uint   `json:"id"`
	Password string `json:"password"`
}

func (l LoginUser) LoginCheck() bool {
	// 进行查找 不应该用id作为账户 然后这里如果find到 应该返回id 没find到 则返回nil or 负数

	user, err := databases.GetUserById(l.ID)
	if err != nil {
		// TODO 日志
		return false
	}

	// 这种判断密码是否相同的逻辑应该放在
	// services里面，不能放在databases，因为不通用
	// TODO 是不是应该加密
	return user.Password == l.Password
}
