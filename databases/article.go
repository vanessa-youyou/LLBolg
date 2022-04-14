package databases

import (
	"LlBlog/models"
	"fmt"
	"strconv"
)

// SearchArticleById 查找文章 byId
func SearchArticleById(article models.ArticleInfo) (bool, *models.ArticleInfo) {
	err := DB.Model(&article).Where("author_id = ?", article.ID).Find(&article).Error
	if err != nil{
		fmt.Println(err)
		return false, &article
	}
	return true, &article
}

// WriteNewArticles 新文章
func WriteNewArticles(a *models.ArticleInfo) bool {
	err := DB.Create(&a).Error
	if err != nil{
		fmt.Println(err)
		return false
	}
	return true
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
		fmt.Println("不是本人在操作")
		return false, nil
	}

	// 删除文章相关评论的点赞数目
	// 1: 要找到所有的此文章的评论的id
	var comments []models.CommentInfo
	err = DB.Model(&models.CommentInfo{}).Where("article_id = ?", a.ID).Find(&comments).Error
	if err != nil{
		return false, err
	}
	for i := 0; i < len(comments); i++ {
		commentId := comments[i].ID
		SetName := strconv.Itoa(int(commentId))
		SetName += "LikeComment:"
		Redis.Del(SetName)
	}

	// 删除此文章相关评论
	var nc models.CommentInfo
	err = DB.Model(&nc).Where("article_id = ?", a.ID).Delete(&models.CommentInfo{}).Error
	if err != nil{
		fmt.Println(err, "数据库删除评论出错")
		return false, err
	}

	// 删除此文章的点赞
	ArticleName := strconv.Itoa(int(a.ID))
	ArticleName += "LikeArticle:"
	Redis.Del(ArticleName)

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
		fmt.Println(err, "数据库更新出错")
		return false, err
	}
	return true, nil
}

// LikeArticle 点赞文章 取消操作 -redis
func LikeArticle(a *models.ArticleInfo, userId uint) (bool, error) {
	// 用set set的名字为文章id 内容为 用户id
	// 1:先查找 有无此 value
	ArticleId := strconv.Itoa(int(a.ID))
	ArticleId += "LikeArticle:"
	UserId := strconv.Itoa(int(userId))
	if !Redis.SIsMember(ArticleId, UserId).Val(){
		// 如果没有 就把这个数据存入
		Redis.SAdd(ArticleId, UserId)
	}else{
		// 如果存在 就删除
		Redis.SRem(ArticleId, UserId)
	}
	return true, nil
}

// NewComment 评论新建
func NewComment(cm *models.CommentInfo) bool {
	var count int
	// 1 检查表中有无文章
	err := DB.Model(&models.ArticleInfo{}).Where("id = ? ", cm.ArticleID).Count(&count).Error
	if err != nil{
		fmt.Println(err)
		return false
	}
	if count == 0{
		return false
	}
	// 存在此文章 新建此评论
	err = DB.Create(&cm).Error
	if err != nil{
		fmt.Println(err)
		return false
	}
	return true
}

// FindAllArticleByUserId 查找个人文章通过操作人ID
func FindAllArticleByUserId(u *models.UserInfo) (bool, error, []models.ArticleInfo) {
	var article []models.ArticleInfo
	err := DB.Model(&article).Where("author_id = ?", u.ID).Find(&article).Error
	if err != nil{
		fmt.Println(err)
		return false, err, nil
	}
	// 应该再在redis中 找到每个文章的点赞量
	for i := 0; i < len(article); i++ {
		// 遍历文章 通过文章id找到 赞的数量 评论的数量
		ArticleName := strconv.Itoa(int(article[i].ID))
		ArticleName += "LikeArticle:"
		likeNum := Redis.SCard(ArticleName).Val()
		fmt.Println("Redis.SCard(ArticleName).Val() is", likeNum)

		// 评论的数量
		var count int
		err = DB.Model(&models.CommentInfo{}).Where("article_id = ?", article[i].ID).Count(&count).Error
		if err != nil{
			fmt.Println(err)
			return false, err, nil
		}

		article[i].LikeNum = int(likeNum)
		article[i].CommentsNum = count
	}

	return true, nil, article
}

// ArticleDetails 文章详情页面（应该要有评论 每个评论的赞）
func ArticleDetails(a *models.ArticleInfo) (bool, []models.CommentInfo) {
	// 1:根据 文章id 找到所有的评论。放到评论的数组里
	var comments []models.CommentInfo
	err := DB.Model(&comments).Where("article_id = ?", a.ID).Find(&comments).Error
	if err != nil{
		fmt.Println("出错啦！err := DB.Model(&comments).Where(\"author_id = ?\", a.ID).Find(&comments).Error")
		fmt.Println(err)
		return false, nil
	}

	// 找到了所有的评论 应该还要知道没条评论的赞
	for i:=1; i<len(comments);i++{
		// 遍历文章 通过文章id找到 赞的数量 评论的数量
		CommentName := strconv.Itoa(int(comments[i].ID))
		CommentName += "LikeComment:"
		likeNum := Redis.SCard(CommentName).Val()
		comments[i].LikeNum = int(likeNum)
	}

	return true, comments
}

// CollectionArticle 收藏文章
func CollectionArticle(coll *models.Collection) (bool, error) {
	var count int
	// 1 检查表中有无文章
	err := DB.Model(&models.ArticleInfo{}).Where("id = ? ", coll.ArticleID).Count(&count).Error
	if err != nil{
		fmt.Println(err)
		return false, err
	}
	if count == 0{
		return false, nil
	}
	// 存在 进行收藏
	err = DB.Create(&coll).Error
	if err != nil{
		fmt.Println(err)
		return false, err
	}
	return true, nil
}

// CancelCollectionArticle 取消收藏
func CancelCollectionArticle(coll *models.Collection) (bool, error) {
	err := DB.Model(&coll).Where("article_id = ? AND user_id = ?", coll.ArticleID, coll.UserID).Delete(&models.Collection{}).Error
	if err != nil{
		fmt.Println(err, "数据库删除出错")
		return false, err
	}
	return true, nil
}

// CreatLabel 创建标签
func CreatLabel(label *models.Label) bool {
	// 1 检查表中有无标签
	var count = 0
	err := DB.Model(&models.Label{}).Where("name = ? ", label.Name).Count(&count).Error
	if err != nil{
		fmt.Println(err)
		return false
	}
	if count != 0{
		return false
	}
	// 不存在此标签 新建
	err = DB.Create(&label).Error
	if err != nil{
		fmt.Println(err)
		return false
	}
	return true
}

// SearchLabel 搜索标签
func SearchLabel(label *models.Label) (bool, []models.Label) {
	var labels []models.Label
	err := DB.Model(&labels).Where("name Like ?", label.Name).Find(&labels).Error
	if err != nil{
		fmt.Println(err)
		return false, nil
	}
	return true, labels
}

// ChooseLabels  为文章添加标签
func ChooseLabels(labels models.LabelReceive, userId uint) (bool, error){
	// 检查是否为本人
	var article models.ArticleInfo
	err := DB.Model(&article).Where("id = ?",labels.ArticleId).Find(&article).Error
	if err != nil{
		return false, err
	}
	if article.AuthorID != userId{
		fmt.Println("不是本人在操作")
		return false, nil
	}

	for i := 0; i < len(labels.LabelId); i++{
		var count = 0
		// 检查 是否有这个标签 没有的话 不可以创建
		err := DB.Model(&models.Label{}).Where("id = ? ", labels.LabelId[i]).Count(&count).Error
		if err != nil{
			fmt.Println(err)
			return false, err
		}
		if count == 0{
			continue
		}

		count = 0
		err = DB.Model(&models.LabelRelation{}).Where("label_id = ? AND article_id = ?",
			labels.LabelId[i], labels.ArticleId).Count(&count).Error
		if err != nil{
			fmt.Println(err)
			return false, err
		}
		if count != 0{
			continue
		}
		// 不存在这个记录 可以创建
		var labelNew = models.LabelRelation{}
		labelNew.LabelId = labels.LabelId[i]
		labelNew.ArticleId = labels.ArticleId

		err = DB.Create(&labelNew).Error
		if err != nil{
			fmt.Println(err)
			return false, err
		}
	}
	return true, nil
}

// FindTheLatestArticles 查找最新的公开文章
func FindTheLatestArticles() (bool, error, []models.ArticleInfo) {
	var article []models.ArticleInfo
	err := DB.Order("id").Model(&article).Where("state = 4").Find(&article).Error
	if err != nil{
		fmt.Println(err)
		return false, err, nil
	}
	// 应该再在redis中 找到每个文章的点赞量
	for i := 0; i < len(article); i++ {
		// 遍历文章 通过文章id找到 赞的数量 评论的数量
		ArticleName := strconv.Itoa(int(article[i].ID))
		ArticleName += "LikeArticle:"
		likeNum := Redis.SCard(ArticleName).Val()
		fmt.Println("Redis.SCard(ArticleName).Val() is", likeNum)

		// 评论的数量
		var count int
		err = DB.Model(&models.CommentInfo{}).Where("article_id = ?", article[i].ID).Count(&count).Error
		if err != nil{
			fmt.Println(err)
			return false, err, nil
		}

		article[i].LikeNum = int(likeNum)
		article[i].CommentsNum = count
	}

	return true, nil, article
}