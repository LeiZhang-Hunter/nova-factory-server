package impl

import (
	"errors"
	"nova-factory-server/app/business/shop/user/dao"
	"nova-factory-server/app/business/shop/user/models"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ShopCartDaoImpl 提供商城用户购物车表的数据访问能力。
type ShopCartDaoImpl struct {
	db        *gorm.DB
	tableName string
}

// NewShopCartDao 创建商城用户购物车 DAO。
func NewShopCartDao(ms *gorm.DB) dao.IShopCartDao {
	return &ShopCartDaoImpl{
		db:        ms,
		tableName: "shop_user_cart",
	}
}

// Set 新增或修改商城用户购物车项。
func (s *ShopCartDaoImpl) Set(c *gin.Context, req *models.CartSetReq) (*models.Cart, error) {
	deptID := baizeContext.GetDeptId(c)
	userID := baizeContext.GetUserId(c)
	selected := int32(1)
	if req.Selected != nil {
		selected = *req.Selected
	}
	status := int32(1)
	if req.Status != nil {
		status = *req.Status
	}

	var result *models.Cart
	err := s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		if req.ID > 0 {
			existing, err := s.getByIDTx(c, tx, req.ID)
			if err != nil {
				return err
			}
			if existing == nil {
				return errors.New("购物车记录不存在")
			}
			if err := tx.Table(s.tableName).
				Where("id = ?", req.ID).
				Where("dept_id = ?", deptID).
				Where("state = ?", commonStatus.NORMAL).
				Updates(map[string]interface{}{
					"user_id":      req.UserID,
					"goods_id":     req.GoodsID,
					"sku_id":       req.SkuID,
					"goods_name":   req.GoodsName,
					"sku_name":     req.SkuName,
					"image_url":    req.ImageURL,
					"retail_price": req.RetailPrice,
					"quantity":     req.Quantity,
					"selected":     selected,
					"status":       status,
					"update_by":    userID,
					"update_time":  gorm.Expr("NOW()"),
				}).Error; err != nil {
				return err
			}
			result, err = s.getByIDTx(c, tx, req.ID)
			return err
		}

		existing, err := s.getByUserIDAndSkuIDTx(c, tx, req.UserID, req.SkuID)
		if err != nil {
			return err
		}
		if existing != nil {
			if err := tx.Table(s.tableName).
				Where("id = ?", existing.ID).
				Where("dept_id = ?", deptID).
				Where("state = ?", commonStatus.NORMAL).
				Updates(map[string]interface{}{
					"goods_id":     req.GoodsID,
					"goods_name":   req.GoodsName,
					"sku_name":     req.SkuName,
					"image_url":    req.ImageURL,
					"retail_price": req.RetailPrice,
					"quantity":     existing.Quantity + req.Quantity,
					"selected":     selected,
					"status":       status,
					"update_by":    userID,
					"update_time":  gorm.Expr("NOW()"),
				}).Error; err != nil {
				return err
			}
			result, err = s.getByIDTx(c, tx, existing.ID)
			return err
		}

		model := &models.Cart{
			ID:          snowflake.GenID(),
			UserID:      req.UserID,
			GoodsID:     req.GoodsID,
			SkuID:       req.SkuID,
			GoodsName:   req.GoodsName,
			SkuName:     req.SkuName,
			ImageURL:    req.ImageURL,
			RetailPrice: req.RetailPrice,
			Quantity:    req.Quantity,
			Selected:    selected,
			Status:      status,
			DeptID:      deptID,
			State:       commonStatus.NORMAL,
		}
		model.SetCreateBy(userID)
		if err := tx.Table(s.tableName).Create(model).Error; err != nil {
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

// GetByID 根据主键查询商城用户购物车项。
func (s *ShopCartDaoImpl) GetByID(c *gin.Context, id int64) (*models.Cart, error) {
	var item models.Cart
	if err := s.db.WithContext(c).Table(s.tableName).
		Where("id = ?", id).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// GetByUserIDAndSkuID 根据用户和SKU查询购物车项。
func (s *ShopCartDaoImpl) GetByUserIDAndSkuID(c *gin.Context, userID int64, skuID string) (*models.Cart, error) {
	return s.getByUserIDAndSkuIDTx(c, s.db, userID, skuID)
}

// List 查询商城用户购物车列表。
func (s *ShopCartDaoImpl) List(c *gin.Context, req *models.CartQuery) (*models.CartListData, error) {
	db := s.db.WithContext(c).Table(s.tableName).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL)
	if req.UserID > 0 {
		db = db.Where("user_id = ?", req.UserID)
	}
	if req.GoodsID != "" {
		db = db.Where("goods_id = ?", req.GoodsID)
	}
	if req.SkuID != "" {
		db = db.Where("sku_id = ?", req.SkuID)
	}
	if req.Selected != nil {
		db = db.Where("selected = ?", *req.Selected)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 20
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]*models.Cart, 0)
	if err := db.Offset(int((req.Page - 1) * req.Size)).
		Limit(int(req.Size)).
		Order("selected DESC, id DESC").
		Find(&rows).Error; err != nil {
		return nil, err
	}
	return &models.CartListData{
		Rows:  rows,
		Total: total,
	}, nil
}

// Remove 软删除商城用户购物车项。
func (s *ShopCartDaoImpl) Remove(c *gin.Context, ids []int64) error {
	return s.db.WithContext(c).Table(s.tableName).
		Where("id IN ?", ids).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL).
		Updates(map[string]interface{}{
			"state":       commonStatus.DELETE,
			"update_by":   baizeContext.GetUserId(c),
			"update_time": gorm.Expr("NOW()"),
		}).Error
}

func (s *ShopCartDaoImpl) getByIDTx(c *gin.Context, tx *gorm.DB, id int64) (*models.Cart, error) {
	var item models.Cart
	if err := tx.WithContext(c).Table(s.tableName).
		Where("id = ?", id).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func (s *ShopCartDaoImpl) getByUserIDAndSkuIDTx(c *gin.Context, tx *gorm.DB, userID int64, skuID string) (*models.Cart, error) {
	var item models.Cart
	if err := tx.WithContext(c).Table(s.tableName).
		Where("user_id = ?", userID).
		Where("sku_id = ?", skuID).
		Where("dept_id = ?", baizeContext.GetDeptId(c)).
		Where("state = ?", commonStatus.NORMAL).
		First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}
