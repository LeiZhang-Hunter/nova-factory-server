package financedao

import (
	"nova-factory-server/app/business/erp/finance/financemodels"

	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/erp/erpbiz"
)

// IFinancePaymentDao ERP 付款单数据访问接口
type IFinancePaymentDao interface {
	Create(c *gin.Context, req *financemodels.FinancePaymentUpsert) (*financemodels.FinancePayment, error)
	Update(c *gin.Context, req *financemodels.FinancePaymentUpsert) (*financemodels.FinancePayment, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*financemodels.FinancePayment, error)
	GetByColumn(c *gin.Context, column string, value any) (*financemodels.FinancePayment, error)
	ListPage(c *gin.Context, req *financemodels.FinancePaymentQuery) (*erpbiz.PageResult[financemodels.FinancePayment], error)
	List(c *gin.Context, req *financemodels.FinancePaymentQuery) (*financemodels.FinancePaymentListData, error)
}
