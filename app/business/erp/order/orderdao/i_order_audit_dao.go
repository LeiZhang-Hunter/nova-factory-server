package orderdao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/erp/order/ordermodels"
)

// IOrderAuditDao ERP订单审核数据访问接口。
type IOrderAuditDao interface {
	Set(c *gin.Context, req *ordermodels.OrderAuditSet) (*ordermodels.OrderAudit, error)
	GetByID(c *gin.Context, id uint64) (*ordermodels.OrderAudit, error)
	List(c *gin.Context, req *ordermodels.OrderAuditQuery) (*ordermodels.OrderAuditListData, error)
	DeleteByIDs(c *gin.Context, ids []uint64) error
	Approve(c *gin.Context, id uint64, remark string, erpOrderID uint64) error
	Reject(c *gin.Context, id uint64, remark string) error
}
