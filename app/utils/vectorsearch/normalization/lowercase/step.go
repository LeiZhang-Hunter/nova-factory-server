package lowercase

import (
	"nova-factory-server/app/utils/gateway/v1/config/cfg"
	"nova-factory-server/app/utils/vectorsearch/normalization/api"
	"strings"
)

const Type = "lowercase"

type Lowercase struct {
	name    string
	enabled bool
}

func init() {
	api.Register(NewLowercase)
}

func NewLowercase() api.Step {
	return &Lowercase{
		name:    Type,
		enabled: true,
	}
}

// Step 对当前字符串做小写归一。
func Step() api.Step {
	return NewLowercase()
}

func (l *Lowercase) Name() string {
	if l.name == "" {
		return Type
	}
	return l.name
}

func (l *Lowercase) Type() string {
	return Type
}

func (l *Lowercase) Apply(ctx *api.Context) error {
	if ctx == nil || !l.enabled {
		return nil
	}
	ctx.Value = strings.ToLower(ctx.Value)
	return nil
}

func (l *Lowercase) Init(config api.InterceptorConfig) error {
	if config.Name != "" {
		l.name = config.Name
	}
	l.enabled = true
	if config.Enabled != nil {
		l.enabled = *config.Enabled
	}
	stepConfig := Config{}
	return cfg.UnpackFromCommonCfg(config.Properties, &stepConfig).Do()
}
