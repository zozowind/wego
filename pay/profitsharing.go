package pay

import (
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/zozowind/wego/core"
	"github.com/zozowind/wego/libs/errmsg"
	"github.com/zozowind/wego/util"
)

/*
参考：https://github.com/blusewang/wxApi-go
*/

const (
	// profitSharingURL 单次分账URL
	profitSharingURL = core.WxPayURL + "/secapi/pay/profitsharing"
	// multiProfitSharingURL 多次分账URL
	multiProfitSharingURL = core.WxPayURL + "/secapi/pay/multiprofitsharing"
	// queryProfitSharingURL 查询分账结果URL
	queryProfitSharingURL = core.WxPayURL + "/secapi/pay/profitsharingquery"
	// addReceiverProfitSharingURL 添加分账接收方URL
	addReceiverProfitSharingURL = core.WxPayURL + "/secapi/pay/profitsharingaddreceiver"
	// removeReceiverProfitSharingURL 删除分账接收方URL
	removeReceiverProfitSharingURL = core.WxPayURL + "/secapi/pay/profitsharingremovereceiver"
	// finishProfitSharingURL 完结分账URL
	finishProfitSharingURL = core.WxPayURL + "/secapi/pay/profitsharingfinish"
	// returnProfitSharingURL 分账回退URL
	returnProfitSharingURL = core.WxPayURL + "/secapi/pay/profitsharingreturn"
	// returnQueryProfitSharingURL 分账回退查询URL
	returnQueryProfitSharingURL = core.WxPayURL + "/secapi/pay/profitsharingreturnquery"
)

var wechatPayAPICert *tls.Certificate

// LoadCertificate 加载微信支付API证书
func LoadCertificate(wechatPayCert string, wechatPayKey string) (err error) {
	cert, err := tls.LoadX509KeyPair(wechatPayCert, wechatPayKey)
	if nil != err {
		err = errors.New("微信支付API证书加载失败")
		return
	}
	wechatPayAPICert = &cert
	return
}

// // RequestMsg .
// type RequestMsg struct {
// 	MchID    string `xml:"mch_id"`    // 微信支付分配的商户号
// 	NonceStr string `xml:"nonce_str"` // 随机字符串
// 	SignType string `xml:"sign_type"` // 签名类型：目前只支持HMAC-SHA256
// 	// Sign    string  `xml:"sign"`//签名：单独赋值
// }

// // CommonRequestMsg .
// type CommonRequestMsg struct {
// 	RequestMsg
// 	AppID string `xml:"appid"` // 微信分配的公众账号ID
// }

// ErrMsg .
type ErrMsg struct {
	ReturnCode string `xml:"return_code"` // 返回状态码：SUCCESS/FAIL 此字段是通信标识，非交易标识
	ReturnMsg  string `xml:"return_msg"`  // 返回信息：如非空，为错误原因
}

// ResultMsg .
type ResultMsg struct {
	ResultCode string `xml:"result_code"`  // 业务结果：SUCCESS：分账申请接收成功，结果通过分账查询接口查询FAIL ：提交业务失败
	ErrCode    string `xml:"err_code"`     // 错误代码：列表详见错误码列表
	ErrCodeDes string `xml:"err_code_des"` // 错误代码描述：结果信息描述
}

// CommonResponseMsg 请求返回的共有的字段
type CommonResponseMsg struct {
	ErrMsg
	ResultMsg
}

// ProfitSharingRequest 分账请求
type ProfitSharingRequest struct {
	// CommonRequestMsg
	MchID    string `xml:"mch_id"`    // 微信支付分配的商户号
	NonceStr string `xml:"nonce_str"` // 随机字符串
	SignType string `xml:"sign_type"` // 签名类型：目前只支持HMAC-SHA256
	AppID    string `xml:"appid"`     // 小程序ID

	Receivers     string `xml:"receivers"`      // 分账接收方列表json
	TransactionID string `xml:"transaction_id"` // 微信支付订单号
	OutOrderNo    string `xml:"out_order_no"`   // 商户系统内部的分账单号

	ReceiverSlice []ProfitSharingReqReceiver `xml:"-"` // receivers字段的结构化
}

// ProfitSharingReqReceiver 分账结果中的接收者
type ProfitSharingReqReceiver struct {
	Type        string `json:"type"`        // 分账接收方类型:MERCHANT_ID：商户ID,PERSONAL_WECHATID：个人微信号,PERSONAL_OPENID：个人openid
	Account     string `json:"account"`     // 参考type
	Amount      int64  `json:"amount"`      // 分账金额，单位为分
	Description string `json:"description"` // 分账描述
}

// ProfitSharingResponse 分账结果
type ProfitSharingResponse struct {
	CommonResponseMsg

	MchID         string `xml:"mch_id"`         // 商户号
	AppID         string `xml:"appid"`          // 公众账号ID
	NonceStr      string `xml:"nonce_str"`      // 随机字符串
	Sign          string `xml:"sign"`           // 签名
	TransactionID string `xml:"transaction_id"` // 微信订单号
	OutOrderNo    string `xml:"out_order_no"`   // 商户分账单号
	OrderID       string `xml:"order_id"`       //微信分账单号
}

// ProfitSharingQueryRequest 查询分账结果请求
type ProfitSharingQueryRequest struct {
	// RequestMsg
	MchID    string `xml:"mch_id"`    // 微信支付分配的商户号
	NonceStr string `xml:"nonce_str"` // 随机字符串
	SignType string `xml:"sign_type"` // 签名类型：目前只支持HMAC-SHA256
	// AppID    string `xml:"appid"`     // 小程序ID

	TransactionID string `xml:"transaction_id"` // 微信订单号
	OutOrderNo    string `xml:"out_order_no"`   // 商户分账单号
}

// ProfitSharingQueryResponse 查询分账结果返回
type ProfitSharingQueryResponse struct {
	CommonResponseMsg
	// return_code为SUCCESS时返回
	MchID    string `xml:"mch_id"`    // 商户号
	NonceStr string `xml:"nonce_str"` //随机字符串
	Sign     string `xml:"sign"`      // 签名
	// return_code和result_code都为SUCCESS时返回
	TransactionID string `xml:"transaction_id"` // 微信订单号
	OutOrderNo    string `xml:"out_order_no"`   // 商户分账单号
	OrderID       string `xml:"order_id"`       // 微信分账单号
	Status        string `xml:"status"`         // 分账单状态
	CloseReason   string `xml:"close_reason"`   // 关闭原因
	Receivers     string `xml:"receivers"`      // 分账接受者列表json

	ReceiverSlice []ProfitSharingQueryResponseReceiver `xml:"-"` // receivers字段的结构化
}

// ProfitSharingQueryResponseReceiver 查询分账结果的接收者
type ProfitSharingQueryResponseReceiver struct {
	ProfitSharingReqReceiver

	Result     string `json:"result"`      // 分账结果
	FinishTime string `json:"finish_time"` // 分账完成时间
	FailReason string `json:"fail_reason"` // 分账失败原因
}

// ProfitSharingAddReceiverRequest 添加分账接收方
type ProfitSharingAddReceiverRequest struct {
	// CommonRequestMsg
	MchID    string `xml:"mch_id"`    // 微信支付分配的商户号
	NonceStr string `xml:"nonce_str"` // 随机字符串
	SignType string `xml:"sign_type"` // 签名类型：目前只支持HMAC-SHA256
	AppID    string `xml:"appid"`     // 小程序ID

	Receiver string `xml:"receiver"` // 分账接收者json

	ReceiverSlice []ProfitSharingAddReceiverRequestReceiver `xml:"-"` // receiver字段的结构化
}

// ProfitSharingAddReceiverRequestReceiver 添加分账接收方的接收者
type ProfitSharingAddReceiverRequestReceiver struct {
	ProfitSharingReqReceiver

	Name           string `json:"name"`            // 分账接收方全称（商户全名或者个人姓名）
	RelationType   string `json:"relation_type"`   // 与分账方的关系类型，分类见util/const.go
	CustomRelation string `json:"custom_relation"` // 自定义的分账关系
}

// ProfitSharingUpdateReceiver .
type ProfitSharingUpdateReceiver struct {
	MchID    string `xml:"mch_id"`    // 商户号
	AppID    string `xml:"appid"`     // 公众账号ID
	NonceStr string `xml:"nonce_str"` // 随机字符串
	Sign     string `xml:"sign"`      // 签名
	Receiver string `xml:"receiver"`  // 分账接收方
}

// ProfitSharingAddReceiverResponse  添加分账接收方返回
type ProfitSharingAddReceiverResponse struct {
	CommonResponseMsg
	// return_code为SUCCESS时返回
	ProfitSharingUpdateReceiver

	ReceiverSlice []ProfitSharingAddReceiverRequestReceiver `xml:"-"` // Receiver字段的结构化
}

// ProfitSharingRemoveReceiverRequest 删除分账接收方
type ProfitSharingRemoveReceiverRequest struct {
	// CommonRequestMsg
	MchID    string `xml:"mch_id"`    // 微信支付分配的商户号
	NonceStr string `xml:"nonce_str"` // 随机字符串
	SignType string `xml:"sign_type"` // 签名类型：目前只支持HMAC-SHA256
	AppID    string `xml:"appid"`     // 小程序ID

	Receiver string `xml:"receiver"` // 分账接收方

	ReceiverSlice []ProfitSharingRemoveReceiverRequestReceiver `xml:"-"` // Receiver字段的结构化
}

// ProfitSharingRemoveReceiverRequestReceiver 删除分账接收方的接收者
type ProfitSharingRemoveReceiverRequestReceiver struct {
	Type    string `json:"type"`
	Account string `json:"account"`
}

// ProfitSharingRemoveReceiverResponse  删除分账接收方返回
type ProfitSharingRemoveReceiverResponse struct {
	CommonResponseMsg
	// return_code为SUCCESS时返回
	ProfitSharingUpdateReceiver

	ReceiverSlice []ProfitSharingRemoveReceiverRequestReceiver `xml:"-"`
}

// ProfitSharingFinishRequest 完结分账请求
type ProfitSharingFinishRequest struct {
	// CommonRequestMsg
	MchID    string `xml:"mch_id"`    // 微信支付分配的商户号
	NonceStr string `xml:"nonce_str"` // 随机字符串
	SignType string `xml:"sign_type"` // 签名类型：目前只支持HMAC-SHA256
	AppID    string `xml:"appid"`     // 小程序ID

	TransactionID string `xml:"transaction_id"` // 微信订单号
	OutOrderNo    string `xml:"out_order_no"`   // 商户分账单号
	Description   string `xml:"description"`    // 分账完结描述
}

// ProfitSharingFinishResponse 完结分账返回
type ProfitSharingFinishResponse struct {
	CommonResponseMsg
	// return_code为SUCCESS时返回
	MchID    string `xml:"mch_id"`    // 商户号
	AppID    string `xml:"appid"`     // 公众账号ID
	NonceStr string `xml:"nonce_str"` // 随机字符串
	Sign     string `xml:"sign"`      // 签名
	// return_code和result_code都为SUCCESS时返回
	TransactionID string `xml:"transaction_id"` // 微信订单号
	OutOrderNo    string `xml:"out_order_no"`   // 商户分账单号
	OrderID       string `xml:"order_id"`       // 微信分账单号
}

// ProfitSharingReturnRequest 分账回退请求
type ProfitSharingReturnRequest struct {
	// CommonRequestMsg
	MchID    string `xml:"mch_id"`    // 微信支付分配的商户号
	NonceStr string `xml:"nonce_str"` // 随机字符串
	SignType string `xml:"sign_type"` // 签名类型：目前只支持HMAC-SHA256
	AppID    string `xml:"appid"`     // 小程序ID

	OrderID           string `xml:"order_id"`            // 微信分账单号，发起分账返回的微信分账单号
	OutOrderNo        string `xml:"out_order_no"`        // 商户分账单号
	OutReturnNo       string `xml:"out_return_no"`       // 商户回退单号
	ReturnAccountType string `xml:"return_account_type"` // 回退方账号类型
	ReturnAccount     string `xml:"return_account"`      // 回退方账号
	ReturnAmount      int64  `xml:"return_amount"`       // 回退金额，单位分
	Description       string `xml:"description"`         // 回退描述
}

// ProfitSharingReturnResponse 分账回退返回
type ProfitSharingReturnResponse struct {
	CommonResponseMsg
	// return_code为SUCCESS时返回
	MchID             string `xml:"mch_id"`              // 商户号
	AppID             string `xml:"appid"`               // 公众账号ID
	NonceStr          string `xml:"nonce_str"`           // 随机字符串
	Sign              string `xml:"sign"`                // 签名
	OrderID           string `xml:"order_id"`            // 微信分账单号
	OutOrderNo        string `xml:"out_order_no"`        // 商户分账单号
	OutReturnNo       string `xml:"out_return_no"`       // 商户回退单号
	ReturnNo          string `xml:"return_no"`           // 微信回退单号
	ReturnAccountType string `xml:"return_account_type"` // 回退方账号类型
	ReturnAccount     string `xml:"return_account"`      // 回退方账号
	ReturnAmount      int64  `xml:"return_amount"`       // 回退金额，单位分
	Description       string `xml:"description"`         // 回退描述
	Result            string `xml:"result"`              // 回退结果
	FailReason        string `xml:"fail_reason"`         // 失败原因
	FinishTime        string `xml:"finish_time"`         // 完成时间
}

// ProfitSharingReturnQueryRequest 分账回退结果查询请求
type ProfitSharingReturnQueryRequest struct {
	// CommonRequestMsg
	MchID    string `xml:"mch_id"`    // 微信支付分配的商户号
	NonceStr string `xml:"nonce_str"` // 随机字符串
	SignType string `xml:"sign_type"` // 签名类型：目前只支持HMAC-SHA256
	AppID    string `xml:"appid"`     // 小程序ID

	OrderID     string `xml:"order_id"`      // 微信分账单号
	OutOrderNo  string `xml:"out_order_no"`  // 商户分账单号
	OutReturnNo string `xml:"out_return_no"` // 商户回退单号
}

// ProfitSharingReturnQueryResponse 分账回退结果返回
type ProfitSharingReturnQueryResponse struct {
	CommonResponseMsg

	// return_code为SUCCESS时返回
	MchID             string `xml:"mch_id"`              // 商户号
	AppID             string `xml:"appid"`               // 公众账号ID
	NonceStr          string `xml:"nonce_str"`           // 随机字符串
	Sign              string `xml:"sign"`                // 签名
	OrderID           string `xml:"order_id"`            // 微信分账单号
	OutOrderNo        string `xml:"out_order_no"`        // 商户分账单号
	OutReturnNo       string `xml:"out_return_no"`       // 商户回退单号
	ReturnNo          string `xml:"return_no"`           // 微信回退单号
	ReturnAccountType string `xml:"return_account_type"` // 回退方类型
	ReturnAccount     string `xml:"return_account"`      // 回退方账号
	ReturnAmount      int64  `xml:"return_amount"`       // 回退金额，单位分
	Description       string `xml:"description"`         // 回退描述
	Result            string `xml:"result"`              // 回退结果
	FailReason        string `xml:"fail_reason"`         // 失败原因
	FinishTime        string `xml:"finish_time"`         // 分账回退完成时间
}

// ProfitSharing 请求单次分账
func ProfitSharing(wb *core.WeBase, req *ProfitSharingRequest) (rsp *ProfitSharingResponse, err error) {
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

	if nil == wechatPayAPICert {
		return nil, errmsg.GetError(errProfitSharingReq, "certificate is not load")
	}
	if len(req.ReceiverSlice) == 0 {
		return nil, errmsg.GetError(errProfitSharingReq, "receiverSlice is not exist")
	}
	raw, err := json.Marshal(req.ReceiverSlice)
	if err != nil {
		return nil, errmsg.GetError(errProfitSharingReq, fmt.Sprintf("receiverSlice marshal err:%s", err.Error()))
	}

	req.Receivers = string(raw)
	params, err := util.StructToURLValue(req, "xml")
	if nil != err {
		return nil, errmsg.GetError(errProfitSharingSign, err.Error())
	}
	params.Add("sign", util.SignSha256(params.Encode(), wb.PayKey))
	data, err := util.HTTPXMLPostWithCertificate(wb.HTTPClient, profitSharingURL, params, wechatPayAPICert)
	if err != nil {
		return nil, errmsg.GetError(errProfitSharingRsp, fmt.Sprintf("profitSharingURL rsp err:%s", err.Error()))
	}
	rsp = &ProfitSharingResponse{}
	err = xml.Unmarshal(data, &rsp)
	if err != nil {
		return nil, errmsg.GetError(errProfitSharingRsp, fmt.Sprintf("profitSharingURL xml.Unmarshal err:%s", err.Error()))
	}
	return
}

// MultiProfitSharing 请求多次分账
func MultiProfitSharing(wb *core.WeBase, req *ProfitSharingRequest) (rsp *ProfitSharingResponse, err error) {
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

	if nil == wechatPayAPICert {
		return nil, errmsg.GetError(errProfitSharingReq, "certificate is not load")
	}
	if len(req.ReceiverSlice) == 0 {
		return nil, errmsg.GetError(errProfitSharingReq, "receiverSlice is not exist")
	}
	raw, err := json.Marshal(req.ReceiverSlice)
	if err != nil {
		return nil, errmsg.GetError(errProfitSharingReq, fmt.Sprintf("receiverSlice marshal err:%s", err.Error()))
	}

	req.Receivers = string(raw)
	params, err := util.StructToURLValue(req, "xml")
	if nil != err {
		return nil, errmsg.GetError(errProfitSharingSign, err.Error())
	}
	params.Add("sign", util.SignSha256(params.Encode(), wb.PayKey))

	data, err := util.HTTPXMLPostWithCertificate(wb.HTTPClient, multiProfitSharingURL, params, wechatPayAPICert)
	if err != nil {
		return nil, errmsg.GetError(errProfitSharingRsp, fmt.Sprintf("multiProfitSharingURL rsp err:%s", err.Error()))
	}
	rsp = &ProfitSharingResponse{}
	err = xml.Unmarshal(data, &rsp)
	if err != nil {
		return nil, errmsg.GetError(errProfitSharingRsp, fmt.Sprintf("multiProfitSharingURL  xml.Unmarshal err:%s", err.Error()))
	}
	return
}

// QueryProfitSharing 查询分账结果
func QueryProfitSharing(wb *core.WeBase, req *ProfitSharingQueryRequest) (rsp *ProfitSharingQueryResponse, err error) {
	// req.RequestMsg = RequestMsg{
	// 	MchID:    wb.PayID,
	// 	NonceStr: util.RandString(32),
	// 	SignType: util.SignTypeHMACSHA256,
	// }
	req.MchID = wb.PayID
	req.NonceStr = util.RandString(32)
	req.SignType = util.SignTypeHMACSHA256

	params, err := util.StructToURLValue(req, "xml")
	if nil != err {
		return nil, errmsg.GetError(errProfitSharingSign, err.Error())
	}
	params.Add("sign", util.SignSha256(params.Encode(), wb.PayKey))

	data, err := util.HTTPXMLPost(wb.HTTPClient, queryProfitSharingURL, params)
	if err != nil {
		return nil, errmsg.GetError(errProfitSharingRsp, fmt.Sprintf("queryProfitSharingURL rsp err:%s", err.Error()))
	}
	rsp = &ProfitSharingQueryResponse{}
	err = xml.Unmarshal(data, &rsp)
	if err != nil {
		return nil, errmsg.GetError(errProfitSharingRsp, fmt.Sprintf("queryProfitSharingURL xml.Unmarshal err:%s", err.Error()))
	}
	err = json.Unmarshal([]byte(rsp.Receivers), &rsp.ReceiverSlice)
	if err != nil {
		return nil, errmsg.GetError(errProfitSharingRsp, fmt.Sprintf("queryProfitSharingURL json.Unmarshal err:%s", err.Error()))
	}
	return
}

// AddReceiverProfitSharing 添加分账接收方
func AddReceiverProfitSharing(wb *core.WeBase, req *ProfitSharingAddReceiverRequest) (rsp *ProfitSharingAddReceiverResponse, err error) {
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

	if len(req.ReceiverSlice) == 0 {
		return nil, errmsg.GetError(errProfitSharingReq, "receiverSlice is not exist")
	}
	raw, err := json.Marshal(req.ReceiverSlice)
	if err != nil {
		return nil, errmsg.GetError(errProfitSharingReq, fmt.Sprintf("receiverSlice marshal err:%s", err.Error()))
	}

	req.Receiver = string(raw)
	params, err := util.StructToURLValue(req, "xml")
	if nil != err {
		return nil, errmsg.GetError(errProfitSharingSign, err.Error())
	}
	params.Add("sign", util.SignSha256(params.Encode(), wb.PayKey))

	data, err := util.HTTPXMLPost(wb.HTTPClient, addReceiverProfitSharingURL, params)
	if err != nil {
		return nil, errmsg.GetError(errProfitSharingRsp, fmt.Sprintf("addReceiverProfitSharingURL rsp err:%s", err.Error()))
	}
	rsp = &ProfitSharingAddReceiverResponse{}
	err = xml.Unmarshal(data, &rsp)
	if err != nil {
		return nil, errmsg.GetError(errProfitSharingRsp, fmt.Sprintf("addReceiverProfitSharingURL  xml.Unmarshal err:%s", err.Error()))
	}
	err = json.Unmarshal([]byte(rsp.Receiver), &rsp.ReceiverSlice)
	if err != nil {
		return nil, errmsg.GetError(errProfitSharingRsp, fmt.Sprintf("addReceiverProfitSharingURL  json.Unmarshal err:%s", err.Error()))
	}
	return
}

// RemoveReceiverProfitSharing 删除分账接收方
func RemoveReceiverProfitSharing(wb *core.WeBase, req *ProfitSharingRemoveReceiverRequest) (rsp *ProfitSharingRemoveReceiverResponse, err error) {
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

	if len(req.ReceiverSlice) == 0 {
		return nil, errmsg.GetError(errProfitSharingReq, "receiverSlice is not exist")
	}
	raw, err := json.Marshal(req.ReceiverSlice)
	if err != nil {
		return nil, errmsg.GetError(errProfitSharingReq, fmt.Sprintf("receiverSlice marshal err:%s", err.Error()))
	}

	req.Receiver = string(raw)
	params, err := util.StructToURLValue(req, "xml")
	if nil != err {
		return nil, errmsg.GetError(errProfitSharingSign, err.Error())
	}
	params.Add("sign", util.SignSha256(params.Encode(), wb.PayKey))

	data, err := util.HTTPXMLPost(wb.HTTPClient, removeReceiverProfitSharingURL, params)
	if err != nil {
		return nil, errmsg.GetError(errProfitSharingRsp, fmt.Sprintf("removeReceiverProfitSharingURL rsp err:%s", err.Error()))
	}
	rsp = &ProfitSharingRemoveReceiverResponse{}
	err = xml.Unmarshal(data, &rsp)
	if err != nil {
		return nil, errmsg.GetError(errProfitSharingRsp, fmt.Sprintf("removeReceiverProfitSharingURL xml.Unmarshal err:%s", err.Error()))
	}
	err = json.Unmarshal([]byte(rsp.Receiver), &rsp.ReceiverSlice)
	if err != nil {
		return nil, errmsg.GetError(errProfitSharingRsp, fmt.Sprintf("removeReceiverProfitSharingURL json.Unmarshal err:%s", err.Error()))
	}
	return
}

// ProfitSharingFinish 完结分账
func ProfitSharingFinish(wb *core.WeBase, req *ProfitSharingFinishRequest) (rsp *ProfitSharingFinishResponse, err error) {
	if nil == wechatPayAPICert {
		return nil, errmsg.GetError(errProfitSharingReq, "certificate is not load")
	}
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
		return nil, errmsg.GetError(errProfitSharingSign, err.Error())
	}
	params.Add("sign", util.SignSha256(params.Encode(), wb.PayKey))
	data, err := util.HTTPXMLPostWithCertificate(wb.HTTPClient, finishProfitSharingURL, params, wechatPayAPICert)
	if err != nil {
		return nil, errmsg.GetError(errProfitSharingRsp, fmt.Sprintf("finishProfitSharingURL rsp err:%s", err.Error()))
	}
	rsp = &ProfitSharingFinishResponse{}
	err = xml.Unmarshal(data, &rsp)
	if err != nil {
		return nil, errmsg.GetError(errProfitSharingRsp, fmt.Sprintf("finishProfitSharingURL xml.Unmarshal err:%s", err.Error()))
	}
	return
}

// ProfitSharingReturn 分账回退
func ProfitSharingReturn(wb *core.WeBase, req *ProfitSharingReturnRequest) (rsp *ProfitSharingReturnResponse, err error) {
	if nil == wechatPayAPICert {
		return nil, errmsg.GetError(errProfitSharingReq, "certificate is not load")
	}
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
		return nil, errmsg.GetError(errProfitSharingSign, err.Error())
	}
	params.Add("sign", util.SignSha256(params.Encode(), wb.PayKey))
	data, err := util.HTTPXMLPostWithCertificate(wb.HTTPClient, returnProfitSharingURL, params, wechatPayAPICert)
	if err != nil {
		return nil, errmsg.GetError(errProfitSharingRsp, fmt.Sprintf("returnProfitSharingURL rsp err:%s", err.Error()))
	}
	rsp = &ProfitSharingReturnResponse{}
	err = xml.Unmarshal(data, &rsp)
	if err != nil {
		return nil, errmsg.GetError(errProfitSharingRsp, fmt.Sprintf("returnProfitSharingURL xml.Unmarshal err:%s", err.Error()))
	}
	return
}

// ReturnQueryReceiverProfitSharing 分账回退结果查询
func ReturnQueryReceiverProfitSharing(wb *core.WeBase, req *ProfitSharingReturnQueryRequest) (rsp *ProfitSharingReturnQueryResponse, err error) {
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
		return nil, errmsg.GetError(errProfitSharingSign, err.Error())
	}
	params.Add("sign", util.SignSha256(params.Encode(), wb.PayKey))

	data, err := util.HTTPXMLPost(wb.HTTPClient, returnQueryProfitSharingURL, params)
	if err != nil {
		return nil, errmsg.GetError(errProfitSharingRsp, fmt.Sprintf("returnQueryProfitSharingURL rsp err:%s", err.Error()))
	}
	rsp = &ProfitSharingReturnQueryResponse{}
	err = xml.Unmarshal(data, &rsp)
	if err != nil {
		return nil, errmsg.GetError(errProfitSharingRsp, fmt.Sprintf("returnQueryProfitSharingURL  xml.Unmarshal err:%s", err.Error()))
	}
	return
}
