package pay

import (
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/zozowind/wego/util"
)

type getSignKeyReq struct {
	MchID    string `xml:"mch_id"`    // 必填， 微信支付分配的商户号
	NonceStr string `xml:"nonce_str"` // 必填， 随机字符串，不长于32位。NOTE: 如果为空则系统会自动生成一个随机字符串。
}

type getSignKeyRsp struct {
	ReturnCode     string `xml:"return_code"`
	ReturnMsg      string `xml:"return_msg"`
	MchID          string `xml:"mch_id"`          //微信支付分配的商户号
	SandboxSignKey string `xml:"sandbox_signkey"` //返回的沙箱密钥
}

func (cli *Client) sandboxSignKey() (err error) {
	if !cli.IsSandbox {
		err = errors.New("current is not sandbox")
		return
	}
	req := &getSignKeyReq{
		MchID:    cli.PayID,
		NonceStr: util.RandString(32),
	}
	params, err := util.StructToURLValue(req, "xml")
	if nil != err {
		return
	}
	params.Add("sign", util.SignMd5(params.Encode(), cli.PayKey))

	data, err := util.HTTPXMLPost(cli.HTTPClient, cli.makeURL(sandboxGetSignKeyURL), params)
	if nil != err {
		return
	}
	rsp := &getSignKeyRsp{}
	err = xml.Unmarshal(data, rsp)
	if nil != err {
		return
	}
	//校验结果
	if rsp.ReturnCode != CodeSuccess {
		err = fmt.Errorf("code: %s, message: %s", rsp.ReturnCode, rsp.ReturnMsg)
		return
	}
	fmt.Printf("%v", rsp)
	cli.SandboxPayKey = rsp.SandboxSignKey
	return

}
