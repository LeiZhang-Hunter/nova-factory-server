package dao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/erp/sale/salemodels"
)

type IApiShopOrderDao interface {
	// Set 新增或修改 ERP 订单及其子表。
	Set(c *gin.Context, req *salemodels.OrderSet) (*salemodels.Order, error)
	// SetWithTx 新增或修改 ERP 订单及其子表（带事务）。
	SetWithTx(c *gin.Context, tx *gorm.DB, req *salemodels.OrderSet) (*salemodels.Order, error)
	// GetByID 查询 ERP 订单详情。
	GetByID(c *gin.Context, id uint64) (*salemodels.Order, error)
	// List 分页查询 ERP 订单。
	List(c *gin.Context, req *salemodels.OrderQuery) (*salemodels.OrderListData, error)
	// DeleteByIDs 删除 ERP 订单。
	DeleteByIDs(c *gin.Context, ids []uint64) error
}
