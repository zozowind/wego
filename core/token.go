package core

import (
	"encoding/json"
	"fmt"

	"github.com/zozowind/wego/util"
)

const (
	accessTokenURL = WxAPIURL + "/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
)

//TokenServer token server interface
type TokenServer interface {
	Token() (string, error)        // 获取 access_token
	RefreshToken() (string, error) //刷新 access_token
}

//AccessToken wechat access token struct
type AccessToken struct {
	Token     string `json:"access_token"`
	ExpiresIn int64  `json:"expires_in"`
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
}

//Token get wechat access token
func (wb *WeBase) Token() (string, error) {
	return wb.TokenServer.Token()
}

//RequestToken request wechat access token
func (wb *WeBase) RequestToken() (*AccessToken, error) {
	acceseToken := &AccessToken{}
	//请求接口
	data, err := util.HTTPGet(wb.HTTPClient, fmt.Sprintf(accessTokenURL, wb.AppID, wb.AppSecret))
	if nil != err {
		return acceseToken, err
	}

	err = json.Unmarshal(data, acceseToken)
	if nil != err {
		return acceseToken, err
	}

	if acceseToken.ErrCode != 0 {
		err = fmt.Errorf("accessToken err: %d, msg: %s", acceseToken.ErrCode, acceseToken.ErrMsg)
		return acceseToken, err
	}

	// 由于网络的延时, access_token 过期时间留有一个缓冲区
	switch {
	case acceseToken.ExpiresIn > 60*30:
		acceseToken.ExpiresIn -= 60 * 5
	case acceseToken.ExpiresIn > 60*2:
		acceseToken.ExpiresIn -= 60
	default:
		err = fmt.Errorf("expires_in too small: %d", acceseToken.ExpiresIn)
		return acceseToken, err
	}

	return acceseToken, nil
}
