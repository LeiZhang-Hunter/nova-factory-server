package impl

import (
	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// IApiShopOrderItemDaoImpl 提供订单商品明细的数据库访问能力。
type IApiShopOrderItemDaoImpl struct {
	db        *gorm.DB
	tableName string
}

// NewIApiShopOrderItemDaoImpl 创建订单商品明细 DAO 实现。
func NewIApiShopOrderItemDaoImpl(db *gorm.DB) dao.IApiShopOrderItemDao {
	return &IApiShopOrderItemDaoImpl{
		db:        db,
		tableName: "shop_order_item",
	}
}

// BatchCreate 批量新增订单商品明细记录。
func (d *IApiShopOrderItemDaoImpl) BatchCreate(c *gin.Context, items []*models.OrderItem) error {
	if len(items) == 0 {
		return nil
	}
	return d.db.WithContext(c).Table(d.tableName).CreateInBatches(items, 100).Error
}

// GetByOrderID 根据订单ID获取商品明细列表。
func (d *IApiShopOrderItemDaoImpl) GetByOrderID(c *gin.Context, orderID int64) ([]*models.OrderItem, error) {
	var items []*models.OrderItem
	err := d.db.WithContext(c).Table(d.tableName).Where("order_id = ? AND state = 0", orderID).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

// GetByOrderNo 根据订单号获取商品明细列表。
func (d *IApiShopOrderItemDaoImpl) GetByOrderNo(c *gin.Context, orderNo string) ([]*models.OrderItem, error) {
	var items []*models.OrderItem
	err := d.db.WithContext(c).Table(d.tableName).Where("order_no = ? AND state = 0", orderNo).Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}
