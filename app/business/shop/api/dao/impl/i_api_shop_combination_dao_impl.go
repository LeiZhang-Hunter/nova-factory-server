package impl

import (
	"errors"
	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/fileUtils"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IApiShopCombinationDaoImpl 拼团商品数据访问实现
type IApiShopCombinationDaoImpl struct {
	db        *gorm.DB
	tableName string
}

// NewIApiShopCombinationDaoImpl 创建拼团商品数据访问对象
func NewIApiShopCombinationDaoImpl(ms *gorm.DB) dao.IApiShopCombinationDao {
	return &IApiShopCombinationDaoImpl{
		db:        ms,
		tableName: "shop_store_combination",
	}
}

// List 查询拼团商品列表（含当前拼团人数）
func (s *IApiShopCombinationDaoImpl) List(c *gin.Context, query *models.CombinationQuery) (*models.CombinationListData, error) {
	db := s.db.WithContext(c).Table(s.tableName+" AS c").
		Select(`c.id, c.product_id, c.mer_id, COALESCE(NULLIF(c.image, ''), g.image_url) AS image, c.images, c.title, g.goods_name AS goods_name, c.attr, c.people, c.info, c.price, g.retail_price AS ot_price, c.sort, c.sales, c.stock, c.is_host, c.is_show, c.is_postage, c.postage, c.start_time, c.stop_time, c.effective_time, c.browse, c.unit_name, c.weight, c.volume, c.num, c.once_num, c.quota, c.quota_show, c.virtual, c.home_module_ids, (SELECT COUNT(*) FROM shop_store_pink WHERE cid = c.id AND status = 1 AND is_refund = 0 AND state = 0) AS pink_count`).
		Joins("LEFT JOIN shop_goods AS g ON c.product_id = g.goods_id").
		Where("c.state = ?", commonStatus.NORMAL)

	if title := strings.TrimSpace(query.Title); title != "" {
		db = db.Where("c.title LIKE ?", "%"+title+"%")
	}
	if query.ProductID > 0 {
		db = db.Where("c.product_id = ?", query.ProductID)
	}
	if query.IsShow != nil {
		db = db.Where("c.is_show = ?", *query.IsShow)
	}
	if query.IsHost != nil {
		db = db.Where("c.is_host = ?", *query.IsHost)
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

	rows := make([]*models.Combination, 0)
	if err := db.Offset(int((query.Page - 1) * query.Size)).
		Limit(int(query.Size)).
		Order("sort ASC").
		Order("id DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		if row != nil {
			row.Image = fileUtils.BuildAbsoluteURL(c, row.Image)
		}
	}

	return &models.CombinationListData{Rows: rows, Total: total}, nil
}

// GetByID 根据主键获取拼团商品
func (s *IApiShopCombinationDaoImpl) GetByID(c *gin.Context, id int64) (*models.Combination, error) {
	var item models.Combination
	if err := s.db.WithContext(c).Table(s.tableName+" AS c").
		Select("c.*, g.goods_name AS goods_name, g.retail_price AS ot_price").
		Joins("LEFT JOIN shop_goods AS g ON c.product_id = g.goods_id").
		Where("c.id = ?", id).
		Where("c.state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	item.Image = fileUtils.BuildAbsoluteURL(c, item.Image)
	return &item, nil
}
