package impl

import (
	"errors"
	"nova-factory-server/app/business/shop/activity/dao"
	"nova-factory-server/app/business/shop/activity/models"
	"nova-factory-server/app/constant/commonStatus"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ShopPinkDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewShopPinkDao(ms *gorm.DB) dao.IShopPinkDao {
	return &ShopPinkDaoImpl{db: ms, tableName: "shop_store_pink"}
}

func (s *ShopPinkDaoImpl) GetByID(c *gin.Context, id int64) (*models.Pink, error) {
	var item models.Pink
	if err := s.baseQuery(c).
		Where("p.id = ?", id).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (s *ShopPinkDaoImpl) List(c *gin.Context, req *models.PinkQuery) (*models.PinkListData, error) {
	db := s.baseQuery(c)
	if orderID := strings.TrimSpace(req.OrderID); orderID != "" {
		db = db.Where("p.order_id = ?", orderID)
	}
	if req.CID > 0 {
		db = db.Where("p.cid = ?", req.CID)
	}
	if req.UID > 0 {
		db = db.Where("p.uid = ?", req.UID)
	}
	if req.Status != nil {
		db = db.Where("p.status = ?", *req.Status)
	}
	if req.IsRefund != nil {
		db = db.Where("p.is_refund = ?", *req.IsRefund)
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}
	if req.Size > 200 {
		req.Size = 200
	}
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]*models.Pink, 0)
	if err := db.Offset(int((req.Page - 1) * req.Size)).Limit(int(req.Size)).Order("p.id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return &models.PinkListData{Rows: rows, Total: total}, nil
}

func (s *ShopPinkDaoImpl) baseQuery(c *gin.Context) *gorm.DB {
	return s.db.WithContext(c).Table(s.tableName+" AS p").
		Select(`p.id, p.uid, p.nickname, p.avatar, p.order_id, p.order_id_key, p.total_num, p.total_price, p.cid, p.pid, p.people, p.price, p.add_time, p.stop_time, p.k_id, p.is_tpl, p.is_refund, p.status, p.is_virtual, p.dept_id, p.state, p.create_time, p.update_time, c.title AS combination_title, c.image AS combination_image`).
		Joins("LEFT JOIN shop_store_combination AS c ON c.id = p.cid").
		Where("p.state = ?", commonStatus.NORMAL)
}
