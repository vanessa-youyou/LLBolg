package errors

var LoginFailed = New(10001, "登陆失败")
var WrongAccountName = New(10002, "注册失败 账户名已被使用")
var RegisteredFailed = New(10003, "注册失败")
var IsNotLogin = New(10004, "尚未登录")

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
