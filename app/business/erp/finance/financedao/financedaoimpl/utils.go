package financedaoimpl

import (
	"strings"

	"nova-factory-server/app/business/erp/finance/financemodels"

	"gorm.io/gorm"
)

func getPageSize(page, size int64) (int64, int64) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if size > 200 {
		size = 200
	}
	return page, size
}

func applyFinanceReceiptFilters(db *gorm.DB, req *financemodels.FinanceReceiptQuery) *gorm.DB {
	if req.No != "" {
		db = db.Where("no LIKE ?", "%"+strings.TrimSpace(req.No)+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.CustomerID > 0 {
		db = db.Where("customer_id = ?", req.CustomerID)
	}
	if req.AccountID > 0 {
		db = db.Where("account_id = ?", req.AccountID)
	}
	if req.FinanceUserID > 0 {
		db = db.Where("finance_user_id = ?", req.FinanceUserID)
	}
	return db
}

func applyFinanceReceiptItemFilters(db *gorm.DB, req *financemodels.FinanceReceiptItemQuery) *gorm.DB {
	if req.ReceiptID > 0 {
		db = db.Where("receipt_id = ?", req.ReceiptID)
	}
	if req.BizType != nil {
		db = db.Where("biz_type = ?", *req.BizType)
	}
	return db
}

func applyFinancePaymentFilters(db *gorm.DB, req *financemodels.FinancePaymentQuery) *gorm.DB {
	if req.No != "" {
		db = db.Where("no LIKE ?", "%"+strings.TrimSpace(req.No)+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.SupplierID > 0 {
		db = db.Where("supplier_id = ?", req.SupplierID)
	}
	if req.AccountID > 0 {
		db = db.Where("account_id = ?", req.AccountID)
	}
	if req.FinanceUserID > 0 {
		db = db.Where("finance_user_id = ?", req.FinanceUserID)
	}
	return db
}

func applyFinancePaymentItemFilters(db *gorm.DB, req *financemodels.FinancePaymentItemQuery) *gorm.DB {
	if req.PaymentID > 0 {
		db = db.Where("payment_id = ?", req.PaymentID)
	}
	if req.BizType != nil {
		db = db.Where("biz_type = ?", *req.BizType)
	}
	return db
}
