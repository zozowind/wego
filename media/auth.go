package media

import (
	"encoding/json"
	"net/url"

	"github.com/zozowind/wego/core"
	"github.com/zozowind/wego/util"
)

const (
	authorizeURL     = core.WxOpenURL + "/connect/oauth2/authorize"
	accessTokenURL   = core.WxAPIURL + "/sns/oauth2/access_token"
	refreshTokenURL  = core.WxAPIURL + "/sns/oauth2/refresh_token"
	responseTypeCode = "code"
	grantTypeRefresh = "refresh_token"
	grantTypeAuth    = "authorization_code"
)

// UserAccessTokenRsp authAccessToken返回数据
type UserAccessTokenRsp struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
}

//AuthCodeURL 获取链接
func (wm *WeMediaClient) AuthCodeURL(redirectURL string, scope string, state string) string {
	params := url.Values{}
	params.Set("appid", wm.AppID)
	params.Set("redirect_uri", redirectURL)
	params.Set("response_type", responseTypeCode)
	params.Set("scope", scope)
	params.Set("state", state)
	return authorizeURL + "?" + params.Encode() + "#wechat_redirect"
}

//GetUserAccessToken 获取用户AccessToken
func (wm *WeMediaClient) GetUserAccessToken(code string) (rsp *UserAccessTokenRsp, err error) {
	params := url.Values{}
	params.Set("appid", wm.AppID)
	params.Set("secret", wm.AppSecret)
	params.Set("code", "code")
	params.Set("grant_type", grantTypeAuth)
	data, err := util.HTTPGet(nil, accessTokenURL+"?"+params.Encode())
	if nil != err {
		return
	}
	rsp = &UserAccessTokenRsp{}
	err = json.Unmarshal(data, rsp)
	return
}

//RefreshUserAccessToken 刷新用户AccessToken
func (wm *WeMediaClient) RefreshUserAccessToken(refreshToken string) (rsp *UserAccessTokenRsp, err error) {
	params := url.Values{}
	params.Set("appid", wm.AppID)
	params.Set("refresh_token", refreshToken)
	params.Set("grant_type", grantTypeRefresh)
	data, err := util.HTTPGet(nil, refreshTokenURL+"?"+params.Encode())
	if nil != err {
		return
	}
	rsp = &UserAccessTokenRsp{}
	err = json.Unmarshal(data, rsp)
	return
}
