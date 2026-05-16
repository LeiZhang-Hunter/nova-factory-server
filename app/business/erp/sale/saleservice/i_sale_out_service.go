package saleservice

import (
	"nova-factory-server/app/business/erp/sale/salemodels"

	"github.com/gin-gonic/gin"
)

// ISaleOutService ERP 销售出库服务接口
type ISaleOutService interface {
	Create(c *gin.Context, req *salemodels.SaleOutUpsert) (*salemodels.SaleOut, error)
	Update(c *gin.Context, req *salemodels.SaleOutUpsert) (*salemodels.SaleOut, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*salemodels.SaleOut, error)
	List(c *gin.Context, req *salemodels.SaleOutQuery) (*salemodels.SaleOutListData, error)
}
