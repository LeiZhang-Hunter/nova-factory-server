package shopserviceimpl

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"nova-factory-server/app/baize"
	"nova-factory-server/app/utils/snowflake"
	"nova-factory-server/app/utils/vectorsearch/goods"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"

	"nova-factory-server/app/business/shop/product/shopdao"
	"nova-factory-server/app/business/shop/product/shopmodels"
	"nova-factory-server/app/business/shop/product/shopretrieval"
	"nova-factory-server/app/business/shop/product/shopservice"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/fileUtils"
	"nova-factory-server/app/utils/retrieval"

	"github.com/gin-gonic/gin"
)

type ShopGoodsServiceImpl struct {
	dao         shopdao.IShopGoodsDao
	vectorDao   shopdao.IShopGoodsVectorDao
	retriever   retrieval.Retriever
	skuDao      shopdao.IShopSkuDao
	categoryDao shopdao.IShopCategoryDao
	cache       cache.Cache
	// metadataExtractor 元数据提取器
	metadataExtractor *goods.MetadataExtractor
}

// 导入商品时的数据库分批大小，避免单次批量过大。
const importBatchSize = 100

// 导出 CSV 时单次分页拉取的行数，兼顾内存与导出吞吐。
const goodsExportBatchSize int64 = 500

// NewShopGoodsService 创建商品服务
func NewShopGoodsService(dao shopdao.IShopGoodsDao, vectorDao shopdao.IShopGoodsVectorDao,
	skuDao shopdao.IShopSkuDao, categoryDao shopdao.IShopCategoryDao, cache cache.Cache) shopservice.IShopGoodsService {
	metadataExtractor := goods.NewMetadataExtractor()
	err := metadataExtractor.Init()
	if err != nil {
		zap.L().Error("metadataExtractor init error", zap.Error(err))
	}
	return &ShopGoodsServiceImpl{
		dao:               dao,
		vectorDao:         vectorDao,
		retriever:         shopretrieval.NewGoodsVectorRetriever(vectorDao),
		skuDao:            skuDao,
		categoryDao:       categoryDao,
		cache:             cache,
		metadataExtractor: metadataExtractor,
	}
}

// Create 创建商品基础信息。
func (s *ShopGoodsServiceImpl) Create(c *gin.Context, req *shopmodels.GoodsUpsert) (*shopmodels.Goods, error) {
	return s.dao.Create(c, req)
}

// Update 更新商品基础信息。
func (s *ShopGoodsServiceImpl) Update(c *gin.Context, req *shopmodels.GoodsUpsert) (*shopmodels.Goods, error) {
	var (
		data *shopmodels.Goods
		err  error
	)
	err = s.dao.Transaction(c, func(txDao shopdao.IShopGoodsDao) error {
		data, err = txDao.Update(c, req)
		if err != nil {
			return err
		}
		if s.vectorDao == nil {
			return nil
		}
		goodsDBID := req.ID
		if data != nil && data.ID > 0 {
			goodsDBID = data.ID
		}
		if goodsDBID <= 0 {
			return nil
		}
		if err = s.vectorDao.UpdateSaleStatusByGoodsID(c, goodsDBID, req.IsOnSale); err != nil {
			return fmt.Errorf("同步商品向量在售状态失败: %w", err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return data, nil
}

// DeleteByIDs 批量删除商品。
func (s *ShopGoodsServiceImpl) DeleteByIDs(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	if s.skuDao != nil {
		goodsIDs := make([]string, 0, len(ids))
		seen := make(map[string]struct{}, len(ids))
		for _, id := range ids {
			goods, err := s.dao.GetByID(c, id)
			if err != nil {
				return err
			}
			if goods == nil || strings.TrimSpace(goods.GoodsID) == "" {
				continue
			}
			if _, ok := seen[goods.GoodsID]; ok {
				continue
			}
			seen[goods.GoodsID] = struct{}{}
			goodsIDs = append(goodsIDs, goods.GoodsID)
		}
		if len(goodsIDs) > 0 {
			skus, err := s.skuDao.ListByGoodsIDs(c, goodsIDs)
			if err != nil {
				return err
			}
			if len(skus) > 0 {
				return errors.New("商品下存在SKU，不允许删除")
			}
		}
	}
	return s.dao.DeleteByIDs(c, ids)
}

// GetByID 按数据库主键查询商品，并补齐分类名、SKU 和媒体绝对地址。
func (s *ShopGoodsServiceImpl) GetByID(c *gin.Context, id int64) (*shopmodels.Goods, error) {
	data, err := s.dao.GetByID(c, id)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	if err = s.attachCategoryNames(c, []*shopmodels.Goods{data}); err != nil {
		return nil, err
	}
	if err = s.attachSkus(c, []*shopmodels.Goods{data}); err != nil {
		return nil, err
	}
	s.normalizeGoodsMediaURLs(c, []*shopmodels.Goods{data})
	return data, nil
}

// GetByGoodsID 按业务商品 ID 查询商品，并补齐分类名、SKU 和媒体绝对地址。
func (s *ShopGoodsServiceImpl) GetByGoodsID(c *gin.Context, goodsID string) (*shopmodels.Goods, error) {
	data, err := s.dao.GetByGoodsID(c, goodsID)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, nil
	}
	if err = s.attachCategoryNames(c, []*shopmodels.Goods{data}); err != nil {
		return nil, err
	}
	if err = s.attachSkus(c, []*shopmodels.Goods{data}); err != nil {
		return nil, err
	}
	s.normalizeGoodsMediaURLs(c, []*shopmodels.Goods{data})
	return data, nil
}

// List 分页查询商品列表，并统一补齐分类名、SKU 与媒体地址。
func (s *ShopGoodsServiceImpl) List(c *gin.Context, req *shopmodels.GoodsQuery) (*shopmodels.GoodsListData, error) {
	data, err := s.dao.List(c, req)
	if err != nil {
		return nil, err
	}
	if data == nil || len(data.Rows) == 0 {
		return data, nil
	}
	if err = s.attachCategoryNames(c, data.Rows); err != nil {
		return nil, err
	}
	if err = s.attachSkus(c, data.Rows); err != nil {
		return nil, err
	}
	s.normalizeGoodsMediaURLs(c, data.Rows)
	return data, nil
}

// ExportCSV 以流式方式导出商品 CSV。
// 这里按分页批次拉取商品，避免一次性将全部商品加载到内存。
func (s *ShopGoodsServiceImpl) ExportCSV(c *gin.Context, req *shopmodels.GoodsQuery, csvWriter *csv.Writer, flush func()) error {
	exportReq := &shopmodels.GoodsQuery{}
	if req != nil {
		*exportReq = *req
	}
	// 导出时强制改用固定分页大小，由服务端主动分页遍历。
	exportReq.Page = 1
	exportReq.Size = goodsExportBatchSize

	firstBatch, err := s.List(c, exportReq)
	if err != nil {
		return err
	}

	if err = csvWriter.Write(goodsCSVHeader()); err != nil {
		return err
	}
	csvWriter.Flush()
	if err = csvWriter.Error(); err != nil {
		return err
	}
	if flush != nil {
		flush()
	}

	// 第一批已经查出，先写入，后续再继续翻页。
	if err = writeGoodsCSVRows(csvWriter, firstBatch.Rows); err != nil {
		return err
	}
	csvWriter.Flush()
	if err = csvWriter.Error(); err != nil {
		return err
	}
	if flush != nil {
		flush()
	}

	written := int64(len(firstBatch.Rows))
	total := firstBatch.Total
	for exportReq.Page++; written < total; exportReq.Page++ {
		// 如果客户端主动断开，则尽快终止导出，避免无效 IO。
		if c.Request.Context().Err() != nil {
			zap.L().Warn("goods csv export canceled by client", zap.Error(c.Request.Context().Err()))
			return c.Request.Context().Err()
		}

		batch, listErr := s.List(c, exportReq)
		if listErr != nil {
			zap.L().Error("query goods export batch fail", zap.Int64("page", exportReq.Page), zap.Error(listErr))
			return listErr
		}
		if batch == nil || len(batch.Rows) == 0 {
			break
		}

		if err = writeGoodsCSVRows(csvWriter, batch.Rows); err != nil {
			return err
		}
		csvWriter.Flush()
		if err = csvWriter.Error(); err != nil {
			return err
		}
		if flush != nil {
			flush()
		}
		written += int64(len(batch.Rows))
	}

	return nil
}

// goodsCSVHeader 返回商品导出 CSV 的表头定义。
func goodsCSVHeader() []string {
	return []string{
		"ID",
		"商品业务ID",
		"商品名称",
		"商品编码",
		"外部ID",
		"分类ID",
		"分类名称",
		"零售价",
		"是否上架",
		"库存数量",
		"销售单位",
		"重量",
		"重量单位",
		"主图",
		"视频",
		"图集",
		"首页模块ID",
		"SKU数量",
		"SKUID",
		"SKU业务ID",
		"SKU名称",
		"SKU编码",
		"SKU外部ID",
		"SKU条码",
		"SKU零售价",
		"SKU库存数量",
		"SKU销售单位",
		"SKU重量",
		"SKU重量单位",
		"SKU主图",
		"SKU视频",
		"SKU图集",
		"SKU描述",
		"描述",
		"创建时间",
		"更新时间",
	}
}

// writeGoodsCSVRows 将商品列表按“商品 x SKU”展开写入 CSV。
// 无 SKU 的商品也会输出一行，仅 SKU 列留空。
func writeGoodsCSVRows(csvWriter *csv.Writer, rows []*shopmodels.Goods) error {
	for _, row := range rows {
		if row == nil {
			continue
		}
		if len(row.Skus) == 0 {
			if err := csvWriter.Write(buildGoodsCSVRecord(row, nil)); err != nil {
				return err
			}
			continue
		}
		for _, sku := range row.Skus {
			if err := csvWriter.Write(buildGoodsCSVRecord(row, sku)); err != nil {
				return err
			}
		}
	}
	return nil
}

// buildGoodsCSVRecord 将商品及可选 SKU 展平为单行 CSV 记录。
func buildGoodsCSVRecord(row *shopmodels.Goods, sku *shopmodels.GoodsSku) []string {
	record := []string{
		strconv.FormatInt(row.ID, 10),
		row.GoodsID,
		row.GoodsName,
		row.GoodsCode,
		row.OuterID,
		strconv.FormatInt(row.ShopCategoryId, 10),
		row.ShopCategoryName,
		strconv.FormatFloat(row.RetailPrice, 'f', -1, 64),
		formatGoodsOnSale(row.IsOnSale),
		strconv.FormatInt(row.Quantity, 10),
		row.Unit,
		strconv.FormatFloat(row.Weight, 'f', -1, 64),
		row.WeightUnit,
		row.ImageURL,
		row.VideoURL,
		strings.Join(row.GalleryImagesArray, " | "),
		row.HomeModuleIDs,
		strconv.Itoa(len(row.Skus)),
	}
	if sku == nil {
		record = append(record, make([]string, 14)...)
	} else {
		record = append(record,
			strconv.FormatUint(sku.ID, 10),
			sku.SkuID,
			sku.SkuName,
			sku.SkuCode,
			sku.OuterID,
			sku.Barcode,
			strconv.FormatFloat(sku.RetailPrice, 'f', -1, 64),
			strconv.FormatInt(sku.Quantity, 10),
			sku.Unit,
			strconv.FormatFloat(sku.Weight, 'f', -1, 64),
			sku.WeightUnit,
			sku.ImageURL,
			sku.VideoURL,
			strings.Join(sku.GalleryImagesArray, " | "),
			sku.Description,
		)
	}
	record = append(record,
		row.Description,
		formatCSVTime(row.CreateTime),
		formatCSVTime(row.UpdateTime),
	)
	return record
}

// formatGoodsOnSale 将上下架状态转为中文文案。
func formatGoodsOnSale(status int32) string {
	if status == 1 {
		return "是"
	}
	return "否"
}

// formatCSVTime 统一格式化时间字段，空时间返回空字符串。
func formatCSVTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// attachCategoryNames 为商品批量补充分类名称。
// 这里会先去重分类 ID，再一次性查询分类，避免 N+1 查询。
func (s *ShopGoodsServiceImpl) attachCategoryNames(c *gin.Context, goodsRows []*shopmodels.Goods) error {
	if len(goodsRows) == 0 || s.categoryDao == nil {
		return nil
	}
	categoryIDSet := make(map[int64]struct{})
	for _, goods := range goodsRows {
		if goods == nil || goods.ShopCategoryId == 0 {
			continue
		}
		categoryIDSet[goods.ShopCategoryId] = struct{}{}
	}
	if len(categoryIDSet) == 0 {
		return nil
	}
	categoryIDs := make([]int64, 0, len(categoryIDSet))
	for id := range categoryIDSet {
		categoryIDs = append(categoryIDs, id)
	}
	sort.Slice(categoryIDs, func(i, j int) bool { return categoryIDs[i] < categoryIDs[j] })
	categories, err := s.categoryDao.ListByIDs(c, categoryIDs)
	if err != nil {
		return err
	}
	categoryNameMap := make(map[int64]string, len(categories))
	for _, category := range categories {
		if category == nil {
			continue
		}
		categoryNameMap[category.ID] = category.CategoryName
	}
	for _, goods := range goodsRows {
		if goods == nil {
			continue
		}
		goods.ShopCategoryName = categoryNameMap[goods.ShopCategoryId]
	}
	return nil
}

// Import 增量导入商品与规格数据
func (s *ShopGoodsServiceImpl) Import(c *gin.Context, records []shopmodels.ImportGoodsRecord) error {
	// 导入采用“先分批、再按批 diff、最后分别批量新增/更新”的流程，
	// 这样可以兼顾导入吞吐与数据库压力。
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
	// 商品业务 ID 优先取 external_id，其次回退到 product_code。
	goodsID := strings.TrimSpace(record.ExternalID)
	if goodsID == "" {
		goodsID = strings.TrimSpace(record.Data.ProductCode)
	}
	if goodsID == "" {
		return nil, errors.New("导入商品缺少external_id或product_code")
	}

	// 聚合 SKU 数据，反推商品层面的库存、售价和重量等摘要字段。
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
		GalleryImages: []string{},
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
	// SKU 业务 ID 允许从多个候选字段回退，尽可能吸收外部数据源的不一致格式。
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
		GoodsID:            goodsID,
		SkuID:              skuID,
		SkuName:            strings.TrimSpace(sku.Skuname),
		SkuCode:            strings.TrimSpace(sku.Skucode),
		OuterID:            strings.TrimSpace(sku.Lcmccode),
		Barcode:            strings.TrimSpace(sku.Barcode),
		Description:        strings.TrimSpace(record.Data.Remark),
		Weight:             sku.Weight,
		WeightUnit:         "kg",
		Unit:               defaultString(strings.TrimSpace(record.Data.UnitName), "件"),
		Quantity:           int64(sku.Size),
		RetailPrice:        pickRetailPrice(sku),
		GalleryImagesArray: []string{},
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
	// 同一批次内以业务 ID 做 key，天然具备去重效果；后出现的数据会覆盖前一条。
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
			// 新商品由服务端生成分布式 ID。
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
	// 图集在库里是 JSON 字符串，这里先序列化再参与比较。
	content, err := json.Marshal(req.GalleryImages)
	if err != nil {
		zap.L().Error("json marsh error", zap.Error(err))
		return false
	}
	return current.GoodsID == req.GoodsID &&
		current.GoodsName == req.GoodsName &&
		current.GoodsCode == req.GoodsCode &&
		current.OuterID == req.OuterID &&
		current.ImageURL == req.ImageURL &&
		current.RetailPrice == req.RetailPrice &&
		current.GalleryImages == string(content) &&
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
	// 先按商品业务 ID 一次性查询全部 SKU，再按 goods_id 回填。
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

// normalizeGoodsMediaURLs 将商品及 SKU 的媒体相对路径补齐为绝对地址，
// 同时把逗号分隔的图集字符串转成数组，便于前端直接使用。
func (s *ShopGoodsServiceImpl) normalizeGoodsMediaURLs(c *gin.Context, rows []*shopmodels.Goods) {
	for _, row := range rows {
		if row == nil {
			continue
		}
		row.ImageURL = fileUtils.BuildAbsoluteURL(c, row.ImageURL)
		row.VideoURL = fileUtils.BuildAbsoluteURL(c, row.VideoURL)
		row.GalleryImagesArray = splitAndNormalizeMediaURLs(c, row.GalleryImages)
		for _, sku := range row.Skus {
			if sku == nil {
				continue
			}
			sku.ImageURL = fileUtils.BuildAbsoluteURL(c, sku.ImageURL)
			sku.VideoURL = fileUtils.BuildAbsoluteURL(c, sku.VideoURL)
			sku.GalleryImagesArray = splitAndNormalizeMediaURLs(c, sku.GalleryImages)
		}
	}
}

// splitAndNormalizeMediaURLs 将逗号分隔的媒体列表切分并补齐绝对地址。
func splitAndNormalizeMediaURLs(c *gin.Context, raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []string{}
	}

	parts := strings.Split(raw, ",")
	urls := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		urls = append(urls, fileUtils.BuildAbsoluteURL(c, part))
	}
	return urls
}
