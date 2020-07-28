package pay

import (
	"github.com/zozowind/wego/libs/errmsg"
)

var (
	errProfitSharingSign = &errmsg.ErrMsg{Code: -3001, Message: "分账相关接口签名错误"}
	errProfitSharingReq  = &errmsg.ErrMsg{Code: -3002, Message: "分账相关接口请求错误"}
	errProfitSharingRsp  = &errmsg.ErrMsg{Code: -3003, Message: "分账相关接口请求返回错误"}

	errPayOrderSign = &errmsg.ErrMsg{Code: -3004, Message: "支付订单相关接口签名错误"}
	errPayOrderReq  = &errmsg.ErrMsg{Code: -3005, Message: "支付订单相关接口请求错误"}
	errPayOrderRsp  = &errmsg.ErrMsg{Code: -3006, Message: "支付订单相关接口请求返回错误"}

	errPayRequestReq = &errmsg.ErrMsg{Code: -5001, Message: "支付接口请求错误"}
)
