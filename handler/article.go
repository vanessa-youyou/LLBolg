package handler

import (
	"LlBlog/core"
	"LlBlog/dao"
	"LlBlog/dto"
	"LlBlog/errors"
	"LlBlog/logic"
	"LlBlog/models"
	"LlBlog/utils"
	"github.com/gin-gonic/gin"
)

// CreatArticle 写文章✓
func CreatArticle(c *gin.Context) (interface{}, error) {
	// 登录验证
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin() {
		return nil, errors.IsNotLogin
	}

	// 获取数据
	var newArticle dto.CreateArticleReq
	err := c.ShouldBind(&newArticle)
	if err != nil {
		return nil, errors.ReceiveParametersError
	}

	// 绑定正确作者id
	newArticle.AuthorID = auth.User.ID
	// 状态有(1:草稿  2:发布  3:发布-审核中 4:发布成功 5:驳回 6:撤销)
	//这里不允许出现撤销 草稿和发布之外的选项
	if newArticle.State != 1 && newArticle.State != 2 && newArticle.State != 6{
		// 默认模式为 发布
		newArticle.State = 2
	}

	// 进行保存
	t, succArticle := logic.WriteNewArticles(&newArticle)
	 if !t{
		 return nil, errors.WriteError
	 }

	return dto.CreateArticleRsp{
		Article: succArticle,
	}, nil
}

// ModifyArticle 修改文章 只有作者可以(修改完成)✓
func ModifyArticle(c *gin.Context) (interface{}, error) {
	// 登录验证
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin() {
		return nil, errors.IsNotLogin
	}

	// 接收数据
	var article dto.UpdateArticleReq
	err := c.ShouldBind(&article)
	if err != nil {
		return nil, errors.ReceiveParametersError
	}

	// 处理 状态参数 （只能为发布/草稿 默认:发布）
	if article.State != 1 && article.State != 2 && article.State != 6{
		// 默认模式为 发布
		article.State = 2
	}

	// 查证 操作人 是否为文章作者
	// 新的查找and 修改 逻辑是(update 文章 where authorID == article.authorID)
	if ! logic.ArticleModify(&article, auth.User.ID){
		utils.Return(c, errors.UpdateError)
	}

	// 成功
	return dto.UpdateArticleRsp{
		Message: "修改成功 这里应该返回文章页面",
	}, nil
}

// DeleteArticle 删除文章（评论还没写 评论应该一起被删）只允许作者本人删除
func DeleteArticle(c *gin.Context) (interface{}, error) {
	// 登录验证
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin() {
		return nil, errors.IsNotLogin
	}
	// 接收数据
	var articleId dto.DeleteArticleReq
	err := c.ShouldBind(&articleId)
	if err != nil {
		return nil, errors.ReceiveParametersError
	}

	// 获得操作人的id
	var userId = auth.User.ID

	// 1:检查操作者是否为 作者本人，是则删除，不是则报错
	if ! logic.RemoveArticle(articleId.ArticleId, userId){
		return nil,errors.DeleteArticleError
	}

	return dto.DeleteArticleRsp{
		Message: "删除成功！ 这里应该返回个人首页面",
	},nil
}


// CreateComment 新建评论
func CreateComment(c *gin.Context) (interface{}, error) {
	// 登录验证
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin() {
		return nil, errors.IsNotLogin
	}

	// 接收数据
	var cm dto.CreateCommentReq
	err := c.ShouldBind(&cm)
	if err != nil {
		return nil, errors.ReceiveParametersError
	}

	// 修改操作人id
	cm.UserId = auth.User.ID
	var count = 0
	if err = dao.DB.Model(&models.ArticleInfo{}).Where("id = ? ", cm.ArticleId).Count(&count).Error; err != nil{
		return nil, err
	}
	if count < 0 {
		return nil, errors.CreatCommentError
	}
	comment := models.CommentInfo{
		UserID:    cm.UserId,
		ArticleID: cm.ArticleId,
		Content:   cm.Context,
		LikeNum:   0,
	}
	if err = dao.DB.Create(&comment).Error; err != nil{
		return nil, errors.CreatCommentError
	}

	return dto.CreateCommentRsp{
		Message: "新建评论成功",
	}, nil
}


// DeleteComment 删除评论(只允许作者本人删除
func DeleteComment(c *gin.Context) (interface{}, error) {
	// 登录验证
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin() {
		return nil, errors.IsNotLogin
	}

	// 获取删除评论时的信息
	var deleteComment dto.DeleteCommentReq
	err := c.ShouldBind(&deleteComment)
	if err != nil {
		return nil, errors.ReceiveParametersError
	}

	var userId = auth.User.ID
	if !logic.RemoveComment(&deleteComment, userId){
		return nil,errors.DeleteCommentError
	}

	// 成功
	return dto.DeleteCommentRsp{
		Message: "删除成功",
	}, nil
}


// PickArticle 点赞文章 用redis存储
func PickArticle(c *gin.Context) (interface{}, error) {
	// 登录验证
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin() {
		return nil, errors.IsNotLogin
	}

	//
	var like dto.PickArticleReq
	err := c.ShouldBind(&like)
	if err != nil {
		return nil, errors.ReceiveParametersError
	}

	// 操作者id
	var userId = auth.User.ID

	if !logic.ArticleLike(&like, userId){
		return nil, errors.PickError

	}

	// 成功
	return dto.PickArticleRsp{
		Message: "点赞/取消点赞 成功 这里应该还在文章页面",
	}, nil

}

// CancelPickArticle 取消点赞文章
func CancelPickArticle(c *gin.Context) (interface{}, error) {
	// 登录验证
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin() {
		return nil, errors.IsNotLogin
	}

	// 获取更改后的 个人信息
	var like dto.CancelPickArticleReq
	err := c.ShouldBind(&like)
	if err != nil {
		return nil, errors.ReceiveParametersError
	}

	// 操作者id
	var userId = auth.User.ID

	if !logic.CancelArticleLike(&like, userId){
		return nil, errors.PickError

	}

	// 成功
	return dto.CancelPickArticleRsp{
		Message: "取消点赞 成功 这里应该还在文章页面",
	}, nil
}

// CommentPick 点赞评论 用redis存储
func CommentPick(c *gin.Context) (interface{}, error) {
	// 登录验证
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin() {
		return nil, errors.IsNotLogin
	}
	// 接收数据
	var cm dto.PickCommentReq
	if err := c.ShouldBind(&cm); err != nil{
		return nil, errors.ReceiveParametersError
	}

	// 操作者id
	var userId = auth.User.ID
	if !logic.CommentLike(&cm, userId){
		return nil, errors.PickError
	}

	// 成功
	return dto.PickCommentRsp{
		Message: "点赞成功",
	}, nil
}

// CancelCommentPick 取消点赞
func CancelCommentPick(c *gin.Context) (interface{}, error) {
	// 登录验证
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin() {
		return nil, errors.IsNotLogin
	}
	// 接收数据
	var cm dto.CancelPickCommentReq
	if err := c.ShouldBind(&cm); err != nil{
		return nil, errors.ReceiveParametersError
	}

	// 操作者id
	var userId = auth.User.ID
	if !logic.CancelCommentLike(&cm, userId){
		return nil, errors.PickError
	}

	// 成功
	return dto.PickCommentRsp{
		Message: "点赞成功",
	}, nil
}

// SearchArticles 查询文章
func SearchArticles(c *gin.Context) (interface{}, error) {
	// 登录验证
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin() {
		return nil, errors.IsNotLogin
	}
	// 接收数据
	var search dto.SearchArticleReq
	if err := c.ShouldBind(&search); err != nil{
		return nil, errors.ReceiveParametersError
	}

	articles, t := logic.SearchArticle(search)
	if !t{
		return nil, errors.SearchERROR
	}
	return dto.SearchArticleRsq{
		Articles: articles,
		Message: "搜索成功",
	}, nil
}

// ArticleDetails 文章详情
func ArticleDetails(c *gin.Context) (interface{}, error) {
	// 登录验证
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin() {
		return nil, errors.IsNotLogin
	}

	// 	获取参数
	var article dto.ArticleDetailsReq
	if err := c.ShouldBind(&article); err != nil{
		return nil, errors.ReceiveParametersError
	}

	t, articles, comments := logic.ArticleDetails(&article)
	if !t{
		return nil,errors.ObtainDetailsError
	}
	return dto.ArticleDetailsRsp{
		Article: articles,
		Comment: comments,
	}, nil

}

// CollectionArticle 收藏
func CollectionArticle(c *gin.Context) (interface{}, error) {
	// 登录验证
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin() {
		return nil, errors.IsNotLogin
	}

	// 	获取参数
	var article dto.CollectionArticleRep
	if err := c.ShouldBind(&article); err != nil{
		return nil, errors.ReceiveParametersError
	}


	var collection models.Collection
	collection.UserID = auth.User.ID
	collection.ArticleID = article.ArticleId

	if !logic.CollectionArticle(&collection){
		return nil, errors.CollectionError

	}

	return dto.CollectionArticleRsq{
		Message: "收藏成功",
	}, nil
}

// CancelCollectionArticle 取消收藏
func CancelCollectionArticle(c *gin.Context) (interface{}, error) {
	// 登录验证
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin() {
		return nil, errors.IsNotLogin
	}

	// 	获取参数
	var article dto.CancelCollectionArticleRep
	if err := c.ShouldBind(&article); err != nil{
		return nil, errors.ReceiveParametersError
	}

	var collection models.Collection
	collection.UserID = auth.User.ID
	collection.ArticleID = article.ArticleId
	if !logic.CancelCollectionArticle(&collection){
		return nil, errors.CollectionError
	}

	return dto.CancelCollectionArticleRsq{
		Message: "取消收藏成功",
	}, nil
}

// CreatLabel 添加标签
func CreatLabel(c *gin.Context) (interface{}, error) {
	// 登录验证
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin() {
		return nil, errors.IsNotLogin
	}

	// 接收数据
	var newLabel dto.CreateLabelRep
	if err := c.ShouldBind(&newLabel); err != nil{
		return nil, errors.ReceiveParametersError
	}
	label := models.Label{
		Name: newLabel.Name,
	}

	var count int
	if err := dao.DB.Model(&models.Label{}).Where("name = ?", label.Name).Count(&count).Error; err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, errors.CreatLabelError
	}

	if err := dao.DB.Create(&label).Error; err != nil {
		return nil, err
	}
	return dto.CreateLabelRsq{
		Message: "创建成功",
	}, nil

}

// SearchLabel 搜索标签
func SearchLabel(c *gin.Context) (interface{}, error) {
	// 登录验证
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin() {
		return nil, errors.IsNotLogin
	}

	// 	获取参数
	var findLabel dto.SearchLabelReq
	if err := c.ShouldBind(&findLabel); err != nil{
		return nil, errors.ReceiveParametersError
	}
	label := models.Label{Name: findLabel.Name}

	var labels []models.Label
	err := dao.DB.Model(&labels).Where("name Like ?", label.Name).Find(&labels).Error
	if err != nil{
		return nil, errors.SearchLabelError
	}

	return dto.SearchLabelRsp{
		Labels: labels,
		Message: "搜索成功",
	}, nil
}

// ChooseLabels 添加文章和 标签的关系
func ChooseLabels(c *gin.Context) (interface{}, error) {
	// 登录验证
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin() {
		return nil, errors.IsNotLogin
	}
	// 	获取参数
	var labelRs dto.LabelReceiveReq
	if err := c.ShouldBind(&labelRs); err != nil{
		return nil, errors.ReceiveParametersError
	}

	// 要检验一下 这个文章是不是本人写的
	if !logic.IsAuthorSelf(labelRs,auth.User.ID){
		return nil, errors.AddLabelError
	}
	return dto.LabelReceiveRsp{
		Message: "为文章添加标签成功",
	}, nil
}

// HomePage 主页
func HomePage(c *gin.Context) (interface{}, error) {
	// 登陆检验
	// 登录验证
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin() {
		return nil, errors.IsNotLogin
	}


	// 获取页数
	var newPage dto.ShowPageReq
	if err := c.ShouldBind(&newPage); err != nil{
		return nil, errors.ReceiveParametersError
	}
	// 通过个人id来获取相关信息(姓名 昵称 自我介绍 等级 性别)
	selfInformation := auth.User

	// 获取公开文章
	t, articlePage := logic.FindTheLatestArticles()
	if !t{
		return nil, errors.ShowPageError
	}

	// 翻页
	// 获取文章后 能知道一共有几页，再按照现在的页面 制作切片 ，需要前端传过来的只有 当前page
	// 设定1面展示 10片文章
	var page dto.Page
	page.PageSize = 10		// 一页展示10片文章
	page.AllPage = len(articlePage)		//  总共有这么多片文章
	beginA := page.NowPage  * 10

	var endA int
	if beginA + 9 > page.AllPage{
		endA = page.AllPage
	}else {
		endA = beginA + 9
	}
	pageArticle := articlePage[beginA : endA]		// pageArticle为 当前页面应该有的文章

	// 返回
	return dto.ShowPageRsp{
		Page:     page,
		Articles: pageArticle,
		SelfPart: *selfInformation,
	}, nil
}
