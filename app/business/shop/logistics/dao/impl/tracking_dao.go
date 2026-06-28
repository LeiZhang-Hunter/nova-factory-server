package impl

import (
	"encoding/json"
	"nova-factory-server/app/business/shop/logistics/dao"
	"nova-factory-server/app/business/shop/logistics/models"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TrackingDaoImpl 物流轨迹数据访问实现
type TrackingDaoImpl struct {
	db    *gorm.DB
	table string
}

// NewTrackingDao 创建物流轨迹 DAO
func NewTrackingDao(db *gorm.DB) dao.ITrackingDao {
	return &TrackingDaoImpl{
		db:    db,
		table: "shop_logistics_tracking_record",
	}
}

// GetSignedRecord 查询已签收的轨迹记录
func (d *TrackingDaoImpl) GetSignedRecord(c *gin.Context, outsid, companyCode string) (*models.TrackingQueryResponse, error) {
	type record struct {
		TraceJSON string `gorm:"column:trace_json"`
	}

	var rec record
	err := d.db.WithContext(c).Table(d.table).
		Select("trace_json").
		Where("outsid = ? AND company_code = ? AND state = 0 AND signed_time IS NOT NULL", outsid, companyCode).
		First(&rec).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	var resp models.TrackingQueryResponse
	if err := json.Unmarshal([]byte(rec.TraceJSON), &resp); err != nil {
		return nil, err
	}
	resp.FromCache = true
	return &resp, nil
}

// SaveSignedRecord 保存已签收的轨迹记录
func (d *TrackingDaoImpl) SaveSignedRecord(c *gin.Context, record *models.TrackingRecordSet) error {
	// 先检查是否已存在
	var count int64
	if err := d.db.WithContext(c).Table(d.table).
		Where("outsid = ? AND company_code = ? AND state = 0", record.Outsid, record.CompanyCode).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil // 已存在，不重复写入
	}

	return d.db.WithContext(c).Table(d.table).Create(map[string]interface{}{
		"id":           snowflake.GenID(),
		"outsid":       record.Outsid,
		"company_code": record.CompanyCode,
		"trace_json":   record.TraceJSON,
		"signed_time":  record.SignedTime,
		"origin_info":  record.OriginInfo,
		"dest_info":    record.DestInfo,
		"state":        0,
	}).Error
}
