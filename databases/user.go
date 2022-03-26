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
	err := DB.Model(&a).Where("id = ? AND author_id = ?", a.ID, userId).Delete(&models.ArticleInfo{}).Error
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
	//li := new(models.GiveLike)
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