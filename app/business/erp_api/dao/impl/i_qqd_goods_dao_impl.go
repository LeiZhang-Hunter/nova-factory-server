package impl

import (
	"context"

	"nova-factory-server/app/business/erp_api/dao"
	"nova-factory-server/app/business/erp_api/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IQQDGoodsDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewIQQDGoodsDaoImpl(db *gorm.DB) dao.IQQDGoodsDao {
	return &IQQDGoodsDaoImpl{
		db:        db,
		tableName: "shop_goods",
	}
}

func (d *IQQDGoodsDaoImpl) List(ctx context.Context, pageNo, pageSize int) ([]models.QQDProductTable, int64, error) {
	var total int64
	db := d.db.WithContext(ctx).Table(d.tableName).
		Where("state = ?", 0).
		Where("is_on_sale = ?", 1)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var products []models.QQDProductTable
	if err := db.Order("goods_id asc").
		Limit(pageSize).
		Offset((pageNo - 1) * pageSize).
		Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (d *IQQDGoodsDaoImpl) UpdateQuantity(ctx context.Context, goodsID string, quantity int64) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var lockedGoods struct {
			GoodsID string `gorm:"column:goods_id"`
		}
		if err := tx.Table(d.tableName).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Select("goods_id").
			Where("goods_id = ?", goodsID).
			First(&lockedGoods).Error; err != nil {
			return err
		}
		return tx.Table(d.tableName).
			Where("goods_id = ?", goodsID).
			Updates(map[string]any{"quantity": quantity}).Error
	})
}
