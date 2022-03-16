package errors

var LoginFailed = New(10001, "登陆失败")

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
