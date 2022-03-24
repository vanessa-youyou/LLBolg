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

// UserLogin 登录检验
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

// UserRegistered 注册
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

// UserInformationUpdate 修改个人信息
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

	// 获取登录人的id
	userN.ID = auth.User.ID
	userN.CustomerName = auth.User.CustomerName

	// 进行更新
	if ! services.UpdateUserInformation(userN){
		utils.Return(c, errors.WrongUpdate)
		return
	}

	// 成功
	utils.Return(c, gin.H{
		"message": "修改成功 这里应该跳转页面到 个人页面",
	})
}

// CreatArticle 写文章
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
	articleN.UserInfoID = auth.User.ID


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

// GiveLike 点赞
func GiveLike(c *gin.Context)  {
	// 登陆检验
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin(){
		utils.Return(c, errors.IsNotLogin)
		return
	}

	// 获取数据
	var ifLike services.Like
	err := c.ShouldBind(&ifLike)
	if err != nil {
		utils.Return(c, err)
		fmt.Println("未接受到传递的信息")
		return
	}

	// 按照id找到这个文章
	article,err := databases.FindArticleById(ifLike.ID)
	if err != nil{
		// 找不到此文章
		utils.Return(c, errors.IsNotLogin)
		return
	}
	// 进行点赞数
	if ifLike.PickIt {
		article.Praise += 1
	}
	if ! services.UpdatArticle(article){
		utils.Return(c, errors.PickError)
		return
	}

	// 成功
	utils.Return(c, gin.H{
		"message": "点赞成功 这里应该还在文章页面",
	})
}

// ModifyArticle 修改文章 只有作者可以
func ModifyArticle(c *gin.Context) {
	// 登陆检验
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin(){
		utils.Return(c, errors.IsNotLogin)
		return
	}
	fmt.Println("在UserInformationUpdate修改个人信息中 已经历过登陆验证 id为",auth.User.ID)
	// 接收数据
	var article models.ArticleInfo
	err := c.ShouldBind(&article)
	if err != nil {
		utils.Return(c, err)
		fmt.Println("未接受到传递的信息")
		return
	}

	// 查证 操作人 是否为文章作者
	var authorId = auth.User.ID
	if ! services.AuthorCheck(authorId, article.ID){
		// 文章作者 和 操作人不匹配
		utils.Return(c, errors.IsNotOneself)
		return
	}

	// 是本人 进行保存
	if ! services.UpdatArticle(&article){
		utils.Return(c, errors.UpdateError)
		return
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
	var articleID1 models.ArticleInfo
	err := c.ShouldBind(&articleID1)
	if err != nil {
		utils.Return(c, err)
		fmt.Println("未接受到传递的信息")
		return
	}
	var articleID = articleID1.ID

	// 查证 操作人 是否为文章作者
	var authorId = auth.User.ID
	if ! services.AuthorCheck(authorId, articleID){
		// 文章作者 和 操作人不匹配
		utils.Return(c, errors.IsNotOneself)
		return
	}

	// 是本人 进行删除
	if ! services.ArticleDelete(articleID){
		utils.Return(c, errors.DeleteArticleError)
		return
	}
	// 成功
	utils.Return(c, gin.H{
		"message": "删除成功 这里应该回到当前页面",
	})
}