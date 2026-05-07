package impl

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"nova-factory-server/app/business/shop/api/dao"
	"nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/snowflake"
	"strings"
	"time"
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
	err := s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {

		existing, err := s.getByUserIDAndSkuIDTx(c, tx, req.UserID, req.SkuID)
		if err != nil {
			return err
		}
		if existing != nil {
			if err := tx.Table(s.tableName).
				Where("id = ?", existing.ID).
				Where("state = ?", commonStatus.NORMAL).
				Updates(map[string]interface{}{
					"goods_id":     req.GoodsID,
					"goods_name":   req.GoodsName,
					"sku_name":     req.SkuName,
					"image_url":    req.ImageURL,
					"retail_price": req.RetailPrice,
					"quantity":     gorm.Expr("quantity + ?", req.Quantity),
					"update_time":  gorm.Expr("NOW()"),
				}).Error; err != nil {
				return err
			}
			result, err = s.getByIDTx(c, tx, existing.ID)
			return err
		}

		model := &models.CartDto{
			ID:          snowflake.GenID(),
			UserID:      req.UserID,
			GoodsID:     req.GoodsID,
			SkuID:       req.SkuID,
			GoodsName:   req.GoodsName,
			SkuName:     req.SkuName,
			ImageURL:    req.ImageURL,
			RetailPrice: req.RetailPrice,
			Quantity:    req.Quantity,
			State:       commonStatus.NORMAL,
			CreateTime:  &now,
		}
		if err := tx.Table(s.tableName).Create(model).Error; err != nil {
			if isDuplicateKeyError(err) {
				existing, getErr := s.getByUserIDAndSkuIDTx(c, tx, req.UserID, req.SkuID)
				if getErr != nil {
					return getErr
				}
				if existing == nil {
					return err
				}
				if updateErr := tx.Table(s.tableName).
					Where("id = ?", existing.ID).
					Where("state = ?", commonStatus.NORMAL).
					Updates(map[string]interface{}{
						"goods_id":     req.GoodsID,
						"goods_name":   req.GoodsName,
						"sku_name":     req.SkuName,
						"image_url":    req.ImageURL,
						"retail_price": req.RetailPrice,
						"quantity":     gorm.Expr("quantity + ?", req.Quantity),
						"update_time":  gorm.Expr("NOW()"),
					}).Error; updateErr != nil {
					return updateErr
				}
				result, err = s.getByIDTx(c, tx, existing.ID)
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

func (s *IApiShopCartDaoImpl) getByUserIDAndSkuIDTx(c *gin.Context, tx *gorm.DB, userID int64, skuID string) (*models.CartDto, error) {
	var item models.CartDto
	if err := tx.WithContext(c).Table(s.tableName).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("user_id = ?", userID).
		Where("sku_id = ?", skuID).
		Where("state = ?", commonStatus.NORMAL).
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
	var item models.CartDto
	if err := tx.WithContext(c).Table(s.tableName).
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
