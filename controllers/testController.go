package controllers

import (
	"LlBlog/core"
	"LlBlog/databases"
	"LlBlog/errors"
	"LlBlog/models"
	"LlBlog/services"
	"LlBlog/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserLogin 登录检验✓
func UserLogin(c *gin.Context) {
	auth := c.MustGet("auth").(core.AuthAuthorization)
	// 数据库那边的操作(要接收的啊承诺书的结构体)
	var userG services.LoginUser

	// 接收数据
	err := c.ShouldBind(&userG)
	if err != nil {
		utils.Return(c, err)
		return
	}

	user, err := databases.GetUserByCustomerName(userG.CustomerName)
	if err != nil {
		// 其实这里跟上面哪个一样 不应该直接把系统错误显示给客户端，不过开发就随便啦
		utils.Return(c, err)
		return
	}

	// 登陆失败报错
	if !userG.LoginCheck(user) {
		utils.Return(c, errors.LoginFailed)
		return
	}

	// 登陆成功
	// 设置cookie
	auth.SetCookie(c, user.ID)
	utils.Return(c, gin.H{
		"user": user.Clear(),
	})
}

// UserRegistered 注册✓
func UserRegistered(c *gin.Context) {
	// add
	var userN models.UserInfo
	// 接收数据
	err := c.ShouldBind(&userN)
	if err != nil {
		utils.Return(c, err)
		return
	}

	// 进行一个查找的动作 看看 账户名字有没有重复
	if !services.RegisteredNameCheck(userN) {
		utils.Return(c, errors.WrongAccountName)
		return
	}

	// 开始add账号进数据库
	if !services.AddAccount(userN) {
		utils.Return(c, errors.RegisteredFailed)
		return
	}

	// 注册成功
	utils.Return(c, gin.H{
		"message": "注册成功 这里应该跳转页面到 个人页面",
	})
}

// UserInformationUpdate 修改个人信息✓
func UserInformationUpdate(c *gin.Context){
	// 登录验证
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin(){
		utils.Return(c, errors.IsNotLogin)
		return
	}
	// 获取更改后的 个人信息
	var userN models.UserInfo
	err := c.ShouldBind(&userN)
	if err != nil {
		utils.Return(c, err)
		fmt.Println("未接受到传递的信息")
		return
	}

	// 获取登录人的id	id CustomerName	Password Label都不能用户自己改
	userN.ID = auth.User.ID
	userN.CustomerName = auth.User.CustomerName
	userN.Password = auth.User.Password
	userN.Label = auth.User.Label

	// 进行更新
	if ! services.UpdateUserInformation(&userN){
		utils.Return(c, errors.WrongUpdate)
		return
	}

	// 成功
	utils.Return(c, gin.H{
		"message": "修改成功 这里应该跳转页面到 个人页面",
	})
}

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
	articleN.AuthorID = auth.User.ID
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

// PasswordUpdate 修改密码（1:旧密码 2：新密码 3：验证新密码）
func PasswordUpdate(c *gin.Context)  {
	// 登陆检验
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin(){
		utils.Return(c, errors.IsNotLogin)
		return
	}
	// 接收数据
	var up = services.UpdatePassword{}
	err := c.ShouldBind(&up)
	if err != nil {
		utils.Return(c, err)
		fmt.Println("未接受到传递的信息")
		return
	}

	var user = auth.User
	// 1: 验证旧密码与新密码是否相等
	if ! up.CheckPassword(user){
		utils.Return(c, errors.UpdatePasswordError)
		return
	}
	utils.Return(c, gin.H{
		"message": "修改成功 这里应该跳转页面到 登录页面",
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

// ShowSelf 展示个人主页
func ShowSelf(c *gin.Context)  {
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

	// 获取文章
	t, articlePage := services.FindAllArticleByUserId(selfInformation)
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
	fmt.Println("len(articlePage) =  " , len(articlePage))
	beginA := page.PageNow  * 10
	fmt.Println("beginA =  " , beginA)

	var endA int
	if beginA + 9 > page.PageNum{
		endA = page.PageNum
		fmt.Println("endA = ", endA)
	}else {
		endA = beginA + 9
		fmt.Println("endA = ", endA)

	}
	pageArticle := articlePage[beginA : endA]		// pageArticle为 当前页面应该有的文章

	// 返回
	utils.Return(c, gin.H{
		"pageMassge" : page,
		"pageArticle" : pageArticle,
		"selfPart" : selfInformation.Clear(),
		"message": "获取成功 这里应该在个人页面",
	})
}

// Upload 七牛云的上传 用post
func Upload(c *gin.Context)  {
	// 登陆检验
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin(){
		utils.Return(c, errors.IsNotLogin)
		return
	}

	file,fileHeader,_ := c.Request.FormFile("file")

	fileSize := fileHeader.Size

	url, err := services.UploadFile(file, fileSize)
	if err != nil{
		utils.Return(c, errors.UploadError)
		return
	}
	// 这里应该 按照操作人id 存入数据库
	if !services.UpdateHeadPortrait(url, auth.User){
		utils.Return(c, errors.UploadError)
		return
	}

	c.JSON(http.StatusOK,gin.H{
		"message":"ok",
		"url":url,
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