package core

import (
	"LlBlog/databases"
	"LlBlog/errors"
	"LlBlog/models"
	"LlBlog/utils/hash"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	LLBlogCookieName = "LLBolg-Auth"
	CookieExpires    = 60 * 60 * 24
)

type cookieInfo struct {
	AccountId  int   `json:"account_id"`
	ExpireTime int64 `json:"expire_time"`
}

type AuthAuthorization struct {
	User    *models.UserInfo
	isLogin bool
}

func (r *AuthAuthorization) CheckLogin() error {
	if !r.isLogin {
		return errors.IsNotLogin
	}
	return nil
}

func (r *AuthAuthorization) IsLogin() bool {
	return r.isLogin
}

func (r *AuthAuthorization) LoadAuthenticationInfo(c *gin.Context) {
	// 阻止解码方法异常传递
	defer func() {
		recover()
	}()

	payload, err := c.Cookie(LLBlogCookieName)
	if err != nil {
		return
	}

	var cookie cookieInfo
	hash.DecodeToken(payload, &cookie)
	if cookie.ExpireTime <= time.Now().Unix() {
		// 过期了
		return
	}

	// 从数据库捞出来
	if ok := r.fetchAccount(cookie.AccountId); ok {
		r.isLogin = true
	}
}

// 设置cookie
func (r *AuthAuthorization) SetCookie(c *gin.Context, aid uint) {
	if aid == 0 {
		c.SetCookie(LLBlogCookieName, "", 99999, "/", "", false, false)
		return
	}
	// 生成加密串
	payload := GenerateToken(int(aid), CookieExpires)
	c.SetCookie(LLBlogCookieName, payload, 99999, "/", "", false, false)
}

// 从数据库查找该用户
func (r *AuthAuthorization) fetchAccount(aid int) bool {
	user, err := databases.GetUserById(uint(aid))
	if err != nil {
		return false
	}
	r.User = user
	return true
}

// 生成token
func GenerateToken(aid int, expire int64) string {
	payload := cookieInfo{
		AccountId:  aid,
		ExpireTime: expire + time.Now().Unix(),
	}
	return hash.GenerateToken(payload, true)
}



