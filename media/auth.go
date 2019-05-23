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
	userInfoURL      = core.WxAPIURL + "/sns/userinfo"
	responseTypeCode = "code"
	grantTypeRefresh = "refresh_token"
	grantTypeAuth    = "authorization_code"
)

// UserAccessTokenRsp authAccessToken返回数据
type UserAccessTokenRsp struct {
	core.WxErrorResponse
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	// UnionID      string `json:"unionid"`
}

//UserInfoRsp 用户信息返回
type UserInfoRsp struct {
	core.WxErrorResponse
	OpenID     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int      `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	HeadImgURL string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	UnionID    string   `json:"unionid"`
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
	params.Set("code", code)
	params.Set("grant_type", grantTypeAuth)
	data, err := util.HTTPGet(nil, accessTokenURL+"?"+params.Encode())
	if nil != err {
		return
	}
	rsp = &UserAccessTokenRsp{}
	err = json.Unmarshal(data, rsp)
	if nil != err {
		return
	}
	err = rsp.Check()
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
	if nil != err {
		return
	}
	err = rsp.Check()
	return
}

//UserInfo 用户信息
func (wm *WeMediaClient) UserInfo(accessToken string, openID string) (rsp *UserInfoRsp, err error) {
	params := url.Values{}
	params.Set("access_token", accessToken)
	params.Set("openid", openID)
	params.Set("lang", "zh_CN")
	data, err := util.HTTPGet(nil, userInfoURL+"?"+params.Encode())
	if nil != err {
		return
	}
	rsp = &UserInfoRsp{}
	err = json.Unmarshal(data, rsp)
	if nil != err {
		return
	}
	err = rsp.Check()
	return
}
