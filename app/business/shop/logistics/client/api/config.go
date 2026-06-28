package api

// Config 快递配置
type Config interface {
	GetData() string

	GetType() string
}
