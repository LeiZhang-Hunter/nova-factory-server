package impl

import (
	"context"

	"nova-factory-server/app/business/erp_api/dao"
	"nova-factory-server/app/business/erp_api/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IQQDGoodsSkuDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewIQQDGoodsSkuDaoImpl(db *gorm.DB) dao.IQQDGoodsSkuDao {
	return &IQQDGoodsSkuDaoImpl{
		db:        db,
		tableName: "shop_goods_sku",
	}
}

func (d *IQQDGoodsSkuDaoImpl) ListByGoodsIDs(ctx context.Context, goodsIDs []string) ([]models.QQDProductSkuTable, error) {
	if len(goodsIDs) == 0 {
		return nil, nil
	}

	var productSkus []models.QQDProductSkuTable
	err := d.db.WithContext(ctx).Table(d.tableName).
		Where("state = ?", 0).
		Where("goods_id IN ?", goodsIDs).
		Order("goods_id asc, sku_id asc").
		Find(&productSkus).Error
	return productSkus, err
}

func (d *IQQDGoodsSkuDaoImpl) UpdateQuantity(ctx context.Context, skuID string, quantity int64, lock bool) error {
	return d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if lock {
			var lockedSku struct {
				SkuID string `gorm:"column:sku_id"`
			}
			if err := tx.Table(d.tableName).
				Clauses(clause.Locking{Strength: "UPDATE"}).
				Select("sku_id").
				Where("sku_id = ?", skuID).
				First(&lockedSku).Error; err != nil {
				return err
			}
		}
		return tx.Table(d.tableName).
			Where("sku_id = ?", skuID).
			Updates(map[string]any{"quantity": quantity}).Error
	})
}
