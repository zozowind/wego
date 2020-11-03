package media

import (
	"encoding/json"
	"net/url"

	"github.com/zozowind/wego/core"
	"github.com/zozowind/wego/util"
)

//CustomerMessage .通用的结构
type CustomerMessage interface {
}

//CustomerCommonReq 公共结构
type CustomerCommonReq struct {
	ToUser  string `json:"touser"`
	MsgType string `json:"msgtype"`
}

//Customer .客服消息定义
const (
	TextMsgtype  = "text"
	ImageMsgtype = "image"

	customerSendURL = core.WxAPIURL + "/cgi-bin/message/custom/send"
)

//CustomersContentReq .文字信息请求参数
type CustomersContentReq struct {
	CustomerCommonReq
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
}

//CustomersImageReq  图片信息请求参数
type CustomersImageReq struct {
	CustomerCommonReq
	Image struct {
		MediaID string `json:"media_id"`
	} `json:"image"`
}

// SendCustomerMessage 获取客服信息发送
func (wm *WeMediaClient) SendCustomerMessage(req CustomerMessage) (err error) {
	token, err := wm.TokenServer.Token()
	if nil != err {
		return
	}
	params := url.Values{}
	params.Set("access_token", token)

	data, err := util.HTTPJsonPost(nil, customerSendURL+"?"+params.Encode(), req)
	rsp := &core.WxErrorResponse{}
	err = json.Unmarshal(data, rsp)
	if nil != err {
		return
	}
	err = rsp.Check()
	return
}
