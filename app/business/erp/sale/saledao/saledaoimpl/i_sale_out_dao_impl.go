package saledaoimpl

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"nova-factory-server/app/business/erp/erpbiz"
	"nova-factory-server/app/business/erp/sale/saledao"
	"nova-factory-server/app/business/erp/sale/salemodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SaleOutDaoImpl 提供数据访问能力。
type SaleOutDaoImpl struct {
	db *gorm.DB
}

// NewSaleOutDao 创建 DAO。
func NewSaleOutDao(db *gorm.DB) saledao.ISaleOutDao {
	return &SaleOutDaoImpl{db: db}
}

func (d *SaleOutDaoImpl) Create(c *gin.Context, req *salemodels.SaleOutUpsert) (*salemodels.SaleOut, error) {
	model := new(salemodels.SaleOut)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	erpbiz.PrepareCreate(model, c)
	if err := d.db.WithContext(c).Table("erp_sale_out").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *SaleOutDaoImpl) Update(c *gin.Context, req *salemodels.SaleOutUpsert) (*salemodels.SaleOut, error) {
	id := erpbiz.GetIntField(req, "ID")
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	model := new(salemodels.SaleOut)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	if err := erpbiz.PrepareUpdate(model, c); err != nil {
		return nil, err
	}
	updates := erpbiz.BuildUpdateMap(model)
	db := d.db.WithContext(c).Table("erp_sale_out").Where("id = ?", id)
	if erpbiz.HasField(model, "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, id)
}

func (d *SaleOutDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	db := d.db.WithContext(c).Table("erp_sale_out").Where("id IN ?", ids)
	if erpbiz.HasField(new(salemodels.SaleOut), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	return db.Updates(map[string]any{
		"state":       commonStatus.DELETE,
		"update_by":   baizeContext.GetUserId(c),
		"update_time": time.Now(),
	}).Error
}

func (d *SaleOutDaoImpl) GetByID(c *gin.Context, id int64) (*salemodels.SaleOut, error) {
	item := new(salemodels.SaleOut)
	db := d.db.WithContext(c).Table("erp_sale_out").Where("id = ?", id)
	if erpbiz.HasField(item, "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	if err := db.First(item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (d *SaleOutDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*salemodels.SaleOut, error) {
	if column == "" {
		return nil, nil
	}
	item := new(salemodels.SaleOut)
	db := d.db.WithContext(c).Table("erp_sale_out").Where(fmt.Sprintf("%s = ?", column), value)
	if erpbiz.HasField(item, "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	if err := db.First(item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return item, nil
}

func (d *SaleOutDaoImpl) ListPage(c *gin.Context, req *salemodels.SaleOutQuery) (*erpbiz.PageResult[salemodels.SaleOut], error) {
	if req == nil {
		req = new(salemodels.SaleOutQuery)
	}
	db := d.db.WithContext(c).Table("erp_sale_out")
	if erpbiz.HasField(new(salemodels.SaleOut), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	db = erpbiz.ApplyFilters(db, req)
	page, size := erpbiz.GetPageSize(req)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]salemodels.SaleOut, 0)
	orderBy := strings.TrimSpace("id DESC")
	if orderBy == "" {
		orderBy = "id DESC"
	}
	if err := db.Order(orderBy).Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*salemodels.SaleOut, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &erpbiz.PageResult[salemodels.SaleOut]{Rows: result, Total: total}, nil
}

func (d *SaleOutDaoImpl) List(c *gin.Context, req *salemodels.SaleOutQuery) (*salemodels.SaleOutListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &salemodels.SaleOutListData{Rows: result.Rows, Total: result.Total}, nil
}
