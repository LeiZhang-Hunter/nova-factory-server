package financedao

import (
	"nova-factory-server/app/business/erp/finance/financemodels"

	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/erp/erpbiz"
)

// IFinanceReceiptDao ERP 收款单数据访问接口
type IFinanceReceiptDao interface {
	Create(c *gin.Context, req *financemodels.FinanceReceiptUpsert) (*financemodels.FinanceReceipt, error)
	Update(c *gin.Context, req *financemodels.FinanceReceiptUpsert) (*financemodels.FinanceReceipt, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*financemodels.FinanceReceipt, error)
	GetByColumn(c *gin.Context, column string, value any) (*financemodels.FinanceReceipt, error)
	ListPage(c *gin.Context, req *financemodels.FinanceReceiptQuery) (*erpbiz.PageResult[financemodels.FinanceReceipt], error)
	List(c *gin.Context, req *financemodels.FinanceReceiptQuery) (*financemodels.FinanceReceiptListData, error)
}
