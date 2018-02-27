package app

import (
	"encoding/json"

	"github.com/zozowind/wego/core"
)

const (
	TemplateUrl = core.WxApiUrl + "/cgi-bin/message/wxopen/template/send?access_token=%s"
)

type TemplateParam struct {
	TemplateId      string                   `json:"template_id"`
	ToUser          string                   `json:"touser"`
	Page            string                   `json:"page"`
	FormId          string                   `json:"form_id"`
	Data            map[string]*TemplateData `json:"data"`
	Color           string                   `json:"color"`
	EmphasisKeyword string                   `json:"emphasis_keyword"`
}

type TemplateData struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

func (this *WeAppClient) SendTemplateMessage(param *TemplateParam) (*core.WxErrorResponse, error) {
	data, err := this.Base.PostWithToken(TemplateUrl, param)
	if nil != err {
		return nil, err
	}
	errRes := &core.WxErrorResponse{}
	err = json.Unmarshal(data, errRes)
	return errRes, err
}
