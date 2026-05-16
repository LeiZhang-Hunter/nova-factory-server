package financeservice

import (
	"nova-factory-server/app/business/erp/finance/financemodels"

	"github.com/gin-gonic/gin"
)

// IFinancePaymentService ERP 付款单服务接口
type IFinancePaymentService interface {
	Create(c *gin.Context, req *financemodels.FinancePaymentUpsert) (*financemodels.FinancePayment, error)
	Update(c *gin.Context, req *financemodels.FinancePaymentUpsert) (*financemodels.FinancePayment, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*financemodels.FinancePayment, error)
	List(c *gin.Context, req *financemodels.FinancePaymentQuery) (*financemodels.FinancePaymentListData, error)
}
