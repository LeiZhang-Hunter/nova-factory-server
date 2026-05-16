package financedao

import (
	"nova-factory-server/app/business/erp/finance/financemodels"

	"github.com/gin-gonic/gin"
)

// IFinanceReceiptItemDao ERP 收款项数据访问接口
type IFinanceReceiptItemDao interface {
	Create(c *gin.Context, req *financemodels.FinanceReceiptItemUpsert) (*financemodels.FinanceReceiptItem, error)
	Update(c *gin.Context, req *financemodels.FinanceReceiptItemUpsert) (*financemodels.FinanceReceiptItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*financemodels.FinanceReceiptItem, error)
	GetByColumn(c *gin.Context, column string, value any) (*financemodels.FinanceReceiptItem, error)
	ListPage(c *gin.Context, req *financemodels.FinanceReceiptItemQuery) (*financemodels.FinanceReceiptItemListData, error)
	List(c *gin.Context, req *financemodels.FinanceReceiptItemQuery) (*financemodels.FinanceReceiptItemListData, error)
}
