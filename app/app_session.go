package app

import (
	"encoding/json"
	"fmt"

	"github.com/zozowind/wego/core"
	"github.com/zozowind/wego/util"
)

const (
	GetSessionUriTemp = "/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
)

type WxGetSessionResponse struct {
	OpenId     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionId    string `json:"unionid"`
}

func (this *WeAppClient) GetSession(code string) (*WxGetSessionResponse, error) {
	res := &WxGetSessionResponse{}

	url := fmt.Sprintf(core.WxApiUrl+GetSessionUriTemp, this.Base.AppId, this.Base.AppSecret, code)
	data, err := util.HttpGet(this.Base.HttpClient, url)
	if err != nil {
		return res, err
	}

	errRes := &core.WxErrorResponse{}
	err = json.Unmarshal(data, errRes)
	if nil != err {
		return res, err
	} else if 0 != errRes.Code {
		return res, fmt.Errorf("code: %d, message: %s", errRes.Code, errRes.Message)
	}
	err = json.Unmarshal(data, res)
	return res, err
}
