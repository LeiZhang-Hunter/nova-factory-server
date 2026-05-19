package masterdao

import (
	"nova-factory-server/app/business/erp/master/mastermodels"

	"github.com/gin-gonic/gin"
)

// IProductVectorDao ERP 产品向量数据访问接口。
type IProductVectorDao interface {
	Upsert(c *gin.Context, product *mastermodels.Product, item *mastermodels.ProductVectorUpsertItem) (*mastermodels.ProductVectorResult, error)
	Search(c *gin.Context, req *mastermodels.ProductVectorSearchReq, vector []float32) (*mastermodels.ProductVectorSearchData, error)
	BatchSearch(c *gin.Context, req *mastermodels.ProductVectorBatchSearchReq, vectors [][]float32) (*mastermodels.ProductVectorBatchSearchData, error)
}
