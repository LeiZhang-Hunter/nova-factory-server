package impl

import (
	"errors"
	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/utils/fileUtils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IApiShopGoodsDaoImpl 商品数据访问实现
type IApiShopGoodsDaoImpl struct {
	db        *gorm.DB
	tableName string
}

// NewIApiShopGoodsDaoImpl  创建商品数据访问对象
func NewIApiShopGoodsDaoImpl(ms *gorm.DB) dao.IApiShopGoodsDao {
	return &IApiShopGoodsDaoImpl{
		db:        ms,
		tableName: "shop_goods",
	}
}

// GetByID 根据 ID 查询
func (s *IApiShopGoodsDaoImpl) GetByID(c *gin.Context, id int64) (*models.Goods, error) {
	var item models.Goods
	if err := s.db.WithContext(c).Table(s.tableName).
		Where("id = ?", id).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	item.ImageURL = fileUtils.BuildAbsoluteURL(c, item.ImageURL)
	return &item, nil
}

// GetByGoodsID 根据商品业务ID查询
func (s *IApiShopGoodsDaoImpl) GetByGoodsID(c *gin.Context, goodsID string) (*models.Goods, error) {
	var item models.Goods
	if err := s.db.WithContext(c).Table(s.tableName).
		Where("goods_id = ?", goodsID).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	item.ImageURL = fileUtils.BuildAbsoluteURL(c, item.ImageURL)
	return &item, nil
}

// List 查询商品列表
func (s *IApiShopGoodsDaoImpl) List(c *gin.Context, query *models.GoodsQuery) (*models.GoodsListData, error) {
	db := s.db.WithContext(c).Table(s.tableName)
	if query.GoodsName != "" {
		db = db.Where("goods_name LIKE ?", "%"+query.GoodsName+"%")
	}
	if query.GoodsCode != "" {
		db = db.Where("goods_code = ?", query.GoodsCode)
	}
	if query.CategoryId > 0 {
		db = db.Where("shop_category_id = ?", query.CategoryId)
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Size <= 0 {
		query.Size = 20
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	rows := make([]*models.Goods, 0)
	orderClause := "id DESC"
	if query.SortBy == "retailPrice" && query.SortOrder != "" {
		orderClause = "retail_price " + query.SortOrder
	}
	if err := db.Offset(int((query.Page - 1) * query.Size)).Limit(int(query.Size)).Order(orderClause).Find(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		if row != nil {
			row.ImageURL = fileUtils.BuildAbsoluteURL(c, row.ImageURL)
		}
	}

	return &models.GoodsListData{
		Rows:  rows,
		Total: total,
	}, nil
}

// ListByUserPurchased 查询用户已购买的商品（用于复购）
func (s *IApiShopGoodsDaoImpl) ListByUserPurchased(c *gin.Context, userID int64, categoryID int64, page, size int64) (*models.GoodsListData, error) {
	// 子查询：获取用户已完成订单的商品ID列表
	subQuery := s.db.Table("shop_order_item").
		Select("DISTINCT goods_id").
		Joins("JOIN shop_order ON shop_order.id = shop_order_item.order_id").
		Where("shop_order.user_id = ?", userID).
		Where("shop_order.status = 3") // 已完成订单

	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}

	db := s.db.WithContext(c).Table(s.tableName).
		Where("goods_id IN (?)", subQuery)

	if categoryID > 0 {
		db = db.Where("shop_category_id = ?", categoryID)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}

	rows := make([]*models.Goods, 0)
	offset := int((page - 1) * size)
	if err := db.Offset(offset).Limit(int(size)).Order("id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		if row != nil {
			row.ImageURL = fileUtils.BuildAbsoluteURL(c, row.ImageURL)
		}
	}

	return &models.GoodsListData{
		Rows:  rows,
		Total: total,
	}, nil
}
