package handler

import (
	"LlBlog/core"
	"LlBlog/dao"
	"LlBlog/dto"
	"LlBlog/errors"
	"LlBlog/logic"
	"LlBlog/models"
	"LlBlog/utils"
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

// Register 注册✓
func Register(c *gin.Context) (interface{}, error) {

	var req dto.RegisterReq
	if err := c.ShouldBind(&req); err != nil {
		return nil, err
	}

	var count int
	if err := dao.DB.Model(&models.UserInfo{}).Where("customer_name = ?", req.CustomerName).Count(&count).Error; err != nil {
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
func UserInformationUpdate(c *gin.Context) (interface{}, error) {
	// 登录验证
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin() {
		return nil, errors.IsNotLogin
	}

	// 获取更改后的 个人信息
	var userMessage dto.UpdateReq
	err := c.ShouldBind(&userMessage)
	if err != nil {
		return nil, errors.ReceiveParametersError
	}

	// 进行更新
	err = dao.DB.Model(auth.User).Updates(models.UserInfo{
		Name:         userMessage.Name,
		Gender:       userMessage.Gender,
		Introduction: userMessage.Introduction,
		HeadPortrait: userMessage.HeadPortrait,
	}).Error

	user,err := dao.GetUserById(auth.User.ID)

	if err != nil {
		return nil, err
	}

	return dto.UpdateRsp{
		User: user.Clear(),
	}, nil

}

// ShowSelf 展示个人主页
func ShowSelf(c *gin.Context) (interface{}, error) {
	// 登陆检验
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin() {
		return nil, errors.IsNotLogin
	}

	// 获取页数
	var page dto.ShowSelfPageReq
	err := c.ShouldBind(&page)
	if err != nil {
		return nil, errors.ReceiveParametersError
	}

	// 通过个人id来获取相关信息(姓名 昵称 自我介绍 等级 性别)
	userMessage, err := dao.GetUserById(auth.User.ID)
	if err != nil{
		return nil, errors.UserObtainError
	}

	// 获取文章
	t, articlePage := logic.FindAllArticleByUserId(auth.User)
	if !t {
		utils.Return(c, errors.ShowPageError)
		return nil, errors.ShowPageError
	}

	// 翻页
	// 获取文章后 能知道一共有几页，再按照现在的页面 制作切片 ，需要前端传过来的只有 当前page
	// 设定1面展示 10片文章
	var Page dto.Page
	Page.PageSize = 10
	Page.NowPage = page.NowPage
	// 一页展示10片文章
	Page.AllPage = len(articlePage) //  总共有这么多片文章
	beginA := Page.NowPage * 10

	var endA int
	if beginA + 9 > Page.AllPage {
		endA = Page.AllPage
	} else {
		endA = beginA + 9
	}
	pageArticle := articlePage[beginA:endA] // pageArticle为 当前页面应该有的文章

	// 返回
	return dto.ShowSelfPageRsp{
		Page: Page,
		Articles: pageArticle,
		SelfPart: userMessage.Clear(),
	}, nil
}

// PasswordUpdate 修改密码
func PasswordUpdate(c *gin.Context) (interface{}, error) {
	// 登陆检验
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin() {
		return nil, errors.IsNotLogin
	}
	// 接收数据
	var pwd = dto.UpdatePasswordReq{}
	err := c.ShouldBind(&pwd)
	if err != nil {
		return nil, errors.ReceiveParametersError

	}

	var user = auth.User
	// 1: 验证旧密码与新密码是否相等
	m := md5.New()
	m.Write([]byte(pwd.OldPassword))
	OldPsw := hex.EncodeToString(m.Sum(nil))
	if user.Password != OldPsw{
		return nil, errors.UpdatePasswordError
	}
	// 修改密码：
	if !dao.UpdatePassword(auth.User, pwd){
		return nil, errors.UpdatePasswordError
	}

	return dto.UpdatePasswordRsp{
		User: user.Clear(),
	}, nil
}

// SearchUSer 查找用户
func SearchUSer(c *gin.Context) (interface{}, error) {
	// 登录验证
	auth := c.MustGet("auth").(core.AuthAuthorization)
	if !auth.IsLogin() {
		return nil, errors.IsNotLogin
	}

	// 	获取参数
	var information dto.SearchUserReq
	if err := c.ShouldBind(&information); err != nil{
		return nil, errors.ReceiveParametersError
	}

	users, t := logic.SearchUser(information)
	if !t{
		return nil, errors.SearchERROR
	}

	return dto.SearchUserRsp{
		Users: users,
	}, nil
}



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
