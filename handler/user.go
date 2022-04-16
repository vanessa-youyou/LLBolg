package handler

import (
	"LlBlog/core"
	"LlBlog/dao"
	"LlBlog/dto"
	"LlBlog/errors"
	"LlBlog/logic"
	"LlBlog/models"
	"crypto/md5"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

// Login 登录检验✓
func Login(c *gin.Context) (interface{}, error) {
	auth := c.MustGet("auth").(core.AuthAuthorization)

	var req dto.LoginReq
	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}

	// 获取账号
	user, err := dao.GetUserByCustomerName(req.CustomerName)
	if err != nil {
		return nil, err
	}

	// 登陆失败报错
	if !logic.LoginCheck(req.Password, user.Password) {
		return nil, errors.LoginFailed
	}

	auth.SetCookie(c, user.ID)
	return dto.LoginRsp{
		User: user.Clear(),
	}, nil
}

// register 注册✓
func Register(c *gin.Context) (interface{}, error) {

	var req dto.RegisterReq
	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}

	var count int
	if err := dao.DB.Model(&models.UserInfo{}).Where("customer_name = ?").Count(&count).Error; err != nil {
		return nil, err
	}

	if count > 0 {
		return nil, errors.WrongAccountName
	}

	m := md5.New()
	m.Write([]byte(req.Password))

	user := models.UserInfo{
		CustomerName: req.CustomerName,
		Password:     hex.EncodeToString(m.Sum(nil)),
	}
	if err := dao.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return dto.RegisterRsp{
		Id: user.ID,
	}, nil
}

// UserInformationUpdate 修改个人信息✓
//func UserInformationUpdate(c *gin.Context) {
//	// 登录验证
//	auth := c.MustGet("auth").(core.AuthAuthorization)
//	if !auth.IsLogin() {
//		utils.Return(c, errors.IsNotLogin)
//		return
//	}
//	// 获取更改后的 个人信息
//	var userN models.UserInfo
//	err := c.ShouldBind(&userN)
//	if err != nil {
//		utils.Return(c, err)
//		fmt.Println("未接受到传递的信息")
//		return
//	}
//
//	// 获取登录人的id	id CustomerName	Password Label都不能用户自己改
//	userN.ID = auth.User.ID
//	userN.CustomerName = auth.User.CustomerName
//	userN.Password = auth.User.Password
//	userN.Label = auth.User.Label
//
//	// 进行更新
//	if !services.UpdateUserInformation(&userN) {
//		utils.Return(c, errors.WrongUpdate)
//		return
//	}
//
//	// 成功
//	utils.Return(c, gin.H{
//		"message": "修改成功 这里应该跳转页面到 个人页面",
//	})
//}
//
//// ShowSelf 展示个人主页
//func ShowSelf(c *gin.Context) {
//	// 登陆检验
//	auth := c.MustGet("auth").(core.AuthAuthorization)
//	if !auth.IsLogin() {
//		utils.Return(c, errors.IsNotLogin)
//		return
//	}
//
//	// 获取页数
//	var page models.Page
//	err := c.ShouldBind(&page)
//	if err != nil {
//		utils.Return(c, err)
//		fmt.Println("未接受到传递的信息")
//		return
//	}
//
//	// 通过个人id来获取相关信息(姓名 昵称 自我介绍 等级 性别)
//	selfInformation := auth.User
//
//	// 获取文章
//	t, articlePage := services.FindAllArticleByUserId(selfInformation)
//	if !t {
//		fmt.Println("获取文章 error")
//		utils.Return(c, errors.ShowPageError)
//		return
//	}
//
//	// 翻页
//	// 获取文章后 能知道一共有几页，再按照现在的页面 制作切片 ，需要前端传过来的只有 当前page
//	// 设定1面展示 10片文章
//	page.PageSize = 10              // 一页展示10片文章
//	page.PageNum = len(articlePage) //  总共有这么多片文章
//	fmt.Println("len(articlePage) =  ", len(articlePage))
//	beginA := page.PageNow * 10
//	fmt.Println("beginA =  ", beginA)
//
//	var endA int
//	if beginA+9 > page.PageNum {
//		endA = page.PageNum
//		fmt.Println("endA = ", endA)
//	} else {
//		endA = beginA + 9
//		fmt.Println("endA = ", endA)
//
//	}
//	pageArticle := articlePage[beginA:endA] // pageArticle为 当前页面应该有的文章
//
//	// 返回
//	utils.Return(c, gin.H{
//		"pageMassge":  page,
//		"pageArticle": pageArticle,
//		"selfPart":    selfInformation.Clear(),
//		"message":     "获取成功 这里应该在个人页面",
//	})
//}
//
//// Upload 七牛云的上传 用post
//func Upload(c *gin.Context) {
//	// 登陆检验
//	auth := c.MustGet("auth").(core.AuthAuthorization)
//	if !auth.IsLogin() {
//		utils.Return(c, errors.IsNotLogin)
//		return
//	}
//
//	file, fileHeader, _ := c.Request.FormFile("file")
//
//	fileSize := fileHeader.Size
//
//	url, err := services.UploadFile(file, fileSize)
//	if err != nil {
//		utils.Return(c, errors.UploadError)
//		return
//	}
//	// 这里应该 按照操作人id 存入数据库
//	if !services.UpdateHeadPortrait(url, auth.User) {
//		utils.Return(c, errors.UploadError)
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{
//		"message": "ok",
//		"url":     url,
//	})
//}
//
//// PasswordUpdate 修改密码（1:旧密码 2：新密码 3：验证新密码）
//func PasswordUpdate(c *gin.Context) {
//	// 登陆检验
//	auth := c.MustGet("auth").(core.AuthAuthorization)
//	if !auth.IsLogin() {
//		utils.Return(c, errors.IsNotLogin)
//		return
//	}
//	// 接收数据
//	var up = services.UpdatePassword{}
//	err := c.ShouldBind(&up)
//	if err != nil {
//		utils.Return(c, err)
//		fmt.Println("未接受到传递的信息")
//		return
//	}
//
//	var user = auth.User
//	// 1: 验证旧密码与新密码是否相等
//	if !up.CheckPassword(user) {
//		utils.Return(c, errors.UpdatePasswordError)
//		return
//	}
//	utils.Return(c, gin.H{
//		"message": "修改成功 这里应该跳转页面到 登录页面",
//	})
//}
