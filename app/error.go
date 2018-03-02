package app

import (
	"github.com/zozowind/wego/libs/errmsg"
)

var (
	errPayNotifyData      = &errmsg.ErrMsg{Code: -2001, Message: "支付通知数据错误", Detail: ""}
	errPayNotifySignCheck = &errmsg.ErrMsg{Code: -2002, Message: "支付通知签名校验错误", Detail: ""}
	errPayNotifyResult    = &errmsg.ErrMsg{Code: -2003, Message: "支付通知结果错误", Detail: ""}
	errUnifiedOrderReq    = &errmsg.ErrMsg{Code: -2004, Message: "订单请求错误", Detail: ""}
	errUnifiedOrderRsp    = &errmsg.ErrMsg{Code: -2005, Message: "订单请求返回错误", Detail: ""}
	errUnifiedOrderResult = &errmsg.ErrMsg{Code: -2006, Message: "订单生成错误", Detail: ""}
	errPayPackage         = &errmsg.ErrMsg{Code: -2007, Message: "支付包生成错误", Detail: ""}
)
