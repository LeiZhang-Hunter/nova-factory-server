package config

type Config interface {
	Status() bool
	Data() string
	OverrideURL() string
	Metadata() map[string]any
}
