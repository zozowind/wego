package app

import (
	"github.com/zozowind/wego/libs/errmsg"
)

var (
	errPayNotifyData      = &errmsg.ErrMsg{Code: -2001, Message: "支付通知数据错误"}
	errPayNotifySignCheck = &errmsg.ErrMsg{Code: -2002, Message: "支付通知签名校验错误"}
	errPayNotifyResult    = &errmsg.ErrMsg{Code: -2003, Message: "支付通知结果错误"}
	errUnifiedOrderReq    = &errmsg.ErrMsg{Code: -2004, Message: "订单请求错误"}
	errUnifiedOrderRsp    = &errmsg.ErrMsg{Code: -2005, Message: "订单请求返回错误"}
	errUnifiedOrderResult = &errmsg.ErrMsg{Code: -2006, Message: "订单生成错误"}
	errPayPackage         = &errmsg.ErrMsg{Code: -2007, Message: "支付包生成错误"}
	errMessageSignCheck   = &errmsg.ErrMsg{Code: -2008, Message: "消息签名校验错误"}
	errMessageMethod      = &errmsg.ErrMsg{Code: -2009, Message: "消息请求方法错误"}
)
