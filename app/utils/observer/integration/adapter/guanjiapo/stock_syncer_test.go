package guanjiapo

import (
	"context"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/observer/integration/event"
	"testing"
)

// ---- mock StockGetReqDataEvent ----

type testStockGetReqData struct {
	page         int64
	pageSize     int64
	skuCode      *string
	goodsCode    *[]string
	whsCode      *string
	isContainWhs *bool
}

func (d *testStockGetReqData) GetPage() int64          { return d.page }
func (d *testStockGetReqData) GetPageSize() int64      { return d.pageSize }
func (d *testStockGetReqData) GetSkuCode() *string     { return d.skuCode }
func (d *testStockGetReqData) GetGoodsCode() *[]string { return d.goodsCode }
func (d *testStockGetReqData) GetWhsCode() *string     { return d.whsCode }
func (d *testStockGetReqData) GetIsContainWhs() *bool  { return d.isContainWhs }

// ---- tests ----

func TestSearchStocks(t *testing.T) {
	ev := new(testEvent)
	ev.cache_ = cache.NewCache()
	var cfg IntegrationConfig
	ev.cfg = &cfg

	req := event.ZStockGetReqEvent{
		Event:                ev,
		StockGetReqDataEvent: &testStockGetReqData{page: 1, pageSize: 10},
	}
	service := New()
	resp, err := service.StockSearcher().SearchStocks(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.GetCode() != 0 {
		t.Errorf("expected code 0, got %d", resp.GetCode())
	}
	t.Logf("total=%d, stocks count=%d", resp.GetTotal(), len(resp.GetStocks()))
}
