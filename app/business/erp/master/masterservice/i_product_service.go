package masterservice

import (
	"nova-factory-server/app/business/erp/master/mastermodels"

	"github.com/gin-gonic/gin"
)

// IProductService ERP 产品服务接口
type IProductService interface {
	Create(c *gin.Context, req *mastermodels.ProductUpsert) (*mastermodels.Product, error)
	Update(c *gin.Context, req *mastermodels.ProductUpsert) (*mastermodels.Product, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*mastermodels.Product, error)
	List(c *gin.Context, req *mastermodels.ProductQuery) (*mastermodels.ProductListData, error)
	GenerateVector(c *gin.Context, req *mastermodels.ProductGenVectorReq) (*mastermodels.ProductVectorResult, error)
	GenerateAllVectors(c *gin.Context, req *mastermodels.ProductGenAllVectorReq) (*mastermodels.ProductVectorTaskData, error)
	GetGenerateAllVectorsProgress(c *gin.Context, taskID string) (*mastermodels.ProductVectorTaskProgress, error)
	ListGenerateAllVectorTasks(c *gin.Context) (*mastermodels.ProductVectorTaskListData, error)
	SearchVector(c *gin.Context, req *mastermodels.ProductVectorSearchReq) (*mastermodels.ProductVectorSearchData, error)
	// BatchSearchVector 批量搜索向量
	BatchSearchVector(c *gin.Context, req *mastermodels.ProductVectorBatchSearchReq) (*mastermodels.ProductVectorBatchSearchData, error)
}
