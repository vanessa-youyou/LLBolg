package services

import (
	"LlBlog/databases"
	"LlBlog/models"
	"crypto/md5"
	"encoding/hex"
)

type LoginUser struct {
	CustomerName string `json:"customer_name" form:"customer_name" `
	Password     string `json:"password" form:"password" `
}

func (l LoginUser) LoginCheck(user *models.UserInfo) bool {
	// 这种判断密码是否相同的逻辑应该放在
	// services里面，不能放在databases，因为不通用
	// TODO 是不是应该加密
	// md5 encryption
	d := []byte(l.Password)
	m := md5.New()
	m.Write(d)
	return user.Password == hex.EncodeToString(m.Sum(nil))
}

func RegisteredNameCheck(u models.UserInfo) bool {
	cn, err := databases.AccountRechecking(u.CustomerName)
	if err != nil {
		// TODO 日志
		return false
	}
	if !cn {
		// 没有相同的账号名 允许创建
		// TODO 日志
		return false
	}
	return true
}

func AddAccount(u models.UserInfo) bool {
	t, err := databases.AccountInsert(&u)
	if err != nil {
		return false
	}
	if !t {
		return false
	}
	return true
}

// UpdateUserInformation 进行更新
func UpdateUserInformation(u models.UserInfo) bool {
	t,err :=databases.UserInformationUpdate(&u)
	if err != nil {
		return false
	}
	return t
}

// 进行文章的保存
func NewArticles(a models.ArticleInfo)  bool {
	t,err := databases.WriteNewArticles(&a)
	if err != nil {
		return false
	}
	return t
}



// Like 点赞的结构体（赞 没有取消）
type Like struct {
	ID uint `json:"article_id" form:"article_id"`
	PickIt bool `json:"pick_it" form:"pick_it" `
}


// UpdatArticle 进行更新
func UpdatArticle (a *models.ArticleInfo) bool {
	t,err := databases.GiveLike(a)
	if err != nil {
		return false
	}
	return t
}

// AuthorCheck 检查操作人是否为作者本人
func AuthorCheck(authorId uint, articleId uint )  bool {
	article,err := databases.FindArticleById(articleId)
	if err != nil {
		return false
	}
	if authorId != article.UserInfoID{
		return false
	}
	return true
}

// ArticleDelete 删除函数
func ArticleDelete(id uint)  bool {
	t,err :=databases.DeleteArticle(id)
	if err != nil {
		return false
	}
	return t
}