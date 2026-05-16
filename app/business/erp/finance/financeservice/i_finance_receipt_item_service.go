package financeservice

import (
	"nova-factory-server/app/business/erp/finance/financemodels"

	"github.com/gin-gonic/gin"
)

// IFinanceReceiptItemService ERP 收款项服务接口
type IFinanceReceiptItemService interface {
	Create(c *gin.Context, req *financemodels.FinanceReceiptItemUpsert) (*financemodels.FinanceReceiptItem, error)
	Update(c *gin.Context, req *financemodels.FinanceReceiptItemUpsert) (*financemodels.FinanceReceiptItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*financemodels.FinanceReceiptItem, error)
	List(c *gin.Context, req *financemodels.FinanceReceiptItemQuery) (*financemodels.FinanceReceiptItemListData, error)
}
