package shopserviceimpl

import (
	"testing"
	"time"

	"nova-factory-server/app/business/shop/shopmodels"

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
		dao:    goodsDao,
		skuDao: skuDao,
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

type mockShopGoodsDao struct {
	byGoodsID    map[string]*shopmodels.Goods
	created      []*shopmodels.GoodsUpsert
	updated      []*shopmodels.GoodsUpsert
	batchCreated []*shopmodels.GoodsUpsert
	batchUpdated []*shopmodels.GoodsUpsert
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
	return nil
}

func (m *mockShopGoodsDao) GetByID(c *gin.Context, id int64) (*shopmodels.Goods, error) {
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
	return nil, nil
}

type mockShopSkuDao struct {
	bySkuID      map[string]*shopmodels.GoodsSku
	created      []*shopmodels.GoodsSkuUpsert
	updated      []*shopmodels.GoodsSkuUpsert
	batchCreated []*shopmodels.GoodsSkuUpsert
	batchUpdated []*shopmodels.GoodsSkuUpsert
}

func (m *mockShopSkuDao) Create(c *gin.Context, req *shopmodels.GoodsSkuUpsert) (*shopmodels.GoodsSku, error) {
	m.created = append(m.created, req)
	return &shopmodels.GoodsSku{
		ID:      101,
		GoodsID: req.GoodsID,
		SkuID:   req.SkuID,
	}, nil
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

func (m *mockShopSkuDao) GetBySkuID(c *gin.Context, skuID string) (*shopmodels.GoodsSku, error) {
	if item, ok := m.bySkuID[skuID]; ok {
		return item, nil
	}
	return nil, nil
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
