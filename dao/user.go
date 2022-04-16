package dao

import (
	"LlBlog/models"
	"errors"

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
func GetUserByCustomerName(cn string) (*models.UserInfo, error) {
	var u models.UserInfo
	if err := DB.Where("customer_name = ?", cn).Find(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

// AccountRechecking 账号查重
func AccountRechecking(c models.UserInfo) (bool, error) {
	var count int
	err := DB.Model(&c).Where("customer_name = ?", c.CustomerName).Count(&count).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 没有相同的账号名 允许创建
		if count == 0 {
			return true, nil
		}
		return true, nil
	}
	if count == 0 {
		return true, nil
	}
	return false, err
}

// AccountInsert 添加账号
// func AccountInsert(u *models.UserInfo) (bool, error) {
// 	d := []byte(u.Password)
// 	m := md5.New()
// 	m.Write(d)
// 	u.Password = hex.EncodeToString(m.Sum(nil))
// 	err := DB.Create(&u).Error
// 	if err != nil {
// 		fmt.Println(err)
// 		return false, err
// 	}
// 	return true, nil
// }
//
// // UserInformationUpdate 个人信息更新
// func UserInformationUpdate(u *models.UserInfo) (bool, error) {
// 	err := DB.Model(&u).Updates(models.UserInfo{
// 		Name:         u.Name,
// 		Gender:       u.Gender,
// 		Introduction: u.Introduction,
// 		Password:     u.Password,
// 		Label:        u.Label,
// 		HeadPortrait: u.HeadPortrait,
// 	}).Error
//
// 	if err != nil {
// 		fmt.Println("数据库更新出错")
// 		return false, err
// 	}
// 	return true, nil
// }
//
// // UpdateHeadPortrait 更新头像
// func UpdateHeadPortrait(url string, u *models.UserInfo) (bool, error) {
// 	err := DB.Model(&u).Updates(models.UserInfo{
// 		HeadPortrait: url,
// 	}).Error
//
// 	if err != nil {
// 		fmt.Println("头像更新出错")
// 		return false, err
// 	}
// 	return true, nil
// }
//
// // AccurateSearch 查询用户
// func AccurateSearch(search models.Search) (bool, error, []models.ArticleInfo) {
// 	var articles []models.ArticleInfo
// 	var err error
// 	fmt.Println(search.SearchWay, "   ", search.Content, "   ", search.Check)
// 	if search.SearchWay && search.Check == "title" {
// 		fmt.Println("1")
// 		err = DB.Model(&models.ArticleInfo{}).Where("title = ?", search.Content).Find(&articles).Error
//
// 	} else if !search.SearchWay && search.Check == "title" {
// 		fmt.Println("2")
// 		err = DB.Model(&models.ArticleInfo{}).Where("title LIKE ?", search.Content+"%").Find(&articles).Error
// 	} else if search.SearchWay && search.Check == "text" {
// 		fmt.Println("3")
// 		err = DB.Model(&models.ArticleInfo{}).Where("text = ?", search.Content).Find(&articles).Error
// 	} else if !search.SearchWay && search.Check == "text" {
// 		fmt.Println("4")
// 		err = DB.Model(&models.ArticleInfo{}).Where("text LIKE ?", search.Content+"%").Find(&articles).Error
// 	}
// 	if err != nil {
// 		fmt.Println("查找文章失败")
// 		return false, err, nil
// 	}
// 	for i := 0; i < len(articles); i++ {
// 		// 遍历文章 通过文章id找到 赞的数量 评论的数量
// 		ArticleName := strconv.Itoa(int(articles[i].ID))
// 		ArticleName += "LikeArticle:"
// 		likeNum := Redis.SCard(ArticleName).Val()
// 		fmt.Println("Redis.SCard(ArticleName).Val() is", likeNum)
//
// 		// 评论的数量
// 		var count int
// 		err = DB.Model(&models.CommentInfo{}).Where("article_id = ?", articles[i].ID).Count(&count).Error
// 		if err != nil {
// 			return false, err, nil
// 		}
//
// 		articles[i].LikeNum = int(likeNum)
// 		articles[i].CommentsNum = count
// 	}
// 	return true, nil, articles
//
// }
//
// // SearchUser 查询用户
// func SearchUser(search models.Search) (bool, error, []models.UserInfo) {
// 	var user []models.UserInfo
// 	var err error
// 	if search.SearchWay {
// 		// 准确查询
// 		err = DB.Model(&models.UserInfo{}).Where("name = ?", search.Content).Find(&user).Error
// 	} else {
// 		err = DB.Model(&models.UserInfo{}).Where("name LIKE ?", search.Content+"%").Find(&user).Error
// 	}
// 	if err != nil {
// 		fmt.Println("查找失败")
// 		return false, err, nil
// 	}
// 	for i := 1; i < len(user); i++ {
// 		// 不允许访客看到 customer_name,password,id
// 		user[i].CustomerName = ""
// 		user[i].Password = ""
// 	}
// 	return true, nil, user
// }
