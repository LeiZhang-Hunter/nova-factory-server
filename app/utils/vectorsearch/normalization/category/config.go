package category

// Rule 定义分类归一规则。
type Rule struct {
	Category string
	Keywords []string
	Mode     MatchMode
}

type Config struct {
	Rules []Rule `yaml:"rules,omitempty"`
}
