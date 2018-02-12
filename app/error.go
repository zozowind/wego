package app

import (
	"github.com/zozowind/wego/libs/errmsg"
)

var (
	errPayNotifyData      = &errmsg.ErrMsg{-2001, "支付通知数据错误", ""}
	errPayNotifySignCheck = &errmsg.ErrMsg{-2002, "支付通知签名校验错误", ""}
	errPayNotifyResult    = &errmsg.ErrMsg{-2003, "支付通知结果错误", ""}
	errUnifiedOrderReq    = &errmsg.ErrMsg{-2004, "订单请求错误", ""}
	errUnifiedOrderRsp    = &errmsg.ErrMsg{-2005, "订单请求返回错误", ""}
	errUnifiedOrderResult = &errmsg.ErrMsg{-2006, "订单生成错误", ""}
	errPayPackage         = &errmsg.ErrMsg{-2007, "支付包生成错误", ""}
)
