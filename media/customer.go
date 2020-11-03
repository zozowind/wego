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
	Touser  string `json:"touser"`
	Msgtype string `json:"msgtype"`
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
	Text CustomersContentInfo `json:"text"`
}

//CustomersImageReq  图片信息请求参数
type CustomersImageReq struct {
	CustomerCommonReq
	Image CustomersImageInfo `json:"image"`
}

//CustomersImageInfo .图片信息请求图片信息
type CustomersImageInfo struct {
	MediaID string `json:"media_id"`
}

//CustomersContentInfo .文本信息请求文本信息
type CustomersContentInfo struct {
	Content string `json:"content"`
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
