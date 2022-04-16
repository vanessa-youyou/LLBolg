package logic

// NewArticles 进行文章的保存
//func NewArticles(a models.ArticleInfo)  bool {
//	err := databases.DB.Create(a).Error
//	if err != nil{
//		fmt.Println(err)
//		return false
//	}
//	return true
//}

// RemoveArticle 检查是否为作者本人 并删除文章
// func RemoveArticle(a *models.ArticleInfo, userId uint) bool {
// 	t, err := databases.ArticleRemove(a, userId)
// 	if err != nil {
// 		fmt.Println(err)
// 		return false
// 	}
// 	return t
// }
//
// // ArticleModify 检验是否为作者本人 并且更新文章
// func ArticleModify(a *models.ArticleInfo, userId uint) bool {
// 	t, err := databases.ModifyArticle(a, userId)
// 	if err != nil {
// 		fmt.Println(err)
// 		return false
// 	}
// 	return t
// }
//
// // IsAuthorSelf 检验作者是否为本人
// func IsAuthorSelf(l models.LabelReceive, userId uint) bool {
// 	t, err := databases.ChooseLabels(l, userId)
// 	if err != nil {
// 		fmt.Println(err)
// 		return false
// 	}
// 	return t
// }
//
// // ArticleLike 点赞文章-redis
// func ArticleLike(a *models.ArticleInfo, userId uint) bool {
// 	t, err := databases.LikeArticle(a, userId)
// 	if err != nil {
// 		fmt.Println(err)
// 		return false
// 	}
// 	return t
// }
//
// // CreatComment 创建评论
// //func CreatComment(cm *models.CommentInfo) bool {
// //	// 找这个文章id是否存在 不存在则false
// //	// 存在则创建
// //	var count = 0
// //	err := databases.DB.Model(&models.ArticleInfo{}).Where("id = ? ", cm.ArticleID).Count(&count).Error
// //	if err != nil{
// //		fmt.Println(err)
// //		return false
// //	}
// //	if count == 0{
// //		return false
// //	}
// //	// 存在此文章 新建此评论
// //	err = databases.DB.Create(&cm).Error
// //	if err != nil{
// //		fmt.Println(err)
// //		return false
// //	}
// //	return true
// //}
//
// // CommentLike 点赞/取消 评论 点赞redis
// func CommentLike(cm *models.CommentInfo, userId uint) bool {
// 	t, err := databases.LikeComment(cm, userId)
// 	if err != nil {
// 		fmt.Println(err)
// 		return false
// 	}
// 	return t
// }
//
// // RemoveComment 删除评论
// func RemoveComment(cm *models.CommentInfo, userId uint) bool {
// 	t, err := databases.CommentDelete(cm, userId)
// 	if err != nil {
// 		fmt.Println(err)
// 		return false
// 	}
// 	return t
// }
//
// // FindAllArticleByUserId 通过操作人id 获取所有文章
// func FindAllArticleByUserId(u *models.UserInfo) (bool, []models.ArticleInfo) {
// 	t, err, articlePage := databases.FindAllArticleByUserId(u)
// 	if err != nil || articlePage == nil {
// 		fmt.Println(err)
// 		return false, nil
// 	}
// 	return t, articlePage
// }
//
// // ArticleDetails 文章详情
// func ArticleDetails(a *models.ArticleInfo) (bool, models.ArticleInfo, []models.CommentInfo) {
// 	// 按照文章id填写 文章
// 	t, a := databases.SearchArticleById(*a)
// 	if !t {
// 		fmt.Println("错在dao中!!!")
// 		return false, *a, nil
// 	}
// 	// 1 找到文章的all 评论，各个评论的赞
// 	t, comments := databases.ArticleDetails(a)
// 	if !t {
// 		fmt.Println("错在dao中")
// 		return false, *a, nil
// 	}
// 	// 2 填上文章的评论数量
// 	a.CommentsNum = len(comments)
// 	// 3 文章的点赞量
// 	// 遍历文章 通过文章id找到 赞的数量 评论的数量
// 	ArticleName := strconv.Itoa(int(a.ID))
// 	ArticleName += "LikeArticle:"
// 	likeNum := databases.Redis.SCard(ArticleName).Val()
// 	fmt.Println("Redis.SCard(ArticleName).Val() is", likeNum)
// 	a.LikeNum = int(likeNum)
// 	return true, *a, comments
// }
//
// // SearchArticle 查找文章
// func SearchArticle(search models.Search) ([]models.ArticleInfo, bool) {
// 	t, err, articles := databases.AccurateSearch(search)
// 	if err != nil || !t {
// 		fmt.Println(err)
// 		return nil, false
// 	}
// 	return articles, t
// }
//
// // CollectionArticle 收藏文章
// func CollectionArticle(collection *models.Collection) bool {
// 	// 在数据库中查找是不是真的有这个文章，有的话就收藏 没有就失败
// 	t, err := databases.CollectionArticle(collection)
// 	if err != nil {
// 		fmt.Println(err)
// 		return false
// 	}
// 	return t
// }
//
// // CancelCollectionArticle 取消收藏
// func CancelCollectionArticle(collection *models.Collection) bool {
// 	t, err := databases.CancelCollectionArticle(collection)
// 	if err != nil {
// 		fmt.Println(err)
// 		return false
// 	}
// 	return t
// }
//
// // FindTheLatestArticles 找 最新的文章
// func FindTheLatestArticles() (bool, []models.ArticleInfo) {
// 	t, err, articlePage := databases.FindTheLatestArticles()
// 	if err != nil || articlePage == nil {
// 		fmt.Println(err)
// 		return false, nil
// 	}
// 	return t, articlePage
// }
