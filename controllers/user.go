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
	"net/http"
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