package impl

import (
	"errors"
	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/constant/commonStatus"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// IApiShopSkuDaoImpl 提供 App 端下单所需的 SKU 数据访问能力。
type IApiShopSkuDaoImpl struct {
	db        *gorm.DB
	tableName string
}

// NewIApiShopSkuDaoImpl 创建 App 端 SKU DAO。
func NewIApiShopSkuDaoImpl(db *gorm.DB) dao.IApiShopSkuDao {
	return &IApiShopSkuDaoImpl{
		db:        db,
		tableName: "shop_goods_sku",
	}
}

// GetByID 根据主键查询 SKU。
func (s *IApiShopSkuDaoImpl) GetByID(c *gin.Context, id int64) (*shopmodels.GoodsSku, error) {
	var item shopmodels.GoodsSku
	if err := getCurrentDB(c, s.db).WithContext(c).
		Table(s.tableName).
		Where("id = ?", id).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// GetByIDForUpdate 在当前事务中按主键锁定 SKU 行。
func (s *IApiShopSkuDaoImpl) GetByIDForUpdate(c *gin.Context, id int64) (*shopmodels.GoodsSku, error) {
	var item shopmodels.GoodsSku
	if err := getCurrentDB(c, s.db).WithContext(c).
		Table(s.tableName).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", id).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// ListByIDs 根据主键列表批量查询 SKU。
func (s *IApiShopSkuDaoImpl) ListByIDs(c *gin.Context, ids []int64) ([]*shopmodels.GoodsSku, error) {
	if len(ids) == 0 {
		return make([]*shopmodels.GoodsSku, 0), nil
	}

	rows := make([]*shopmodels.GoodsSku, 0, len(ids))
	if err := getCurrentDB(c, s.db).WithContext(c).
		Table(s.tableName).
		Where("id IN ?", ids).
		Where("state = ?", commonStatus.NORMAL).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *IApiShopSkuDaoImpl) ListByGoodsIDs(c *gin.Context, goodsIDs []int64) ([]*shopmodels.GoodsSku, error) {
	if len(goodsIDs) == 0 {
		return make([]*shopmodels.GoodsSku, 0), nil
	}

	rows := make([]*shopmodels.GoodsSku, 0)
	if err := getCurrentDB(c, s.db).WithContext(c).
		Table(s.tableName).
		Where("goods_id IN ?", goodsIDs).
		Where("state = ?", commonStatus.NORMAL).
		Order("id ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// DeductStock 原子扣减 SKU 库存。
func (s *IApiShopSkuDaoImpl) DeductStock(c *gin.Context, id int64, quantity int64) error {
	if quantity <= 0 {
		return errors.New("扣减库存数量必须大于0")
	}
	result := getCurrentDB(c, s.db).WithContext(c).
		Table(s.tableName).
		Where("id = ?", id).
		Where("quantity >= ?", quantity).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]interface{}{
			"quantity":    gorm.Expr("quantity - ?", quantity),
			"update_time": gorm.Expr("NOW()"),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("库存不足")
	}
	return nil
}

// RestoreStock 原子回补 SKU 库存。
func (s *IApiShopSkuDaoImpl) RestoreStock(c *gin.Context, id int64, quantity int64) error {
	if quantity <= 0 {
		return errors.New("回补库存数量必须大于0")
	}
	result := getCurrentDB(c, s.db).WithContext(c).
		Table(s.tableName).
		Where("id = ?", id).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]interface{}{
			"quantity":    gorm.Expr("quantity + ?", quantity),
			"update_time": gorm.Expr("NOW()"),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("商品规格不存在")
	}
	return nil
}

func (s *IApiShopSkuDaoImpl) ListBySkuIDs(c *gin.Context, ids []int64) ([]*shopmodels.GoodsSku, error) {
	if len(ids) == 0 {
		return make([]*shopmodels.GoodsSku, 0), nil
	}

	rows := make([]*shopmodels.GoodsSku, 0, len(ids))
	if err := getCurrentDB(c, s.db).WithContext(c).
		Table(s.tableName).
		Where("sku_id IN ?", ids).
		Where("state = ?", commonStatus.NORMAL).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// GetBySkuID 根据 sku id查询 SKU。
func (s *IApiShopSkuDaoImpl) GetBySkuID(c *gin.Context, skuId int64) (*shopmodels.GoodsSku, error) {
	var item shopmodels.GoodsSku
	if err := getCurrentDB(c, s.db).WithContext(c).
		Table(s.tableName).
		Where("sku_id = ?", skuId).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// GetBySkuIDForUpdate 在当前事务中按主键锁定 SKU 行。
func (s *IApiShopSkuDaoImpl) GetBySkuIDForUpdate(c *gin.Context, id int64) (*shopmodels.GoodsSku, error) {
	var item shopmodels.GoodsSku
	if err := getCurrentDB(c, s.db).WithContext(c).
		Table(s.tableName).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("sku_id = ?", id).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// DeductStockBySkuId 原子扣减 SKU by sku id 库存。
func (s *IApiShopSkuDaoImpl) DeductStockBySkuId(c *gin.Context, id int64, quantity int64) error {
	if quantity <= 0 {
		return errors.New("扣减库存数量必须大于0")
	}
	result := getCurrentDB(c, s.db).WithContext(c).
		Table(s.tableName).
		Where("sku_id = ?", id).
		Where("quantity >= ?", quantity).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]interface{}{
			"quantity":    gorm.Expr("quantity - ?", quantity),
			"update_time": gorm.Expr("NOW()"),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("库存不足")
	}
	return nil
}
