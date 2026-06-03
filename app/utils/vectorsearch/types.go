package vectorsearch

// LabeledValue 表示带标签的文本片段，常用于将结构化字段拼装为 embedding 文本。
type LabeledValue struct {
	Label string
	Value string
}

// LabeledContentProvider 提供可用于向量文本构建的结构化字段。
type LabeledContentProvider interface {
	VectorSearchLabeledValues() []LabeledValue
}

// ProcessedQuery 表示对用户原始检索词做完清洗、分词、扩展后的结构化结果。
type ProcessedQuery struct {
	Original       string
	Normalized     string
	Tokens         []string
	ExpandedTokens []string
	Keywords       []string
	CategoryTerms  []string
	SpecTerms      []string
	CodeTerms      []string
	EmbeddingText  string
	HybridText     string
	IsCodeLike     bool
	IsShortQuery   bool
}
