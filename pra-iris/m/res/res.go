package res

type Res struct {
	Code         int
	Error_msg_en error
	Error_msg_zh error
}

func MyRes(code int, error_msg_en error, error_msg_zh error) Res{
	return Res{code, error_msg_en, error_msg_zh}
}
