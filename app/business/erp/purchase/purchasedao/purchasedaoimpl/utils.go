package purchasedaoimpl

import (
	"strings"

	"nova-factory-server/app/business/erp/purchase/purchasemodels"

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

func applyPurchaseOrderFilters(db *gorm.DB, req *purchasemodels.PurchaseOrderQuery) *gorm.DB {
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
	return db
}

func applyPurchaseOrderItemFilters(db *gorm.DB, req *purchasemodels.PurchaseOrderItemQuery) *gorm.DB {
	if req.OrderID > 0 {
		db = db.Where("order_id = ?", req.OrderID)
	}
	if req.ProductID > 0 {
		db = db.Where("product_id = ?", req.ProductID)
	}
	return db
}

func applyPurchaseInFilters(db *gorm.DB, req *purchasemodels.PurchaseInQuery) *gorm.DB {
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
	if req.OrderID > 0 {
		db = db.Where("order_id = ?", req.OrderID)
	}
	return db
}

func applyPurchaseInItemFilters(db *gorm.DB, req *purchasemodels.PurchaseInItemQuery) *gorm.DB {
	if req.WarehouseID > 0 {
		db = db.Where("warehouse_id = ?", req.WarehouseID)
	}
	if req.ProductID > 0 {
		db = db.Where("product_id = ?", req.ProductID)
	}
	if req.InID > 0 {
		db = db.Where("in_id = ?", req.InID)
	}
	return db
}

func applyPurchaseReturnFilters(db *gorm.DB, req *purchasemodels.PurchaseReturnQuery) *gorm.DB {
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
	if req.OrderID > 0 {
		db = db.Where("order_id = ?", req.OrderID)
	}
	return db
}

func applyPurchaseReturnItemFilters(db *gorm.DB, req *purchasemodels.PurchaseReturnItemQuery) *gorm.DB {
	if req.WarehouseID > 0 {
		db = db.Where("warehouse_id = ?", req.WarehouseID)
	}
	if req.ProductID > 0 {
		db = db.Where("product_id = ?", req.ProductID)
	}
	if req.ReturnID > 0 {
		db = db.Where("return_id = ?", req.ReturnID)
	}
	return db
}
