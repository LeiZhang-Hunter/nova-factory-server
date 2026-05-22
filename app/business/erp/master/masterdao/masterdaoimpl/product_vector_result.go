package masterdaoimpl

import (
	"context"
	"fmt"
	"sort"

	"nova-factory-server/app/business/erp/master/mastermodels"
	searchutil "nova-factory-server/app/utils/vectorsearch"

	"github.com/gin-gonic/gin"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

// buildEmptyProductVectorBatchSearchData 构造空结果，保持批量查询返回结构稳定。
func buildEmptyProductVectorBatchSearchData(queries []string) *mastermodels.ProductVectorBatchSearchData {
	rows := make([]*mastermodels.ProductVectorBatchSearchItem, 0, len(queries))
	for _, query := range queries {
		rows = append(rows, &mastermodels.ProductVectorBatchSearchItem{
			Query: query,
			Rows:  make([]*mastermodels.ProductVectorSearchItem, 0),
			Total: 0,
		})
	}
	return &mastermodels.ProductVectorBatchSearchData{
		Rows:  rows,
		Total: int64(len(rows)),
	}
}

// parseProductVectorSearchResultSet 将 Milvus 返回结果转换为业务检索结果。
func parseProductVectorSearchResultSet(resultSet milvusclient.ResultSet) (*mastermodels.ProductVectorSearchData, error) {
	if resultSet.Err != nil {
		return nil, resultSet.Err
	}

	rows := make([]*mastermodels.ProductVectorSearchItem, 0, resultSet.ResultCount)
	for i := 0; i < resultSet.ResultCount; i++ {
		// 先保留 Milvus 的原始得分，后续应用层重排会基于该分数继续融合字段特征。
		item := &mastermodels.ProductVectorSearchItem{
			Score: resultSet.Scores[i],
		}
		var err error
		if item.ProductID, err = getProductVectorInt64ColumnValue(&resultSet, productVectorPKField, i); err != nil {
			return nil, err
		}
		if item.Name, err = getProductVectorStringColumnValue(&resultSet, productVectorNameField, i); err != nil {
			return nil, err
		}
		if item.BarCode, err = getProductVectorStringColumnValue(&resultSet, productVectorBarCodeField, i); err != nil {
			return nil, err
		}
		if item.CategoryId, err = getProductVectorInt64ColumnValue(&resultSet, productVectorCategoryIDField, i); err != nil {
			return nil, err
		}
		if item.CategoryName, err = getProductVectorStringColumnValue(&resultSet, productVectorCategoryNameField, i); err != nil {
			return nil, err
		}
		if item.UnitId, err = getProductVectorInt64ColumnValue(&resultSet, productVectorUnitIDField, i); err != nil {
			return nil, err
		}
		if item.UnitName, err = getProductVectorStringColumnValue(&resultSet, productVectorUnitNameField, i); err != nil {
			return nil, err
		}
		if item.Standard, err = getProductVectorStringColumnValue(&resultSet, productVectorStandardField, i); err != nil {
			return nil, err
		}
		if item.Remark, err = getProductVectorStringColumnValue(&resultSet, productVectorRemarkField, i); err != nil {
			return nil, err
		}
		expiryDay, err := getProductVectorInt64ColumnValue(&resultSet, productVectorExpiryDayField, i)
		if err != nil {
			return nil, err
		}
		item.ExpiryDay = int32(expiryDay)
		if item.Weight, err = getProductVectorFloat64ColumnValue(&resultSet, productVectorWeightField, i); err != nil {
			return nil, err
		}
		if item.PurchasePrice, err = getProductVectorFloat64ColumnValue(&resultSet, productVectorPurchasePriceField, i); err != nil {
			return nil, err
		}
		if item.SalePrice, err = getProductVectorFloat64ColumnValue(&resultSet, productVectorSalePriceField, i); err != nil {
			return nil, err
		}
		if item.MinPrice, err = getProductVectorFloat64ColumnValue(&resultSet, productVectorMinPriceField, i); err != nil {
			return nil, err
		}
		if item.Content, err = getProductVectorStringColumnValue(&resultSet, productVectorContentField, i); err != nil {
			return nil, err
		}
		rows = append(rows, item)
	}

	return &mastermodels.ProductVectorSearchData{
		Rows:  rows,
		Total: int64(len(rows)),
	}, nil
}

// buildProductVectorRequestContext 从 gin 上下文中提取请求上下文。
func buildProductVectorRequestContext(c *gin.Context) context.Context {
	if c != nil && c.Request != nil {
		return c.Request.Context()
	}
	return context.Background()
}

// normalizeProductVectorSearchLimit 统一处理搜索数量上限，避免查询参数失控。
func normalizeProductVectorSearchLimit(limit int) int {
	if limit <= 0 {
		return productVectorSearchDefaultLimit
	}
	if limit > productVectorSearchMaxLimit {
		return productVectorSearchMaxLimit
	}
	return limit
}

// resolveProductVectorSearchCandidateLimit 计算候选召回数量，为应用层重排预留足够样本。
func resolveProductVectorSearchCandidateLimit(limit int) int {
	limit = normalizeProductVectorSearchLimit(limit)
	candidateLimit := limit * productVectorSearchCandidateMultiplier
	if candidateLimit < productVectorSearchMinCandidates {
		candidateLimit = productVectorSearchMinCandidates
	}
	if candidateLimit > productVectorSearchMaxCandidates {
		candidateLimit = productVectorSearchMaxCandidates
	}
	if candidateLimit < limit {
		return limit
	}
	return candidateLimit
}

// rerankProductVectorSearchRows 使用字段加权与动态阈值对 Milvus 候选结果做二次排序。
func rerankProductVectorSearchRows(query string, rows []*mastermodels.ProductVectorSearchItem, limit int) []*mastermodels.ProductVectorSearchItem {
	if len(rows) == 0 {
		return make([]*mastermodels.ProductVectorSearchItem, 0)
	}

	// 将业务结果映射为通用候选结构，交由 app/utils/vectorsearch 做统一重排。
	processedQuery := searchutil.ProcessQuery(query)
	candidates := make([]searchutil.RankCandidate, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		candidates = append(candidates, searchutil.RankCandidate{
			ID:        row.ProductID,
			Title:     row.Name,
			Code:      row.BarCode,
			Category:  row.CategoryName,
			Unit:      row.UnitName,
			Standard:  row.Standard,
			Remark:    row.Remark,
			Content:   row.Content,
			BaseScore: row.Score,
		})
	}
	if len(candidates) == 0 {
		return make([]*mastermodels.ProductVectorSearchItem, 0)
	}

	// 重排后的分数会覆盖原始分数，便于前端或调用方直接按最终得分展示。
	ranked := searchutil.RerankCandidates(processedQuery, candidates, normalizeProductVectorSearchLimit(limit))
	result := make([]*mastermodels.ProductVectorSearchItem, 0, len(ranked))
	for _, item := range ranked {
		if item.Index < 0 || item.Index >= len(rows) || rows[item.Index] == nil {
			continue
		}
		row := rows[item.Index]
		row.Score = item.Score
		result = append(result, row)
	}

	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Score == result[j].Score {
			return result[i].ProductID < result[j].ProductID
		}
		return result[i].Score > result[j].Score
	})
	return result
}

// getProductVectorInt64ColumnValue 读取 int64 类型列值。
func getProductVectorInt64ColumnValue(resultSet *milvusclient.ResultSet, fieldName string, idx int) (int64, error) {
	if resultSet == nil {
		return 0, fmt.Errorf("搜索结果为空")
	}
	col := resultSet.GetColumn(fieldName)
	if col == nil {
		return 0, fmt.Errorf("搜索结果缺少字段 %s", fieldName)
	}
	return col.GetAsInt64(idx)
}

// getProductVectorStringColumnValue 读取字符串类型列值。
func getProductVectorStringColumnValue(resultSet *milvusclient.ResultSet, fieldName string, idx int) (string, error) {
	if resultSet == nil {
		return "", fmt.Errorf("搜索结果为空")
	}
	col := resultSet.GetColumn(fieldName)
	if col == nil {
		return "", fmt.Errorf("搜索结果缺少字段 %s", fieldName)
	}
	return col.GetAsString(idx)
}

// getProductVectorFloat64ColumnValue 读取浮点类型列值。
func getProductVectorFloat64ColumnValue(resultSet *milvusclient.ResultSet, fieldName string, idx int) (float64, error) {
	if resultSet == nil {
		return 0, fmt.Errorf("搜索结果为空")
	}
	col := resultSet.GetColumn(fieldName)
	if col == nil {
		return 0, fmt.Errorf("搜索结果缺少字段 %s", fieldName)
	}
	return col.GetAsDouble(idx)
}
