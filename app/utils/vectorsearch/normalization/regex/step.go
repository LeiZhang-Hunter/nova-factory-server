package regex

import (
	"fmt"
	"nova-factory-server/app/utils/gateway/v1/config/cfg"
	"nova-factory-server/app/utils/vectorsearch/normalization/api"
	"nova-factory-server/app/utils/vectorsearch/normalization/util"
	"regexp"
)

const Type = "regex"

type Regex struct {
	name        string
	enabled     bool
	pattern     *regexp.Regexp
	replacement string
}

func init() {
	api.Register(NewRegex)
}

func NewRegex() api.Step {
	return &Regex{
		name:    Type,
		enabled: true,
	}
}

// Step 使用正则对当前字符串做替换归一。
func Step(name string, pattern *regexp.Regexp, replacement string) api.Step {
	stepName := util.NormalizeWhitespace(name)
	if stepName == "" {
		stepName = "regex_replace"
	}
	return &Regex{
		name:        stepName,
		enabled:     true,
		pattern:     pattern,
		replacement: replacement,
	}
}

func (r *Regex) Name() string {
	if r.name == "" {
		return Type
	}
	return r.name
}

func (r *Regex) Type() string {
	return Type
}

func (r *Regex) Apply(ctx *api.Context) error {
	if ctx == nil || !r.enabled || r.pattern == nil {
		return nil
	}
	before := ctx.Value
	after := r.pattern.ReplaceAllString(before, r.replacement)
	if before == after {
		return nil
	}
	ctx.Value = after
	ctx.AddMatch(api.Match{
		Step:    r.Name(),
		Kind:    "regex_replace",
		Pattern: r.pattern.String(),
		Input:   before,
		Output:  after,
	})
	return nil
}

func (r *Regex) Init(config api.InterceptorConfig) error {
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
	if stepConfig.Pattern == "" {
		return fmt.Errorf("regex pattern is empty")
	}
	pattern, err := regexp.Compile(stepConfig.Pattern)
	if err != nil {
		return fmt.Errorf("compile regex pattern failed: %w", err)
	}
	r.pattern = pattern
	r.replacement = stepConfig.Replacement
	return nil
}
