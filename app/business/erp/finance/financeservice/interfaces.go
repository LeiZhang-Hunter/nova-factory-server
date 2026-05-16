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

// IFinancePaymentItemService ERP 付款项服务接口
type IFinancePaymentItemService interface {
	Create(c *gin.Context, req *financemodels.FinancePaymentItemUpsert) (*financemodels.FinancePaymentItem, error)
	Update(c *gin.Context, req *financemodels.FinancePaymentItemUpsert) (*financemodels.FinancePaymentItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*financemodels.FinancePaymentItem, error)
	List(c *gin.Context, req *financemodels.FinancePaymentItemQuery) (*financemodels.FinancePaymentItemListData, error)
}

// IFinanceReceiptService ERP 收款单服务接口
type IFinanceReceiptService interface {
	Create(c *gin.Context, req *financemodels.FinanceReceiptUpsert) (*financemodels.FinanceReceipt, error)
	Update(c *gin.Context, req *financemodels.FinanceReceiptUpsert) (*financemodels.FinanceReceipt, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*financemodels.FinanceReceipt, error)
	List(c *gin.Context, req *financemodels.FinanceReceiptQuery) (*financemodels.FinanceReceiptListData, error)
}

// IFinanceReceiptItemService ERP 收款项服务接口
type IFinanceReceiptItemService interface {
	Create(c *gin.Context, req *financemodels.FinanceReceiptItemUpsert) (*financemodels.FinanceReceiptItem, error)
	Update(c *gin.Context, req *financemodels.FinanceReceiptItemUpsert) (*financemodels.FinanceReceiptItem, error)
	DeleteByIDs(c *gin.Context, ids []int64) error
	GetByID(c *gin.Context, id int64) (*financemodels.FinanceReceiptItem, error)
	List(c *gin.Context, req *financemodels.FinanceReceiptItemQuery) (*financemodels.FinanceReceiptItemListData, error)
}
