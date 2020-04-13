package profitsharing

import "github.com/zozowind/wego/libs/errmsg"

const (
	errProfitSharingSign = &errmsg.ErrMsg{Code: -2010, Message: "分账相关接口签名错误"}
	errProfitSharingReq  = &errmsg.ErrMsg{Code: -2011, Message: "分账相关接口请求错误"}
	errProfitSharingRsp  = &errmsg.ErrMsg{Code: -2012, Message: "分账相关接口请求返回错误"}
)
