package goods

import (
	"nova-factory-server/app/utils/gateway/v1/config/cfg"
	"nova-factory-server/app/utils/vectorsearch/normalization"
	normalizationapi "nova-factory-server/app/utils/vectorsearch/normalization/api"
	replacepkg "nova-factory-server/app/utils/vectorsearch/normalization/replace"
	"nova-factory-server/app/utils/vectorsearch/normalization/shopcategory"
	"nova-factory-server/app/utils/vectorsearch/normalization/whitespace"
)

// NormalizationEnhancer 商品归一化处理器
type NormalizationEnhancer struct {
	normalPipeline *normalization.Pipeline
}

// NewGoodsNormalizationEnhancer 商品归一化
func NewGoodsNormalizationEnhancer() *NormalizationEnhancer {
	return &NormalizationEnhancer{}
}

func (g *NormalizationEnhancer) Init() error {
	replaceProperties, err := cfg.Pack(replacepkg.Config{
		Old: "x",
		New: "*",
	})
	if err != nil {
		return err
	}

	pipeline := normalization.NewPipeline(normalizationapi.Config{
		Interceptors: []*normalizationapi.InterceptorConfig{
			{Type: whitespace.Type},
			{Type: shopcategory.Type},
			{Type: replacepkg.Type, Properties: replaceProperties},
		},
	})
	g.normalPipeline = pipeline
	return nil
}

// Normalization 归一化
func (g *NormalizationEnhancer) Normalization(input string) string {
	result, err := g.normalPipeline.Normalize(input)
	if err != nil {
		return ""
	}
	return result.Value
}
