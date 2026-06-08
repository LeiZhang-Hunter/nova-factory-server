package api

import (
	"strings"
	"unicode"
)

// Match 记录某个步骤触发的归一化命中信息。
type Match struct {
	Step     string
	Kind     string
	Pattern  string
	Input    string
	Output   string
	Category string
}

// Category 表示命中的分类结果。
type Category struct {
	Name string `json:"name"`
	ID   int64  `json:"id"`
}

// Result 表示 pipeline 执行后的归一化结果。
type Result struct {
	Original   string
	Value      string
	Categories []Category
	Matches    []Match
	Metadata   map[string][]string
}

// Context 是 pipeline 在步骤间传递的上下文。
type Context struct {
	Original string
	Value    string

	categories []Category
	matches    []Match
	metadata   map[string][]string
}

// NewContext 创建归一化上下文。
func NewContext(input string) *Context {
	return &Context{
		Original: input,
		Value:    input,
		metadata: make(map[string][]string),
	}
}

// AddCategory 追加分类标签并去重。
func (c *Context) AddCategory(name string, id int64) {
	name = normalizeWhitespace(name)
	if name == "" {
		return
	}
	c.categories = appendUniqueCategory(c.categories, Category{
		Name: name,
		ID:   id,
	})
}

// AddMetadata 追加归一化元数据并去重。
func (c *Context) AddMetadata(key, value string) {
	key = normalizeWhitespace(key)
	value = normalizeWhitespace(value)
	if key == "" || value == "" {
		return
	}
	c.metadata[key] = appendUnique(c.metadata[key], value)
}

// AddMatch 记录命中信息。
func (c *Context) AddMatch(match Match) {
	match.Step = normalizeWhitespace(match.Step)
	match.Kind = normalizeWhitespace(match.Kind)
	match.Pattern = strings.TrimSpace(match.Pattern)
	c.matches = append(c.matches, match)
}

// Result 输出最终快照，避免暴露内部切片。
func (c *Context) Result() Result {
	result := Result{
		Original:   c.Original,
		Value:      c.Value,
		Categories: append([]Category(nil), c.categories...),
		Matches:    append([]Match(nil), c.matches...),
	}
	if len(c.metadata) > 0 {
		result.Metadata = make(map[string][]string, len(c.metadata))
		for key, values := range c.metadata {
			result.Metadata[key] = append([]string(nil), values...)
		}
	}
	return result
}

func appendUnique(values []string, value string) []string {
	for _, item := range values {
		if item == value {
			return values
		}
	}
	return append(values, value)
}

func appendUniqueCategory(values []Category, value Category) []Category {
	for _, item := range values {
		if item.Name == value.Name && item.ID == value.ID {
			return values
		}
	}
	return append(values, value)
}

func normalizeWhitespace(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	var builder strings.Builder
	builder.Grow(len(value))
	lastSpace := false
	for _, r := range value {
		if unicode.IsSpace(r) {
			if lastSpace {
				continue
			}
			builder.WriteByte(' ')
			lastSpace = true
			continue
		}
		builder.WriteRune(r)
		lastSpace = false
	}
	return strings.TrimSpace(builder.String())
}
