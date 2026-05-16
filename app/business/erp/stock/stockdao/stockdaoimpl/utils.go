package stockdaoimpl

import (
	"errors"
	"strings"
	"time"

	"nova-factory-server/app/business/erp/stock/stockmodels"

	"gorm.io/gorm"
)

func parseTime(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	layouts := []string{"2006-01-02 15:04:05", time.RFC3339, "2006-01-02"}
	for _, layout := range layouts {
		if parsed, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return &parsed, nil
		}
	}
	return nil, errors.New("时间格式错误")
}

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

func applyStockCheckFilters(db *gorm.DB, req *stockmodels.StockCheckQuery) *gorm.DB {
	if req.No != "" {
		db = db.Where("no LIKE ?", "%"+strings.TrimSpace(req.No)+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	return db
}

func applyStockInFilters(db *gorm.DB, req *stockmodels.StockInQuery) *gorm.DB {
	if req.No != "" {
		db = db.Where("no LIKE ?", "%"+strings.TrimSpace(req.No)+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.SupplierID > 0 {
		db = db.Where("supplier_id = ?", req.SupplierID)
	}
	return db
}

func applyStockMoveFilters(db *gorm.DB, req *stockmodels.StockMoveQuery) *gorm.DB {
	if req.No != "" {
		db = db.Where("no LIKE ?", "%"+strings.TrimSpace(req.No)+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	return db
}

func applyStockOutFilters(db *gorm.DB, req *stockmodels.StockOutQuery) *gorm.DB {
	if req.No != "" {
		db = db.Where("no LIKE ?", "%"+strings.TrimSpace(req.No)+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.CustomerID > 0 {
		db = db.Where("customer_id = ?", req.CustomerID)
	}
	return db
}

func applyStockInItemFilters(db *gorm.DB, req *stockmodels.StockInItemQuery) *gorm.DB {
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

func applyStockCheckItemFilters(db *gorm.DB, req *stockmodels.StockCheckItemQuery) *gorm.DB {
	if req.WarehouseID > 0 {
		db = db.Where("warehouse_id = ?", req.WarehouseID)
	}
	if req.ProductID > 0 {
		db = db.Where("product_id = ?", req.ProductID)
	}
	if req.CheckID > 0 {
		db = db.Where("check_id = ?", req.CheckID)
	}
	return db
}

func applyStockMoveItemFilters(db *gorm.DB, req *stockmodels.StockMoveItemQuery) *gorm.DB {
	if req.ProductID > 0 {
		db = db.Where("product_id = ?", req.ProductID)
	}
	if req.MoveID > 0 {
		db = db.Where("move_id = ?", req.MoveID)
	}
	return db
}

func applyStockOutItemFilters(db *gorm.DB, req *stockmodels.StockOutItemQuery) *gorm.DB {
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

func applyStockRecordFilters(db *gorm.DB, req *stockmodels.StockRecordQuery) *gorm.DB {
	if req.WarehouseID > 0 {
		db = db.Where("warehouse_id = ?", req.WarehouseID)
	}
	if req.ProductID > 0 {
		db = db.Where("product_id = ?", req.ProductID)
	}
	if req.BizType != nil {
		db = db.Where("biz_type = ?", *req.BizType)
	}
	return db
}

func applyStockFilters(db *gorm.DB, req *stockmodels.StockQuery) *gorm.DB {
	if req.WarehouseID > 0 {
		db = db.Where("warehouse_id = ?", req.WarehouseID)
	}
	if req.ProductID > 0 {
		db = db.Where("product_id = ?", req.ProductID)
	}
	return db
}
