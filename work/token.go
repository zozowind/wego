package work

import (
	"encoding/json"
	"fmt"

	"github.com/zozowind/wego/core"
	"github.com/zozowind/wego/util"
)

const (
	accessTokenURL = WxWorkAPIURL + "/cgi-bin/gettoken?corpid=%s&corpsecret=%s"
)

//RequestToken request wechat access token
func (w *WeWorkClient) RequestToken() (*core.AccessToken, error) {
	acceseToken := &core.AccessToken{}
	//请求接口
	data, err := util.HTTPGet(w.HTTPClient, fmt.Sprintf(accessTokenURL, w.AppID, w.AppSecret))
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
