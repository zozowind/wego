package media

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"net/url"
	"regexp"
	"time"

	"github.com/zozowind/wego/util"
)

//WeJSConfig jsconfig
type WeJSConfig struct {
	NonceStr    string `json:"noncestr" query:"noncestr"`
	JsAPITicket string `json:"jsapi_ticket" query:"jsapi_ticket"`
	Timestamp   int64  `json:"timestamp" query:"timestamp"`
	URL         string `json:"url" query:"url"`
	Signature   string `json:"signature" query:"-"`
}

func (conf *WeJSConfig) sign() error {
	value, err := util.StructToURLValue(conf, "query")
	if nil != err {
		return err
	}
	str, err := url.QueryUnescape(value.Encode())
	if nil != err {
		return err
	}
	h := sha1.New()
	io.WriteString(h, str)
	conf.Signature = hex.EncodeToString(h.Sum(nil))
	return nil
}

//JsConfig 生成JSconfig
func (wm *WeMediaClient) JsConfig(url string) (*WeJSConfig, error) {
	//处理url
	var re = regexp.MustCompile(`#.*`)
	url = re.ReplaceAllString(url, "")
	ticket, err := wm.TicketServer.Ticket()
	if nil != err {
		return nil, err
	}

	conf := &WeJSConfig{
		NonceStr:    util.RandString(16),
		JsAPITicket: ticket,
		Timestamp:   time.Now().Unix(),
		URL:         url,
	}
	err = conf.sign()
	return conf, err
}
