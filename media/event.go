package media

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/zozowind/wego/core"
	"github.com/zozowind/wego/media/message"
	"github.com/zozowind/wego/util"
)

//EventHandler 事件处理
func (wm *WeMediaClient) EventHandler(r *http.Request) (response []byte, err error) {
	if r.FormValue("echostr") != "" {
		response = []byte(r.FormValue("echostr"))
		return
	}

	msg, random, err := wm.eventMessage(r)
	if nil != err {
		return
	}

	var reply *message.Reply
	if nil != wm.eventHandlers {
		handler, ok := wm.eventHandlers[msg.Event]
		if ok {
			reply, err = handler(msg)
		} else {
			if nil == wm.defaultEventHandler {
				wm.defaultEventHandler = defaultEventHandler
			}
			reply, err = wm.defaultEventHandler(msg)
		}
	}

	err = buildResponse(reply, msg.FromUserName, msg.ToUserName)
	if nil != err {
		return
	}

	if wm.MessageConfig.Mode == core.MessageModePlain {
		response, err = xml.Marshal(reply.MsgData)
	} else {
		var responseRawXMLMsg []byte
		responseRawXMLMsg, err = xml.Marshal(reply.MsgData)
		//安全模式下对消息进行加密
		var encryptedMsg []byte
		encryptedMsg, err = util.EncryptMsg(random, responseRawXMLMsg, wm.AppID, wm.MessageConfig.EncodingAESKey)
		if err != nil {
			return
		}
		//TODO 如果获取不到timestamp nonce 则自己生成
		timestamp := time.Now().Unix()
		nonce := util.RandString(8)
		msgSignature := util.StrSortSha1Sign([]string{
			wm.MessageConfig.Token,
			fmt.Sprintf("%d", timestamp),
			nonce,
			string(encryptedMsg),
		})
		replyMsg := message.ResponseEncryptedXMLMsg{
			EncryptedMsg: string(encryptedMsg),
			MsgSignature: msgSignature,
			Timestamp:    timestamp,
			Nonce:        nonce,
		}
		response, err = xml.Marshal(replyMsg)
	}
	return
}

//OnEvent 注册公众号的事件处理方法
func (wm *WeMediaClient) OnEvent(e message.EventType, fn func(*message.MixMessage) (*message.Reply, error)) {
	if wm.eventHandlers == nil {
		wm.eventHandlers = map[message.EventType]func(*message.MixMessage) (*message.Reply, error){}
	}
	wm.eventHandlers[e] = fn
}

//SetDefaultEventHandler 设置默认的事件处理
func (wm *WeMediaClient) SetDefaultEventHandler(fn func(*message.MixMessage) (*message.Reply, error)) {
	wm.defaultEventHandler = fn
}

func (wm *WeMediaClient) eventMessage(r *http.Request) (msg *message.MixMessage, random []byte, err error) {
	// safe mode
	var rawXMLMsgBytes []byte
	if wm.MessageConfig.Mode == core.MessageModePlain {
		rawXMLMsgBytes, err = ioutil.ReadAll(r.Body)
		if err != nil {
			//消息内容读取失败
			return
		}
	} else {
		encryptedXMLMsg := &message.EncryptedXMLMsg{}
		err = xml.NewDecoder(r.Body).Decode(encryptedXMLMsg)
		if err != nil {
			//消息内容解析失败
			return
		}

		//验证签名
		sign := r.FormValue("msg_signature")
		checkSign := util.StrSortSha1Sign([]string{
			wm.MessageConfig.Token,
			r.FormValue("timestamp"),
			r.FormValue("nonce"),
			encryptedXMLMsg.EncryptedMsg,
		})
		if sign != checkSign {
			err = errors.New("sign check fail")
			return
		}
		random, rawXMLMsgBytes, err = util.DecryptMsg(wm.AppID, encryptedXMLMsg.EncryptedMsg, wm.MessageConfig.EncodingAESKey)
		if nil != err {
			return
		}
	}

	//解析消息结构体
	msg = &message.MixMessage{}
	err = xml.Unmarshal(rawXMLMsgBytes, msg)
	return
}

func defaultEventHandler(msg *message.MixMessage) (reply *message.Reply, err error) {
	fmt.Printf("%v", msg)
	return
}

func buildResponse(reply *message.Reply, toUser string, fromUser string) (err error) {
	if nil == reply {
		return
	}
	msgType := reply.MsgType
	switch msgType {
	case message.MsgTypeText:
	case message.MsgTypeImage:
	case message.MsgTypeVoice:
	case message.MsgTypeVideo:
	case message.MsgTypeMusic:
	case message.MsgTypeNews:
	case message.MsgTypeTransfer:
	default:
		err = message.ErrUnsupportReply
		return
	}

	reply.MsgData.SetToUserName(toUser)
	reply.MsgData.SetFromUserName(fromUser)
	reply.MsgData.SetMsgType(msgType)
	reply.MsgData.SetCreateTime(time.Now())
	return
}
