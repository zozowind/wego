package menu

//Button 菜单按钮
type Button struct {
	Type       string    `json:"type,omitempty"`
	Name       string    `json:"name,omitempty"`
	Key        string    `json:"key,omitempty"`
	URL        string    `json:"url,omitempty"`
	MediaID    string    `json:"media_id,omitempty"`
	AppID      string    `json:"appid,omitempty"`
	PagePath   string    `json:"pagepath,omitempty"`
	SubButtons []*Button `json:"sub_button,omitempty"`
}

//NewSubButton 二级菜单
func NewSubButton(name string, subButtons []*Button) *Button {
	return &Button{
		Name:       name,
		SubButtons: subButtons,
	}
}

//NewClickButton clict Btn类型
func NewClickButton(name, key string) *Button {
	return &Button{
		Type: "click",
		Key:  key,
		Name: name,
	}
}

//NewViewButton view 类型Btn
func NewViewButton(name, url string) *Button {
	return &Button{
		Type: "view",
		Name: name,
		URL:  url,
	}
}

//NewScanCodePushButton 扫码推事件
func NewScanCodePushButton(name, key string) *Button {
	return &Button{
		Type: "scancode_push",
		Name: name,
		Key:  key,
	}
}

//NewScanCodeWaitMsgButton  设置 扫码推事件且弹出"消息接收中"提示框
func NewScanCodeWaitMsgButton(name, key string) *Button {
	return &Button{
		Type: "scancode_waitmsg",
		Name: name,
		Key:  key,
	}
}

//NewPicSysPhotoButton  设置弹出系统拍照发图按钮
func NewPicSysPhotoButton(name, key string) *Button {
	return &Button{
		Type: "pic_sysphoto",
		Name: name,
		Key:  key,
	}
}

//NewPicPhotoOrAlbumButton  设置弹出系统拍照发图按钮
func NewPicPhotoOrAlbumButton(name, key string) *Button {
	return &Button{
		Type: "pic_photo_or_album",
		Name: name,
		Key:  key,
	}
}

//NewPicWeixinButton  设置弹出微信相册发图器类型按钮
func NewPicWeixinButton(name, key string) *Button {
	return &Button{
		Type: "pic_weixin",
		Name: name,
		Key:  key,
	}
}

//NewLocationSelectButton  设置 弹出地理位置选择器 类型按钮
func NewLocationSelectButton(name, key string) *Button {
	return &Button{
		Type: "location_select",
		Name: name,
		Key:  key,
	}
}

//NewMediaIDButton  设置 下发消息(除文本消息) 类型按钮
func NewMediaIDButton(name, mediaID string) *Button {
	return &Button{
		Type:    "media_id",
		Name:    name,
		MediaID: mediaID,
	}
}

//NewViewLimitedButton  设置 跳转图文消息URL 类型按钮
func NewViewLimitedButton(name, mediaID string) *Button {
	return &Button{
		Type:    "view_limited",
		Name:    name,
		MediaID: mediaID,
	}
}

//NewMiniprogramButton  设置 跳转小程序 类型按钮 (公众号后台必须已经关联小程序)
func NewMiniprogramButton(name, url, appID, pagePath string) *Button {
	return &Button{
		Type:     "miniprogram",
		Name:     name,
		URL:      url,
		AppID:    appID,
		PagePath: pagePath,
	}
}
