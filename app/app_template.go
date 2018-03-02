package app

import (
	"encoding/json"

	"github.com/zozowind/wego/core"
)

const (
	templateURL = core.WxAPIURL + "/cgi-bin/message/wxopen/template/send?access_token=%s"
)

//TemplateParam template message request struct
type TemplateParam struct {
	TemplateID      string                   `json:"template_id"`
	ToUser          string                   `json:"touser"`
	Page            string                   `json:"page"`
	FormID          string                   `json:"form_id"`
	Data            map[string]*TemplateData `json:"data"`
	Color           string                   `json:"color"`
	EmphasisKeyword string                   `json:"emphasis_keyword"`
}

//TemplateData template message data struct
type TemplateData struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

//SendTemplateMessage send template message
func (client *WeAppClient) SendTemplateMessage(param *TemplateParam) (*core.WxErrorResponse, error) {
	data, err := client.Base.PostWithToken(templateURL, param)
	if nil != err {
		return nil, err
	}
	errRes := &core.WxErrorResponse{}
	err = json.Unmarshal(data, errRes)
	return errRes, err
}
