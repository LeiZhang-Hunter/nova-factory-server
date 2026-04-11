package shopdaoimpl

import (
	"errors"
	"nova-factory-server/app/business/shop/product/shopdao"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/constant/commonStatus"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ShopGoodsDaoImpl struct {
	db        *gorm.DB
	tableName string
}

// NewShopGoodsDao 创建商品数据访问对象
func NewShopGoodsDao(ms *gorm.DB) shopdao.IShopGoodsDao {
	return &ShopGoodsDaoImpl{
		db:        ms,
		tableName: "shop_goods",
	}
}

func (s *ShopGoodsDaoImpl) Create(c *gin.Context, req *shopmodels.GoodsUpsert) (*shopmodels.Goods, error) {
	model := buildGoodsModel(req)
	if err := s.db.WithContext(c).Table(s.tableName).Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

// BatchCreate 批量新增商品
func (s *ShopGoodsDaoImpl) BatchCreate(c *gin.Context, reqs []*shopmodels.GoodsUpsert, batchSize int) error {
	if len(reqs) == 0 {
		return nil
	}
	if batchSize <= 0 {
		batchSize = len(reqs)
	}
	models := make([]*shopmodels.Goods, 0, len(reqs))
	now := time.Now()
	for _, req := range reqs {
		model := buildGoodsModel(req)
		model.CreateBy = 0
		model.UpdateBy = 0
		model.CreateTime = &now
		model.UpdateTime = &now
		models = append(models, model)
	}
	return s.db.WithContext(c).Table(s.tableName).CreateInBatches(models, batchSize).Error
}

func (s *ShopGoodsDaoImpl) Update(c *gin.Context, req *shopmodels.GoodsUpsert) (*shopmodels.Goods, error) {
	updates := buildGoodsUpdates(req)
	if err := s.db.WithContext(c).Table(s.tableName).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
		return nil, err
	}
	return s.GetByID(c, int64(req.ID))
}

// BatchUpdate 批量更新商品
func (s *ShopGoodsDaoImpl) BatchUpdate(c *gin.Context, reqs []*shopmodels.GoodsUpsert, batchSize int) error {
	if len(reqs) == 0 {
		return nil
	}
	if batchSize <= 0 {
		batchSize = len(reqs)
	}
	return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		for start := 0; start < len(reqs); start += batchSize {
			end := start + batchSize
			if end > len(reqs) {
				end = len(reqs)
			}
			for _, req := range reqs[start:end] {
				if err := tx.Table(s.tableName).Where("id = ?", req.ID).Updates(buildGoodsUpdates(req)).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (s *ShopGoodsDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return s.db.WithContext(c).Table(s.tableName).Where("id IN ?", ids).Delete(nil).Error
}

func (s *ShopGoodsDaoImpl) GetByID(c *gin.Context, id int64) (*shopmodels.Goods, error) {
	var item shopmodels.Goods
	if err := s.db.WithContext(c).Table(s.tableName).Where("id = ?", id).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// GetByGoodsID 根据商品业务ID查询商品
func (s *ShopGoodsDaoImpl) GetByGoodsID(c *gin.Context, goodsID string) (*shopmodels.Goods, error) {
	var item shopmodels.Goods
	if err := s.db.WithContext(c).Table(s.tableName).Where("goods_id = ?", goodsID).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// ListByGoodsIDs 根据商品业务ID批量查询商品
func (s *ShopGoodsDaoImpl) ListByGoodsIDs(c *gin.Context, goodsIDs []string) ([]*shopmodels.Goods, error) {
	if len(goodsIDs) == 0 {
		return []*shopmodels.Goods{}, nil
	}
	rows := make([]*shopmodels.Goods, 0)
	if err := s.db.WithContext(c).Table(s.tableName).Where("goods_id IN ?", goodsIDs).Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *ShopGoodsDaoImpl) List(c *gin.Context, req *shopmodels.GoodsQuery) (*shopmodels.GoodsListData, error) {
	db := s.db.WithContext(c).Table(s.tableName)
	if req.GoodsName != "" {
		db = db.Where("goods_name LIKE ?", "%"+req.GoodsName+"%")
	}
	if req.GoodsCode != "" {
		db = db.Where("goods_code = ?", req.GoodsCode)
	}
	if req.OuterID != "" {
		db = db.Where("outer_id = ?", req.OuterID)
	}
	if req.IsOnSale != nil {
		db = db.Where("is_on_sale = ?", req.IsOnSale)
	}

	db = db.Where("state = ?", commonStatus.NORMAL)
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
	rows := make([]*shopmodels.Goods, 0)
	if err := db.Offset(int((req.Page - 1) * req.Size)).Limit(int(req.Size)).Order("id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return &shopmodels.GoodsListData{
		Rows:  rows,
		Total: total,
	}, nil
}

func buildGoodsModel(req *shopmodels.GoodsUpsert) *shopmodels.Goods {
	return &shopmodels.Goods{
		GoodsID:       req.GoodsID,
		GoodsName:     req.GoodsName,
		GoodsCode:     req.GoodsCode,
		OuterID:       req.OuterID,
		ImageURL:      req.ImageURL,
		RetailPrice:   req.RetailPrice,
		GalleryImages: req.GalleryImages,
		VideoURL:      req.VideoURL,
		Description:   req.Description,
		Weight:        req.Weight,
		WeightUnit:    req.WeightUnit,
		Unit:          req.Unit,
		IsOnSale:      req.IsOnSale,
		Quantity:      req.Quantity,
	}
}

func buildGoodsUpdates(req *shopmodels.GoodsUpsert) map[string]interface{} {
	return map[string]interface{}{
		"goods_id":       req.GoodsID,
		"goods_name":     req.GoodsName,
		"goods_code":     req.GoodsCode,
		"outer_id":       req.OuterID,
		"image_url":      req.ImageURL,
		"retail_price":   req.RetailPrice,
		"gallery_images": req.GalleryImages,
		"video_url":      req.VideoURL,
		"description":    req.Description,
		"weight":         req.Weight,
		"weight_unit":    req.WeightUnit,
		"unit":           req.Unit,
		"is_on_sale":     req.IsOnSale,
		"quantity":       req.Quantity,
	}
}
