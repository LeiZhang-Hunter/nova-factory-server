package impl

import (
	"errors"
	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"
	shopConstant "nova-factory-server/app/constant/shop"
	"nova-factory-server/app/utils/snowflake"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// IApiShopCartDaoImpl 购物车，这个购物车是用来给APp端用的
type IApiShopCartDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewIApiShopCartDaoImpl(db *gorm.DB) dao.IApiShopCartDao {
	return &IApiShopCartDaoImpl{
		db:        db,
		tableName: "shop_user_cart",
	}
}

func (s *IApiShopCartDaoImpl) Save(c *gin.Context, req *models.CartSetData) (*models.CartDto, error) {
	var result *models.CartDto
	now := time.Now()
	state := req.State
	if state == 0 {
		state = shopConstant.CartStateNormal
	}
	err := s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {

		existing, err := s.getByMergeIdentityTx(c, tx, req.UserID, req.SkuID, req.ProductType, req.ActivityID, state)
		if err != nil {
			return err
		}
		if existing != nil {
			if err := tx.Table(s.tableName).
				Where("id = ?", existing.ID).
				Where("state = ?", state).
				Updates(map[string]interface{}{
					"goods_id":     req.GoodsID,
					"goods_name":   req.GoodsName,
					"sku_name":     req.SkuName,
					"image_url":    req.ImageURL,
					"retail_price": req.RetailPrice,
					"quantity":     gorm.Expr("quantity + ?", req.Quantity),
					"product_type": req.ProductType,
					"activity_id":  req.ActivityID,
					"pink_id":      req.PinkID,
					"update_time":  gorm.Expr("NOW()"),
				}).Error; err != nil {
				return err
			}
			result, err = s.getByIDAndStateTx(c, tx, existing.ID, state)
			return err
		}

		model := &models.CartDto{
			ID:          snowflake.GenID(),
			UserID:      req.UserID,
			GoodsID:     req.GoodsID,
			SkuID:       int64(req.SkuID),
			GoodsName:   req.GoodsName,
			SkuName:     req.SkuName,
			ImageURL:    req.ImageURL,
			RetailPrice: req.RetailPrice,
			Quantity:    req.Quantity,
			ProductType: req.ProductType,
			ActivityID:  req.ActivityID,
			PinkID:      req.PinkID,
			State:       state,
			CreateTime:  &now,
		}
		if err := tx.Table(s.tableName).Create(model).Error; err != nil {
			if isDuplicateKeyError(err) {
				existing, getErr := s.getByMergeIdentityTx(c, tx, req.UserID, req.SkuID, req.ProductType, req.ActivityID, state)
				if getErr != nil {
					return getErr
				}
				if existing == nil {
					return err
				}
				if updateErr := tx.Table(s.tableName).
					Where("id = ?", existing.ID).
					Where("state = ?", state).
					Updates(map[string]interface{}{
						"goods_id":     req.GoodsID,
						"goods_name":   req.GoodsName,
						"sku_name":     req.SkuName,
						"image_url":    req.ImageURL,
						"retail_price": req.RetailPrice,
						"quantity":     gorm.Expr("quantity + ?", req.Quantity),
						"product_type": req.ProductType,
						"activity_id":  req.ActivityID,
						"pink_id":      req.PinkID,
						"update_time":  gorm.Expr("NOW()"),
					}).Error; updateErr != nil {
					return updateErr
				}
				result, err = s.getByIDAndStateTx(c, tx, existing.ID, state)
				return err
			}
			return err
		}
		result = model
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// List 查询购物车列表（不分页）。
func (s *IApiShopCartDaoImpl) List(c *gin.Context, userID int64) ([]*models.CartDto, error) {
	db := s.db.WithContext(c).Table(s.tableName).Where("user_id = ?", userID).Where("state = ?", shopConstant.CartStateNormal)

	rows := make([]*models.CartDto, 0)
	if err := db.Order("update_time DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// ListByIDs 按购物车ID列表查询当前用户的购物车商品。
func (s *IApiShopCartDaoImpl) ListByIDs(c *gin.Context, userID int64, ids []int64) ([]*models.CartDto, error) {
	return s.ListByIDsAndState(c, userID, ids, shopConstant.CartStateNormal)
}

func (s *IApiShopCartDaoImpl) ListByIDsAndState(c *gin.Context, userID int64, ids []int64, state int32) ([]*models.CartDto, error) {
	if len(ids) == 0 {
		return make([]*models.CartDto, 0), nil
	}

	rows := make([]*models.CartDto, 0)
	if err := getCurrentDB(c, s.db).WithContext(c).
		Table(s.tableName).
		Where("user_id = ?", userID).
		Where("id IN ?", ids).
		Where("state = ?", state).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *IApiShopCartDaoImpl) DeleteByIds(c *gin.Context, userID int64, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	return getCurrentDB(c, s.db).WithContext(c).Table(s.tableName).
		Where("user_id = ?", userID).
		Where("id IN ?", ids).
		Delete(nil).Error
}

func (s *IApiShopCartDaoImpl) getByMergeIdentityTx(c *gin.Context, tx *gorm.DB, userID int64, skuID uint64, productType int32, activityID int64, state int32) (*models.CartDto, error) {
	var item models.CartDto
	if err := tx.WithContext(c).Table(s.tableName).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ?", userID).
		Where("sku_id = ?", skuID).
		Where("product_type = ?", productType).
		Where("activity_id = ?", activityID).
		Where("state = ?", state).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// isDuplicateKeyError 判断是否为唯一键冲突（主要用于并发下 insert 竞争场景）。
func isDuplicateKeyError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return mysqlErr.Number == 1062
	}
	return strings.Contains(err.Error(), "Duplicate entry")
}

func (s *IApiShopCartDaoImpl) getByIDTx(c *gin.Context, tx *gorm.DB, id int64) (*models.CartDto, error) {
	return s.getByIDAndStateTx(c, tx, id, shopConstant.CartStateNormal)
}

func (s *IApiShopCartDaoImpl) getByIDAndStateTx(c *gin.Context, tx *gorm.DB, id int64, state int32) (*models.CartDto, error) {
	var item models.CartDto
	if err := tx.WithContext(c).Table(s.tableName).
		Where("id = ?", id).
		Where("state = ?", state).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// Remove 软删除商城用户购物车项。
func (s *IApiShopCartDaoImpl) Remove(c *gin.Context, ids []string) error {
	return s.db.WithContext(c).Table(s.tableName).
		Where("id IN ?", ids).
		Delete(&models.CartDto{}).Error
}
