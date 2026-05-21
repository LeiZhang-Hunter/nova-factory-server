package retrieval

import (
	"context"

	"nova-factory-server/app/utils/retrieval/schema"
)

// Retriever 定义统一检索器接口。
//
// 实现方负责完成 query 预处理、向量生成、召回与结果映射。
type Retriever interface {
	Retrieve(ctx context.Context, query string, opts ...Option) ([]*schema.Document, error)
}
