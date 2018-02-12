package app

import "encoding/json"

type Watermark struct {
	Timestamp int    `json:"timestamp"`
	AppId     string `json:"appid"`
}

type UserDecryptInfo struct {
	OpenId    string    `json:"openId"`
	Nickname  string    `json:"nickName"`
	Gender    int       `json:"gender"`
	Province  string    `json:"province"`
	City      string    `json:"city"`
	Country   string    `json:"country"`
	Avatar    string    `json:"avatarUrl"`
	Language  string    `json:"language"`
	UnionId   string    `json:"unionId"`
	Watermark Watermark `json:"watermark"`
}

type GroupDecryptInfo struct {
	OpenGId string `json:"openGId"`
}

func (this *WeAppClient) DecryptUserInfo(encryptedData string, iv string, sessionKey string) (*UserDecryptInfo, error) {
	info := &UserDecryptInfo{}
	data, err := this.Decrypt(encryptedData, iv, sessionKey)
	if nil == err {
		err = json.Unmarshal([]byte(data), info)
	}
	return info, err
}

func (this *WeAppClient) DecryptGroupInfo(encryptedData string, iv string, sessionKey string) (*GroupDecryptInfo, error) {
	info := &GroupDecryptInfo{}
	data, err := this.Decrypt(encryptedData, iv, sessionKey)
	if nil == err {
		err = json.Unmarshal([]byte(data), info)
	}
	return info, err
}

func (this *WeAppClient) Decrypt(encryptedData string, iv string, sessionKey string) (string, error) {
	pc := WxBizDataCrypt{
		AppID:      this.Base.AppId,
		SessionKey: sessionKey,
	}
	//返回json string格式
	data, err := pc.Decrypt(encryptedData, iv, true)
	if nil != err {
		return "", err
	}

	return data.(string), err
}
