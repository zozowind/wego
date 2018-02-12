package core

type TokenServer interface {
	Get() (string, error)                 // 获取 access_token
	Refresh(token string) (string, error) // 刷新 access_token
}
