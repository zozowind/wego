package core

import (
	"github.com/zozowind/wego/libs/errmsg"
)

var (
	errSuccess     = &errmsg.ErrMsg{Code: 0, Message: "成功", Detail: ""}
	errGetToken    = &errmsg.ErrMsg{Code: -100, Message: "获取Token错误", Detail: ""}
	errNetwork     = &errmsg.ErrMsg{Code: -101, Message: "网络请求错误", Detail: ""}
	errResultParse = &errmsg.ErrMsg{Code: -102, Message: "结果解析错误", Detail: ""}
)
