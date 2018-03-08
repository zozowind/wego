package work

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/zozowind/wego/core"
	"github.com/zozowind/wego/libs/errmsg"
)

const (
	sendMessageURL = WxWorkAPIURL + "/cgi-bin/message/send?access_token=%s"
	//MessageTargetUser message target key user
	MessageTargetUser = "user"
	//MessageTargetParty message target key party
	MessageTargetParty = "party"
	//MessageTargetTag message target key tag
	MessageTargetTag = "tag"
	//MessageTargetAll message target key all
	MessageTargetAll = "all"
)

// PushMessage  push message interface
type PushMessage interface {
	setMessageType()
	setToUser(string)
	setToParty(string)
	setToTag(string)
	setAgentID(int)
	setSafe(int)
}

//BaseMessage push message base struct
type BaseMessage struct {
	ToUser  string `json:"touser"`
	ToParty string `json:"toparty"`
	ToTag   string `json:"totag"`
	MsgType string `json:"msgtype"`
	AgentID int    `json:"agentid"`
	Safe    int    `json:"safe"`
}

//WeWorkMessageResponse wework push message response struct
type WeWorkMessageResponse struct {
	core.WxErrorResponse
	InvalidUser  string `json:"invaliduser"` // 不区分大小写，返回的列表都统一转为小写
	InvalidParty string `json:"invalidparty"`
	InvalidTag   string `json:"invalidtag"`
}

func (m *BaseMessage) setToUser(t string) {
	m.ToUser = t
}

func (m *BaseMessage) setToParty(t string) {
	m.ToParty = t
}

func (m *BaseMessage) setToTag(t string) {
	m.ToTag = t
}

func (m *BaseMessage) setAgentID(id int) {
	m.AgentID = id
}

func (m *BaseMessage) setSafe(safe int) {
	m.Safe = 1
	if safe == 0 {
		m.Safe = 0
	}
}

// InMessageContent message data for text message
type InMessageContent struct {
	Content string `json:"content"`
}

// InMessageMedia message data for image | voice | file  message
type InMessageMedia struct {
	MediaID string `json:"media_id"`
}

// InMessageInfo message data info for card | video message
type InMessageInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// InMessageVideo message data for video message
type InMessageVideo struct {
	InMessageMedia
	InMessageInfo
}

// InMessageCard message data for textcard message
type InMessageCard struct {
	InMessageInfo
	URL    string `json:"url"`
	BtnTxt string `json:"btntxt"`
}

// TextMessage push text message struct
type TextMessage struct {
	BaseMessage
	Text InMessageContent `json:"text"`
}

func (m *TextMessage) setMessageType() {
	m.MsgType = "text"
}

// ImageMessage push image message struct
type ImageMessage struct {
	BaseMessage
	Image InMessageMedia `json:"image"`
}

func (m *ImageMessage) setMessageType() {
	m.MsgType = MediaTypeImage
}

// VoiceMessage push voice message struct
type VoiceMessage struct {
	BaseMessage
	Voice InMessageMedia `json:"voice"`
}

func (m *VoiceMessage) setMessageType() {
	m.MsgType = MediaTypeVoice
}

// VideoMessage push video message struct
type VideoMessage struct {
	BaseMessage
	Video InMessageVideo `json:"video"`
}

func (m *VideoMessage) setMessageType() {
	m.MsgType = MediaTypeVideo
}

// FileMessage push file message struct
type FileMessage struct {
	BaseMessage
	File InMessageMedia `json:"file"`
}

func (m *FileMessage) setMessageType() {
	m.MsgType = MediaTypeFile
}

// CardMessage push textcard message struct
type CardMessage struct {
	BaseMessage
	TextCard InMessageCard `json:"textcard"`
}

func (m *CardMessage) setMessageType() {
	m.MsgType = "textcard"
}

//InMessageArticle article info
type InMessageArticle struct {
	InMessageCard
	PicURL string `json:"picurl"`
}

//InMessageMpArticle  info of article which storing in wework
type InMessageMpArticle struct {
	Title            string `json:"title"`
	ThumbMediaID     string `json:"thumb_media_id"`
	Author           string `json:"author"`
	ContentSourceURL string `json:"content_source_url"`
	Content          string `json:"content"`
	Digest           string `json:"digest"`
}

// NewsMessage mutiple article message, include 1-8 articles
type NewsMessage struct {
	BaseMessage
	Articles []InMessageArticle
}

func (m *NewsMessage) setMessageType() {
	m.MsgType = "news"
}

//MpNewsMessage mutiple article message, include 1-8 articles which store in we work
type MpNewsMessage struct {
	BaseMessage
	Articles []InMessageMpArticle
}

func (m *MpNewsMessage) setMessageType() {
	m.MsgType = "mpnews"
}

func (w *WeWorkClient) completeMessage(toType string, ids []string, safe int, message PushMessage) (PushMessage, error) {
	//设置消息类型
	switch toType {
	case MessageTargetUser:
		message.setToUser(strings.Join(ids, "|"))
	case MessageTargetParty:
		message.setToParty(strings.Join(ids, "|"))
	case MessageTargetTag:
		message.setToTag(strings.Join(ids, "|"))
	case MessageTargetAll:
		message.setToUser("@all")
	default:
		return message, fmt.Errorf("wrong error type %s", toType)
	}
	message.setMessageType()
	message.setSafe(safe)
	message.setAgentID(w.AgentID)
	return message, nil
}

//SendMessage send message to target
func (w *WeWorkClient) SendMessage(toType string, ids []string, safe int, message PushMessage) (*WeWorkMessageResponse, error) {
	message, err := w.completeMessage(toType, ids, safe, message)
	if nil != err {
		return nil, errmsg.GetError(errSendMessage, err.Error())
	}

	data, err := w.PostWithToken(sendMessageURL, message)
	if nil != err {
		return nil, errmsg.GetError(errSendMessage, err.Error())
	}

	res := &WeWorkMessageResponse{}
	err = json.Unmarshal(data, res)
	if nil != err {
		return res, errmsg.GetError(errSendMessage, err.Error())
	}

	err = res.Check()
	if nil != err {
		return res, errmsg.GetError(errDownloadMedia, err.Error())
	}
	return res, nil
}
