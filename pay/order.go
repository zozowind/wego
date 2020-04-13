package pay

import (
	"encoding/xml"
	"fmt"

	"github.com/zozowind/wego/core"
	"github.com/zozowind/wego/libs/errmsg"
	"github.com/zozowind/wego/util"
)

const (
	// queryPayOrderURL 查询订单URL
	queryPayOrderURL = core.WxPayURL + "/pay/orderquery"
	// closePayOrderURL 关闭订单URL
	closePayOrderURL = core.WxPayURL + "/pay/closeorder"
	// refundPayOrderURL 申请退款URL
	refundPayOrderURL = core.WxPayURL + "/secapi/pay/refund"
	// queryRefundPayOrderURL 查询退款URL
	queryRefundPayOrderURL = core.WxPayURL + "/pay/refundquery"
)

// OrderQueryRequest 订单查询请求
type OrderQueryRequest struct {
	// CommonRequestMsg
	MchID string `xml:"mch_id"` // 微信支付分配的商户号
	// NonceStr string `xml:"nonce_str"` // 随机字符串
	// SignType string `xml:"sign_type"` // 签名类型：目前只支持HMAC-SHA256
	AppID string `xml:"appid"` // 小程序ID

	// 下面这些参数至少提供一个
	TransactionID string `xml:"transaction_id"` // 微信的订单号，优先使用
	OutTradeNo    string `xml:"out_trade_no"`   // 商户系统内部的订单号，当没提供transaction_id时需要传这个。

	// 可选参数
	NonceStr string `xml:"nonce_str"` // 随机字符串，不长于32位。NOTE: 如果为空则系统会自动生成一个随机字符串。
	SignType string `xml:"sign_type"` // 签名类型，目前支持HMAC-SHA256和MD5，默认为MD5
}

// OrderQueryResponse 订单查询返回
type OrderQueryResponse struct {
	CommonResponseMsg

	// 必选返回
	TradeState     string `xml:"trade_state"`      // 交易状态
	TradeStateDesc string `xml:"trade_state_desc"` // 对当前查询订单状态的描述和下一步操作的指引
	OpenID         string `xml:"openid"`           // 用户在商户appid下的唯一标识
	TransactionID  string `xml:"transaction_id"`   // 微信支付订单号
	OutTradeNo     string `xml:"out_trade_no"`     // 商户系统的订单号，与请求一致。
	TradeType      string `xml:"trade_type"`       // 调用接口提交的交易类型，取值如下：JSAPI，NATIVE，APP，MICROPAY，详细说明见参数规定
	BankType       string `xml:"bank_type"`        // 银行类型，采用字符串类型的银行标识
	TotalFee       int64  `xml:"total_fee"`        // 订单总金额，单位为分
	CashFee        int64  `xml:"cash_fee"`         // 现金支付金额订单现金支付金额，详见支付金额
	TimeEnd        string `xml:"time_end"`         // 订单支付时间，格式为yyyyMMddHHmmss，如2009年12月25日9点10分10秒表示为20091225091010。其他详见时间规则

	// 下面字段都是可选返回的(详细见微信支付文档), 为空值表示没有返回, 程序逻辑里需要判断
	DeviceInfo         string `xml:"device_info"`          // 微信支付分配的终端设备号
	IsSubscribe        *bool  `xml:"is_subscribe"`         // 用户是否关注公众账号
	SubOpenID          string `xml:"sub_openid"`           // 用户在子商户appid下的唯一标识
	SubIsSubscribe     *bool  `xml:"sub_is_subscribe"`     // 用户是否关注子公众账号
	SettlementTotalFee *int64 `xml:"settlement_total_fee"` // 应结订单金额=订单金额-非充值代金券金额，应结订单金额<=订单金额。
	FeeType            string `xml:"fee_type"`             // 货币类型，符合ISO 4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型
	CashFeeType        string `xml:"cash_fee_type"`        // 货币类型，符合ISO 4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型
	Detail             string `xml:"detail"`               // 商品详情
	Attach             string `xml:"attach"`               // 附加数据，原样返回
}

// CloseOrderRequest 关闭订单请求
type CloseOrderRequest struct {
	// CommonRequestMsg
	MchID    string `xml:"mch_id"`    // 微信支付分配的商户号
	NonceStr string `xml:"nonce_str"` // 随机字符串
	SignType string `xml:"sign_type"` // 签名类型：目前只支持HMAC-SHA256
	AppID    string `xml:"appid"`     // 小程序ID

	// 必选参数
	OutTradeNo string `xml:"out_trade_no"` // 商户系统内部订单号
}

// CloseOrderResponse 关闭订单返回
type CloseOrderResponse struct {
	CommonResponseMsg
	ResultMsg string `xml:"result_msg"` // 业务结果描述

	AppID    string `xml:"appid"`     // 小程序ID
	MchID    string `xml:"mch_id"`    // 商户号
	NonceStr string `xml:"nonce_str"` // 随机字符串
	Sign     string `xml:"sign"`      // 签名
}

// RefundPayOrderRequest 申请退款请求
type RefundPayOrderRequest struct {
	// CommonRequestMsg
	MchID    string `xml:"mch_id"`    // 微信支付分配的商户号
	NonceStr string `xml:"nonce_str"` // 随机字符串
	SignType string `xml:"sign_type"` // 签名类型：目前只支持HMAC-SHA256
	AppID    string `xml:"appid"`     // 小程序ID

	// 必选参数, TransactionID 和 OutTradeNo 二选一即可.
	TransactionID string `xml:"transaction_id"` // 微信生成的订单号，在支付通知中有返回
	OutTradeNo    string `xml:"out_trade_no"`   // 商户侧传给微信的订单号
	OutRefundNo   string `xml:"out_refund_no"`  // 商户系统内部的退款单号，商户系统内部唯一，同一退款单号多次请求只退一笔
	TotalFee      int64  `xml:"total_fee"`      // 订单总金额，单位为分，只能为整数，详见支付金额
	RefundFee     int64  `xml:"refund_fee"`     // 退款总金额，订单总金额，单位为分，只能为整数，详见支付金额
	// 可选参数
	RefundFeeType string `xml:"refund_fee_type"` // 货币类型，符合ISO 4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型
	RefundDesc    string `xml:"refund_desc"`     // 若商户传入，会在下发给用户的退款消息中体现退款原因
	RefundAccount string `xml:"refund_account"`  // 退款资金来源
	NotifyURL     string `xml:"notify_url"`      // 退款结果通知URL
}

// RefundPayOrderResponse 申请退款返回
type RefundPayOrderResponse struct {
	CommonResponseMsg

	// 必选返回
	AppID         string `xml:"appid"`          // 小程序ID
	MchID         string `xml:"mch_id"`         // 商户号
	NonceStr      string `xml:"nonce_str"`      // 随机字符串
	Sign          string `xml:"sign"`           // 签名
	TransactionID string `xml:"transaction_id"` // 微信订单号
	OutTradeNo    string `xml:"out_trade_no"`   // 商户系统内部的订单号
	OutRefundNo   string `xml:"out_refund_no"`  // 商户退款单号
	RefundID      string `xml:"refund_id"`      // 微信退款单号
	RefundFee     int64  `xml:"refund_fee"`     // 退款总金额,单位为分,可以做部分退款
	TotalFee      int64  `xml:"total_fee"`      // 订单总金额，单位为分，只能为整数，详见支付金额
	CashFee       int64  `xml:"cash_fee"`       // 现金支付金额，单位为分，只能为整数，详见支付金额

	// 下面字段都是可选返回的(详细见微信支付文档), 为空值表示没有返回, 程序逻辑里需要判断
	SettlementRefundFee *int64 `xml:"settlement_refund_fee"` // 退款金额=申请退款金额-非充值代金券退款金额，退款金额<=申请退款金额
	SettlementTotalFee  *int64 `xml:"settlement_total_fee"`  // 应结订单金额=订单金额-非充值代金券金额，应结订单金额<=订单金额。
	FeeType             string `xml:"fee_type"`              // 订单金额货币类型，符合ISO 4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型
	CashFeeType         string `xml:"cash_fee_type"`         // 货币类型，符合ISO 4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型
	CashRefundFee       *int64 `xml:"cash_refund_fee"`       // 现金退款金额，单位为分，只能为整数，详见支付金额
}

// QueryRefundPayOrderRequest 查询退款请求
type QueryRefundPayOrderRequest struct {
	// CommonRequestMsg
	MchID    string `xml:"mch_id"`    // 微信支付分配的商户号
	NonceStr string `xml:"nonce_str"` // 随机字符串
	SignType string `xml:"sign_type"` // 签名类型：目前只支持HMAC-SHA256
	AppID    string `xml:"appid"`     // 小程序ID

	// 必选参数, 四选一
	TransactionID string `xml:"transaction_id"` // 微信订单号
	OutTradeNo    string `xml:"out_trade_no"`   // 商户订单号
	OutRefundNo   string `xml:"out_refund_no"`  // 商户退款单号
	RefundID      string `xml:"refund_id"`      // 微信退款单号
}

// QueryRefundPayOrderResponse 查询退款返回
type QueryRefundPayOrderResponse struct {
	CommonResponseMsg

	// 必选返回
	TransactionID string `xml:"transaction_id"` // 微信订单号
	OutTradeNo    string `xml:"out_trade_no"`   // 商户系统内部的订单号
	TotalFee      int64  `xml:"total_fee"`      // 订单总金额，单位为分，只能为整数，详见支付金额
	CashFee       int64  `xml:"cash_fee"`       // 现金支付金额，单位为分，只能为整数，详见支付金额
	RefundCount   int    `xml:"refund_count"`   // 退款笔数

	// 下面字段都是可选返回的(详细见微信支付文档), 为空值表示没有返回, 程序逻辑里需要判断
	SettlementTotalFee *int64 `xml:"settlement_total_fee"` // 应结订单金额=订单金额-非充值代金券金额，应结订单金额<=订单金额。
	FeeType            string `xml:"fee_type"`             // 订单金额货币类型，符合ISO 4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型
	CashFeeType        string `xml:"cash_fee_type"`        // 现金支付货币类型
}

// QueryPayOrder 查询支付订单
func QueryPayOrder(wb *core.WeBase, req *OrderQueryRequest) (rsp *OrderQueryResponse, err error) {
	// req.CommonRequestMsg = CommonRequestMsg{
	// 	RequestMsg: RequestMsg{
	// 		MchID:    wb.PayID,
	// 		NonceStr: util.RandString(32),
	// 		SignType: util.SignTypeHMACSHA256,
	// 	},
	// 	AppID: wb.AppID,
	// }
	req.MchID = wb.PayID
	req.NonceStr = util.RandString(32)
	req.SignType = util.SignTypeHMACSHA256
	req.AppID = wb.AppID

	params, err := util.StructToURLValue(req, "xml")
	if nil != err {
		return nil, errmsg.GetError(errPayOrderSign, err.Error())
	}
	params.Add("sign", util.SignSha256(params.Encode(), wb.PayKey))

	data, err := util.HTTPXMLPost(wb.HTTPClient, queryPayOrderURL, params)
	if err != nil {
		return nil, errmsg.GetError(errPayOrderRsp, fmt.Sprintf("queryPayOrderURL rsp err:%s", err.Error()))
	}
	rsp = &OrderQueryResponse{}
	err = xml.Unmarshal(data, &rsp)
	if err != nil {
		return nil, errmsg.GetError(errPayOrderRsp, fmt.Sprintf("queryPayOrderURL xml.Unmarshal err:%s", err.Error()))
	}
	return
}

// ClosePayOrder 关闭支付订单
func ClosePayOrder(wb *core.WeBase, req *CloseOrderRequest) (rsp *CloseOrderResponse, err error) {
	// req.CommonRequestMsg = CommonRequestMsg{
	// 	RequestMsg: RequestMsg{
	// 		MchID:    wb.PayID,
	// 		NonceStr: util.RandString(32),
	// 		SignType: util.SignTypeHMACSHA256,
	// 	},
	// 	AppID: wb.AppID,
	// }
	req.MchID = wb.PayID
	req.NonceStr = util.RandString(32)
	req.SignType = util.SignTypeHMACSHA256
	req.AppID = wb.AppID

	params, err := util.StructToURLValue(req, "xml")
	if nil != err {
		return nil, errmsg.GetError(errPayOrderSign, err.Error())
	}
	params.Add("sign", util.SignSha256(params.Encode(), wb.PayKey))

	data, err := util.HTTPXMLPost(wb.HTTPClient, closePayOrderURL, params)
	if err != nil {
		return nil, errmsg.GetError(errPayOrderRsp, fmt.Sprintf("closePayOrderURL rsp err:%s", err.Error()))
	}
	rsp = &CloseOrderResponse{}
	err = xml.Unmarshal(data, &rsp)
	if err != nil {
		return nil, errmsg.GetError(errPayOrderRsp, fmt.Sprintf("closePayOrderURL xml.Unmarshal err:%s", err.Error()))
	}
	return
}

// RefundPayOrder 申请退款
func RefundPayOrder(wb *core.WeBase, req *RefundPayOrderRequest) (rsp *RefundPayOrderResponse, err error) {
	// req.CommonRequestMsg = CommonRequestMsg{
	// 	RequestMsg: RequestMsg{
	// 		MchID:    wb.PayID,
	// 		NonceStr: util.RandString(32),
	// 		SignType: util.SignTypeHMACSHA256,
	// 	},
	// 	AppID: wb.AppID,
	// }
	req.MchID = wb.PayID
	req.NonceStr = util.RandString(32)
	req.SignType = util.SignTypeHMACSHA256
	req.AppID = wb.AppID

	params, err := util.StructToURLValue(req, "xml")
	if nil != err {
		return nil, errmsg.GetError(errPayOrderSign, err.Error())
	}
	params.Add("sign", util.SignSha256(params.Encode(), wb.PayKey))

	data, err := util.HTTPXMLPostWithCertificate(wb.HTTPClient, refundPayOrderURL, params, wechatPayAPICert)
	if err != nil {
		return nil, errmsg.GetError(errPayOrderRsp, fmt.Sprintf("refundPayOrderURL rsp err:%s", err.Error()))
	}
	rsp = &RefundPayOrderResponse{}
	err = xml.Unmarshal(data, &rsp)
	if err != nil {
		return nil, errmsg.GetError(errPayOrderRsp, fmt.Sprintf("refundPayOrderURL xml.Unmarshal err:%s,data:%s", err.Error(), string(data)))
	}
	if rsp.ReturnCode != util.CodeSuccess {
		return nil, errmsg.GetError(errPayOrderRsp, fmt.Sprintf("refundPayOrderURL rsp ReturnCode fail:%s,data:%s", rsp.ReturnMsg, string(data)))
	}
	if rsp.ResultCode != util.CodeSuccess {
		return nil, errmsg.GetError(errPayOrderRsp, fmt.Sprintf("refundPayOrderURL rsp ResultCode fail:%s,data:%s", rsp.ErrCode, string(data)))
	}
	return
}

// QueryRefundPayOrder 查询申请退款
func QueryRefundPayOrder(wb *core.WeBase, req *QueryRefundPayOrderRequest) (rsp *QueryRefundPayOrderResponse, err error) {
	// req.CommonRequestMsg = CommonRequestMsg{
	// 	RequestMsg: RequestMsg{
	// 		MchID:    wb.PayID,
	// 		NonceStr: util.RandString(32),
	// 		SignType: util.SignTypeHMACSHA256,
	// 	},
	// 	AppID: wb.AppID,
	// }
	req.MchID = wb.PayID
	req.NonceStr = util.RandString(32)
	req.SignType = util.SignTypeHMACSHA256
	req.AppID = wb.AppID

	params, err := util.StructToURLValue(req, "xml")
	if nil != err {
		return nil, errmsg.GetError(errPayOrderSign, err.Error())
	}
	params.Add("sign", util.SignSha256(params.Encode(), wb.PayKey))

	data, err := util.HTTPXMLPost(wb.HTTPClient, queryRefundPayOrderURL, params)
	if err != nil {
		return nil, errmsg.GetError(errPayOrderRsp, fmt.Sprintf("queryRefundPayOrderURL rsp err:%s", err.Error()))
	}
	rsp = &QueryRefundPayOrderResponse{}
	err = xml.Unmarshal(data, &rsp)
	if err != nil {
		return nil, errmsg.GetError(errPayOrderRsp, fmt.Sprintf("queryRefundPayOrderURL xml.Unmarshal err:%s", err.Error()))
	}
	return
}
