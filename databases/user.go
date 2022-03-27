package databases

import (
	"LlBlog/models"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
)

// databases应该关注于单纯的通用逻辑
// 例如通过id获取记录这种

func GetUserById(id uint) (*models.UserInfo, error) {
	var u models.UserInfo
	if err := DB.Where("id = ?", id).Find(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// GetUserByCustomerName 通过用户名查找 个人信息
func GetUserByCustomerName(cn string) (*models.UserInfo, error){
	var u models.UserInfo
	if err := DB.Where("customer_name = ?", cn).Find(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// AccountRechecking 账号查重
func AccountRechecking(c models.UserInfo) (bool, error){
	var count int
	err := DB.Model(&c).Where("customer_name = ?", c.CustomerName).Count(&count).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 没有相同的账号名 允许创建
		if count == 0{
			return true, nil
		}
		return true, nil
	}
	if count == 0{
		return true, nil
	}
	return false, err
}

// AccountInsert 添加账号
func AccountInsert(u *models.UserInfo)  (bool, error){
	d := []byte(u.Password)
	m := md5.New()
	m.Write(d)
	u.Password = hex.EncodeToString(m.Sum(nil))
	err := DB.Create(&u).Error
	if err != nil{
		return false, err
	}
	return true, nil
}

// UserInformationUpdate 个人信息更新
func UserInformationUpdate(u *models.UserInfo) (bool, error) {
	err := DB.Model(&u).Updates(models.UserInfo{
		Name: u.Name,
		Gender: u.Gender,
		Introduction: u.Introduction,
		Password: u.Password,
		Label: u.Label,
	}).Error

	if err != nil{
		fmt.Println("数据库更新出错")
		return false, err
	}
	return true, nil
}

// WriteNewArticles 新文章
func WriteNewArticles(a *models.ArticleInfo) (bool, error) {
	err := DB.Create(&a).Error
	if err != nil{
		return false, err
	}
	return true, nil
}

// ArticleRemove 删除文章（删除前需要查证操作人是不是作者）
func ArticleRemove(a *models.ArticleInfo, userId uint) (bool, error) {
	// 如果文章id = 要改的文章id 而且 文章作者 等于 操作者
	// 防止传过来的a.AuthorID 是假的 重新通过文章id 查一下文章
	err := DB.Model(&models.ArticleInfo{}).Where("id = ? ", a.ID).Find(&a).Error
	if err != nil{
		return false, err
	}
	if a.AuthorID != userId{
		fmt.Println(a.AuthorID,"   ", userId)
		fmt.Println("不是本人在操作")
		return false, nil
	}
	var cl models.CommentLike
	// 删除此文章相关的评论的点赞
	err = DB.Model(&cl).Where("article_id = ?", a.ID).Delete(&models.CommentLike{}).Error
	if err != nil{
		fmt.Println("数据库删除点赞出错")
		return false, err
	}

	// 删除此文章相关评论
	var nc models.CommentInfo
	err = DB.Model(&nc).Where("article_id = ?", a.ID).Delete(&models.CommentInfo{}).Error
	if err != nil{
		fmt.Println("数据库删除评论出错")
		return false, err
	}

	// 删除此文章的点赞
	var al models.GiveLike
	err = DB.Model(&al).Where("article_id = ?", a.ID).Delete(&models.GiveLike{}).Error
	if err != nil{
		fmt.Println("数据库删除文章的赞出错")
		return false, err
	}

	// 删除文章
	err = DB.Model(&a).Where("id = ? AND author_id = ?", a.ID, userId).Delete(&models.ArticleInfo{}).Error
	if err != nil{
		fmt.Println("数据库删除出错")
		return false, err
	}
	return true, nil
}

// ModifyArticle 修改文章（修改前需要查证操作人是不是作者）
func ModifyArticle(a *models.ArticleInfo, userId uint) (bool, error) {
	// 如果文章id = 要改的文章id 而且 文章作者 等于 操作者
	err := DB.Model(&a).Where("id = ? AND author_id = ?", a.ID, userId).Updates(models.ArticleInfo{
		Title: a.Title,
		Text: a.Text,
		AuthorID: userId,
		Original: a.Original,
		PlacedTop: a.PlacedTop,
		State: a.State,
	}).Error
	if err != nil{
		fmt.Println("数据库更新出错")
		return false, err
	}
	return true, nil
}

// ArticlePick 点赞/取消点赞操作：
func ArticlePick(a *models.ArticleInfo, userId uint) (bool, error) {
	var count int = 0
	// 1 检查表中 有无文章id =a.id 用户名id = userid 的 有就删除 else 创建
	err := DB.Model(&models.GiveLike{}).Where("user_id = ? AND article_id = ?", userId, a.ID).Count(&count).Error

	// 没有 创建
	var L models.GiveLike
	L.ArticleID = a.ID
	L.UserID = userId
	if count == 0{
		err := DB.Create(&L).Error
		if err != nil{
			return false, err
		}
		return true, nil
	}else{
		err = DB.Model(L).Where("user_id = ? AND article_id = ?", userId, a.ID).Delete(&models.GiveLike{}).Error
		if err != nil{
			return false, err
		}
		return true, nil

	}
}

// NewComment 按照文章id查文章
func NewComment(cm *models.CommentInfo) (bool, error) {
	// 1 查找文章
	var count int = 0
	// 1 检查表中
	err := DB.Model(&models.ArticleInfo{}).Where("id = ? ", cm.ArticleID).Count(&count).Error
	if err != nil{
		return false, err
	}
	if count == 0{
		return false, nil
	}
	// 存在此文章 新建此评论
	err = DB.Create(&cm).Error
	if err != nil{
		return false, err
	}
	return true, nil
}

// CommentPick 评论点赞 点赞/取消点赞操作：
func CommentPick(cm *models.CommentInfo, userId uint) (bool, error) {
	var count int = 0
	// 1 检查表中 有无评论id =a.id 用户名id = userid 的 有就删除 else 创建
	err := DB.Model(&models.CommentLike{}).Where("user_id = ? AND comment_id = ? AND article_id = ?", userId, cm.ID, cm.ArticleID).Count(&count).Error

	// 没有 创建
	var g models.CommentLike
	g.UserID = userId
	g.CommentID = cm.ID
	g.ArticleID = cm.ArticleID
	if count == 0{
		err := DB.Create(&g).Error
		if err != nil{
			return false, err
		}
		return true, nil
	}else{
		err = DB.Model(g).Where("user_id = ? AND comment_id = ? AND article_id",  userId, cm.ID, cm.ArticleID).Delete(&models.CommentLike{}).Error
		if err != nil{
			return false, err
		}
		return true, nil

	}
}

// CommentDelete 删除评论
func CommentDelete(cm *models.CommentInfo, userId uint) (bool, error) {
	// 1检验操作人是不是评论人员：
	if cm.UserID != userId{
		return false, nil
	}
	// 删除所有的赞 where commentId = cm.Id
	var g models.CommentLike
	err := DB.Model(g).Where("comment_id = ?",cm.ID).Delete(&models.CommentLike{}).Error
	if err != nil{
		return false, err
	}
	// 删除此条评论
	err = DB.Model(cm).Where("id = ?",cm.ID).Delete(&models.CommentInfo{}).Error
	if err != nil{
		return false, err
	}
	return true, nil
}