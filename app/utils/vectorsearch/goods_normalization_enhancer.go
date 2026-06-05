package vectorsearch

import (
	"fmt"
	"nova-factory-server/app/utils/gateway/v1/config/cfg"
	"nova-factory-server/app/utils/vectorsearch/normalization"
	normalizationapi "nova-factory-server/app/utils/vectorsearch/normalization/api"
	"nova-factory-server/app/utils/vectorsearch/normalization/regex"
	replacepkg "nova-factory-server/app/utils/vectorsearch/normalization/replace"
	"nova-factory-server/app/utils/vectorsearch/normalization/shopcategory"
	"nova-factory-server/app/utils/vectorsearch/normalization/whitespace"
)

const goodsSearchSpecPattern = `(?i)^.*?((?:d)?\d+(?:\.\d+)?\*\d+(?:\.\d+)?\*\d+(?:\.\d+)?|(?:d)?\d+(?:\.\d+)?\*\d+(?:\.\d+)?(?:[-\s]+(?:d)?\d+(?:\.\d+)?\*\d+(?:\.\d+)?)?).*$`

// GoodsNormalizationEnhancer 商品归一化处理器
type GoodsNormalizationEnhancer struct {
	normalPipeline *normalization.Pipeline
}

// NewGoodsNormalizationEnhancer 商品归一化
func NewGoodsNormalizationEnhancer() *GoodsNormalizationEnhancer {

	return &GoodsNormalizationEnhancer{}
}

func (g *GoodsNormalizationEnhancer) Init() error {
	replaceProperties, err := cfg.Pack(replacepkg.Config{
		Old: "x",
		New: "*",
	})
	if err != nil {
		return err
	}
	regexProperties, err := cfg.Pack(regex.Config{
		Pattern:     goodsSearchSpecPattern,
		Replacement: "$1",
	})
	if err != nil {
		return err
	}
	pipeline := normalization.NewPipeline(normalizationapi.Config{
		Interceptors: []*normalizationapi.InterceptorConfig{
			{Type: whitespace.Type},
			{Type: shopcategory.Type},
			{Type: replacepkg.Type, Properties: replaceProperties},
			{Type: regex.Type, Properties: regexProperties},
		},
	})
	g.normalPipeline = pipeline
	return nil
}

// Normalization 归一化
func (g *GoodsNormalizationEnhancer) Normalization(input string) string {
	result, err := g.normalPipeline.Normalize(input)
	if err != nil {
		return ""
	}
	formatStr := result.Value
	fmt.Println(formatStr)
	return result.Value
}
