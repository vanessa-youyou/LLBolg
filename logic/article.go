package logic

import (
	"LlBlog/dao"
	"LlBlog/dto"
	"LlBlog/models"
	"fmt"
	"strconv"
)

// WriteNewArticles 新增文章
func WriteNewArticles(article *dto.CreateArticleReq) (bool, models.ArticleInfo) {
	a := models.ArticleInfo{
		Title:       article.Title,
		Text:        article.Text,
		AuthorID:    article.AuthorID,
		Original:    article.Original,
		PlacedTop:   article.PlacedTop,
		State:       article.State,
		LikeNum:     0,
		CommentsNum: 0,
	}
	newArticle, t := dao.WriteNewArticles(&a)
	if !t{
		return false, newArticle
	}
	return true, newArticle
}

// RemoveArticle 检查是否为作者本人 并删除文章
func RemoveArticle(aId uint, userId uint) bool {
	t, err := dao.ArticleRemove(aId, userId)
	if err != nil {
		return false
	}
	return t
}

// ArticleModify 检验是否为作者本人 并且更新文章
func ArticleModify(article *dto.UpdateArticleReq, userId uint) bool {

	a  := models.ArticleInfo{
		ID:          article.ArticleId,
		Title:       article.Title,
		Text:        article.Text,
		AuthorID:    userId,
		Original:    article.Original,
		PlacedTop:   article.PlacedTop,
		State:       article.State,
	}
	t, err := dao.ModifyArticle(&a, userId)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return t
}

// IsAuthorSelf 检验作者是否为本人
func IsAuthorSelf(l dto.LabelReceiveReq, userId uint) bool {
	t, err := dao.ChooseLabels(l, userId)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return t
}

// ArticleLike 点赞文章-redis
func ArticleLike(comment *dto.PickArticleReq, userId uint) bool {
	cm := models.GiveLike{
		UserID:    userId,
		ArticleID: comment.ArticleId,
	}
	t, err := dao.LikeArticle(&cm, userId)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return t
}

// CancelArticleLike 取消点赞文章
func CancelArticleLike(comment *dto.CancelPickArticleReq, userId uint) bool {
	cm := models.GiveLike{
		UserID:    userId,
		ArticleID: comment.ArticleId,
	}
	t, err := dao.CancelLikeArticle(&cm, userId)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return t
}

// CommentLike 点赞评论 点赞redis
func CommentLike(cm *dto.PickCommentReq, userId uint) bool {
	comment := models.CommentLike{
		UserID:    userId,
		CommentID: cm.CommentId,
		ArticleID: cm.ArticleId,
	}
	t, err := dao.LikeComment(&comment, userId)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return t
}

// CancelCommentLike 取消点赞评论
func CancelCommentLike(cm *dto.CancelPickCommentReq, userId uint) bool {
	comment := models.CommentLike{
		UserID:    userId,
		CommentID: cm.CommentId,
		ArticleID: cm.ArticleId,
	}
	t, err := dao.CancelLikeComment(&comment, userId)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return t
}


// RemoveComment 删除评论
func RemoveComment(cm *dto.DeleteCommentReq, userId uint) bool {
	comment := models.CommentInfo{
		ID:        	cm.CommentId,
		UserID:    	userId,
		ArticleID: 	cm.ArticleId,
	}
	t := dao.CommentDelete(&comment, userId)

	return t
}

// FindAllArticleByUserId 通过操作人id 获取所有文章
func FindAllArticleByUserId(u *models.UserInfo) (bool, []models.ArticleInfo) {
	t, err, articlePage := dao.FindAllArticleByUserId(u)
	if err != nil || articlePage == nil {
		fmt.Println(err)
		return false, nil
	}
	return t, articlePage
}

// ArticleDetails 文章详情
func ArticleDetails(arrticle *dto.ArticleDetailsReq) (bool, models.ArticleInfo, []models.CommentInfo) {
	a := models.ArticleInfo{ID: arrticle.ArticleId}
	// 按照文章id填写 文章
	var t bool
	t, a = dao.SearchArticleById(&a)
	if !t {
		fmt.Println("此文章不存在 / 不可访问状态")
		return false, a, nil
	}
	// 1 找到文章的all 评论，各个评论的赞
	t, comments := dao.ArticleDetails(&a)
	if !t {
		fmt.Println("错在dao中 评论，各个评论的赞 查询失败")
		return false, a, nil
	}
	// 2 填上文章的评论数量
	a.CommentsNum = len(comments)
	// 3 文章的点赞量
	// 遍历文章 通过文章id找到 赞的数量 评论的数量
	ArticleName := strconv.Itoa(int(a.ID))
	ArticleName += "LikeArticle:"
	likeNum := dao.Redis.SCard(ArticleName).Val()
	// fmt.Println("Redis.SCard(ArticleName).Val() is", likeNum)
	a.LikeNum = int(likeNum)
	return true, a, comments
}

// SearchArticle 查找文章
func SearchArticle(search dto.SearchArticleReq) ([]models.ArticleInfo, bool) {
	t, err, articles := dao.AccurateSearch(search)
	if err != nil || !t {
		fmt.Println(err)
		return nil, false
	}
	return articles, t
}

// CollectionArticle 收藏文章
func CollectionArticle(collection *models.Collection) bool {
	// 在数据库中查找是不是真的有这个文章，有的话就收藏 没有就失败
	t, err := dao.CollectionArticle(collection)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return t
}

// CancelCollectionArticle 取消收藏
func CancelCollectionArticle(collection *models.Collection) bool {
	t, err := dao.CancelCollectionArticle(collection)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return t
}

// FindTheLatestArticles 找 最新的文章
func FindTheLatestArticles() (bool, []models.ArticleInfo) {
	t, err, articlePage := dao.FindTheLatestArticles()
	if err != nil || articlePage == nil {
		fmt.Println(err)
		return false, nil
	}
	return t, articlePage
}
