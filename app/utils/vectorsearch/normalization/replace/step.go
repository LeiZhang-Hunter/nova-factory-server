package replace

import (
	"fmt"
	"strings"

	"nova-factory-server/app/utils/gateway/v1/config/cfg"
	"nova-factory-server/app/utils/vectorsearch/normalization/api"
	"nova-factory-server/app/utils/vectorsearch/normalization/util"
)

const Type = "replace"

type Replace struct {
	name    string
	enabled bool
	old     string
	new     string
}

func init() {
	api.Register(NewReplace)
}

func NewReplace() api.Step {
	return &Replace{
		name:    Type,
		enabled: true,
	}
}

func (r *Replace) Name() string {
	if r.name == "" {
		return Type
	}
	return r.name
}

func (r *Replace) Type() string {
	return Type
}

func (r *Replace) Apply(ctx *api.Context) error {
	if ctx == nil || !r.enabled || r.old == "" {
		return nil
	}
	before := ctx.Value
	after := strings.ReplaceAll(before, r.old, r.new)
	if before == after {
		return nil
	}
	ctx.Value = after
	ctx.AddMatch(api.Match{
		Step:    r.Name(),
		Kind:    "replace",
		Pattern: r.old,
		Input:   before,
		Output:  after,
	})
	return nil
}

func (r *Replace) Init(config api.InterceptorConfig) error {
	if config.Name != "" {
		r.name = util.NormalizeWhitespace(config.Name)
	}
	r.enabled = true
	if config.Enabled != nil {
		r.enabled = *config.Enabled
	}

	stepConfig := Config{}
	if err := cfg.UnpackFromCommonCfg(config.Properties, &stepConfig).Do(); err != nil {
		return err
	}
	if stepConfig.Old == "" {
		return fmt.Errorf("replace old value is empty")
	}
	r.old = stepConfig.Old
	r.new = stepConfig.New
	return nil
}
