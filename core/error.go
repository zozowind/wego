package core

import (
	"github.com/zozowind/wego/libs/errmsg"
)

var (
	errSuccess     = &errmsg.ErrMsg{Code: 0, Message: "成功"}
	errGetToken    = &errmsg.ErrMsg{Code: -100, Message: "获取Token错误"}
	errNetwork     = &errmsg.ErrMsg{Code: -101, Message: "网络请求错误"}
	errResultParse = &errmsg.ErrMsg{Code: -102, Message: "结果解析错误"}
)
