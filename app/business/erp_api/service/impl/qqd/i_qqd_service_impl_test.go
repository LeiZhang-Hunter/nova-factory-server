package qqd

import (
	"context"
	"errors"
	"net/url"
	"sync"
	"testing"
	"time"

	qqddaoimpl "nova-factory-server/app/business/erp_api/dao/impl"
	"nova-factory-server/app/business/erp_api/models"
	qqdservice "nova-factory-server/app/business/erp_api/service/qqd"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestIssueTokenStoresTokensInCacheAndRejectsReusedRefreshToken(t *testing.T) {
	service := newTestService(t, nil)
	ctx := context.Background()

	code := createTestAuthCode(t, service)
	initial, err := service.IssueToken(ctx, "app-key", "app-secret", code, "authorization_code", "")
	if err != nil {
		t.Fatalf("issue token: %v", err)
	}
	if initial.Token == "" || initial.RefreshToken == "" {
		t.Fatalf("expected token response, got %#v", initial)
	}
	if !service.ValidAccessToken(ctx, initial.Token, "app-key") {
		t.Fatal("expected access token to be valid")
	}

	refreshed, err := service.IssueToken(ctx, "app-key", "app-secret", "", "refresh_token", initial.RefreshToken)
	if err != nil {
		t.Fatalf("refresh token: %v", err)
	}
	if refreshed.Token == "" || refreshed.Token == initial.Token {
		t.Fatalf("expected new access token, got %#v", refreshed)
	}

	if _, err := service.IssueToken(ctx, "app-key", "app-secret", "", "refresh_token", initial.RefreshToken); err == nil {
		t.Fatal("expected reused refresh token to fail")
	}
}

func TestProductListReadsProductAndSkuTables(t *testing.T) {
	db := openTestStockDB(t)
	service := newTestService(t, db)

	response, err := service.ProductList(context.Background(), qqdservice.ProductListRequest{})
	if err != nil {
		t.Fatalf("product list: %v", err)
	}
	if response["totalresults"] != 3 {
		t.Fatalf("expected three products, got %#v", response)
	}
	items, ok := response["productinfo"].([]map[string]any)
	if !ok || len(items) != 3 {
		t.Fatalf("expected productinfo item, got %#v", response["productinfo"])
	}
	item := items[0]
	if item["productid"] != "1001" {
		t.Fatalf("expected first goods id 1001, got %#v", item["productid"])
	}
	skus, ok := item["skus"].([]map[string]any)
	if !ok || len(skus) != 1 {
		t.Fatalf("expected sku info, got %#v", item["skus"])
	}
	if skus[0]["skuid"] != "1001" || skus[0]["quantity"] != float64(8) {
		t.Fatalf("unexpected sku fields: %#v", skus[0])
	}
}

func TestProductStockUpdateUpdatesMySQLOnly(t *testing.T) {
	db := openTestStockDB(t)
	service := newTestService(t, db)

	response, err := service.ProductStockUpdate(context.Background(), qqdservice.ProductStockUpdateRequest{
		ProductID:  "1001",
		ProductQty: "99",
	})
	if err != nil {
		t.Fatalf("product stock update without skus: %v", err)
	}
	if response["productid"] != "1001" {
		t.Fatalf("expected response productid, got %#v", response)
	}
	assertGoodsQuantity(t, db, "1001", 99)
	assertSkuQuantity(t, db, "1001", 99)

	if _, err := service.ProductStockUpdate(context.Background(), qqdservice.ProductStockUpdateRequest{ProductQty: "1"}); err == nil {
		t.Fatal("expected missing productid to fail")
	}
	if _, err := service.ProductStockUpdate(context.Background(), qqdservice.ProductStockUpdateRequest{ProductID: "1001"}); err == nil {
		t.Fatal("expected missing productqty without skus to fail")
	}
}

func TestAddProductsUpdatesGoodsAndSkuTables(t *testing.T) {
	db := openTestStockDB(t)
	service := newTestService(t, db)

	goodsInfos := []map[string]any{
		{"goodsid": "2001", "quantity": 7},
		{"goodsid": "2002", "skus": []any{
			map[string]any{"skuid": "2003", "quantity": 11},
			map[string]any{"skuid": "2004", "quantity": 12},
		}},
	}

	stored, err := service.AddProducts(context.Background(), goodsInfos)
	if err != nil {
		t.Fatalf("add products: %v", err)
	}
	if len(stored) != len(goodsInfos) {
		t.Fatalf("expected returned goods length %d, got %d", len(goodsInfos), len(stored))
	}

	assertGoodsQuantity(t, db, "2001", 7)
	assertSkuQuantity(t, db, "2001", 7)
	assertGoodsQuantity(t, db, "2002", 23)
	assertSkuQuantity(t, db, "2003", 11)
	assertSkuQuantity(t, db, "2004", 12)
}

func newTestService(t *testing.T, db *gorm.DB) qqdservice.Service {
	t.Helper()

	service := &IQQDServiceImpl{
		cfg: models.QQDConfig{
			AppKey:          "app-key",
			AppSecret:       "app-secret",
			SelfMallAccount: "mall-account",
			CodeTTL:         "10m",
			TokenTTL:        "1h",
			RefreshTokenTTL: "24h",
		},
		cache: newTestCache(),
	}
	if db != nil {
		service.goodsDao = qqddaoimpl.NewIQQDGoodsDaoImpl(db)
		service.goodsSkuDao = qqddaoimpl.NewIQQDGoodsSkuDaoImpl(db)
	}
	return service
}

type testCache struct {
	mu   sync.Mutex
	data map[string]testCacheItem
}

type testCacheItem struct {
	value    string
	deadline time.Time
}

func newTestCache() *testCache {
	return &testCache{data: make(map[string]testCacheItem)}
}

func (c *testCache) Set(ctx context.Context, key string, val string, expiration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item := testCacheItem{value: val}
	if expiration > 0 {
		item.deadline = time.Now().Add(expiration)
	}
	c.data[key] = item
}

func (c *testCache) Get(ctx context.Context, key string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	item, ok := c.data[key]
	if !ok || item.expired() {
		delete(c.data, key)
		return "", errors.New("cache: nil")
	}
	return item.value, nil
}

func (c *testCache) Del(ctx context.Context, keys ...string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, key := range keys {
		delete(c.data, key)
	}
}

func (c *testCache) Exists(ctx context.Context, keys ...string) int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	var count int64
	for _, key := range keys {
		item, ok := c.data[key]
		if ok && !item.expired() {
			count++
			continue
		}
		delete(c.data, key)
	}
	return count
}

func (c *testCache) HSet(ctx context.Context, key string, values ...any) {}
func (c *testCache) Expire(ctx context.Context, key string, expiration time.Duration) bool {
	return false
}
func (c *testCache) HGet(ctx context.Context, key, field string) string { return "" }
func (c *testCache) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64) {
	return nil, 0
}
func (c *testCache) JudgmentAndHSet(ctx context.Context, rk, key string, gs any) {}
func (c *testCache) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) bool {
	return false
}
func (c *testCache) Publish(ctx context.Context, channel string, message interface{}) {}
func (c *testCache) Subscribe(ctx context.Context, channels ...string) *redis.PubSub  { return nil }
func (c *testCache) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) *redis.StringSliceCmd {
	return nil
}
func (c *testCache) ZAdd(ctx context.Context, key string, members ...redis.Z) *redis.IntCmd {
	return nil
}
func (c *testCache) MGet(ctx context.Context, keys []string) *redis.SliceCmd { return nil }

func (i testCacheItem) expired() bool {
	return !i.deadline.IsZero() && time.Now().After(i.deadline)
}

func createTestAuthCode(t *testing.T, service qqdservice.Service) string {
	t.Helper()

	callback, err := service.CreateAuthorizationCallback(context.Background(), "app-key", "app-secret", "http://example.com/callback", "")
	if err != nil {
		t.Fatalf("create auth callback: %v", err)
	}
	parsed, err := url.Parse(callback)
	if err != nil {
		t.Fatalf("parse callback: %v", err)
	}
	code := parsed.Query().Get("code")
	if code == "" {
		t.Fatal("expected authorization code")
	}
	return code
}

func openTestStockDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.Exec(`
CREATE TABLE shop_goods (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  goods_id varchar(64) NOT NULL,
  goods_name varchar(255) NOT NULL,
  outer_id varchar(64),
  image_url varchar(255),
  retail_price decimal(24,6) DEFAULT 0,
  description text,
  quantity int DEFAULT 0,
  is_on_sale int DEFAULT 1,
  create_time datetime DEFAULT NULL,
  update_time datetime DEFAULT NULL,
  state tinyint DEFAULT 0,
  shop_category_id bigint DEFAULT 0,
  home_module_ids varchar(255) DEFAULT ''
)`).Error; err != nil {
		t.Fatalf("create shop_goods: %v", err)
	}
	if err := db.Exec(`
CREATE TABLE shop_goods_sku (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  goods_id varchar(64) NOT NULL,
  sku_id varchar(64) NOT NULL,
  sku_name varchar(255),
  outer_id varchar(64),
  retail_price decimal(24,6) DEFAULT 0,
  quantity int DEFAULT 0,
  create_time datetime DEFAULT NULL,
  update_time datetime DEFAULT NULL,
  state tinyint DEFAULT 0
)`).Error; err != nil {
		t.Fatalf("create shop_goods_sku: %v", err)
	}

	now := time.Date(2026, 6, 5, 10, 0, 0, 0, time.Local)
	if err := db.Table("shop_goods").Create(map[string]any{
		"goods_id":        "1001",
		"goods_name":      "test product",
		"outer_id":        "OUT1001",
		"image_url":       "https://example.com/1001.jpg",
		"retail_price":    12.5,
		"description":     "desc",
		"quantity":        8,
		"is_on_sale":      1,
		"create_time":     now,
		"update_time":     now,
		"state":           0,
		"home_module_ids": "",
	}).Error; err != nil {
		t.Fatalf("insert shop_goods: %v", err)
	}
	if err := db.Table("shop_goods_sku").Create(map[string]any{
		"goods_id":     "1001",
		"sku_id":       "1001",
		"sku_name":     "red",
		"outer_id":     "SOUT1001",
		"retail_price": 13.5,
		"quantity":     8,
		"create_time":  now,
		"update_time":  now,
		"state":        0,
	}).Error; err != nil {
		t.Fatalf("insert shop_goods_sku: %v", err)
	}
	for _, goods := range []map[string]any{
		{"goods_id": "2001", "goods_name": "product 2001", "quantity": 0, "is_on_sale": 1, "state": 0, "home_module_ids": ""},
		{"goods_id": "2002", "goods_name": "product 2002", "quantity": 0, "is_on_sale": 1, "state": 0, "home_module_ids": ""},
	} {
		if err := db.Table("shop_goods").Create(goods).Error; err != nil {
			t.Fatalf("insert shop_goods: %v", err)
		}
	}
	for _, sku := range []map[string]any{
		{"goods_id": "2001", "sku_id": "2001", "quantity": 0, "state": 0},
		{"goods_id": "2002", "sku_id": "2003", "quantity": 0, "state": 0},
		{"goods_id": "2002", "sku_id": "2004", "quantity": 0, "state": 0},
	} {
		if err := db.Table("shop_goods_sku").Create(sku).Error; err != nil {
			t.Fatalf("insert shop_goods_sku: %v", err)
		}
	}
	return db
}

func assertGoodsQuantity(t *testing.T, db *gorm.DB, goodsID string, expected int64) {
	t.Helper()

	var goods struct {
		Quantity int64 `gorm:"column:quantity"`
	}
	if err := db.Table("shop_goods").Where("goods_id = ? AND state = ?", goodsID, 0).First(&goods).Error; err != nil {
		t.Fatalf("load goods %s: %v", goodsID, err)
	}
	if goods.Quantity != expected {
		t.Fatalf("expected goods %s quantity %d, got %d", goodsID, expected, goods.Quantity)
	}
}

func assertSkuQuantity(t *testing.T, db *gorm.DB, skuID string, expected int64) {
	t.Helper()

	var sku struct {
		Quantity int64 `gorm:"column:quantity"`
	}
	if err := db.Table("shop_goods_sku").Where("sku_id = ? AND state = ?", skuID, 0).First(&sku).Error; err != nil {
		t.Fatalf("load sku %s: %v", skuID, err)
	}
	if sku.Quantity != expected {
		t.Fatalf("expected sku %s quantity %d, got %d", skuID, expected, sku.Quantity)
	}
}
