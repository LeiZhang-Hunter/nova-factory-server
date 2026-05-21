package masterdaoimpl

import (
	"fmt"
	"strings"

	"nova-factory-server/app/business/erp/master/mastermodels"
	searchutil "nova-factory-server/app/utils/vectorsearch"

	"github.com/milvus-io/milvus/client/v2/entity"
)

type productVectorSearchPlan struct {
	Query         string
	SearchText    string
	ExactCodeText string
	Vector        entity.Vector
	IsCodeLike    bool
	Intent        string
	Processed     *searchutil.ProcessedQuery
}

// productVectorSearchStrategy 定义查询意图策略接口。
type productVectorSearchStrategy interface {
	Name() string
	Match(processed *searchutil.ProcessedQuery) bool
	Build(plan productVectorSearchPlan, processed *searchutil.ProcessedQuery) productVectorSearchPlan
}

// productVectorCodeSearchStrategy 处理编码/条码/SKU 类高精度查询。
type productVectorCodeSearchStrategy struct{}

func (s productVectorCodeSearchStrategy) Name() string { return "code" }

func (s productVectorCodeSearchStrategy) Match(processed *searchutil.ProcessedQuery) bool {
	return processed != nil && processed.IsCodeLike
}

func (s productVectorCodeSearchStrategy) Build(plan productVectorSearchPlan, processed *searchutil.ProcessedQuery) productVectorSearchPlan {
	plan.Intent = s.Name()
	plan.IsCodeLike = true
	if processed != nil && len(processed.CodeTerms) > 0 {
		plan.ExactCodeText = processed.CodeTerms[0]
		plan.SearchText = strings.TrimSpace(strings.Join([]string{
			processed.CodeTerms[0],
			processed.Normalized,
			plan.SearchText,
		}, " "))
	}
	return plan
}

// productVectorSpecSearchStrategy 处理带规格信息的商品查询。
type productVectorSpecSearchStrategy struct{}

func (s productVectorSpecSearchStrategy) Name() string { return "spec" }

func (s productVectorSpecSearchStrategy) Match(processed *searchutil.ProcessedQuery) bool {
	return processed != nil && len(processed.SpecTerms) > 0
}

func (s productVectorSpecSearchStrategy) Build(plan productVectorSearchPlan, processed *searchutil.ProcessedQuery) productVectorSearchPlan {
	plan.Intent = s.Name()
	if processed != nil {
		parts := make([]string, 0, 2+len(processed.SpecTerms))
		parts = append(parts, processed.Normalized)
		parts = append(parts, processed.SpecTerms...)
		parts = append(parts, plan.SearchText)
		plan.SearchText = strings.TrimSpace(strings.Join(parts, " "))
	}
	return plan
}

// productVectorCategorySearchStrategy 处理分类类查询，增强分类词在检索文本中的权重。
type productVectorCategorySearchStrategy struct{}

func (s productVectorCategorySearchStrategy) Name() string { return "category" }

func (s productVectorCategorySearchStrategy) Match(processed *searchutil.ProcessedQuery) bool {
	return processed != nil && len(processed.CategoryTerms) > 0
}

func (s productVectorCategorySearchStrategy) Build(plan productVectorSearchPlan, processed *searchutil.ProcessedQuery) productVectorSearchPlan {
	plan.Intent = s.Name()
	if processed != nil {
		parts := make([]string, 0, 2+len(processed.CategoryTerms))
		parts = append(parts, processed.Normalized)
		parts = append(parts, processed.CategoryTerms...)
		parts = append(parts, plan.SearchText)
		plan.SearchText = strings.TrimSpace(strings.Join(parts, " "))
	}
	return plan
}

// productVectorSemanticSearchStrategy 兜底处理普通语义检索。
type productVectorSemanticSearchStrategy struct{}

func (s productVectorSemanticSearchStrategy) Name() string { return "semantic" }

func (s productVectorSemanticSearchStrategy) Match(processed *searchutil.ProcessedQuery) bool {
	return true
}

func (s productVectorSemanticSearchStrategy) Build(plan productVectorSearchPlan, processed *searchutil.ProcessedQuery) productVectorSearchPlan {
	plan.Intent = s.Name()
	return plan
}

// productVectorSearchStrategies 返回按优先级排列的查询策略集合。
func productVectorSearchStrategies() []productVectorSearchStrategy {
	return []productVectorSearchStrategy{
		productVectorCodeSearchStrategy{},
		productVectorSpecSearchStrategy{},
		productVectorCategorySearchStrategy{},
		productVectorSemanticSearchStrategy{},
	}
}

// buildProductVectorSearchPlans 将批量查询参数转换为统一检索计划，隔离参数校验与 query 意图分析。
func buildProductVectorSearchPlans(req *mastermodels.ProductVectorBatchSearchReq, vectors [][]float32) ([]productVectorSearchPlan, error) {
	plans := make([]productVectorSearchPlan, 0, len(vectors))
	strategies := productVectorSearchStrategies()
	for idx, vector := range vectors {
		if len(vector) == 0 {
			return nil, fmt.Errorf("第%d条产品搜索向量为空", idx+1)
		}
		query := strings.TrimSpace(req.Queries[idx])
		if query == "" {
			return nil, fmt.Errorf("第%d条产品搜索文本为空", idx+1)
		}

		searchText := query
		if len(req.SearchTexts) == len(req.Queries) {
			if candidate := strings.TrimSpace(req.SearchTexts[idx]); candidate != "" {
				searchText = candidate
			}
		}

		processed := searchutil.ProcessQuery(query)
		plan := productVectorSearchPlan{
			Query:      query,
			SearchText: searchText,
			Vector:     entity.FloatVector(vector),
			Processed:  processed,
		}
		for _, strategy := range strategies {
			if !strategy.Match(processed) {
				continue
			}
			plan = strategy.Build(plan, processed)
			break
		}
		plans = append(plans, plan)
	}
	return plans, nil
}
