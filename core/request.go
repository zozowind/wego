package core

import (
	"encoding/json"
	"fmt"

	"github.com/zozowind/wego/util"
)

const WxApiUrl = "https://api.weixin.qq.com"
const WxPayUrl = "https://api.mch.weixin.qq.com"

type WxErrorResponse struct {
	Code    int    `json:"errcode"`
	Message string `json:"errmsg"`
}

func (this *WeBase) PostWithToken(urlTemp string, param interface{}) ([]byte, error) {
	data := []byte{}
	token, err := this.Token()
	if nil != err {
		return data, err
	}
	url := fmt.Sprintf(urlTemp, token)
	data, err = util.HttpJsonPost(this.HttpClient, url, param)
	if nil != err {
		return data, err
	}
	errRes := &WxErrorResponse{}
	err = json.Unmarshal(data, errRes)
	if nil != err {
		return data, err
	}
	if nil == err && 0 != errRes.Code {
		return data, fmt.Errorf("code: %d, message: %s", errRes.Code, errRes.Message)
	}
	return data, nil
}
