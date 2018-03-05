package app

import "encoding/json"

// Watermark Watermark struct
type Watermark struct {
	Timestamp int    `json:"timestamp"`
	AppID     string `json:"appid"`
}

//UserDecryptInfo decypt wechat user info struct
type UserDecryptInfo struct {
	OpenID    string    `json:"openId"`
	Nickname  string    `json:"nickName"`
	Gender    int       `json:"gender"`
	Province  string    `json:"province"`
	City      string    `json:"city"`
	Country   string    `json:"country"`
	Avatar    string    `json:"avatarUrl"`
	Language  string    `json:"language"`
	UnionID   string    `json:"unionId"`
	Watermark Watermark `json:"watermark"`
}

//GroupDecryptInfo decypt wechat group info struct
type GroupDecryptInfo struct {
	OpenGId string `json:"openGId"`
}

//DecryptUserInfo decypt wechat user info
func (client *WeAppClient) DecryptUserInfo(encryptedData string, iv string, sessionKey string) (*UserDecryptInfo, error) {
	info := &UserDecryptInfo{}
	data, err := client.decrypt(encryptedData, iv, sessionKey)
	if nil == err {
		err = json.Unmarshal([]byte(data), info)
	}
	return info, err
}

//DecryptGroupInfo decypt wechat group info
func (client *WeAppClient) DecryptGroupInfo(encryptedData string, iv string, sessionKey string) (*GroupDecryptInfo, error) {
	info := &GroupDecryptInfo{}
	data, err := client.decrypt(encryptedData, iv, sessionKey)
	if nil == err {
		err = json.Unmarshal([]byte(data), info)
	}
	return info, err
}

func (client *WeAppClient) decrypt(encryptedData string, iv string, sessionKey string) (string, error) {
	pc := WxBizDataCrypt{
		AppID:      client.AppID,
		SessionKey: sessionKey,
	}
	//返回json string格式
	data, err := pc.Decrypt(encryptedData, iv, true)
	if nil != err {
		return "", err
	}

	return data.(string), err
}
