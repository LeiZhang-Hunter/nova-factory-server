package shopDaoImpl

import (
	"nova-factory-server/app/business/shop/shopDao"
	"nova-factory-server/app/business/shop/shopModels"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ShopGoodsDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewShopGoodsDao(ms *gorm.DB) shopDao.IShopGoodsDao {
	return &ShopGoodsDaoImpl{
		db:        ms,
		tableName: "shop_goods",
	}
}

func (s *ShopGoodsDaoImpl) Create(c *gin.Context, req *shopModels.GoodsUpsert) (*shopModels.Goods, error) {
	model := &shopModels.Goods{
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
	if err := s.db.WithContext(c).Table(s.tableName).Create(model).Error; err != nil {
		return nil, err
	}
	return model, nil
}

func (s *ShopGoodsDaoImpl) Update(c *gin.Context, req *shopModels.GoodsUpsert) (*shopModels.Goods, error) {
	updates := map[string]interface{}{
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
	if err := s.db.WithContext(c).Table(s.tableName).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
		return nil, err
	}
	return s.GetByID(c, int64(req.ID))
}

func (s *ShopGoodsDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return s.db.WithContext(c).Table(s.tableName).Where("id IN ?", ids).Delete(nil).Error
}

func (s *ShopGoodsDaoImpl) GetByID(c *gin.Context, id int64) (*shopModels.Goods, error) {
	var item shopModels.Goods
	if err := s.db.WithContext(c).Table(s.tableName).Where("id = ?", id).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (s *ShopGoodsDaoImpl) List(c *gin.Context, req *shopModels.GoodsQuery) (*shopModels.GoodsListData, error) {
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
	if req.IsOnSale == 0 || req.IsOnSale == 1 {
		db = db.Where("is_on_sale = ?", req.IsOnSale)
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
	rows := make([]*shopModels.Goods, 0)
	if err := db.Offset(int((req.Page - 1) * req.Size)).Limit(int(req.Size)).Order("id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	return &shopModels.GoodsListData{
		Rows:  rows,
		Total: total,
	}, nil
}
