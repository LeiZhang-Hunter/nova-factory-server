package category

import (
	"nova-factory-server/app/utils/gateway/v1/config/cfg"
	"nova-factory-server/app/utils/vectorsearch/normalization/api"
	"nova-factory-server/app/utils/vectorsearch/normalization/util"
	"strings"
)

const Type = "product_type"

// MatchMode 定义分类归一时的匹配方式。
type MatchMode int

const (
	MatchExact MatchMode = iota
	MatchContains
	MatchPrefix
	MatchSuffix
)

type Category struct {
	name     string
	enabled  bool
	filtered []Rule
}

func init() {
	api.Register(NewCategory)
}

func NewCategory() api.Step {
	return &Category{
		name:    Type,
		enabled: true,
	}
}

func (c *Category) Name() string {
	if c.name == "" {
		return Type
	}
	return c.name
}

func (c *Category) Type() string {
	return Type
}

func (c *Category) Apply(ctx *api.Context) error {
	if ctx == nil || !c.enabled || len(c.filtered) == 0 {
		return nil
	}
	value := util.NormalizeWhitespace(ctx.Value)
	if value == "" {
		ctx.Value = ""
		return nil
	}
	ctx.Value = value
	for _, rule := range c.filtered {
		for _, keyword := range rule.Keywords {
			keyword = util.NormalizeWhitespace(keyword)
			if keyword == "" || !matchKeyword(value, keyword, rule.Mode) {
				continue
			}
			ctx.AddCategory(rule.Category)
			ctx.AddMetadata("category", rule.Category)
			ctx.AddMatch(api.Match{
				Step:     c.Name(),
				Kind:     "category",
				Pattern:  keyword,
				Input:    value,
				Output:   value,
				Category: rule.Category,
			})
			break
		}
	}
	return nil
}

func (c *Category) Init(config api.InterceptorConfig) error {
	if config.Name != "" {
		c.name = config.Name
	}
	c.enabled = true
	if config.Enabled != nil {
		c.enabled = *config.Enabled
	}
	stepConfig := Config{}
	unpack := cfg.UnpackFromCommonCfg(config.Properties, &stepConfig)
	err := unpack.Do()
	if err != nil {
		return err
	}
	c.filtered = stepConfig.Rules
	return nil
}

func matchKeyword(value, keyword string, mode MatchMode) bool {
	switch mode {
	case MatchExact:
		return value == keyword
	case MatchPrefix:
		return strings.HasPrefix(value, keyword)
	case MatchSuffix:
		return strings.HasSuffix(value, keyword)
	case MatchContains:
		fallthrough
	default:
		return strings.Contains(value, keyword)
	}
}
