package interceptor

type ExtensionConfig struct {
	Order       int      `yaml:"order,omitempty" default:"900"`
	BelongTo    []string `yaml:"belongTo,omitempty"`
	IgnoreRetry bool     `yaml:"ignoreRetry,omitempty" default:"true"`
}
