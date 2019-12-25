package menu

//SetReq 设置菜单请求数据
type SetReq struct {
	Button    []*Button  `json:"button,omitempty"`
	MatchRule *MatchRule `json:"matchrule,omitempty"`
}

//DeleteConditionalReq 删除个性化菜单请求数据
type DeleteConditionalReq struct {
	MenuID int64 `json:"menuid"`
}

//TryMatchReq 菜单匹配请求
type TryMatchReq struct {
	UserID string `json:"user_id"`
}

//MenusRes 菜单返回
type MenusRes struct {
	Menu            *ConditionalMenuRes   `json:"menu"`
	ConditionalMenu []*ConditionalMenuRes `json:"conditionalmenu,omitempty"`
}

//ConditionalMenuRes 个性化菜单返回结果
type ConditionalMenuRes struct {
	Button    []Button  `json:"button"`
	MatchRule MatchRule `json:"matchrule"`
	MenuID    int64     `json:"menuid"`
}

//MatchRule 个性化菜单规则
type MatchRule struct {
	TagID              string `json:"tag_id,omitempty"`
	GroupID            string `json:"group_id,omitempty"`
	Sex                string `json:"sex,omitempty"`
	Country            string `json:"country,omitempty"`
	Province           string `json:"province,omitempty"`
	City               string `json:"city,omitempty"`
	ClientPlatformType string `json:"client_platform_type,omitempty"`
	Language           string `json:"language,omitempty"`
}
