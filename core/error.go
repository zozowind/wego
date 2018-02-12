package core

import (
	"github.com/zozowind/wego/libs/errmsg"
)

var (
	errSuccess     = &errmsg.ErrMsg{0, "成功", ""}
	errGetToken    = &errmsg.ErrMsg{-100, "获取Token错误", ""}
	errNetwork     = &errmsg.ErrMsg{-101, "网络请求错误", ""}
	errResultParse = &errmsg.ErrMsg{-102, "结果解析错误", ""}
)
