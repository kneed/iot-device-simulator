package app

const (
	Success           = "000000"
	Error             = "500"
	InvalidParams     = "400"
	ObjectNotExist    = "-100002"
	ErrorExistAlready = "-100001"
)


var MsgFlags = map[string]string{
	Success:           "success",
	Error:             "sorry, Something went wrong",
	InvalidParams:     "请求参数错误",
	ErrorExistAlready: "对象已存在",
	ObjectNotExist:    "对象不存在",
}

// GetMsg get error information based on Code
func GetMsg(code string) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return "未定义的错误"
}
