package daoimpl

import (
	"gorm.io/gorm"
	"nova-factory-server/app/business/shop/finance/models"
	"strings"
)

func applyAccountFilters(db *gorm.DB, req *models.AccountQuery) *gorm.DB {
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
