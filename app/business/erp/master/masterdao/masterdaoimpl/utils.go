package masterdaoimpl

import (
	"strings"

	"nova-factory-server/app/business/erp/master/mastermodels"

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

func applyProductFilters(db *gorm.DB, req *mastermodels.ProductQuery) *gorm.DB {
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+strings.TrimSpace(req.Name)+"%")
	}
	if req.Code != "" {
		db = db.Where("product_code = ?", req.Code)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	return db
}

func applyProductCategoryFilters(db *gorm.DB, req *mastermodels.ProductCategoryQuery) *gorm.DB {
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+strings.TrimSpace(req.Name)+"%")
	}
	if req.Code != "" {
		db = db.Where("code LIKE ?", "%"+strings.TrimSpace(req.Code)+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.ParentID > 0 {
		db = db.Where("parent_id = ?", req.ParentID)
	}
	return db
}

func applyProductUnitFilters(db *gorm.DB, req *mastermodels.ProductUnitQuery) *gorm.DB {
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+strings.TrimSpace(req.Name)+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	return db
}

func applyCustomerFilters(db *gorm.DB, req *mastermodels.CustomerQuery) *gorm.DB {
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+strings.TrimSpace(req.Name)+"%")
	}
	if req.Code != "" {
		db = db.Where("code LIKE ?", "%"+strings.TrimSpace(req.Code)+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	return db
}

func applySupplierFilters(db *gorm.DB, req *mastermodels.SupplierQuery) *gorm.DB {
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+strings.TrimSpace(req.Name)+"%")
	}
	if req.Code != "" {
		db = db.Where("code LIKE ?", "%"+strings.TrimSpace(req.Code)+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	return db
}

func applyWarehouseFilters(db *gorm.DB, req *mastermodels.WarehouseQuery) *gorm.DB {
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+strings.TrimSpace(req.Name)+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.DefaultStatus != nil {
		db = db.Where("default_status = ?", *req.DefaultStatus)
	}
	return db
}

func applyAccountFilters(db *gorm.DB, req *mastermodels.AccountQuery) *gorm.DB {
	if req.No != "" {
		db = db.Where("no LIKE ?", "%"+strings.TrimSpace(req.No)+"%")
	}
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+strings.TrimSpace(req.Name)+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.DefaultStatus != nil {
		db = db.Where("default_status = ?", *req.DefaultStatus)
	}
	return db
}
