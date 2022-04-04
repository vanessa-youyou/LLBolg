package controllers

import (
	"LlBlog/core"
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
	// 状态有(1:草稿 draft 2:发布 published 3:发布-审核中 published-review 4:发布成功 5:驳回)
	//这里不允许出现 草稿和发布之外的选项
	if articleN.State != "draft" && articleN.State != "published"{
		// 默认模式为 发布
		articleN.State = "published"
	}

	// 进行保存
	if ! services.NewArticles(articleN){
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
	if article.State != "draft" && article.State != "published"{
		// 默认模式为 发布
		article.State = "published"
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

	if !services.CreatComment(&cm){
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
	fmt.Println("1111111进入controller")
	// 接收数据
	var search models.Search
	err := c.ShouldBind(&search)
	if err != nil {
		utils.Return(c, err)
		fmt.Println("未接受到传递的信息")
		return
	}
	fmt.Println("接收数据完毕")
	fmt.Println(search.SearchWay,"  22 ", search.Content,"  33", search.Content)

	articles, t := services.SearchArticle(search)
	fmt.Println("services完毕")
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
		utils.Return(c, errors.ObtainDetailsError)
		return
	}
	utils.Return(c, gin.H{
		"message": " 查找成功 这里应该还在文章页面",
		"article": articles,
		"comment": comments,
	})
}