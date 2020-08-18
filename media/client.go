package media

import (
	"github.com/zozowind/wego/core"
	"github.com/zozowind/wego/media/message"
)

//WeMediaClient wechat media client struct
type WeMediaClient struct {
	core.WeBase
	TicketServer  TicketServer
	MessageConfig *core.MessageConfig
	AuthHosts     []string //业务认证域名

	messageHandlers map[message.MsgType]func(*message.MixMessage) (*message.Reply, error)
	defaultHandler  func(*message.MixMessage) (*message.Reply, error)
}
