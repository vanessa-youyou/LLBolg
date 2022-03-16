package Services

import "LlBlog/07/Databases"

type LoginUser struct {
	ID uint	`json:"id"`
	Password string `json:"password"`
}

func (l LoginUser)LoginCheck() bool {
	// 进行查找 不应该用id作为账户 然后这里如果find到 应该返回id 没find到 则返回nil or 负数
	b := Databases.UserLoginCheck(l.ID, l.Password)
	return b
}