package controllers

import (
	"LlBlog/core"
	"LlBlog/databases"
	"LlBlog/errors"
	"LlBlog/models"
	"LlBlog/services"
	"LlBlog/utils"
	"fmt"
	"github.com/gin-gonic/gin"
)

// CreatArticle 写文章✓
func CreatArticle(c *gin.Context){
	// 登陆检验
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin(){
		utils.Return(c, errors.IsNotLogin)
		return
	}
	// 获取数据
	var articleN models.ArticleInfo
	err := c.ShouldBind(&articleN)
	if err != nil {
		utils.Return(c, err)
		fmt.Println("未接受到传递的信息")
		return
	}

	// 绑定正确作者id
	articleN.AuthorID = auth.User.ID
	// 状态有(1:草稿  2:发布  3:发布-审核中 4:发布成功 5:驳回 6:撤销)
	//这里不允许出现撤销 草稿和发布之外的选项
	if articleN.State != 1 && articleN.State != 2 && articleN.State != 6{
		// 默认模式为 发布
		articleN.State = 2
	}

	// 进行保存
	t := databases.WriteNewArticles(&articleN)

	if !t{
		utils.Return(c, errors.WriteError)
		return
	}

	// 成功
	utils.Return(c, gin.H{
		"message": "写入成功 这里应该跳转页面到 个人页面",
	})
}

// ModifyArticle 修改文章 只有作者可以(修改完成)✓
func ModifyArticle(c *gin.Context) {
	// 登陆检验
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin(){
		utils.Return(c, errors.IsNotLogin)
		return
	}
	// 接收数据
	var article models.ArticleInfo
	err := c.ShouldBind(&article)
	if err != nil {
		utils.Return(c, err)
		fmt.Println("未接受到传递的信息")
		return
	}

	// 查证 操作人 是否为文章作者
	var userId = auth.User.ID

	// 处理 状态参数 （只能为发布/草稿 默认:发布）
	if article.State != 1 && article.State != 2 && article.State != 6{
		// 默认模式为 发布
		article.State = 2
	}

	// 新的查找and 修改 逻辑是(update 文章 where authorID == article.authorID)
	if ! services.ArticleModify(&article, userId){
		utils.Return(c, errors.UpdateError)
		//	return
	}

	// 成功
	utils.Return(c, gin.H{
		"message": "修改成功 这里应该还在文章页面",
	})
}

// DeleteArticle 删除文章（评论还没写 评论应该一起被删）只允许作者本人删除
func DeleteArticle(c *gin.Context)  {
	// 登陆检验
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin(){
		utils.Return(c, errors.IsNotLogin)
		return
	}
	// 接收数据
	var article models.ArticleInfo
	err := c.ShouldBind(&article)
	if err != nil {
		utils.Return(c, err)
		fmt.Println("未接受到传递的信息")
		return
	}

	// 获得操作人的id
	var userId = auth.User.ID

	// 1:检查操作者是否为 作者本人，是则删除，不是则报错
	if ! services.RemoveArticle(&article, userId){
		utils.Return(c, errors.DeleteArticleError)
		return
	}

	utils.Return(c, gin.H{
		"message": "删除成功 这里应该回到当前页面",
	})
}

// WriteComment 新建评论
func WriteComment(c *gin.Context)  {
	// 登陆检验
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin(){
		utils.Return(c, errors.IsNotLogin)
		return
	}

	// 接收数据
	var cm models.CommentInfo
	err := c.ShouldBind(&cm)
	if err != nil {
		utils.Return(c, err)
		fmt.Println("未接受到传递的信息")
		return
	}

	// 修改操作人id
	cm.UserID = auth.User.ID

	t := databases.NewComment(&cm)
	if !t{
		utils.Return(c, errors.CreatCommentError)
		return
	}

	utils.Return(c, gin.H{
		"message": "评论成功 这里应该跳转页面到 文章页面",
	})
}

// DeleteComment 删除评论(只允许作者本人删除
func DeleteComment(c *gin.Context)  {
	// 登陆检验
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin(){
		utils.Return(c, errors.IsNotLogin)
		return
	}

	// 接收数据
	var cm models.CommentInfo
	err := c.ShouldBind(&cm)
	if err != nil {
		utils.Return(c, err)
		fmt.Println("未接受到传递的信息")
		return
	}

	var userId = auth.User.ID
	if !services.RemoveComment(&cm, userId){
		utils.Return(c, errors.DeleteCommentError)
		return
	}

	// 成功
	utils.Return(c, gin.H{
		"message": "成功删除评论 这里应该还在文章页面",
	})
}

// PickArticle 点赞文章 用redis存储
func PickArticle(c *gin.Context)  {
	// 登陆检验
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin(){
		utils.Return(c, errors.IsNotLogin)
		return
	}
	// 接收数据
	var article models.ArticleInfo
	err := c.ShouldBind(&article)
	if err != nil {
		utils.Return(c, err)
		fmt.Println("未接受到传递的信息")
		return
	}
	// 操作者id
	var userId = auth.User.ID

	if !services.ArticleLike(&article, userId){
		utils.Return(c, errors.PickError)
		return
	}

	// 成功
	utils.Return(c, gin.H{
		"message": "点赞/取消点赞 成功 这里应该还在文章页面",
	})

}

// CommentPick 点赞评论 用redis存储
func CommentPick(c *gin.Context)  {
	// 登陆检验
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin(){
		utils.Return(c, errors.IsNotLogin)
		return
	}
	// 接收数据
	var cm models.CommentInfo
	err := c.ShouldBind(&cm)
	if err != nil {
		utils.Return(c, err)
		fmt.Println("未接受到传递的信息")
		return
	}
	// 操作者id
	var userId = auth.User.ID
	if !services.CommentLike(&cm, userId){
		utils.Return(c, errors.PickError)
		return
	}

	// 成功
	utils.Return(c, gin.H{
		"message": "点赞/取消点赞 成功 这里应该还在文章页面",
	})
}

// SearchArticles 查询文章
func SearchArticles(c *gin.Context)  {
	// 登陆检验
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin(){
		utils.Return(c, errors.IsNotLogin)
		return
	}
	// 接收数据
	var search models.Search
	err := c.ShouldBind(&search)
	if err != nil {
		utils.Return(c, err)
		return
	}

	articles, t := services.SearchArticle(search)
	if !t{
		utils.Return(c, errors.SearchERROR)
		return
	}
	// 一起返回
	utils.Return(c, gin.H{
		"message": " 查找成功 这里应该还在文章页面",
		"articles": articles,
	})
}

// SearchUSer 查找用户
func SearchUSer(c *gin.Context){
	// 登陆检验
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin(){
		utils.Return(c, errors.IsNotLogin)
		return
	}

	// 	获取参数
	var information models.Search
	err := c.ShouldBind(&information)
	if err != nil {
		utils.Return(c, err)
		fmt.Println("未接受到传递的信息")
		return
	}
	users, t := services.SearchUser(information)
	if !t{
		utils.Return(c, errors.SearchERROR)
		return
	}

	utils.Return(c, gin.H{
		"message": " 查找成功 这里应该还在文章页面",
		"Users": users,
	})
}

// ArticleDetails 文章详情
func ArticleDetails(c *gin.Context)  {
	// 登陆检验
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin(){
		utils.Return(c, errors.IsNotLogin)
		return
	}

	// 	获取参数
	var article models.ArticleInfo
	err := c.ShouldBind(&article)
	if err != nil {
		utils.Return(c, err)
		fmt.Println("未接受到传递的信息")
		return
	}

	t, articles, comments := services.ArticleDetails(&article)
	if !t{
		fmt.Println("出错拉！错误在service里")
		utils.Return(c, errors.ObtainDetailsError)
		return
	}
	utils.Return(c, gin.H{
		"message": " 查找成功 这里应该还在文章页面",
		"article": articles,
		"comment": comments,
	})

}

// CollectionArticle 收藏
func CollectionArticle(c *gin.Context){
	// 登陆检验
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin(){
		utils.Return(c, errors.IsNotLogin)
		return
	}

	// 	获取参数
	var article models.ArticleInfo
	err := c.ShouldBind(&article)
	if err != nil {
		utils.Return(c, err)
		fmt.Println("未接受到传递的信息")
		return
	}

	var collection models.Collection
	collection.UserID = auth.User.ID
	collection.ArticleID = article.ID

	if !services.CollectionArticle(&collection){
		utils.Return(c, errors.CollectionError)
		return
	}
	utils.Return(c, gin.H{
		"message": " 收藏成功 这里应该还在文章页面",

	})
}

// CancelCollectionArticle 取消收藏
func CancelCollectionArticle(c *gin.Context){
	// 登陆检验
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin(){
		utils.Return(c, errors.IsNotLogin)
		return
	}

	// 	获取参数
	var article models.ArticleInfo
	err := c.ShouldBind(&article)
	if err != nil {
		utils.Return(c, err)
		fmt.Println("未接受到传递的信息")
		return
	}

	var collection models.Collection
	collection.UserID = auth.User.ID
	collection.ArticleID = article.ID
	if !services.CancelCollectionArticle(&collection){
		utils.Return(c, errors.CollectionError)
		return
	}

	utils.Return(c, gin.H{
		"message": " 取消成功 这里应该还在文章页面",
	})
}

// CreatLabel 添加标签
func CreatLabel(c *gin.Context)  {
	// 登陆检验
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin(){
		utils.Return(c, errors.IsNotLogin)
		return
	}
	// 	获取参数
	var label models.Label
	err := c.ShouldBind(&label)
	if err != nil {
		utils.Return(c, err)
		fmt.Println("未接受到传递的信息")
		return
	}

	if !databases.CreatLabel(&label){
		utils.Return(c, errors.CreatLabelError)
		return
	}

	utils.Return(c, gin.H{
		"message": " 创建标签成功 这里应该在文章页面",
	})

}

// SearchLabel 搜索标签
func SearchLabel(c *gin.Context)  {
	// 登陆检验
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin(){
		utils.Return(c, errors.IsNotLogin)
		return
	}
	// 	获取参数
	var label models.Label
	err := c.ShouldBind(&label)
	if err != nil {
		utils.Return(c, err)
		fmt.Println("未接受到传递的信息")
		return
	}

	t, labels := databases.SearchLabel(&label)
	if !t{
		utils.Return(c, errors.CreatLabelError)
		return
	}

	utils.Return(c, gin.H{
		"Labels": labels,
		"message": "搜索成功！",
	})
}

// ChooseLabels 添加文章和 标签的关系
func ChooseLabels(c *gin.Context)  {
	// 登陆检验
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin(){
		utils.Return(c, errors.IsNotLogin)
		return
	}
	// 	获取参数
	var labelRs models.LabelReceive
	var err error
	if err = c.ShouldBind(&labelRs); err != nil {
		utils.Return(c, err)
		fmt.Println("未接受到传递的信息")
		return
	}

	// 要检验一下 这个文章是不是本人写的
	if !services.IsAuthorSelf(labelRs,auth.User.ID){
		fmt.Println("IsAuthorSelf(labelRs,auth.User.ID){")
		utils.Return(c, err)
		return
	}
	utils.Return(c, gin.H{
		"message": "添加成功！",
	})
}

// HomePage 主页
func HomePage(c *gin.Context)  {
	// 登陆检验
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin(){
		utils.Return(c, errors.IsNotLogin)
		return
	}

	// 获取页数
	var page models.Page
	err := c.ShouldBind(&page)
	if err != nil {
		utils.Return(c, err)
		fmt.Println("未接受到传递的信息")
		return
	}

	// 通过个人id来获取相关信息(姓名 昵称 自我介绍 等级 性别)
	selfInformation := auth.User

	// 获取公开文章
	t, articlePage := services.FindTheLatestArticles()
	if !t{
		fmt.Println("获取文章 error")
		utils.Return(c, errors.ShowPageError)
		return
	}

	// 翻页
	// 获取文章后 能知道一共有几页，再按照现在的页面 制作切片 ，需要前端传过来的只有 当前page
	// 设定1面展示 10片文章
	page.PageSize = 10		// 一页展示10片文章
	page.PageNum = len(articlePage)		//  总共有这么多片文章
	beginA := (page.PageNow - 1)  * 10

	var endA int
	if beginA + 9 > page.PageNum{
		endA = page.PageNum
	}else {
		endA = beginA + 9
	}
	pageArticle := articlePage[beginA : endA]		// pageArticle为 当前页面应该有的文章

	// 返回
	utils.Return(c, gin.H{
		"pageMessage" : page,
		"pageArticle" : pageArticle,
		"selfPart" : selfInformation.Clear(),
		"message": "获取成功 这里应该在个人页面",
	})
}
