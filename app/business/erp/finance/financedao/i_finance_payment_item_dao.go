package financedao

import (
	"nova-factory-server/app/business/erp/finance/financemodels"

	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/erp/erpbiz"
)

// IFinancePaymentItemDao ERP 付款项数据访问接口
type IFinancePaymentItemDao interface {
	Create(c *gin.Context, req *financemodels.FinancePaymentItemUpsert) (*financemodels.FinancePaymentItem, error)
	Update(c *gin.Context, req *financemodels.FinancePaymentItemUpsert) (*financemodels.FinancePaymentItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*financemodels.FinancePaymentItem, error)
	GetByColumn(c *gin.Context, column string, value any) (*financemodels.FinancePaymentItem, error)
	ListPage(c *gin.Context, req *financemodels.FinancePaymentItemQuery) (*erpbiz.PageResult[financemodels.FinancePaymentItem], error)
	List(c *gin.Context, req *financemodels.FinancePaymentItemQuery) (*financemodels.FinancePaymentItemListData, error)
}
