package errors

var LoginFailed = New(10001, "登陆失败")
var WrongAccountName =New(1002,"注册失败 账户名已被使用")
var RegisteredFailed = New(1003,"注册失败")

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
