package guanjiapo

import (
	"context"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/observer/integration/config"
	"nova-factory-server/app/utils/observer/integration/event"
	"testing"

	"gorm.io/gorm"
)

// ---- mock Event (event.Event) ----

type testEvent struct {
	cfg    config.Config
	cache_ cache.Cache
}

func (e *testEvent) Config() config.Config       { return e.cfg }
func (e *testEvent) Action() event.EventType     { return event.EventProductCreated }
func (e *testEvent) GetCache() cache.Cache       { return e.cache_ }
func (e *testEvent) GetCallback() event.Callback { return nil }
func (e *testEvent) GetDB() *gorm.DB             { return nil }
func (e *testEvent) GetTransaction() bool        { return false }

// ---- mock Base (event.Base) ----

type testBase struct{}

func (b *testBase) Metadata() map[string]any { return nil }
func (b *testBase) Ptr() any                 { return b }

// ---- mock ProductGetReqDataEvent ----

type testProductGetReqData struct {
	page       int64
	pageSize   int64
	returnType *int64
	goodsCode  *[]string
	goodsName  *[]string
}

func (d *testProductGetReqData) GetPage() int64          { return d.page }
func (d *testProductGetReqData) GetPageSize() int64      { return d.pageSize }
func (d *testProductGetReqData) GetReturnType() *int64   { return d.returnType }
func (d *testProductGetReqData) GetGoodsCode() *[]string { return d.goodsCode }
func (d *testProductGetReqData) GetGoodsName() *[]string { return d.goodsName }

// ---- mock ProductRelationQueryReqDataEvent ----

type testProductRelationQueryReqData struct {
	page      int64
	pageSize  int64
	goodsCode *[]string
	goodsName *[]string
}

func (d *testProductRelationQueryReqData) GetPage() int64          { return d.page }
func (d *testProductRelationQueryReqData) GetPageSize() int64      { return d.pageSize }
func (d *testProductRelationQueryReqData) GetGoodsCode() *[]string { return d.goodsCode }
func (d *testProductRelationQueryReqData) GetGoodsName() *[]string { return d.goodsName }

// ---- mock ZProductUpdateReqEvent ----

type testProdUpdateReqData struct {
	goodsID string
	remark  string
}

func (d *testProdUpdateReqData) GetGoodsID() string { return d.goodsID }
func (d *testProdUpdateReqData) GetRemark() string  { return d.remark }

type testZProductUpdateReq struct {
	*testEvent
	*testBase
	items *[]event.ZProductUpdateReqData
}

func (r *testZProductUpdateReq) GetItems() *[]event.ZProductUpdateReqData { return r.items }

// ---- config ----

type IntegrationConfig struct{}

func (i *IntegrationConfig) GetOverrideURL() string      { return "" }
func (i *IntegrationConfig) GetMetadata() map[string]any { return nil }
func (i *IntegrationConfig) GetData() string {
	return "{\"systemName\":\"管家婆一代\",\"credentials\":{\"appKey\":\"000220630104019715\",\"appSecret\":\"e373a270f0fd4b6e8c57df2890ff9637\",\"selfmallaccount\":\"青岛质德工业设备有限公司\"}}"
}
func (i *IntegrationConfig) GetType() string  { return "gjp_v1" }
func (i *IntegrationConfig) GetStatus() *bool { a := true; return &a }

// ---- tests ----

func TestSearchProducts(t *testing.T) {
	ev := new(testEvent)
	ev.cache_ = cache.NewCache()
	var cfg IntegrationConfig
	ev.cfg = &cfg

	req := event.ZProductGetReqEvent{
		Event:                  ev,
		ProductGetReqDataEvent: &testProductGetReqData{page: 1, pageSize: 10},
	}
	service := New()
	data, err := service.ProductSearcher().SearchProducts(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if len(data.GetGoods()) != 10 {
		t.Errorf("pageSize=10, expected 10 goods, got %d", len(data.GetGoods()))
	}
}

func TestUpdateProductRemark(t *testing.T) {
	ev := new(testEvent)
	ev.cache_ = cache.NewCache()
	var cfg IntegrationConfig
	ev.cfg = &cfg

	bs := new(testBase)
	items := []event.ZProductUpdateReqData{
		&testProdUpdateReqData{goodsID: "1001", remark: "测试备注"},
	}
	req := &testZProductUpdateReq{
		testEvent: ev,
		testBase:  bs,
		items:     &items,
	}

	service := New()
	resp, err := service.ProductSearcher().UpdateProductRemark(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.GetCode() != 0 {
		t.Errorf("expected code 0, got %d", resp.GetCode())
	}
}

func TestProductRelationQuery(t *testing.T) {
	ev := new(testEvent)
	ev.cache_ = cache.NewCache()
	var cfg IntegrationConfig
	ev.cfg = &cfg

	req := event.ZProductRelationQueryReqEvent{
		Event:                            ev,
		ProductRelationQueryReqDataEvent: &testProductRelationQueryReqData{page: 1, pageSize: 10},
	}
	service := New()
	resp, err := service.ProductSearcher().ProductRelationQuery(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.GetCode() != 0 {
		t.Errorf("expected code 0, got %d", resp.GetCode())
	}
	t.Logf("total=%d, goodsrelation count=%d", resp.GetTotal(), len(resp.GetGoodsRelation()))
}
