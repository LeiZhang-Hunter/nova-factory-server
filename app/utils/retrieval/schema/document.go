package schema

// Document 表示统一的检索结果文档。
//
// 不同检索器可以把各自的业务结果映射成统一文档结构，
// 并把实现相关的附加信息放到 Metadata 中。
type Document struct {
	ID       string         `json:"id"`
	Content  string         `json:"content"`
	Score    float64        `json:"score"`
	Metadata map[string]any `json:"metadata,omitempty"`
}
