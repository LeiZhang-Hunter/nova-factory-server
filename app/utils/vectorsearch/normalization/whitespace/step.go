package whitespace

import (
	"nova-factory-server/app/utils/gateway/v1/config/cfg"
	"nova-factory-server/app/utils/vectorsearch/normalization/api"
	"nova-factory-server/app/utils/vectorsearch/normalization/util"
)

const Type = "whitespace"

type Whitespace struct {
	name    string
	enabled bool
}

func init() {
	api.Register(NewWhitespace)
}

func NewWhitespace() api.Step {
	return &Whitespace{
		name:    Type,
		enabled: true,
	}
}

// Step 统一压缩并裁剪空白字符。
func Step() api.Step {
	return NewWhitespace()
}

func (w *Whitespace) Name() string {
	if w.name == "" {
		return Type
	}
	return w.name
}

func (w *Whitespace) Type() string {
	return Type
}

func (w *Whitespace) Apply(ctx *api.Context) error {
	if ctx == nil || !w.enabled {
		return nil
	}
	ctx.Value = util.NormalizeWhitespace(ctx.Value)
	return nil
}

func (w *Whitespace) Init(config api.InterceptorConfig) error {
	if config.Name != "" {
		w.name = config.Name
	}
	w.enabled = true
	if config.Enabled != nil {
		w.enabled = *config.Enabled
	}
	stepConfig := Config{}
	return cfg.UnpackFromCommonCfg(config.Properties, &stepConfig).Do()
}
