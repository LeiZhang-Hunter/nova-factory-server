package config

type Config interface {
	GetStatus() *bool
	GetData() string
	GetOverrideURL() string
	GetMetadata() map[string]any
}
