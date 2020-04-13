package profitsharing

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

// ProfitSharingRequest 分账请求
type ProfitSharingRequest struct {
	MchID    string `xml:"mch_id"`
	AppID    string `xml:"appid"`
	NonceStr string `xml:"nonce_str"`
	SignType string `xml:"sign_type"`
	// Sign          string                     `xml:"sign"`//单独赋值
	TransactionID string                     `xml:"transaction_id"`
	OutOrderNo    string                     `xml:"out_order_no"`
	Receivers     string                     `xml:"receivers"`
	ReceiverSlice []ProfitSharingReqReceiver `xml:"-"`
}

// ProfitSharingReqReceiver 分账结果中的接收者
type ProfitSharingReqReceiver struct {
	Type        string `json:"type"`
	Account     string `json:"account"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
}

// ProfitSharingResponse 分账结果
type ProfitSharingResponse struct {
	CommonResponseMsg

	MchID         string `xml:"mch_id"`
	AppID         string `xml:"appid"`
	NonceStr      string `xml:"nonce_str"`
	Sign          string `xml:"sign"`
	TransactionID string `xml:"transaction_id"`
	OutOrderNo    string `xml:"out_order_no"`
	OrderID       string `xml:"order_id"`
}

// ErrMsg .
type ErrMsg struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
}

// ResultMsg .
type ResultMsg struct {
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
}

// CommonResponseMsg 请求返回的共有的字段
type CommonResponseMsg struct {
	ErrMsg
	ResultMsg
}

// ProfitSharingQueryRequest 查询分账结果请求
type ProfitSharingQueryRequest struct {
	MchID         string `xml:"mch_id"`
	TransactionID string `xml:"transaction_id"`
	OutOrderNo    string `xml:"out_order_no"`
	NonceStr      string `xml:"nonce_str"`
	SignType      string `xml:"sign_type"`
	// Sign          string                     `xml:"sign"`//单独赋值
}

// ProfitSharingQueryResponse 查询分账结果返回
type ProfitSharingQueryResponse struct {
	CommonResponseMsg
	// return_code为SUCCESS时返回
	MchID    string `xml:"mch_id"`
	NonceStr string `xml:"nonce_str"`
	Sign     string `xml:"sign"`
	// return_code和result_code都为SUCCESS时返回
	TransactionID string                               `xml:"transaction_id"`
	OutOrderNo    string                               `xml:"out_order_no"`
	OrderID       string                               `xml:"order_id"`
	Status        string                               `xml:"status"`
	CloseReason   string                               `xml:"close_reason"`
	Receivers     string                               `xml:"receivers"`
	ReceiverSlice []ProfitSharingQueryResponseReceiver `xml:"-"`
}

// ProfitSharingQueryResponseReceiver 查询分账结果的接收者
type ProfitSharingQueryResponseReceiver struct {
	Type        string `json:"type"`
	Account     string `json:"account"`
	Amount      int64  `json:"amount"`
	Description string `json:"description"`
	Result      string `json:"result"`
	FinishTime  string `json:"finish_time"`
	FailReason  string `json:"fail_reason"`
}

// ProfitSharingAddReceiverRequest 添加分账接收方
type ProfitSharingAddReceiverRequest struct {
	MchID    string `xml:"mch_id"`
	AppID    string `xml:"appid"`
	NonceStr string `xml:"nonce_str"`
	SignType string `xml:"sign_type"`
	// Sign     string `xml:"sign"`
	Receiver      string                                    `xml:"receiver"`
	ReceiverSlice []ProfitSharingAddReceiverRequestReceiver `xml:"-"`
}

// ProfitSharingAddReceiverRequestReceiver 添加分账接收方的接收者
type ProfitSharingAddReceiverRequestReceiver struct {
	Type           string `json:"type"`
	Account        string `json:"account"`
	Description    string `json:"description"`
	Name           string `json:"name"`
	RelationType   string `json:"relation_type"`
	CustomRelation string `json:"custom_relation"`
}

// ProfitSharingAddReceiverResponse  添加分账接收方返回
type ProfitSharingAddReceiverResponse struct {
	CommonResponseMsg
	// return_code为SUCCESS时返回
	MchID         string                               `xml:"mch_id"`
	AppID         string                               `xml:"appid"`
	NonceStr      string                               `xml:"nonce_str"`
	Sign          string                               `xml:"sign"`
	Receiver      string                               `xml:"receiver"`
	ReceiverSlice []ProfitSharingQueryResponseReceiver `xml:"-"`
}

// ProfitSharing 请求单次分账
func (wb *WeBase) ProfitSharing(req *ProfitSharingRequest) (rsp *ProfitSharingResponse, err error) {
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
	req.AppID = wb.AppID
	req.MchID = wb.PayID
	req.NonceStr = util.RandString(32)
	req.SignType = util.SignTypeHMACSHA256
	params, err := util.StructToURLValue(req, "xml")
	if nil != err {
		return nil, errmsg.GetError(errProfitSharingSign, err.Error())
	}
	params.Add("sign", util.SignSha256(params.Encode(), client.PayKey))
	data, err := util.HTTPXMLPostWithCertificate(client.HTTPClient, profitSharingURL, params, wechatPayAPICert)
	if err != nil {
		return nil, errmsg.GetError(errProfitSharingRsp, fmt.Sprintf("profitSharingURL rsp err:%s", err.Error()))
	}
	rsp = &ProfitSharingResponse{}
	err = xml.Unmarshal(data, &rsp)
	return
}

// MultiProfitSharing 请求多次分账
func (wb *WeBase) MultiProfitSharing(req *ProfitSharingRequest) (rsp *ProfitSharingResponse, err error) {
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
	req.AppID = wb.AppID
	req.MchID = wb.PayID
	req.NonceStr = util.RandString(32)
	req.SignType = util.SignTypeHMACSHA256
	params, err := util.StructToURLValue(req, "xml")
	if nil != err {
		return nil, errmsg.GetError(errProfitSharingSign, err.Error())
	}
	params.Add("sign", util.SignSha256(params.Encode(), client.PayKey))

	data, err := util.HTTPXMLPostWithCertificate(client.HTTPClient, multiProfitSharingURL, params, wechatPayAPICert)
	if err != nil {
		return nil, errmsg.GetError(errProfitSharingRsp, fmt.Sprintf("multiProfitSharingURL rsp err:%s", err.Error()))
	}
	rsp = &ProfitSharingResponse{}
	err = xml.Unmarshal(data, &rsp)
	return
}

// QueryProfitSharing 查询分账结果
func (wb *WeBase) QueryProfitSharing(transactionID string, outOrderNo string) (rsp *ProfitSharingQueryResponse, err error) {
	req := &ProfitSharingQueryRequest{
		MchID:         wb.PayID,
		TransactionID: transactionID,
		OutOrderNo:    outOrderNo,
		NonceStr:      util.RandString(32),
		SignType:      util.SignTypeHMACSHA256,
	}
	params, err := util.StructToURLValue(req, "xml")
	if nil != err {
		return nil, errmsg.GetError(errProfitSharingSign, err.Error())
	}
	params.Add("sign", util.SignSha256(params.Encode(), client.PayKey))

	data, err := util.HTTPXMLPost(client.HTTPClient, queryProfitSharingURL, params)
	if err != nil {
		return nil, errmsg.GetError(errProfitSharingRsp, fmt.Sprintf("queryProfitSharingURL rsp err:%s", err.Error()))
	}
	rsp = &ProfitSharingQueryResponse{}
	err = xml.Unmarshal(data, &rsp)
	return
}

// AddReceiverProfitSharing 添加分账接收方
func (wb *WeBase) AddReceiverProfitSharing(req *ProfitSharingAddReceiverRequest) (rsp *ProfitSharingAddReceiverResponse, err error) {
	req.MchID = wb.PayID
	req.AppID = wb.AppID
	req.NonceStr = util.RandString(32)
	req.SignType = util.SignTypeHMACSHA256
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
	params.Add("sign", util.SignSha256(params.Encode(), client.PayKey))

	data, err := util.HTTPXMLPost(client.HTTPClient, addReceiverProfitSharingURL, params)
	if err != nil {
		return nil, errmsg.GetError(errProfitSharingRsp, fmt.Sprintf("addReceiverProfitSharingURL rsp err:%s", err.Error()))
	}
	rsp = &ProfitSharingAddReceiverResponse{}
	err = xml.Unmarshal(data, &rsp)
	return
}

// RemoveReceiverProfitSharing 删除分账接收方
func (wb *WeBase) RemoveReceiverProfitSharing(req *ProfitSharingAddReceiverRequest) (rsp *ProfitSharingAddReceiverResponse, err error) {
	req.MchID = wb.PayID
	req.AppID = wb.AppID
	req.NonceStr = util.RandString(32)
	req.SignType = util.SignTypeHMACSHA256
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
	params.Add("sign", util.SignSha256(params.Encode(), client.PayKey))

	data, err := util.HTTPXMLPost(client.HTTPClient, addReceiverProfitSharingURL, params)
	if err != nil {
		return nil, errmsg.GetError(errProfitSharingRsp, fmt.Sprintf("addReceiverProfitSharingURL rsp err:%s", err.Error()))
	}
	rsp = &ProfitSharingAddReceiverResponse{}
	err = xml.Unmarshal(data, &rsp)
	return
}
