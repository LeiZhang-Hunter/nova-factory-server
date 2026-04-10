package shopdaoimpl

import (
	"errors"
	"nova-factory-server/app/business/shop/shopdao"
	"nova-factory-server/app/business/shop/shopmodels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ShopSkuDaoImpl struct {
	db        *gorm.DB
	tableName string
}

// NewShopSkuDao 创建商品规格数据访问对象
func NewShopSkuDao(ms *gorm.DB) shopdao.IShopSkuDao {
	return &ShopSkuDaoImpl{
		db:        ms,
		tableName: "shop_goods_sku",
	}
}

func (s *ShopSkuDaoImpl) Create(c *gin.Context, req *shopmodels.GoodsSkuUpsert) (*shopmodels.GoodsSku, error) {
	model := buildSkuModel(req)
	if err := s.db.WithContext(c).Table(s.tableName).Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

// BatchCreate 批量新增商品规格
func (s *ShopSkuDaoImpl) BatchCreate(c *gin.Context, reqs []*shopmodels.GoodsSkuUpsert, batchSize int) error {
	if len(reqs) == 0 {
		return nil
	}
	if batchSize <= 0 {
		batchSize = len(reqs)
	}
	models := make([]*shopmodels.GoodsSku, 0, len(reqs))
	for _, req := range reqs {
		models = append(models, buildSkuModel(req))
	}
	return s.db.WithContext(c).Table(s.tableName).CreateInBatches(models, batchSize).Error
}

func (s *ShopSkuDaoImpl) Update(c *gin.Context, req *shopmodels.GoodsSkuUpsert) (*shopmodels.GoodsSku, error) {
	updates := buildSkuUpdates(req)
	if err := s.db.WithContext(c).Table(s.tableName).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
		return nil, err
	}
	return s.GetByID(c, int64(req.ID))
}

// BatchUpdate 批量更新商品规格
func (s *ShopSkuDaoImpl) BatchUpdate(c *gin.Context, reqs []*shopmodels.GoodsSkuUpsert, batchSize int) error {
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
				if err := tx.Table(s.tableName).Where("id = ?", req.ID).Updates(buildSkuUpdates(req)).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (s *ShopSkuDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return s.db.WithContext(c).Table(s.tableName).Where("id IN ?", ids).Delete(nil).Error
}

func (s *ShopSkuDaoImpl) GetByID(c *gin.Context, id int64) (*shopmodels.GoodsSku, error) {
	var item shopmodels.GoodsSku
	if err := s.db.WithContext(c).Table(s.tableName).Where("id = ?", id).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// GetBySkuID 根据规格业务ID查询商品规格
func (s *ShopSkuDaoImpl) GetBySkuID(c *gin.Context, skuID string) (*shopmodels.GoodsSku, error) {
	var item shopmodels.GoodsSku
	if err := s.db.WithContext(c).Table(s.tableName).Where("sku_id = ?", skuID).First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

// ListBySkuIDs 根据规格业务ID批量查询商品规格
func (s *ShopSkuDaoImpl) ListBySkuIDs(c *gin.Context, skuIDs []string) ([]*shopmodels.GoodsSku, error) {
	if len(skuIDs) == 0 {
		return []*shopmodels.GoodsSku{}, nil
	}
	rows := make([]*shopmodels.GoodsSku, 0)
	if err := s.db.WithContext(c).Table(s.tableName).Where("sku_id IN ?", skuIDs).Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *ShopSkuDaoImpl) List(c *gin.Context, req *shopmodels.GoodsSkuQuery) (*shopmodels.GoodsSkuListData, error) {
	db := s.db.WithContext(c).Table(s.tableName)
	if req.GoodsID != "" {
		db = db.Where("goods_id = ?", req.GoodsID)
	}
	if req.SkuName != "" {
		db = db.Where("sku_name LIKE ?", "%"+req.SkuName+"%")
	}
	if req.SkuCode != "" {
		db = db.Where("sku_code = ?", req.SkuCode)
	}
	if req.OuterID != "" {
		db = db.Where("outer_id = ?", req.OuterID)
	}
	if req.Barcode != "" {
		db = db.Where("barcode = ?", req.Barcode)
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
	rows := make([]*shopmodels.GoodsSku, 0)
	if err := db.Offset(int((req.Page - 1) * req.Size)).Limit(int(req.Size)).Order("id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return &shopmodels.GoodsSkuListData{
		Rows:  rows,
		Total: total,
	}, nil
}

func buildSkuModel(req *shopmodels.GoodsSkuUpsert) *shopmodels.GoodsSku {
	return &shopmodels.GoodsSku{
		GoodsID:       req.GoodsID,
		SkuID:         req.SkuID,
		SkuName:       req.SkuName,
		SkuCode:       req.SkuCode,
		OuterID:       req.OuterID,
		Barcode:       req.Barcode,
		ImageURL:      req.ImageURL,
		RetailPrice:   req.RetailPrice,
		GalleryImages: req.GalleryImages,
		VideoURL:      req.VideoURL,
		Description:   req.Description,
		Weight:        req.Weight,
		WeightUnit:    req.WeightUnit,
		Unit:          req.Unit,
		Quantity:      req.Quantity,
	}
}

func buildSkuUpdates(req *shopmodels.GoodsSkuUpsert) map[string]interface{} {
	return map[string]interface{}{
		"goods_id":       req.GoodsID,
		"sku_id":         req.SkuID,
		"sku_name":       req.SkuName,
		"sku_code":       req.SkuCode,
		"outer_id":       req.OuterID,
		"barcode":        req.Barcode,
		"image_url":      req.ImageURL,
		"retail_price":   req.RetailPrice,
		"gallery_images": req.GalleryImages,
		"video_url":      req.VideoURL,
		"description":    req.Description,
		"weight":         req.Weight,
		"weight_unit":    req.WeightUnit,
		"unit":           req.Unit,
		"quantity":       req.Quantity,
	}
}
