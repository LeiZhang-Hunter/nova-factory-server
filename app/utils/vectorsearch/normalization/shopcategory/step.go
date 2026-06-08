package shopcategory

import (
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/constant/shop"
	"nova-factory-server/app/utils/store"
	"nova-factory-server/app/utils/vectorsearch/normalization/api"
	"nova-factory-server/app/utils/vectorsearch/normalization/util"
	"sort"
	"strconv"
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
	name    string
	enabled bool
	cache   store.IShopCategoryStore
}

type matchedCategory struct {
	id      int64
	name    string
	pattern string
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

func (c *Category) Init(config api.InterceptorConfig) error {
	if config.Name != "" {
		c.name = config.Name
	}
	c.enabled = true
	if config.Enabled != nil {
		c.enabled = *config.Enabled
	}
	c.cache = store.GetStore(shop.ShopCategoryStoreName)
	return nil
}

func (c *Category) Apply(ctx *api.Context) error {
	if ctx == nil || !c.enabled {
		return c.applyFromCache(ctx)
	}
	value := util.NormalizeWhitespace(ctx.Value)
	if value == "" {
		ctx.Value = ""
		return nil
	}
	ctx.Value = value
	categories := c.matchCategoriesFromCache(value)
	for _, matched := range categories {
		c.appendCategoryMatch(ctx, matched, matched.name)
	}
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

func (c *Category) applyFromCache(ctx *api.Context) error {
	if ctx == nil || !c.enabled {
		return nil
	}
	value := util.NormalizeWhitespace(ctx.Value)
	if value == "" {
		ctx.Value = ""
		return nil
	}
	ctx.Value = value
	for _, matched := range c.matchCategoriesFromCache(value) {
		c.appendCategoryMatch(ctx, matched, matched.name)
	}
	return nil
}

func (c *Category) appendCategoryMatch(ctx *api.Context, category matchedCategory, pattern string) {
	category.name = util.NormalizeWhitespace(category.name)
	pattern = util.NormalizeWhitespace(pattern)
	if category.name == "" {
		return
	}
	ctx.AddCategory(category.name, category.id)
	ctx.AddMetadata("category", category.name)
	if category.id > 0 {
		ctx.AddMetadata("category_id", strconv.FormatInt(category.id, 10))
	}
	ctx.AddMatch(api.Match{
		Step:     c.Name(),
		Kind:     "category",
		Pattern:  pattern,
		Input:    ctx.Value,
		Output:   ctx.Value,
		Category: category.name,
	})
}

func (c *Category) matchCategoriesFromCache(value string) []matchedCategory {
	if c.cache == nil {
		return nil
	}
	rows, ok := c.cache.Get()
	if !ok {
		return nil
	}
	matches := make([]matchedCategory, 0)
	seen := make(map[string]struct{})
	collectMatchedCategories(rows, value, &matches, seen)
	sort.SliceStable(matches, func(i, j int) bool {
		return len([]rune(matches[i].name)) > len([]rune(matches[j].name))
	})
	return matches
}

func collectMatchedCategories(rows []store.ShopCategoryData, value string, matches *[]matchedCategory, seen map[string]struct{}) {
	for _, row := range rows {
		if row == nil {
			continue
		}
		category, ok := row.(*shopmodels.CategoryInfo)
		if ok {
			name := util.NormalizeWhitespace(category.CategoryName)
			if name != "" && matchKeyword(value, name, MatchContains) {
				key := strconv.FormatInt(category.ID, 10) + ":" + name
				if _, exists := seen[key]; !exists {
					seen[key] = struct{}{}
					*matches = append(*matches, matchedCategory{
						id:      category.ID,
						name:    category.CategoryName,
						pattern: name,
					})
				}
			}
		}
		collectMatchedCategories(row.ChildrenData(), value, matches, seen)
	}
}
