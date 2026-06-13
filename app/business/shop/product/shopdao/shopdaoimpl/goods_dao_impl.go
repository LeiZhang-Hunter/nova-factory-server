package shopdaoimpl

import (
	"encoding/json"
	"errors"
	"go.uber.org/zap"
	homeDao "nova-factory-server/app/business/shop/home/dao"
	homeModels "nova-factory-server/app/business/shop/home/models"
	"nova-factory-server/app/business/shop/product/shopdao"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/fileUtils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ShopGoodsDaoImpl struct {
	db        *gorm.DB
	itemDao   homeDao.IShopHomeModuleItemDao
	tableName string
}

// NewShopGoodsDao 创建商品数据访问对象
const goodsBusinessType = "goods"

func NewShopGoodsDao(ms *gorm.DB, itemDao homeDao.IShopHomeModuleItemDao) shopdao.IShopGoodsDao {
	return &ShopGoodsDaoImpl{
		db:        ms,
		itemDao:   itemDao,
		tableName: "shop_goods",
	}
}

func (s *ShopGoodsDaoImpl) Transaction(c *gin.Context, fn func(txDao shopdao.IShopGoodsDao) error) error {
	return s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		txDao := &ShopGoodsDaoImpl{
			db:        tx,
			itemDao:   s.itemDao,
			tableName: s.tableName,
		}
		return fn(txDao)
	})
}

func (s *ShopGoodsDaoImpl) Create(c *gin.Context, req *shopmodels.GoodsUpsert) (*shopmodels.Goods, error) {
	model := buildGoodsModel(req)
	if err := s.db.WithContext(c).Table(s.tableName).Create(model).Error; err != nil {
		return nil, err
	}
	if err := s.itemDao.SyncBusinessModules(c, &homeModels.HomeModuleItemBusinessSync{
		BusinessType: goodsBusinessType,
		LinkID:       model.ID,
		ModuleIDs:    req.HomeModuleIDs,
		ItemName:     req.GoodsName,
		ItemSubTitle: req.Description,
		ItemImage:    req.ImageURL,
		Sort:         0,
		Status:       int8(req.IsOnSale),
	}); err != nil {
		return nil, err
	}
	return s.GetByID(c, model.ID)
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
	updates := buildGoodsUpdates(c, req)
	if err := s.db.WithContext(c).Table(s.tableName).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
		return nil, err
	}
	if err := s.itemDao.SyncBusinessModules(c, &homeModels.HomeModuleItemBusinessSync{
		BusinessType: goodsBusinessType,
		LinkID:       req.ID,
		ModuleIDs:    req.HomeModuleIDs,
		ItemName:     req.GoodsName,
		ItemSubTitle: req.Description,
		ItemImage:    req.ImageURL,
		Sort:         0,
		Status:       int8(req.IsOnSale),
	}); err != nil {
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
				if err := tx.Table(s.tableName).Where("id = ?", req.ID).Updates(buildGoodsUpdates(c, req)).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (s *ShopGoodsDaoImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if err := s.itemDao.DeleteByBusiness(c, goodsBusinessType, ids); err != nil {
		return err
	}
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
	if err := s.attachHomeModuleIDs(c, []*shopmodels.Goods{&item}); err != nil {
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
	if err := s.attachHomeModuleIDs(c, []*shopmodels.Goods{&item}); err != nil {
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
	if err := s.attachHomeModuleIDs(c, rows); err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *ShopGoodsDaoImpl) List(c *gin.Context, req *shopmodels.GoodsQuery) (*shopmodels.GoodsListData, error) {
	db := s.db.WithContext(c).Table(s.tableName)
	if req.ID > 0 {
		db = db.Where("id = ?", req.ID)
	}
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
	if req.CategoryId > 0 {
		db = db.Where("shop_category_id = ?", req.CategoryId)
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

	// 构建排序
	orderClause := "id DESC"
	if req.SortBy == "retailPrice" && req.SortOrder != "" {
		orderClause = "retail_price " + req.SortOrder
	}
	if err := db.Offset(int((req.Page - 1) * req.Size)).Limit(int(req.Size)).Order(orderClause).Find(&rows).Error; err != nil {
		return nil, err
	}
	if err := s.attachHomeModuleIDs(c, rows); err != nil {
		return nil, err
	}
	return &shopmodels.GoodsListData{
		Rows:  rows,
		Total: total,
	}, nil
}

func buildGoodsModel(req *shopmodels.GoodsUpsert) *shopmodels.Goods {
	content, err := json.Marshal(req.GalleryImages)
	if err != nil {
		zap.L().Error("buildGoodsModel json.Marshal", zap.Error(err))
	}
	return &shopmodels.Goods{
		GoodsID:       req.GoodsID,
		GoodsName:     req.GoodsName,
		GoodsCode:     req.GoodsCode,
		OuterID:       req.OuterID,
		ImageURL:      req.ImageURL,
		RetailPrice:   req.RetailPrice,
		GalleryImages: string(content),
		VideoURL:      req.VideoURL,
		Description:   req.Description,
		Weight:        req.Weight,
		WeightUnit:    req.WeightUnit,
		Unit:          req.Unit,
		IsOnSale:      req.IsOnSale,
		Quantity:      req.Quantity,
		HomeModuleIDs: strings.Join(req.HomeModuleIDs, ","),
	}
}

func buildGoodsUpdates(c *gin.Context, req *shopmodels.GoodsUpsert) map[string]interface{} {
	var galleryImagesArray []string
	for _, v := range req.GalleryImages {
		path, err := fileUtils.NormalizeResourcePath(v)
		if err != nil {
			zap.L().Error("normalizeResourcePath error", zap.Error(err))
			continue
		}
		galleryImagesArray = append(galleryImagesArray, path)
	}

	var galleryImagesStr string
	if len(galleryImagesArray) > 0 {
		galleryImagesStr = strings.Join(galleryImagesArray, ",")
	}

	imageURL, err := fileUtils.NormalizeResourcePath(req.ImageURL)
	if err != nil {
		zap.L().Error("normalizeResourcePath error", zap.Error(err))
	}

	videoUrl, err := fileUtils.NormalizeResourcePath(req.ImageURL)
	if err != nil {
		zap.L().Error("normalizeResourcePath error", zap.Error(err))
	}

	return map[string]interface{}{
		"goods_id":         req.GoodsID,
		"goods_name":       req.GoodsName,
		"goods_code":       req.GoodsCode,
		"outer_id":         req.OuterID,
		"image_url":        imageURL,
		"retail_price":     req.RetailPrice,
		"gallery_images":   galleryImagesStr,
		"video_url":        videoUrl,
		"description":      req.Description,
		"weight":           req.Weight,
		"weight_unit":      req.WeightUnit,
		"unit":             req.Unit,
		"is_on_sale":       req.IsOnSale,
		"quantity":         req.Quantity,
		"shop_category_id": req.ShopCategoryId,
		"home_module_ids":  strings.Join(req.HomeModuleIDs, ","),
	}
}

func (s *ShopGoodsDaoImpl) UpdateStockByGoodsID(c *gin.Context, goodsID string, quantity int64) error {
	return s.db.WithContext(c).Table(s.tableName).
		Where("goods_id = ?", goodsID).
		Update("quantity", quantity).Error
}

func (s *ShopGoodsDaoImpl) UpsertByGoodsID(c *gin.Context, goodsID string, updates map[string]any) error {
	var count int64
	s.db.WithContext(c).Table(s.tableName).Where("goods_id = ?", goodsID).Count(&count)
	if count == 0 {
		updates["goods_id"] = goodsID
		return s.db.WithContext(c).Table(s.tableName).Create(updates).Error
	}
	return s.db.WithContext(c).Table(s.tableName).Where("goods_id = ?", goodsID).Updates(updates).Error
}

func (s *ShopGoodsDaoImpl) attachHomeModuleIDs(c *gin.Context, rows []*shopmodels.Goods) error {
	linkIDs := make([]int64, 0, len(rows))
	for _, row := range rows {
		if row == nil {
			continue
		}
		linkIDs = append(linkIDs, row.ID)
	}
	moduleMap, err := s.itemDao.ListModuleIDsByBusiness(c, goodsBusinessType, linkIDs)
	if err != nil {
		return err
	}
	for _, row := range rows {
		if row == nil {
			continue
		}
		row.HomeModuleIDs = moduleMap[row.ID]
	}
	return nil
}
