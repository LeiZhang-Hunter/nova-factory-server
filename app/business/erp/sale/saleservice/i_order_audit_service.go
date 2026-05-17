package saleservice

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/erp/sale/salemodels"
)

// IOrderAuditService ERP订单审核服务接口。
type IOrderAuditService interface {
	Set(c *gin.Context, req *salemodels.OrderAuditSet) (*salemodels.OrderAudit, error)
	Import(c *gin.Context, req *salemodels.OrderAuditImportReq) (*salemodels.OrderAuditImportResult, error)
	GetByID(c *gin.Context, id uint64) (*salemodels.OrderAudit, error)
	List(c *gin.Context, req *salemodels.OrderAuditQuery) (*salemodels.OrderAuditListData, error)
	DeleteByIDs(c *gin.Context, ids []uint64) error
	Approve(c *gin.Context, req *salemodels.OrderAuditAction) (*salemodels.OrderAudit, error)
	Reject(c *gin.Context, req *salemodels.OrderAuditAction) (*salemodels.OrderAudit, error)
}
