package orderservice

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/erp/order/ordermodels"
)

// IOrderAuditService ERP订单审核服务接口。
type IOrderAuditService interface {
	Set(c *gin.Context, req *ordermodels.OrderAuditSet) (*ordermodels.OrderAudit, error)
	Import(c *gin.Context, req *ordermodels.OrderAuditImportReq) (*ordermodels.OrderAuditImportResult, error)
	GetByID(c *gin.Context, id uint64) (*ordermodels.OrderAudit, error)
	List(c *gin.Context, req *ordermodels.OrderAuditQuery) (*ordermodels.OrderAuditListData, error)
	DeleteByIDs(c *gin.Context, ids []uint64) error
	Approve(c *gin.Context, req *ordermodels.OrderAuditAction) (*ordermodels.OrderAudit, error)
	Reject(c *gin.Context, req *ordermodels.OrderAuditAction) (*ordermodels.OrderAudit, error)
}
