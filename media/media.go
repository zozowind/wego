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

	eventHandlers       map[message.EventType]func(*message.MixMessage) (*message.Reply, error)
	defaultEventHandler func(*message.MixMessage) (*message.Reply, error)
}
