package saledaoimpl

import (
	"strings"

	"nova-factory-server/app/business/erp/sale/salemodels"

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

func applySaleOutFilters(db *gorm.DB, req *salemodels.SaleOutQuery) *gorm.DB {
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
	if req.OrderID > 0 {
		db = db.Where("order_id = ?", req.OrderID)
	}
	return db
}

func applySaleOutItemFilters(db *gorm.DB, req *salemodels.SaleOutItemQuery) *gorm.DB {
	if req.WarehouseID > 0 {
		db = db.Where("warehouse_id = ?", req.WarehouseID)
	}
	if req.ProductID > 0 {
		db = db.Where("product_id = ?", req.ProductID)
	}
	if req.OutID > 0 {
		db = db.Where("out_id = ?", req.OutID)
	}
	return db
}

func applySaleReturnFilters(db *gorm.DB, req *salemodels.SaleReturnQuery) *gorm.DB {
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
	if req.OrderID > 0 {
		db = db.Where("order_id = ?", req.OrderID)
	}
	return db
}

func applySaleReturnItemFilters(db *gorm.DB, req *salemodels.SaleReturnItemQuery) *gorm.DB {
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
