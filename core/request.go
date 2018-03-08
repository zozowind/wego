package core

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/zozowind/wego/util"
)

const (
	//WxAPIURL wechat api url
	WxAPIURL = "https://api.weixin.qq.com"
	//WxPayURL wechat api url
	WxPayURL = "https://api.mch.weixin.qq.com"
)

//WxErrorResponse wechat common api response
type WxErrorResponse struct {
	Code    int    `json:"errcode"`
	Message string `json:"errmsg"`
}

//Check check wx response if has error
func (res *WxErrorResponse) Check() error {
	if res.Code != 0 {
		return fmt.Errorf("code: %d, message: %s", res.Code, res.Message)
	}
	return nil
}

//PostResponseWithToken wechat request with token
func (wb *WeBase) PostResponseWithToken(urlTemp string, param interface{}) ([]byte, error) {
	data := []byte{}
	token, err := wb.Token()
	if nil != err {
		return data, err
	}
	url := fmt.Sprintf(urlTemp, token)
	data, err = util.HTTPJsonPost(wb.HTTPClient, url, param)
	return data, err
}

//GetResponseWithToken wechat request with token
func (wb *WeBase) GetResponseWithToken(urlTemp string, params url.Values) ([]byte, error) {
	data := []byte{}
	token, err := wb.Token()
	if nil != err {
		return data, err
	}
	url := fmt.Sprintf(urlTemp, token)
	sep := "?"
	if strings.Contains(url, sep) {
		sep = "&"
	}
	url += sep + params.Encode()
	data, err = util.HTTPGet(wb.HTTPClient, url)
	return data, err
}

//PostWithToken post with token
func (wb *WeBase) PostWithToken(urlTemp string, param interface{}) ([]byte, error) {
	data, err := wb.PostResponseWithToken(urlTemp, param)
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
