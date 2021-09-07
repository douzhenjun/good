package response

import "github.com/kataras/iris/v12/mvc"

type response struct {
	Code int         `json:"errorno"`
	Data interface{} `json:"data,omitempty"`
	*errMsg
}

type errMsg struct {
	En string `json:"error_msg_en,omitempty"`
	Zh string `json:"error_msg_zh,omitempty"`
}

func NewMsg(en, zh string) *errMsg {
	return &errMsg{en, zh}
}

func (e *errMsg) Error() string {
	return e.En
}

var _ error = (*errMsg)(nil)

func IsMsg(err error) *errMsg {
	if msg, ok := err.(*errMsg); ok {
		return msg
	}
	return nil
}

func Success(data interface{}) mvc.Result {
	return Response(CodeSuccess, data, EmptyMsg)
}

/*
Fail 返回指定的错误信息
*/
func Fail(msg *errMsg) mvc.Result {
	return Response(CodeFail, nil, msg)
}

func FailMsg(en string, zh string) mvc.Result {
	return Fail(NewMsg(en, zh))
}

/*
Error 返回error错误信息
*/
func Error(err error) mvc.Result {
	if msg := IsMsg(err); msg != nil {
		return Fail(msg)
	}
	msg := err.Error()
	return Fail(NewMsg(msg, msg))
}

func Response(code int, data interface{}, msg *errMsg) mvc.Result {
	return mvc.Response{Object: response{code, data, msg}}
}
