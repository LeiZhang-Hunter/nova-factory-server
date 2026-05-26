package impl

import (
	"errors"
	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/constant/commonStatus"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IApiShopPinkDaoImpl 拼团记录数据访问实现
type IApiShopPinkDaoImpl struct {
	db        *gorm.DB
	tableName string
}

// NewIApiShopPinkDaoImpl 创建拼团记录数据访问对象
func NewIApiShopPinkDaoImpl(ms *gorm.DB) dao.IApiShopPinkDao {
	return &IApiShopPinkDaoImpl{
		db:        ms,
		tableName: "shop_store_pink",
	}
}

// List 查询拼团记录列表（含拼团商品标题和图片）
func (s *IApiShopPinkDaoImpl) List(c *gin.Context, query *models.PinkQuery) (*models.PinkListData, error) {
	db := s.db.WithContext(c).Table(s.tableName+" AS p").
		Select(`p.id, p.uid, p.nickname, p.avatar, p.order_id, p.order_id_key, p.total_num, p.total_price, p.cid, p.pid, p.people, p.price, p.add_time, p.stop_time, p.k_id, p.is_tpl, p.is_refund, p.status, p.is_virtual, c.title AS combination_title, c.image AS combination_image`).
		Joins("LEFT JOIN shop_store_combination AS c ON c.id = p.cid").
		Where("p.state = ?", commonStatus.NORMAL)

	if orderID := strings.TrimSpace(query.OrderID); orderID != "" {
		db = db.Where("p.order_id = ?", orderID)
	}
	if query.CID > 0 {
		db = db.Where("p.cid = ?", query.CID)
	}
	if query.UID > 0 {
		db = db.Where("p.uid = ?", query.UID)
	}
	if query.Status != nil {
		db = db.Where("p.status = ?", *query.Status)
	}
	if query.IsRefund != nil {
		db = db.Where("p.is_refund = ?", *query.IsRefund)
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Size <= 0 {
		query.Size = 20
	}
	if query.Size > 200 {
		query.Size = 200
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	rows := make([]*models.Pink, 0)
	if err := db.Offset(int((query.Page - 1) * query.Size)).
		Limit(int(query.Size)).
		Order("p.id DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	return &models.PinkListData{Rows: rows, Total: total}, nil
}

// GetByID 根据主键获取拼团记录
func (s *IApiShopPinkDaoImpl) GetByID(c *gin.Context, id int64) (*models.Pink, error) {
	var item models.Pink
	if err := s.db.WithContext(c).Table(s.tableName+" AS p").
		Select(`p.id, p.uid, p.nickname, p.avatar, p.order_id, p.order_id_key, p.total_num, p.total_price, p.cid, p.pid, p.people, p.price, p.add_time, p.stop_time, p.k_id, p.is_tpl, p.is_refund, p.status, p.is_virtual, c.title AS combination_title, c.image AS combination_image`).
		Joins("LEFT JOIN shop_store_combination AS c ON c.id = p.cid").
		Where("p.id = ?", id).
		Where("p.state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// CountMembers 统计团内成员数（含子团）
func (s *IApiShopPinkDaoImpl) CountMembers(c *gin.Context, pinkID int64) (int64, error) {
	var count int64
	err := s.db.WithContext(c).Table(s.tableName).
		Where("id = ? OR k_id = ?", pinkID, pinkID).
		Where("state = ?", commonStatus.NORMAL).
		Count(&count).Error
	return count, err
}
