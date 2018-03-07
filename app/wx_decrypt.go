package app

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/zozowind/wego/libs/errmsg"
)

var (
	illegalAesKey   = &errmsg.ErrMsg{Code: -41001, Message: "sessionKey length is error"}
	illegalIv       = &errmsg.ErrMsg{Code: -41002, Message: "iv length is error"}
	illegalBuffer   = &errmsg.ErrMsg{Code: -41003, Message: "illegalBuffer"}
	errDecodeBase64 = &errmsg.ErrMsg{Code: -41004, Message: "errDecodeBase64"}
)

// WxBizDataCrypt represents an active WxBizDataCrypt object
type WxBizDataCrypt struct {
	AppID      string
	SessionKey string
}

// Decrypt Weixin APP's AES Data
// If isJSON is true, Decrypt return JSON type.
// If isJSON is false, Decrypt return map type.
func (wxCrypt *WxBizDataCrypt) Decrypt(encryptedData string, iv string, isJSON bool) (interface{}, error) {
	if len(wxCrypt.SessionKey) != 24 {
		return nil, illegalAesKey
	}
	aesKey, err := base64.StdEncoding.DecodeString(wxCrypt.SessionKey)
	if err != nil {
		return nil, errmsg.GetError(errDecodeBase64, err.Error())
	}

	if len(iv) != 24 {
		return nil, illegalIv
	}
	aesIV, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		return nil, errmsg.GetError(errDecodeBase64, err.Error())
	}

	aesCipherText, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, errmsg.GetError(errDecodeBase64, err.Error())
	}
	aesPlantText := make([]byte, len(aesCipherText))

	aesBlock, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, errmsg.GetError(illegalBuffer, err.Error())
	}

	mode := cipher.NewCBCDecrypter(aesBlock, aesIV)
	mode.CryptBlocks(aesPlantText, aesCipherText)
	aesPlantText = PKCS7UnPadding(aesPlantText)

	var decrypted map[string]interface{}
	aesPlantText = []byte(strings.Replace(string(aesPlantText), "\a", "", -1))
	err = json.Unmarshal([]byte(aesPlantText), &decrypted)
	if err != nil {
		return nil, errmsg.GetError(illegalBuffer, err.Error())
	}

	if decrypted["watermark"].(map[string]interface{})["appid"] != wxCrypt.AppID {
		return nil, errmsg.GetError(illegalBuffer, "appId is not match")
	}

	if isJSON == true {
		return string(aesPlantText), nil
	}

	return decrypted, nil
}

// PKCS7UnPadding return unpadding []Byte plantText
func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unPadding := int(plantText[length-1])
	if unPadding < 1 || unPadding > 32 {
		unPadding = 0
	}
	return plantText[:(length - unPadding)]
}
