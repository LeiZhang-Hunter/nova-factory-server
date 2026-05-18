package saledao

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/erp/sale/salemodels"
)

// IOrderAuditDao ERP订单审核数据访问接口。
type IOrderAuditDao interface {
	Set(c *gin.Context, req *salemodels.OrderAuditSet) (*salemodels.OrderAudit, error)
	GetByID(c *gin.Context, id uint64) (*salemodels.OrderAudit, error)
	GetByIDWithTx(c *gin.Context, tx *gorm.DB, id uint64) (*salemodels.OrderAudit, error)
	List(c *gin.Context, req *salemodels.OrderAuditQuery) (*salemodels.OrderAuditListData, error)
	DeleteByIDs(c *gin.Context, ids []uint64) error
	Approve(c *gin.Context, id uint64, remark string, erpOrderID uint64) error
	ApproveWithTx(c *gin.Context, tx *gorm.DB, id uint64, remark string, erpOrderID uint64) error
	Reject(c *gin.Context, id uint64, remark string) error
	RejectWithTx(c *gin.Context, tx *gorm.DB, id uint64, remark string) error
}
