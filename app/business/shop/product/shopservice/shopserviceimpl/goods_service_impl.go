package shopserviceimpl

import (
	"errors"
	"nova-factory-server/app/baize"
	"nova-factory-server/app/utils/snowflake"
	"strings"
	"time"

	"nova-factory-server/app/business/shop/product/shopdao"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/business/shop/product/shopservice"

	"github.com/gin-gonic/gin"
)

type ShopGoodsServiceImpl struct {
	dao    shopdao.IShopGoodsDao
	skuDao shopdao.IShopSkuDao
}

const importBatchSize = 100

// NewShopGoodsService 创建商品服务
func NewShopGoodsService(dao shopdao.IShopGoodsDao, skuDao shopdao.IShopSkuDao) shopservice.IShopGoodsService {
	return &ShopGoodsServiceImpl{
		dao:    dao,
		skuDao: skuDao,
	}
}

func (s *ShopGoodsServiceImpl) Create(c *gin.Context, req *shopmodels.GoodsUpsert) (*shopmodels.Goods, error) {
	return s.dao.Create(c, req)
}

func (s *ShopGoodsServiceImpl) Update(c *gin.Context, req *shopmodels.GoodsUpsert) (*shopmodels.Goods, error) {
	return s.dao.Update(c, req)
}

func (s *ShopGoodsServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	return s.dao.DeleteByIDs(c, ids)
}

func (s *ShopGoodsServiceImpl) GetByID(c *gin.Context, id int64) (*shopmodels.Goods, error) {
	data, err := s.dao.GetByID(c, id)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	if err = s.attachSkus(c, []*shopmodels.Goods{data}); err != nil {
		return nil, err
	}
	return data, nil
}

func (s *ShopGoodsServiceImpl) List(c *gin.Context, req *shopmodels.GoodsQuery) (*shopmodels.GoodsListData, error) {
	data, err := s.dao.List(c, req)
	if err != nil {
		return nil, err
	}
	if data == nil || len(data.Rows) == 0 {
		return data, nil
	}
	if err = s.attachSkus(c, data.Rows); err != nil {
		return nil, err
	}
	return data, nil
}

// Import 增量导入商品与规格数据
func (s *ShopGoodsServiceImpl) Import(c *gin.Context, records []shopmodels.ImportGoodsRecord) error {
	for _, batch := range splitImportRecords(records, importBatchSize) {
		goodsMap, skuMap, err := buildImportBatch(batch)
		if err != nil {
			return err
		}

		goodsCreates, goodsUpdates, err := s.diffGoods(c, goodsMap)
		if err != nil {
			return err
		}
		if err = s.dao.BatchCreate(c, goodsCreates, importBatchSize); err != nil {
			return err
		}
		if err = s.dao.BatchUpdate(c, goodsUpdates, importBatchSize); err != nil {
			return err
		}

		skuCreates, skuUpdates, err := s.diffSkus(c, skuMap)
		if err != nil {
			return err
		}
		if err = s.skuDao.BatchCreate(c, skuCreates, importBatchSize); err != nil {
			return err
		}
		if err = s.skuDao.BatchUpdate(c, skuUpdates, importBatchSize); err != nil {
			return err
		}
	}
	return nil
}

// buildGoodsUpsert 组装商品导入参数
func buildGoodsUpsert(record shopmodels.ImportGoodsRecord) (*shopmodels.GoodsUpsert, error) {
	goodsID := strings.TrimSpace(record.ExternalID)
	if goodsID == "" {
		goodsID = strings.TrimSpace(record.Data.ProductCode)
	}
	if goodsID == "" {
		return nil, errors.New("导入商品缺少external_id或product_code")
	}
	quantity := int64(0)
	retailPrice := 0.0
	weight := 0.0
	for _, sku := range record.Data.Skus {
		if sku.Size > 0 {
			quantity += int64(sku.Size)
		}
		if retailPrice <= 0 {
			retailPrice = pickRetailPrice(sku)
		}
		if weight <= 0 && sku.Weight > 0 {
			weight = sku.Weight
		}
	}
	now := time.Now()
	return &shopmodels.GoodsUpsert{
		GoodsID:       goodsID,
		GoodsName:     strings.TrimSpace(record.Data.ProductName),
		GoodsCode:     strings.TrimSpace(record.Data.ProductCode),
		OuterID:       strings.TrimSpace(record.ExternalID),
		Description:   strings.TrimSpace(record.Data.Remark),
		Weight:        weight,
		WeightUnit:    "kg",
		Unit:          defaultString(strings.TrimSpace(record.Data.UnitName), "件"),
		IsOnSale:      1,
		Quantity:      quantity,
		RetailPrice:   retailPrice,
		GalleryImages: "[]",
		BaseEntity: baize.BaseEntity{
			CreateBy:   1,
			UpdateBy:   1,
			CreateTime: &now,
			UpdateTime: &now,
		},
	}, nil
}

// buildGoodsSkuUpsert 组装商品规格导入参数
func buildGoodsSkuUpsert(goodsID string, record shopmodels.ImportGoodsRecord, sku shopmodels.ImportGoodsSkuRawData) (*shopmodels.GoodsSkuUpsert, bool) {
	skuID := strings.TrimSpace(sku.Skuid)
	if skuID == "" {
		skuID = strings.TrimSpace(sku.Skucode)
	}
	if skuID == "" {
		skuID = strings.TrimSpace(sku.Barcode)
	}
	if skuID == "" {
		return nil, false
	}
	now := time.Now()
	return &shopmodels.GoodsSkuUpsert{
		GoodsID:       goodsID,
		SkuID:         skuID,
		SkuName:       strings.TrimSpace(sku.Skuname),
		SkuCode:       strings.TrimSpace(sku.Skucode),
		OuterID:       strings.TrimSpace(sku.Lcmccode),
		Barcode:       strings.TrimSpace(sku.Barcode),
		Description:   strings.TrimSpace(record.Data.Remark),
		Weight:        sku.Weight,
		WeightUnit:    "kg",
		Unit:          defaultString(strings.TrimSpace(record.Data.UnitName), "件"),
		Quantity:      int64(sku.Size),
		RetailPrice:   pickRetailPrice(sku),
		GalleryImages: "[]",
		BaseEntity: baize.BaseEntity{
			CreateBy:   1,
			UpdateBy:   1,
			CreateTime: &now,
			UpdateTime: &now,
		},
	}, true
}

// pickRetailPrice 依次选择导入数据中的有效零售价
func pickRetailPrice(sku shopmodels.ImportGoodsSkuRawData) float64 {
	prices := []float64{sku.Price, sku.Price2, sku.Price3, sku.Price4, sku.Price5}
	for _, price := range prices {
		if price > 0 {
			return price
		}
	}
	return 0
}

// defaultString 返回非空字符串，空值时使用默认值
func defaultString(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

// splitImportRecords 按批次拆分导入记录
func splitImportRecords(records []shopmodels.ImportGoodsRecord, batchSize int) [][]shopmodels.ImportGoodsRecord {
	if len(records) == 0 {
		return nil
	}
	if batchSize <= 0 {
		batchSize = len(records)
	}
	result := make([][]shopmodels.ImportGoodsRecord, 0, (len(records)+batchSize-1)/batchSize)
	for start := 0; start < len(records); start += batchSize {
		end := start + batchSize
		if end > len(records) {
			end = len(records)
		}
		result = append(result, records[start:end])
	}
	return result
}

// buildImportBatch 构建单批次导入的商品与规格数据
func buildImportBatch(records []shopmodels.ImportGoodsRecord) (map[string]*shopmodels.GoodsUpsert, map[string]*shopmodels.GoodsSkuUpsert, error) {
	goodsMap := make(map[string]*shopmodels.GoodsUpsert, len(records))
	skuMap := make(map[string]*shopmodels.GoodsSkuUpsert)
	for _, record := range records {
		goodsUpsert, err := buildGoodsUpsert(record)
		if err != nil {
			return nil, nil, err
		}
		goodsMap[goodsUpsert.GoodsID] = goodsUpsert
		for _, sku := range record.Data.Skus {
			skuUpsert, ok := buildGoodsSkuUpsert(goodsUpsert.GoodsID, record, sku)
			if !ok {
				continue
			}
			skuMap[skuUpsert.SkuID] = skuUpsert
		}
	}
	return goodsMap, skuMap, nil
}

// diffGoods 比对单批次商品数据的新增和更新列表
func (s *ShopGoodsServiceImpl) diffGoods(c *gin.Context, goodsMap map[string]*shopmodels.GoodsUpsert) ([]*shopmodels.GoodsUpsert, []*shopmodels.GoodsUpsert, error) {
	if len(goodsMap) == 0 {
		return nil, nil, nil
	}
	existingRows, err := s.dao.ListByGoodsIDs(c, mapKeys(goodsMap))
	if err != nil {
		return nil, nil, err
	}
	existingMap := make(map[string]*shopmodels.Goods, len(existingRows))
	for _, row := range existingRows {
		existingMap[row.GoodsID] = row
	}
	creates := make([]*shopmodels.GoodsUpsert, 0)
	updates := make([]*shopmodels.GoodsUpsert, 0)
	for goodsID, req := range goodsMap {
		current := existingMap[goodsID]
		if current == nil {
			req.ID = snowflake.GenID()
			creates = append(creates, req)
			continue
		}
		if !sameGoods(current, req) {
			req.ID = current.ID
			updates = append(updates, req)
		}
	}
	return creates, updates, nil
}

// diffSkus 比对单批次规格数据的新增和更新列表
func (s *ShopGoodsServiceImpl) diffSkus(c *gin.Context, skuMap map[string]*shopmodels.GoodsSkuUpsert) ([]*shopmodels.GoodsSkuUpsert, []*shopmodels.GoodsSkuUpsert, error) {
	if len(skuMap) == 0 {
		return nil, nil, nil
	}
	existingRows, err := s.skuDao.ListBySkuIDs(c, mapKeys(skuMap))
	if err != nil {
		return nil, nil, err
	}
	existingMap := make(map[string]*shopmodels.GoodsSku, len(existingRows))
	for _, row := range existingRows {
		existingMap[row.SkuID] = row
	}
	creates := make([]*shopmodels.GoodsSkuUpsert, 0)
	updates := make([]*shopmodels.GoodsSkuUpsert, 0)
	for skuID, req := range skuMap {
		current := existingMap[skuID]
		if current == nil {
			creates = append(creates, req)
			continue
		}
		if !sameSku(current, req) {
			req.ID = current.ID
			updates = append(updates, req)
		}
	}
	return creates, updates, nil
}

// sameGoods 判断商品数据是否发生变化
func sameGoods(current *shopmodels.Goods, req *shopmodels.GoodsUpsert) bool {
	if current == nil || req == nil {
		return false
	}
	return current.GoodsID == req.GoodsID &&
		current.GoodsName == req.GoodsName &&
		current.GoodsCode == req.GoodsCode &&
		current.OuterID == req.OuterID &&
		current.ImageURL == req.ImageURL &&
		current.RetailPrice == req.RetailPrice &&
		current.GalleryImages == req.GalleryImages &&
		current.VideoURL == req.VideoURL &&
		current.Description == req.Description &&
		current.Weight == req.Weight &&
		current.WeightUnit == req.WeightUnit &&
		current.Unit == req.Unit &&
		current.IsOnSale == req.IsOnSale &&
		current.Quantity == req.Quantity
}

// sameSku 判断商品规格数据是否发生变化
func sameSku(current *shopmodels.GoodsSku, req *shopmodels.GoodsSkuUpsert) bool {
	if current == nil || req == nil {
		return false
	}
	return current.GoodsID == req.GoodsID &&
		current.SkuID == req.SkuID &&
		current.SkuName == req.SkuName &&
		current.SkuCode == req.SkuCode &&
		current.OuterID == req.OuterID &&
		current.Barcode == req.Barcode &&
		current.ImageURL == req.ImageURL &&
		current.RetailPrice == req.RetailPrice &&
		current.GalleryImages == req.GalleryImages &&
		current.VideoURL == req.VideoURL &&
		current.Description == req.Description &&
		current.Weight == req.Weight &&
		current.WeightUnit == req.WeightUnit &&
		current.Unit == req.Unit &&
		current.Quantity == req.Quantity
}

// mapKeys 获取映射中的全部键
func mapKeys[T any](data map[string]T) []string {
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	return keys
}

// goodsIDsFromRows 提取商品列表中的商品业务ID
func goodsIDsFromRows(rows []*shopmodels.Goods) []string {
	goodsIDs := make([]string, 0, len(rows))
	seen := make(map[string]struct{}, len(rows))
	for _, row := range rows {
		if row == nil || row.GoodsID == "" {
			continue
		}
		if _, ok := seen[row.GoodsID]; ok {
			continue
		}
		seen[row.GoodsID] = struct{}{}
		goodsIDs = append(goodsIDs, row.GoodsID)
	}
	return goodsIDs
}

// attachSkus 为商品挂载规格列表
func (s *ShopGoodsServiceImpl) attachSkus(c *gin.Context, rows []*shopmodels.Goods) error {
	if len(rows) == 0 {
		return nil
	}
	skus, err := s.skuDao.ListByGoodsIDs(c, goodsIDsFromRows(rows))
	if err != nil {
		return err
	}
	skuMap := make(map[string][]*shopmodels.GoodsSku, len(rows))
	for _, sku := range skus {
		skuMap[sku.GoodsID] = append(skuMap[sku.GoodsID], sku)
	}
	for _, row := range rows {
		row.Skus = skuMap[row.GoodsID]
		if row.Skus == nil {
			row.Skus = make([]*shopmodels.GoodsSku, 0)
		}
	}
	return nil
}
