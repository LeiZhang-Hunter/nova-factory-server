package shopserviceimpl

import (
	"bytes"
	"encoding/csv"
	"errors"
	"testing"
	"time"

	"nova-factory-server/app/business/shop/product/shopdao"
	"nova-factory-server/app/business/shop/product/shopmodels"

	"github.com/gin-gonic/gin"
)

func TestBuildGoodsUpsert(t *testing.T) {
	record := shopmodels.ImportGoodsRecord{
		ExternalID: "ext-1",
		Data: shopmodels.ImportGoodsRawData{
			ProductCode: "code-1",
			ProductName: "商品A",
			Remark:      "描述",
			UnitName:    "箱",
			Skus: []shopmodels.ImportGoodsSkuRawData{
				{Price: 12.5, Size: 3, Weight: 1.2},
				{Price2: 10, Size: 2},
			},
		},
		SyncedAt: time.Now(),
	}

	req, err := buildGoodsUpsert(record)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if req.GoodsID != "ext-1" {
		t.Fatalf("unexpected goods id: %s", req.GoodsID)
	}
	if req.Quantity != 5 {
		t.Fatalf("unexpected quantity: %d", req.Quantity)
	}
	if req.RetailPrice != 12.5 {
		t.Fatalf("unexpected retail price: %v", req.RetailPrice)
	}
	if req.Unit != "箱" {
		t.Fatalf("unexpected unit: %s", req.Unit)
	}
}

func TestImportIncrementalUpsert(t *testing.T) {
	goodsDao := &mockShopGoodsDao{
		byGoodsID: map[string]*shopmodels.Goods{
			"ext-1": {
				ID:      7,
				GoodsID: "ext-1",
			},
		},
	}
	skuDao := &mockShopSkuDao{
		bySkuID: map[string]*shopmodels.GoodsSku{
			"sku-1": {
				ID:      11,
				SkuID:   "sku-1",
				GoodsID: "ext-1",
			},
		},
	}
	service := &ShopGoodsServiceImpl{
		dao:         goodsDao,
		skuDao:      skuDao,
		categoryDao: &mockShopCategoryDao{},
	}
	records := []shopmodels.ImportGoodsRecord{
		{
			ExternalID: "ext-1",
			Data: shopmodels.ImportGoodsRawData{
				ProductCode: "code-1",
				ProductName: "商品A",
				Remark:      "描述A",
				Skus: []shopmodels.ImportGoodsSkuRawData{
					{Skuid: "sku-1", Skuname: "规格A", Price: 19.9, Size: 5},
					{Skuid: "sku-2", Skuname: "规格B", Price: 29.9, Size: 6},
				},
			},
		},
	}

	if err := service.Import(&gin.Context{}, records); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(goodsDao.batchCreated) != 0 {
		t.Fatalf("unexpected goods create count: %d", len(goodsDao.batchCreated))
	}
	if len(goodsDao.batchUpdated) != 1 {
		t.Fatalf("unexpected goods update count: %d", len(goodsDao.batchUpdated))
	}
	if len(skuDao.batchUpdated) != 1 {
		t.Fatalf("unexpected sku update count: %d", len(skuDao.batchUpdated))
	}
	if len(skuDao.batchCreated) != 1 {
		t.Fatalf("unexpected sku create count: %d", len(skuDao.batchCreated))
	}
	if skuDao.batchCreated[0].GoodsID != "ext-1" {
		t.Fatalf("unexpected created sku goods id: %s", skuDao.batchCreated[0].GoodsID)
	}
}

func TestListAttachSkus(t *testing.T) {
	goodsDao := &mockShopGoodsDao{
		listData: &shopmodels.GoodsListData{
			Rows: []*shopmodels.Goods{
				{GoodsID: "g-1", GoodsName: "商品1", ShopCategoryId: 101},
				{GoodsID: "g-2", GoodsName: "商品2", ShopCategoryId: 102},
			},
			Total: 2,
		},
	}
	skuDao := &mockShopSkuDao{
		skusByGoodsID: map[string][]*shopmodels.GoodsSku{
			"g-1": {
				{SkuID: "s-1", GoodsID: "g-1"},
				{SkuID: "s-2", GoodsID: "g-1"},
			},
			"g-2": {
				{SkuID: "s-3", GoodsID: "g-2"},
			},
		},
	}
	categoryDao := &mockShopCategoryDao{
		byID: map[int64]*shopmodels.Category{
			101: {ID: 101, CategoryName: "分类A"},
			102: {ID: 102, CategoryName: "分类B"},
		},
	}
	service := &ShopGoodsServiceImpl{
		dao:         goodsDao,
		skuDao:      skuDao,
		categoryDao: categoryDao,
	}

	data, err := service.List(&gin.Context{}, &shopmodels.GoodsQuery{})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(data.Rows[0].Skus) != 2 {
		t.Fatalf("unexpected sku count for goods 1: %d", len(data.Rows[0].Skus))
	}
	if len(data.Rows[1].Skus) != 1 {
		t.Fatalf("unexpected sku count for goods 2: %d", len(data.Rows[1].Skus))
	}
	if data.Rows[0].ShopCategoryName != "分类A" {
		t.Fatalf("unexpected category name for goods 1: %s", data.Rows[0].ShopCategoryName)
	}
	if data.Rows[1].ShopCategoryName != "分类B" {
		t.Fatalf("unexpected category name for goods 2: %s", data.Rows[1].ShopCategoryName)
	}
}

func TestGetByIDAttachSkus(t *testing.T) {
	goodsDao := &mockShopGoodsDao{
		byID: map[int64]*shopmodels.Goods{
			1: {ID: 1, GoodsID: "g-1", GoodsName: "商品1", ShopCategoryId: 101},
		},
	}
	skuDao := &mockShopSkuDao{
		skusByGoodsID: map[string][]*shopmodels.GoodsSku{
			"g-1": {
				{SkuID: "s-1", GoodsID: "g-1"},
				{SkuID: "s-2", GoodsID: "g-1"},
			},
		},
	}
	categoryDao := &mockShopCategoryDao{
		byID: map[int64]*shopmodels.Category{
			101: {ID: 101, CategoryName: "分类A"},
		},
	}
	service := &ShopGoodsServiceImpl{
		dao:         goodsDao,
		skuDao:      skuDao,
		categoryDao: categoryDao,
	}

	data, err := service.GetByID(&gin.Context{}, 1)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if data == nil {
		t.Fatal("expected goods data")
	}
	if len(data.Skus) != 2 {
		t.Fatalf("unexpected sku count: %d", len(data.Skus))
	}
	if data.ShopCategoryName != "分类A" {
		t.Fatalf("unexpected category name: %s", data.ShopCategoryName)
	}
}

func TestUpdateSyncsVectorSaleStatus(t *testing.T) {
	goodsDao := &mockShopGoodsDao{}
	vectorDao := &mockShopGoodsVectorDao{}
	service := &ShopGoodsServiceImpl{
		dao:       goodsDao,
		vectorDao: vectorDao,
	}
	req := &shopmodels.GoodsUpsert{
		ID:       12,
		GoodsID:  "g-12",
		IsOnSale: 0,
	}

	data, err := service.Update(&gin.Context{}, req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if data == nil || data.ID != 12 {
		t.Fatalf("unexpected update result: %#v", data)
	}
	if goodsDao.transactionCalls != 1 {
		t.Fatalf("expected transaction once, got %d", goodsDao.transactionCalls)
	}
	if len(goodsDao.updated) != 1 {
		t.Fatalf("expected one goods update, got %d", len(goodsDao.updated))
	}
	if len(vectorDao.updateSaleStatusCalls) != 1 {
		t.Fatalf("expected one vector sync, got %d", len(vectorDao.updateSaleStatusCalls))
	}
	call := vectorDao.updateSaleStatusCalls[0]
	if call.goodsDBID != 12 || call.isOnSale != 0 {
		t.Fatalf("unexpected vector sync call: %#v", call)
	}
}

func TestUpdateRollbackWhenVectorSyncFails(t *testing.T) {
	goodsDao := &mockShopGoodsDao{}
	vectorDao := &mockShopGoodsVectorDao{
		updateSaleStatusErr: errors.New("milvus unavailable"),
	}
	service := &ShopGoodsServiceImpl{
		dao:       goodsDao,
		vectorDao: vectorDao,
	}
	req := &shopmodels.GoodsUpsert{
		ID:       18,
		GoodsID:  "g-18",
		IsOnSale: 1,
	}

	data, err := service.Update(&gin.Context{}, req)
	if err == nil {
		t.Fatal("expected err")
	}
	if data != nil {
		t.Fatalf("expected nil data when vector sync fails, got %#v", data)
	}
	if err.Error() != "同步商品向量在售状态失败: milvus unavailable" {
		t.Fatalf("unexpected err: %v", err)
	}
	if goodsDao.transactionCalls != 1 {
		t.Fatalf("expected transaction once, got %d", goodsDao.transactionCalls)
	}
	if len(goodsDao.updated) != 1 {
		t.Fatalf("expected one goods update before rollback, got %d", len(goodsDao.updated))
	}
	if len(vectorDao.updateSaleStatusCalls) != 1 {
		t.Fatalf("expected one vector sync attempt, got %d", len(vectorDao.updateSaleStatusCalls))
	}
}

func TestDeleteByIDsRejectWhenGoodsHasSkus(t *testing.T) {
	goodsDao := &mockShopGoodsDao{
		byID: map[int64]*shopmodels.Goods{
			1: {ID: 1, GoodsID: "g-1"},
			2: {ID: 2, GoodsID: "g-2"},
		},
	}
	skuDao := &mockShopSkuDao{
		skusByGoodsID: map[string][]*shopmodels.GoodsSku{
			"g-2": {
				{SkuID: "s-1", GoodsID: "g-2"},
			},
		},
	}
	service := &ShopGoodsServiceImpl{
		dao:    goodsDao,
		skuDao: skuDao,
	}

	err := service.DeleteByIDs(&gin.Context{}, []int64{1, 2})
	if err == nil {
		t.Fatal("expected err")
	}
	if err.Error() != "商品下存在SKU，不允许删除" {
		t.Fatalf("unexpected err: %v", err)
	}
	if goodsDao.deleteCalls != 0 {
		t.Fatalf("expected delete not called, got %d", goodsDao.deleteCalls)
	}
}

func TestDeleteByIDsAllowWhenGoodsHasNoSkus(t *testing.T) {
	goodsDao := &mockShopGoodsDao{
		byID: map[int64]*shopmodels.Goods{
			3: {ID: 3, GoodsID: "g-3"},
			4: {ID: 4, GoodsID: "g-4"},
		},
	}
	skuDao := &mockShopSkuDao{
		skusByGoodsID: map[string][]*shopmodels.GoodsSku{},
	}
	service := &ShopGoodsServiceImpl{
		dao:    goodsDao,
		skuDao: skuDao,
	}

	err := service.DeleteByIDs(&gin.Context{}, []int64{3, 4})
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if goodsDao.deleteCalls != 1 {
		t.Fatalf("expected delete called once, got %d", goodsDao.deleteCalls)
	}
	if len(goodsDao.deletedIDs) != 2 || goodsDao.deletedIDs[0] != 3 || goodsDao.deletedIDs[1] != 4 {
		t.Fatalf("unexpected deleted ids: %#v", goodsDao.deletedIDs)
	}
}

func TestWriteGoodsCSVRowsIncludesSkuDetails(t *testing.T) {
	now := time.Date(2026, 5, 14, 12, 0, 0, 0, time.UTC)
	rows := []*shopmodels.Goods{
		{
			ID:               1,
			GoodsID:          "g-1",
			GoodsName:        "商品1",
			GoodsCode:        "code-1",
			OuterID:          "outer-1",
			ShopCategoryId:   101,
			ShopCategoryName: "分类A",
			RetailPrice:      99.5,
			IsOnSale:         1,
			Quantity:         30,
			Unit:             "件",
			Weight:           2.5,
			WeightUnit:       "kg",
			ImageURL:         "https://example.com/goods.png",
			VideoURL:         "https://example.com/goods.mp4",
			GalleryImagesArray: []string{
				"https://example.com/g1.png",
				"https://example.com/g2.png",
			},
			HomeModuleIDs: "home-1",
			Description:   "商品描述",
			BaseEntity:    shopmodels.Goods{}.BaseEntity,
			Skus: []*shopmodels.GoodsSku{
				{
					SkuID:              "sku-1",
					SkuName:            "规格1",
					SkuCode:            "sku-code-1",
					OuterID:            "sku-outer-1",
					Barcode:            "barcode-1",
					RetailPrice:        19.9,
					Quantity:           10,
					Unit:               "盒",
					Weight:             1.1,
					WeightUnit:         "kg",
					ImageURL:           "https://example.com/sku1.png",
					VideoURL:           "https://example.com/sku1.mp4",
					GalleryImagesArray: []string{"https://example.com/sku1-g1.png"},
					Description:        "规格1描述",
				},
				{
					SkuID:       "sku-2",
					SkuName:     "规格2",
					SkuCode:     "sku-code-2",
					RetailPrice: 29.9,
					Quantity:    20,
				},
			},
		},
		{
			ID:         2,
			GoodsID:    "g-2",
			GoodsName:  "商品2",
			Quantity:   5,
			BaseEntity: shopmodels.Goods{}.BaseEntity,
		},
	}
	rows[0].CreateTime = &now
	rows[0].UpdateTime = &now
	rows[1].CreateTime = &now
	rows[1].UpdateTime = &now

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	if err := writer.Write(goodsCSVHeader()); err != nil {
		t.Fatalf("write header err: %v", err)
	}
	if err := writeGoodsCSVRows(writer, rows); err != nil {
		t.Fatalf("write rows err: %v", err)
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		t.Fatalf("flush err: %v", err)
	}

	reader := csv.NewReader(bytes.NewReader(buf.Bytes()))
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("read csv err: %v", err)
	}
	if len(records) != 4 {
		t.Fatalf("unexpected record count: %d", len(records))
	}
	if records[1][18] != "sku-1" || records[1][19] != "规格1" {
		t.Fatalf("unexpected first sku columns: %v", records[1][18:20])
	}
	if records[2][18] != "sku-2" || records[2][19] != "规格2" {
		t.Fatalf("unexpected second sku columns: %v", records[2][18:20])
	}
	if records[3][1] != "g-2" {
		t.Fatalf("unexpected goods id for goods without sku: %s", records[3][1])
	}
	if records[3][18] != "" || records[3][19] != "" {
		t.Fatalf("expected empty sku columns for goods without sku: %v", records[3][18:20])
	}
}

type mockShopGoodsDao struct {
	byID             map[int64]*shopmodels.Goods
	byGoodsID        map[string]*shopmodels.Goods
	created          []*shopmodels.GoodsUpsert
	updated          []*shopmodels.GoodsUpsert
	batchCreated     []*shopmodels.GoodsUpsert
	batchUpdated     []*shopmodels.GoodsUpsert
	deletedIDs       []int64
	deleteCalls      int
	listData         *shopmodels.GoodsListData
	transactionCalls int
}

func (m *mockShopGoodsDao) Create(c *gin.Context, req *shopmodels.GoodsUpsert) (*shopmodels.Goods, error) {
	m.created = append(m.created, req)
	return &shopmodels.Goods{
		ID:      100,
		GoodsID: req.GoodsID,
	}, nil
}

func (m *mockShopGoodsDao) BatchCreate(c *gin.Context, reqs []*shopmodels.GoodsUpsert, batchSize int) error {
	m.batchCreated = append(m.batchCreated, reqs...)
	return nil
}

func (m *mockShopGoodsDao) Transaction(c *gin.Context, fn func(txDao shopdao.IShopGoodsDao) error) error {
	m.transactionCalls++
	return fn(m)
}

func (m *mockShopGoodsDao) Update(c *gin.Context, req *shopmodels.GoodsUpsert) (*shopmodels.Goods, error) {
	m.updated = append(m.updated, req)
	return &shopmodels.Goods{
		ID:      req.ID,
		GoodsID: req.GoodsID,
	}, nil
}

func (m *mockShopGoodsDao) BatchUpdate(c *gin.Context, reqs []*shopmodels.GoodsUpsert, batchSize int) error {
	m.batchUpdated = append(m.batchUpdated, reqs...)
	return nil
}

func (m *mockShopGoodsDao) DeleteByIDs(c *gin.Context, ids []int64) error {
	m.deleteCalls++
	m.deletedIDs = append(m.deletedIDs, ids...)
	return nil
}

func (m *mockShopGoodsDao) GetByID(c *gin.Context, id int64) (*shopmodels.Goods, error) {
	if item, ok := m.byID[id]; ok {
		return item, nil
	}
	return nil, nil
}

func (m *mockShopGoodsDao) GetByGoodsID(c *gin.Context, goodsID string) (*shopmodels.Goods, error) {
	if item, ok := m.byGoodsID[goodsID]; ok {
		return item, nil
	}
	return nil, nil
}

func (m *mockShopGoodsDao) ListByGoodsIDs(c *gin.Context, goodsIDs []string) ([]*shopmodels.Goods, error) {
	rows := make([]*shopmodels.Goods, 0, len(goodsIDs))
	for _, goodsID := range goodsIDs {
		if item, ok := m.byGoodsID[goodsID]; ok {
			rows = append(rows, item)
		}
	}
	return rows, nil
}

func (m *mockShopGoodsDao) List(c *gin.Context, req *shopmodels.GoodsQuery) (*shopmodels.GoodsListData, error) {
	return m.listData, nil
}

func (m *mockShopGoodsDao) UpdateStockByGoodsID(c *gin.Context, goodsID string, quantity int64) error {
	return nil
}

func (m *mockShopGoodsDao) UpsertByGoodsID(c *gin.Context, goodsID string, updates map[string]any) error {
	return nil
}

type mockShopSkuDao struct {
	bySkuID       map[string]*shopmodels.GoodsSku
	created       []*shopmodels.GoodsSkuUpsert
	updated       []*shopmodels.GoodsSkuUpsert
	batchCreated  []*shopmodels.GoodsSkuUpsert
	batchUpdated  []*shopmodels.GoodsSkuUpsert
	skusByGoodsID map[string][]*shopmodels.GoodsSku
}

type mockShopCategoryDao struct {
	byID map[int64]*shopmodels.Category
}

type mockShopGoodsVectorDao struct {
	updateSaleStatusCalls []mockGoodsVectorSaleStatusCall
	updateSaleStatusErr   error
}

type mockGoodsVectorSaleStatusCall struct {
	goodsDBID int64
	isOnSale  int32
}

func (m *mockShopCategoryDao) Create(c *gin.Context, req *shopmodels.CategoryUpsert) (*shopmodels.Category, error) {
	return nil, nil
}

func (m *mockShopCategoryDao) Update(c *gin.Context, req *shopmodels.CategoryUpsert) (*shopmodels.Category, error) {
	return nil, nil
}

func (m *mockShopCategoryDao) DeleteByIDs(c *gin.Context, ids []int64) error {
	return nil
}

func (m *mockShopCategoryDao) GetByID(c *gin.Context, id int64) (*shopmodels.Category, error) {
	if item, ok := m.byID[id]; ok {
		return item, nil
	}
	return nil, nil
}

func (m *mockShopCategoryDao) All(c *gin.Context) ([]*shopmodels.Category, error) {
	rows := make([]*shopmodels.Category, 0, len(m.byID))
	for _, item := range m.byID {
		rows = append(rows, item)
	}
	return rows, nil
}

func (m *mockShopCategoryDao) ListByIDs(c *gin.Context, ids []int64) ([]*shopmodels.Category, error) {
	rows := make([]*shopmodels.Category, 0, len(ids))
	for _, id := range ids {
		if item, ok := m.byID[id]; ok {
			rows = append(rows, item)
		}
	}
	return rows, nil
}

func (m *mockShopCategoryDao) List(c *gin.Context, req *shopmodels.CategoryQuery) (*shopmodels.CategoryListData, error) {
	return nil, nil
}

func (m *mockShopGoodsVectorDao) Upsert(c *gin.Context, goods *shopmodels.Goods,
	items []*shopmodels.GoodsVectorUpsertItem) (*shopmodels.GoodsVectorResult, error) {
	return nil, nil
}

func (m *mockShopGoodsVectorDao) UpdateSaleStatusByGoodsID(c *gin.Context, goodsDBID int64, isOnSale int32) error {
	m.updateSaleStatusCalls = append(m.updateSaleStatusCalls, mockGoodsVectorSaleStatusCall{
		goodsDBID: goodsDBID,
		isOnSale:  isOnSale,
	})
	return m.updateSaleStatusErr
}

func (m *mockShopGoodsVectorDao) DeleteBySkuIDs(c *gin.Context, skuIDs []int64) error {
	return nil
}

func (m *mockShopGoodsVectorDao) Search(c *gin.Context, req *shopmodels.GoodsVectorSearchReq,
	vector []float32, fallbackWithoutMetadata bool) (*shopmodels.GoodsVectorSearchData, error) {
	return nil, nil
}

func (m *mockShopGoodsVectorDao) BatchSearch(c *gin.Context, req *shopmodels.GoodsVectorBatchSearchReq,
	vectors [][]float32, fallbackWithoutMetadata bool) (*shopmodels.GoodsVectorBatchSearchData, error) {
	return nil, nil
}

func (m *mockShopSkuDao) Create(c *gin.Context, req *shopmodels.GoodsSkuUpsert) (*shopmodels.GoodsSku, error) {
	m.created = append(m.created, req)
	return &shopmodels.GoodsSku{
		ID:      101,
		GoodsID: req.GoodsID,
		SkuID:   req.SkuID,
	}, nil
}

func (m *mockShopSkuDao) Transaction(c *gin.Context, fn func(txDao shopdao.IShopSkuDao) error) error {
	return fn(m)
}

func (m *mockShopSkuDao) BatchCreate(c *gin.Context, reqs []*shopmodels.GoodsSkuUpsert, batchSize int) error {
	m.batchCreated = append(m.batchCreated, reqs...)
	return nil
}

func (m *mockShopSkuDao) Update(c *gin.Context, req *shopmodels.GoodsSkuUpsert) (*shopmodels.GoodsSku, error) {
	m.updated = append(m.updated, req)
	return &shopmodels.GoodsSku{
		ID:      req.ID,
		GoodsID: req.GoodsID,
		SkuID:   req.SkuID,
	}, nil
}

func (m *mockShopSkuDao) BatchUpdate(c *gin.Context, reqs []*shopmodels.GoodsSkuUpsert, batchSize int) error {
	m.batchUpdated = append(m.batchUpdated, reqs...)
	return nil
}

func (m *mockShopSkuDao) DeleteByIDs(c *gin.Context, ids []int64) error {
	return nil
}

func (m *mockShopSkuDao) GetByID(c *gin.Context, id int64) (*shopmodels.GoodsSku, error) {
	return nil, nil
}

func (m *mockShopSkuDao) ListByIDs(c *gin.Context, ids []int64) ([]*shopmodels.GoodsSku, error) {
	return []*shopmodels.GoodsSku{}, nil
}

func (m *mockShopSkuDao) GetBySkuID(c *gin.Context, skuID string) (*shopmodels.GoodsSku, error) {
	if item, ok := m.bySkuID[skuID]; ok {
		return item, nil
	}
	return nil, nil
}

func (m *mockShopSkuDao) ListByGoodsIDs(c *gin.Context, goodsIDs []string) ([]*shopmodels.GoodsSku, error) {
	rows := make([]*shopmodels.GoodsSku, 0)
	for _, goodsID := range goodsIDs {
		rows = append(rows, m.skusByGoodsID[goodsID]...)
	}
	return rows, nil
}

func (m *mockShopSkuDao) ListBySkuIDs(c *gin.Context, skuIDs []string) ([]*shopmodels.GoodsSku, error) {
	rows := make([]*shopmodels.GoodsSku, 0, len(skuIDs))
	for _, skuID := range skuIDs {
		if item, ok := m.bySkuID[skuID]; ok {
			rows = append(rows, item)
		}
	}
	return rows, nil
}

func (m *mockShopSkuDao) List(c *gin.Context, req *shopmodels.GoodsSkuQuery) (*shopmodels.GoodsSkuListData, error) {
	return nil, nil
}

func (m *mockShopSkuDao) UpdateStockBySkuID(c *gin.Context, skuID string, quantity int64) error {
	return nil
}

func (m *mockShopSkuDao) SumStockByGoodsID(c *gin.Context, goodsID string) (int64, error) {
	return 0, nil
}

func (m *mockShopSkuDao) UpsertBySkuID(c *gin.Context, skuID string, updates map[string]any) error {
	return nil
}
