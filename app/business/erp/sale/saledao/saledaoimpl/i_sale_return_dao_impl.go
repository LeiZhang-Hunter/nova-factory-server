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

// SaleReturnDaoImpl 提供数据访问能力。
type SaleReturnDaoImpl struct {
	db *gorm.DB
}

// NewSaleReturnDao 创建 DAO。
func NewSaleReturnDao(db *gorm.DB) saledao.ISaleReturnDao {
	return &SaleReturnDaoImpl{db: db}
}

func (d *SaleReturnDaoImpl) Create(c *gin.Context, req *salemodels.SaleReturnUpsert) (*salemodels.SaleReturn, error) {
	model := new(salemodels.SaleReturn)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	erpbiz.PrepareCreate(model, c)
	if err := d.db.WithContext(c).Table("erp_sale_return").Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (d *SaleReturnDaoImpl) Update(c *gin.Context, req *salemodels.SaleReturnUpsert) (*salemodels.SaleReturn, error) {
	id := erpbiz.GetIntField(req, "ID")
	if id <= 0 {
		return nil, errors.New("id不能为空")
	}
	model := new(salemodels.SaleReturn)
	if err := erpbiz.CopyStruct(model, req); err != nil {
		return nil, err
	}
	if err := erpbiz.PrepareUpdate(model, c); err != nil {
		return nil, err
	}
	updates := erpbiz.BuildUpdateMap(model)
	db := d.db.WithContext(c).Table("erp_sale_return").Where("id = ?", id)
	if erpbiz.HasField(model, "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	if err := db.Updates(updates).Error; err != nil {
		return nil, err
	}
	return d.GetByID(c, id)
}

func (d *SaleReturnDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	db := d.db.WithContext(c).Table("erp_sale_return").Where("id IN ?", ids)
	if erpbiz.HasField(new(salemodels.SaleReturn), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	return db.Updates(map[string]any{
		"state":       commonStatus.DELETE,
		"update_by":   baizeContext.GetUserId(c),
		"update_time": time.Now(),
	}).Error
}

func (d *SaleReturnDaoImpl) GetByID(c *gin.Context, id int64) (*salemodels.SaleReturn, error) {
	item := new(salemodels.SaleReturn)
	db := d.db.WithContext(c).Table("erp_sale_return").Where("id = ?", id)
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

func (d *SaleReturnDaoImpl) GetByColumn(c *gin.Context, column string, value any) (*salemodels.SaleReturn, error) {
	if column == "" {
		return nil, nil
	}
	item := new(salemodels.SaleReturn)
	db := d.db.WithContext(c).Table("erp_sale_return").Where(fmt.Sprintf("%s = ?", column), value)
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

func (d *SaleReturnDaoImpl) ListPage(c *gin.Context, req *salemodels.SaleReturnQuery) (*erpbiz.PageResult[salemodels.SaleReturn], error) {
	if req == nil {
		req = new(salemodels.SaleReturnQuery)
	}
	db := d.db.WithContext(c).Table("erp_sale_return")
	if erpbiz.HasField(new(salemodels.SaleReturn), "State") {
		db = db.Where("state = ?", commonStatus.NORMAL)
	}
	db = erpbiz.ApplyFilters(db, req)
	page, size := erpbiz.GetPageSize(req)
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]salemodels.SaleReturn, 0)
	orderBy := strings.TrimSpace("id DESC")
	if orderBy == "" {
		orderBy = "id DESC"
	}
	if err := db.Order(orderBy).Offset(int((page - 1) * size)).Limit(int(size)).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make([]*salemodels.SaleReturn, 0, len(rows))
	for i := range rows {
		item := rows[i]
		result = append(result, &item)
	}
	return &erpbiz.PageResult[salemodels.SaleReturn]{Rows: result, Total: total}, nil
}

func (d *SaleReturnDaoImpl) List(c *gin.Context, req *salemodels.SaleReturnQuery) (*salemodels.SaleReturnListData, error) {
	result, err := d.ListPage(c, req)
	if err != nil {
		return nil, err
	}
	return &salemodels.SaleReturnListData{Rows: result.Rows, Total: result.Total}, nil
}
