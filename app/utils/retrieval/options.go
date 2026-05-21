package retrieval

import "github.com/cloudwego/eino/components/embedding"

// Options 描述一次检索请求的可选配置。
//
// 不同检索器可以按自己的语义解释这些字段，
// 未使用的字段允许忽略。
type Options struct {
	// Index 是检索器使用的索引，不同检索器中的索引可能有不同含义。
	Index *string

	// SubIndex 是检索器使用的子索引，不同检索器中的子索引可能有不同含义。
	SubIndex *string

	// TopK 是检索返回文档数量上限。
	TopK *int

	// ScoreThreshold 是检索结果的分数阈值。
	ScoreThreshold *float64

	// Embedding 用于生成查询向量。
	Embedding embedding.Embedder

	// DSLInfo 用于传递检索器专属 DSL 信息。
	DSLInfo map[string]interface{}
}

// Option 是 Options 的函数式配置项。
type Option func(*Options)

// ApplyOptions 将函数式参数合并为最终配置。
func ApplyOptions(opts ...Option) Options {
	var result Options
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		opt(&result)
	}
	return result
}

// WithIndex 设置索引。
func WithIndex(index string) Option {
	return func(opts *Options) {
		opts.Index = &index
	}
}

// WithSubIndex 设置子索引。
func WithSubIndex(subIndex string) Option {
	return func(opts *Options) {
		opts.SubIndex = &subIndex
	}
}

// WithTopK 设置检索文档数量上限。
func WithTopK(topK int) Option {
	return func(opts *Options) {
		opts.TopK = &topK
	}
}

// WithScoreThreshold 设置相似度阈值。
func WithScoreThreshold(threshold float64) Option {
	return func(opts *Options) {
		opts.ScoreThreshold = &threshold
	}
}

// WithEmbedding 设置向量生成组件。
func WithEmbedding(emb embedding.Embedder) Option {
	return func(opts *Options) {
		opts.Embedding = emb
	}
}

// WithDSLInfo 设置 DSL 信息。
func WithDSLInfo(dsl map[string]any) Option {
	return func(opts *Options) {
		opts.DSLInfo = dsl
	}
}
