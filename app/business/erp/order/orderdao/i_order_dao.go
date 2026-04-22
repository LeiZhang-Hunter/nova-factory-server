package orderdao

import (
	"nova-factory-server/app/business/erp/order/ordermodels"

	"github.com/gin-gonic/gin"
)

// IOrderDao ERP订单数据访问接口。
type IOrderDao interface {
	// Set 新增或修改 ERP 订单及其子表。
	Set(c *gin.Context, req *ordermodels.OrderSet) (*ordermodels.Order, error)
	// GetByID 查询 ERP 订单详情。
	GetByID(c *gin.Context, id uint64) (*ordermodels.Order, error)
	// List 分页查询 ERP 订单。
	List(c *gin.Context, req *ordermodels.OrderQuery) (*ordermodels.OrderListData, error)
	// DeleteByIDs 删除 ERP 订单。
	DeleteByIDs(c *gin.Context, ids []uint64) error
}

// IOrderAuditDao ERP订单审核数据访问接口。
type IOrderAuditDao interface {
	Set(c *gin.Context, req *ordermodels.OrderAuditSet) (*ordermodels.OrderAudit, error)
	GetByID(c *gin.Context, id uint64) (*ordermodels.OrderAudit, error)
	List(c *gin.Context, req *ordermodels.OrderAuditQuery) (*ordermodels.OrderAuditListData, error)
	DeleteByIDs(c *gin.Context, ids []uint64) error
	Approve(c *gin.Context, id uint64, remark string, erpOrderID uint64) error
	Reject(c *gin.Context, id uint64, remark string) error
}
