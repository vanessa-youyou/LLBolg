package controllers

import (
	"LlBlog/errors"
	"LlBlog/models"
	"LlBlog/services"
	"LlBlog/utils"
	"fmt"
	"github.com/gin-gonic/gin"
)

// UserLogin 登录检验
func UserLogin(c *gin.Context) {
	// 数据库那边的操作(要接收的啊承诺书的结构体)
	var userG services.LoginUser

	// 接收数据
	err := c.ShouldBind(&userG)
	if err != nil {
		utils.Return(c, err)
		return
	}

	fmt.Printf("%#v\n", userG)
	// 进行一个查找的动作

	// 登陆失败报错
	if !userG.LoginCheck() {
		utils.Return(c, errors.LoginFailed)
		return
	}
	// 登陆成功
	utils.Return(c, gin.H{
		"message": "登陆成功 这里应该跳转页面到 个人页面",
	})
}

func UserRegistered(c *gin.Context)  {
	// add
	var userN models.UserInfo
	fmt.Println("进入了controller1")
	// 接收数据
	err := c.ShouldBind(&userN)
	if err != nil {
		fmt.Printf("%#v\n", userN)
		utils.Return(c, err)
		return
	}
	fmt.Printf("%#v\n", userN)
	// 进行一个查找的动作 看看 账户名字有没有重复
	if !services.RegisteredNameCheck(userN){
		utils.Return(c, errors.WrongAccountName)
		return
	}
	// 开始add账号进数据库
	if !services.AddAccount(userN){
		utils.Return(c, errors.RegisteredFailed)
		return
	}

	// 注册成功
	utils.Return(c, gin.H{
		"message": "注册成功 这里应该跳转页面到 个人页面",
	})

}
