package work

import (
	"github.com/zozowind/wego/libs/errmsg"
)

var (
	errSendMessage   = &errmsg.ErrMsg{Code: -4001, Message: "发送消息错误"}
	errUploadMedia   = &errmsg.ErrMsg{Code: -4002, Message: "素材文件上传错误"}
	errDownloadMedia = &errmsg.ErrMsg{Code: -4003, Message: "素材文件下载错误"}
)
