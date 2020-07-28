package pay

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"

	"github.com/zozowind/wego/core"
	"github.com/zozowind/wego/util"
)

const (
	orderUnifiedURL      = "/pay/unifiedorder"
	orderQueryURL        = "/pay/orderquery"
	orderCloseURL        = "/pay/closeorder"
	refundURL            = "/secapi/pay/refund"
	refundQueryURL       = "/pay/refundquery"
	billDownloadURL      = "/pay/downloadbill"
	fundFlowDownloadURL  = "/pay/downloadfundflow"
	itilReportURL        = "/payitil/report"
	batchQueryCommentURL = "/billcommentsp/batchquerycomment"
	sandboxGetSignKeyURL = "/pay/getsignkey"

	//CodeSuccess 成功
	CodeSuccess = "SUCCESS"

	//TradeTypeJSAPI JSAPI支付（或小程序支付）
	TradeTypeJSAPI = "JSAPI"
	//TradeTypeNative Native支付
	TradeTypeNative = "NATIVE"
	//TradeTypeAPP app支付
	TradeTypeAPP = "APP"
	//TradeTypeMWeb H5支付
	TradeTypeMWeb = "MWEB"
	//TradeTypeMicroPay 付款码支付
	TradeTypeMicroPay = "MICROPAY"
	// CurrencyCNY 境内商户号仅支持人民币
	CurrencyCNY = "CNY"

	//TradeStateSuccess SUCCESS—支付成功
	TradeStateSuccess = "SUCCESS"
	//TradeStateRefund REFUND—转入退款
	TradeStateRefund = "REFUND"
	//TradeStateNopay NOTPAY—未支付
	TradeStateNopay = "NOTPAY"
	//TradeStateClosed CLOSED—已关闭
	TradeStateClosed = "CLOSED"
	//TradeStateRevoked REVOKED—已撤销（付款码支付）
	TradeStateRevoked = "REVOKED"
	//TradeStateUserPaying USERPAYING--用户支付中（付款码支付）
	TradeStateUserPaying = "USERPAYING"
	//TradeStatePayError PAYERROR--支付失败(其他原因，如银行返回失败)
	TradeStatePayError = "PAYERROR"
	//CouponTypeCash 充值代金券
	CouponTypeCash = "CASH"
	//CouponTypeNoCash 非充值优惠券
	CouponTypeNoCash = "NO_CASH"
	//RefundSourceUnsettled 未结算资金退款（默认使用未结算资金退款）
	RefundSourceUnsettled = "REFUND_SOURCE_UNSETTLED_FUNDS"
	//RefundSourceRecharge 可用余额退款
	RefundSourceRecharge = "REFUND_SOURCE_RECHARGE_FUNDS"

	//RefundChannelOriginal ORIGINAL—原路退款
	RefundChannelOriginal = "ORIGINAL"
	//RefundChannelBalance BALANCE—退回到余额
	RefundChannelBalance = "BALANCE"
	//RefundChannelOtherBalance OTHER_BALANCE—原账户异常退到其他余额账户
	RefundChannelOtherBalance = "OTHER_BALANCE"
	//RefundChannelOtherBankcard OTHER_BANKCARD—原银行卡异常退到其他银行卡
	RefundChannelOtherBankcard = "OTHER_BANKCARD"
	//RefundStatusSuccess SUCCESS—退款成功
	RefundStatusSuccess = "SUCCESS"
	//RefundStatusRefundClose REFUNDCLOSE—退款关闭。
	RefundStatusRefundClose = "REFUNDCLOSE"
	//RefundStatusProcessing PROCESSING—退款处理中
	RefundStatusProcessing = "PROCESSING"
	//RefundStatusChange CHANGE—退款异常，退款到银行发现用户的卡作废或者冻结了，导致原路退款银行卡失败，可前往商户平台（pay.weixin.qq.com）-交易中心，手动处理此笔退款。
	RefundStatusChange = "CHANGE"
	//TarTypeGZIP 压缩形式
	TarTypeGZIP = "GZIP"
	//BillTypeAll ALL（默认值），返回当日所有订单信息（不含充值退款订单）
	BillTypeAll = "ALL"
	//BillTypeSuccess SUCCESS，返回当日成功支付的订单（不含充值退款订单）
	BillTypeSuccess = "SUCCESS"
	//BillTypeRefund REFUND，返回当日退款订单（不含充值退款订单）
	BillTypeRefund = "REFUND"
	//BillTypeRechargeRefund RECHARGE_REFUND，返回当日充值退款订单
	BillTypeRechargeRefund = "RECHARGE_REFUND"
	//AccountTypeBasic Basic  基本账户
	AccountTypeBasic = "Basic"
	//AccountTypeOperation Operation
	AccountTypeOperation = "Operation"
	//AccountTypeFees Fees 手续费账户
	AccountTypeFees = "Fees"
)

//Client 支付客户端
type Client struct {
	PayID         string //支付账号，一般为商户账号
	PayKey        string //支付key
	HTTPClient    *http.Client
	IsSandbox     bool
	SandboxPayKey string
}

//InitClient 初始化支付
func InitClient(payID string, payKey string, isSandbox bool) (cli *Client, err error) {
	cli = &Client{
		PayID:     payID,
		PayKey:    payKey,
		IsSandbox: isSandbox,
	}
	if cli.IsSandbox {
		err = cli.sandboxSignKey()
	}
	return
}

//Req 支付接口
type Req interface {
	setMchID(string)
	setNonceStr(string)
	setSignType(string)
}

//CommonReq 通用参数
type CommonReq struct {
	AppID    string `xml:"appid"`     // 必填， 微信分配的小程序ID
	MchID    string `xml:"mch_id"`    // 必填， 微信支付分配的商户号
	NonceStr string `xml:"nonce_str"` // 必填， 随机字符串，不长于32位。NOTE: 如果为空则系统会自动生成一个随机字符串。
	SignType string `xml:"sign_type"` // 选填， 签名类型，默认为MD5，支持HMAC-SHA256和MD5
}

func (req *CommonReq) setMchID(str string) {
	req.MchID = str
}
func (req *CommonReq) setNonceStr(str string) {
	req.NonceStr = str
}
func (req *CommonReq) setSignType(str string) {
	req.SignType = str
}

//Rsp 支付返回数据接口
type Rsp interface {
	getReturnCode() string
	getReturnMsg() string
	getResultCode() string
	getErrCode() string
	getErrCodeDes() string
	getSign() string
	setSign(str string)
}

//CommonRsp 通用返回
type CommonRsp struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	//以下字段在return_code为SUCCESS的时候有返回
	AppID      string `xml:"appid"`     //微信分配的小程序ID
	MchID      string `xml:"mch_id"`    //微信支付分配的商户号
	NonceStr   string `xml:"nonce_str"` // 随机字符串，不长于32位。NOTE: 如果为空则系统会自动生成一个随机字符串。
	Sign       string `xml:"sign"`      //签名
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
}

func (rsp *CommonRsp) getReturnCode() string {
	return rsp.ReturnCode
}

func (rsp *CommonRsp) getReturnMsg() string {
	return rsp.ReturnMsg
}

func (rsp *CommonRsp) getResultCode() string {
	return rsp.ResultCode
}

func (rsp *CommonRsp) getErrCode() string {
	return rsp.ErrCode
}

func (rsp *CommonRsp) getErrCodeDes() string {
	return rsp.ErrCodeDes
}
func (rsp *CommonRsp) getSign() string {
	return rsp.Sign
}
func (rsp *CommonRsp) setSign(str string) {
	rsp.Sign = str
}

func (cli *Client) makeURL(url string) string {
	s := ""
	if cli.IsSandbox {
		s = "/sandboxnew"
	}
	return core.WxPayURL + s + url
}

func (cli *Client) getSignKey() string {
	key := cli.PayKey
	if cli.IsSandbox {
		key = cli.SandboxPayKey
	}
	return key
}

func (cli *Client) request(url string, req Req, rsp Rsp) (err error) {
	params, err := cli.makeParams(req)
	if nil != err {
		return
	}

	data, err := util.HTTPXMLPost(cli.HTTPClient, cli.makeURL(url), params)
	if nil != err {
		return
	}

	err = xml.Unmarshal(data, rsp)
	if nil != err {
		return
	}
	err = cli.checkRsp(req, rsp)
	return
}

func (cli *Client) makeParams(req Req) (params url.Values, err error) {
	req.setMchID(cli.PayID)
	req.setNonceStr(util.RandString(32))
	//基于安全统一使用SignTypeHMACSHA256
	if cli.IsSandbox {
		req.setSignType(util.SignTypeMD5)
	} else {
		req.setSignType(util.SignTypeHMACSHA256)
	}
	params, err = util.StructToURLValue(req, "xml")
	if nil != err {
		return
	}
	if cli.IsSandbox {
		params.Add("sign", util.SignMd5(params.Encode(), cli.getSignKey()))
	} else {
		params.Add("sign", util.SignSha256(params.Encode(), cli.getSignKey()))
	}
	return
}

func (cli *Client) checkRsp(req Req, rsp Rsp) (err error) {
	//校验结果
	if rsp.getReturnCode() != CodeSuccess {
		err = fmt.Errorf("return code: %s, message: %s", rsp.getReturnCode(), rsp.getReturnMsg())
		return
	}

	if rsp.getResultCode() != CodeSuccess {
		err = fmt.Errorf("result_code:%s, err_code: %s, err_code_des: %s", rsp.getResultCode(), rsp.getErrCode(), rsp.getErrCodeDes())
		return
	}

	//校验签名
	sign := rsp.getSign()
	rsp.setSign("")
	resData, err := util.StructToURLValue(rsp, "xml")
	if nil != err {
		return
	}

	var ok bool
	if cli.IsSandbox {
		ok = util.CheckSignMd5(resData.Encode(), cli.getSignKey(), sign)
	} else {
		ok = util.CheckSingSha256(resData.Encode(), cli.getSignKey(), sign)
	}

	if !ok {
		err = fmt.Errorf(fmt.Sprintf("sign invalid, response: %#v, signStr: %s", rsp, resData))
	}
	return
}
