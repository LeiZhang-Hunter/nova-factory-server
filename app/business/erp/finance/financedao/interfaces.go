package financedao

import (
	"nova-factory-server/app/business/erp/finance/financemodels"

	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/erp/erpcrud"
)

// IFinancePaymentDao ERP 付款单数据访问接口
type IFinancePaymentDao interface {
	Create(c *gin.Context, req *financemodels.FinancePaymentUpsert) (*financemodels.FinancePayment, error)
	Update(c *gin.Context, req *financemodels.FinancePaymentUpsert) (*financemodels.FinancePayment, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*financemodels.FinancePayment, error)
	GetByColumn(c *gin.Context, column string, value any) (*financemodels.FinancePayment, error)
	ListPage(c *gin.Context, req *financemodels.FinancePaymentQuery) (*erpcrud.PageResult[financemodels.FinancePayment], error)
	List(c *gin.Context, req *financemodels.FinancePaymentQuery) (*financemodels.FinancePaymentListData, error)
}

// IFinancePaymentItemDao ERP 付款项数据访问接口
type IFinancePaymentItemDao interface {
	Create(c *gin.Context, req *financemodels.FinancePaymentItemUpsert) (*financemodels.FinancePaymentItem, error)
	Update(c *gin.Context, req *financemodels.FinancePaymentItemUpsert) (*financemodels.FinancePaymentItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*financemodels.FinancePaymentItem, error)
	GetByColumn(c *gin.Context, column string, value any) (*financemodels.FinancePaymentItem, error)
	ListPage(c *gin.Context, req *financemodels.FinancePaymentItemQuery) (*erpcrud.PageResult[financemodels.FinancePaymentItem], error)
	List(c *gin.Context, req *financemodels.FinancePaymentItemQuery) (*financemodels.FinancePaymentItemListData, error)
}

// IFinanceReceiptDao ERP 收款单数据访问接口
type IFinanceReceiptDao interface {
	Create(c *gin.Context, req *financemodels.FinanceReceiptUpsert) (*financemodels.FinanceReceipt, error)
	Update(c *gin.Context, req *financemodels.FinanceReceiptUpsert) (*financemodels.FinanceReceipt, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*financemodels.FinanceReceipt, error)
	GetByColumn(c *gin.Context, column string, value any) (*financemodels.FinanceReceipt, error)
	ListPage(c *gin.Context, req *financemodels.FinanceReceiptQuery) (*erpcrud.PageResult[financemodels.FinanceReceipt], error)
	List(c *gin.Context, req *financemodels.FinanceReceiptQuery) (*financemodels.FinanceReceiptListData, error)
}

// IFinanceReceiptItemDao ERP 收款项数据访问接口
type IFinanceReceiptItemDao interface {
	Create(c *gin.Context, req *financemodels.FinanceReceiptItemUpsert) (*financemodels.FinanceReceiptItem, error)
	Update(c *gin.Context, req *financemodels.FinanceReceiptItemUpsert) (*financemodels.FinanceReceiptItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*financemodels.FinanceReceiptItem, error)
	GetByColumn(c *gin.Context, column string, value any) (*financemodels.FinanceReceiptItem, error)
	ListPage(c *gin.Context, req *financemodels.FinanceReceiptItemQuery) (*erpcrud.PageResult[financemodels.FinanceReceiptItem], error)
	List(c *gin.Context, req *financemodels.FinanceReceiptItemQuery) (*financemodels.FinanceReceiptItemListData, error)
}
