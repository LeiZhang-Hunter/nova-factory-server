package shopDaoImpl

import (
	"nova-factory-server/app/business/shop/shopDao"
	"nova-factory-server/app/business/shop/shopModels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ShopSkuDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewShopSkuDao(ms *gorm.DB) shopDao.IShopSkuDao {
	return &ShopSkuDaoImpl{
		db:        ms,
		tableName: "shop_goods_sku",
	}
}

func (s *ShopSkuDaoImpl) Create(c *gin.Context, req *shopModels.GoodsSkuUpsert) (*shopModels.GoodsSku, error) {
	model := &shopModels.GoodsSku{
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
	if err := s.db.WithContext(c).Table(s.tableName).Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (s *ShopSkuDaoImpl) Update(c *gin.Context, req *shopModels.GoodsSkuUpsert) (*shopModels.GoodsSku, error) {
	updates := map[string]interface{}{
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
	if err := s.db.WithContext(c).Table(s.tableName).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
		return nil, err
	}
	return s.GetByID(c, int64(req.ID))
}

func (s *ShopSkuDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return s.db.WithContext(c).Table(s.tableName).Where("id IN ?", ids).Delete(nil).Error
}

func (s *ShopSkuDaoImpl) GetByID(c *gin.Context, id int64) (*shopModels.GoodsSku, error) {
	var item shopModels.GoodsSku
	if err := s.db.WithContext(c).Table(s.tableName).Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *ShopSkuDaoImpl) List(c *gin.Context, req *shopModels.GoodsSkuQuery) (*shopModels.GoodsSkuListData, error) {
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
	rows := make([]*shopModels.GoodsSku, 0)
	if err := db.Offset(int((req.Page - 1) * req.Size)).Limit(int(req.Size)).Order("id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return &shopModels.GoodsSkuListData{
		Rows:  rows,
		Total: total,
	}, nil
}
