package errors

// 用户模块错误	1000~1999

var LoginFailed = New(10001, "登陆失败")
var WrongAccountName = New(10002, "注册失败 账户名已被使用")
var RegisteredFailed = New(10003, "注册失败")
var IsNotLogin = New(10004, "尚未登录")
var WrongUpdate = New(1005,"修改失败")


// 文章模块	2000~2999

var WriteError = New(2001,"写入失败")
var PickError = New(2002,"点赞失败")
var UpdateError = New(2003,"修改失败")
var DeleteArticleError = New(2004,"删除文章失败")
var IsNotOneself = New(2005,"不是本人在操作")


func New(code int, msg string) ErrorBase {
	return ErrorBase{
		Code: code,
		Msg:  msg,
	}
}

type ErrorBase struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (b ErrorBase) Error() string {
	return b.Msg
}
