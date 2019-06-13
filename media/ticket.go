package media

import (
	"encoding/json"
	"fmt"

	"github.com/zozowind/wego/core"
	"github.com/zozowind/wego/util"
)

const (
	jsTicketURL = core.WxAPIURL + "/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"
)

//JsAPITicket jsapiticket
type JsAPITicket struct {
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	Ticket    string `json:"ticket"`
	ExpiresIn int64  `json:"expires_in"`
}

//TicketServer ticket server interface
type TicketServer interface {
	Ticket() (string, error)        // 获取 Ticket
	RefreshTicket() (string, error) //刷新 Ticket
}

// Ticket get wechat jsapi_ticket
func (wm *WeMediaClient) Ticket() (string, error) {
	return wm.TicketServer.Ticket()
}

//RequestTicket request wechat jsapi_ticket
func (wm *WeMediaClient) RequestTicket() (*JsAPITicket, error) {
	jsTicket := &JsAPITicket{}

	token, err := wm.Token()
	if nil != err {
		return jsTicket, err
	}
	//请求接口
	data, err := util.HTTPGet(wm.HTTPClient, fmt.Sprintf(jsTicketURL, token))
	if nil != err {
		return jsTicket, err
	}

	err = json.Unmarshal(data, jsTicket)
	if nil != err {
		return jsTicket, err
	}

	if jsTicket.ErrCode != 0 {
		err = fmt.Errorf("jsAPITicket err: %d, msg: %s", jsTicket.ErrCode, jsTicket.ErrMsg)
		return jsTicket, err
	}

	// 由于网络的延时, js_ticket 过期时间留有一个缓冲区
	switch {
	case jsTicket.ExpiresIn > 60*30:
		jsTicket.ExpiresIn -= 60 * 5
	case jsTicket.ExpiresIn > 60*2:
		jsTicket.ExpiresIn -= 60
	default:
		err = fmt.Errorf("expires_in too small: %d", jsTicket.ExpiresIn)
		return jsTicket, err
	}

	return jsTicket, nil
}
