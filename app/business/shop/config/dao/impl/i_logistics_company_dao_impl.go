package impl

import (
	"errors"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/snowflake"
	"strings"
	"time"

	"nova-factory-server/app/business/shop/config/dao"
	"nova-factory-server/app/business/shop/config/models"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LogisticsCompanyDaoImpl 提供 ERP 物流公司表的数据访问能力。
type LogisticsCompanyDaoImpl struct {
	db    *gorm.DB
	table string
}

// NewLogisticsCompanyDao 创建 ERP 物流公司 DAO。
func NewLogisticsCompanyDao(db *gorm.DB) dao.ILogisticsCompanyDao {
	return &LogisticsCompanyDaoImpl{
		db:    db,
		table: "shop_logistics_company",
	}
}

// Create 新增 ERP 物流公司记录。
func (l *LogisticsCompanyDaoImpl) Create(c *gin.Context, req *models.LogisticsCompanyUpsert) (*models.LogisticsCompany, error) {
	model := &models.LogisticsCompany{
		ID:      snowflake.GenID(),
		Code:    strings.TrimSpace(req.Code),
		Name:    strings.TrimSpace(req.Name),
		Company: strings.TrimSpace(req.Company),

		ProvinceName: strings.TrimSpace(req.ProvinceName),
		ProvinceCode: strings.TrimSpace(req.ProvinceCode),

		CityName: strings.TrimSpace(req.CityName),
		CityCode: strings.TrimSpace(req.CityCode),

		DistrictName: strings.TrimSpace(req.DistrictName),
		DistrictCode: strings.TrimSpace(req.DistrictCode),

		StreetName: strings.TrimSpace(req.StreetName),
		StreetCode: strings.TrimSpace(req.StreetCode),

		ShortName:    strings.TrimSpace(req.ShortName),
		ContactName:  strings.TrimSpace(req.ContactName),
		ContactPhone: strings.TrimSpace(req.ContactPhone),
		Address:      strings.TrimSpace(req.Address),
		Remark:       strings.TrimSpace(req.Remark),
		Sort:         req.Sort,
		Status:       req.Status,
		DeptID:       baizeContext.GetDeptId(c),
		State:        commonStatus.NORMAL,
	}
	model.SetCreateBy(baizeContext.GetUserId(c))
	if err := l.db.WithContext(c).Table(l.table).Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

// Update 修改 ERP 物流公司记录。
func (l *LogisticsCompanyDaoImpl) Update(c *gin.Context, req *models.LogisticsCompanyUpsert) (*models.LogisticsCompany, error) {
	model := &models.LogisticsCompany{
		ID:      req.ID,
		Code:    strings.TrimSpace(req.Code),
		Name:    strings.TrimSpace(req.Name),
		Company: strings.TrimSpace(req.Company),

		ProvinceName: strings.TrimSpace(req.ProvinceName),
		ProvinceCode: strings.TrimSpace(req.ProvinceCode),

		CityName: strings.TrimSpace(req.CityName),
		CityCode: strings.TrimSpace(req.CityCode),

		DistrictName: strings.TrimSpace(req.DistrictName),
		DistrictCode: strings.TrimSpace(req.DistrictCode),

		StreetName: strings.TrimSpace(req.StreetName),
		StreetCode: strings.TrimSpace(req.StreetCode),

		ShortName:    strings.TrimSpace(req.ShortName),
		ContactName:  strings.TrimSpace(req.ContactName),
		ContactPhone: strings.TrimSpace(req.ContactPhone),
		Address:      strings.TrimSpace(req.Address),
		Remark:       strings.TrimSpace(req.Remark),
		Sort:         req.Sort,
		Status:       req.Status,
	}
	model.SetUpdateBy(baizeContext.GetUserId(c))
	if err := l.db.WithContext(c).Table(l.table).
		Where("id = ?", req.ID).
		Where("state = 0").
		Select("code", "name", "short_name", "contact_name", "contact_phone", "address", "website", "remark", "sort", "status", "update_by", "update_time").
		Updates(model).Error; err != nil {
		return nil, err
	}
	return l.GetByID(c, int64(req.ID))
}

// DeleteByIDs 软删除 ERP 物流公司记录。
func (l *LogisticsCompanyDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	now := time.Now()
	return l.db.WithContext(c).Table(l.table).Where("id IN ?", ids).Updates(map[string]interface{}{
		"state":       -1,
		"update_by":   baizeContext.GetUserId(c),
		"update_time": now,
	}).Error
}

// GetByID 根据主键查询 ERP 物流公司记录。
func (l *LogisticsCompanyDaoImpl) GetByID(c *gin.Context, id int64) (*models.LogisticsCompany, error) {
	var item models.LogisticsCompany
	if err := l.db.WithContext(c).Table(l.table).
		Where("id = ?", id).
		Where("state = 0").
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// GetByCode 根据编码查询 ERP 物流公司记录。
func (l *LogisticsCompanyDaoImpl) GetByCode(c *gin.Context, code string) (*models.LogisticsCompany, error) {
	var item models.LogisticsCompany
	if err := l.db.WithContext(c).Table(l.table).
		Where("code = ?", strings.TrimSpace(code)).
		Where("state = 0").
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// GetByName 根据名称查询 ERP 物流公司记录。
func (l *LogisticsCompanyDaoImpl) GetByName(c *gin.Context, name string) (*models.LogisticsCompany, error) {
	var item models.LogisticsCompany
	if err := l.db.WithContext(c).Table(l.table).
		Where("name = ?", strings.TrimSpace(name)).
		Where("state = 0").
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// List 分页查询 ERP 物流公司记录。
func (l *LogisticsCompanyDaoImpl) List(c *gin.Context, req *models.LogisticsCompanyQuery) (*models.LogisticsCompanyListData, error) {
	db := l.db.WithContext(c).Table(l.table).Where("state = 0")
	if req.Code != "" {
		db = db.Where("code LIKE ?", "%"+strings.TrimSpace(req.Code)+"%")
	}
	if req.Name != "" {
		db = db.Where("name LIKE ?", "%"+strings.TrimSpace(req.Name)+"%")
	}
	if req.Status != nil {
		db = db.Where("status = ?", req.Status)
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]*models.LogisticsCompany, 0)
	if err := db.Order("sort ASC, id DESC").
		Offset(int((req.Page - 1) * req.Size)).
		Limit(int(req.Size)).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return &models.LogisticsCompanyListData{
		Rows:  rows,
		Total: total,
	}, nil
}
