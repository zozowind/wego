package core

const (
	//MessageModePlain 明文
	MessageModePlain = iota //0
	//MessageModeCompa 兼容
	MessageModeCompa //1
	//MessageModeEncry 加密
	MessageModeEncry //2
	//MessageTypeJSON json格式
	MessageTypeJSON = "json"
	//MessageTypeXML  xml格式
	MessageTypeXML = "xml"
)

//MessageConfig 消息接收设置
type MessageConfig struct {
	Token          string
	EncodingAESKey string
	Mode           int
	Type           string
}

// WxBizMsgCrypt 消息解密加密器
type WxBizMsgCrypt struct {
	Token          string
	EncodingAesKey string
	AppID          string
}

// Decrypt 消息解密
func (wxCrypt *WxBizMsgCrypt) Decrypt() {

}

// Encrypt 消息加密
func (wxCrypt *WxBizMsgCrypt) Encrypt() {

}
