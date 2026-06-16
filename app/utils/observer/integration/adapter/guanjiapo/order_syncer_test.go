package guanjiapo

import (
	"context"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/observer/integration/event"
	"testing"
)

// ---- mock ZOrderStatusSyncReqData ----

type testOrderStatusSyncData struct {
	tid          string
	status       string
	refundStatus string
}

func (d *testOrderStatusSyncData) GetTid() string          { return d.tid }
func (d *testOrderStatusSyncData) GetStatus() string       { return d.status }
func (d *testOrderStatusSyncData) GetRefundStatus() string { return d.refundStatus }

// ---- mock ZOrderStatusSyncReqEvent ----

type testOrderStatusSyncReq struct {
	*testEvent
	*testBase
	orders *[]event.ZOrderStatusSyncReqData
}

func (r *testOrderStatusSyncReq) GetOrders() *[]event.ZOrderStatusSyncReqData { return r.orders }

// ---- mock ZAfterSaleOrderDetail / ZAfterSaleOrderExDetail ----

type testAfterSaleDetail struct {
	oid            string
	eshopGoodsName string
	eshopSkuName   string
	backQty        int
	backTotal      float64
	outerIid       string
}

func (d *testAfterSaleDetail) GetOid() string            { return d.oid }
func (d *testAfterSaleDetail) GetEshopGoodsName() string { return d.eshopGoodsName }
func (d *testAfterSaleDetail) GetEshopSkuName() string   { return d.eshopSkuName }
func (d *testAfterSaleDetail) GetBackQty() int           { return d.backQty }
func (d *testAfterSaleDetail) GetBackTotal() float64     { return d.backTotal }
func (d *testAfterSaleDetail) GetOuterIid() string       { return d.outerIid }

type testAfterSaleExDetail struct {
	oid            string
	eshopGoodsName string
	eshopSkuName   string
	exchangeQty    int
	backTotal      float64
	outerIid       string
}

func (d *testAfterSaleExDetail) GetOid() string            { return d.oid }
func (d *testAfterSaleExDetail) GetEshopGoodsName() string { return d.eshopGoodsName }
func (d *testAfterSaleExDetail) GetEshopSkuName() string   { return d.eshopSkuName }
func (d *testAfterSaleExDetail) GetExchangeQty() int       { return d.exchangeQty }
func (d *testAfterSaleExDetail) GetBackTotal() float64     { return d.backTotal }
func (d *testAfterSaleExDetail) GetOuterIid() string       { return d.outerIid }

// ---- mock ZAfterSaleOrderSyncReqData ----

type testAfterSaleOrderSyncData struct {
	rtid           string
	tid            string
	total          float64
	privilege      float64
	postFee        float64
	created        string
	aftSaleType    string
	reasonCode     string
	logistBillCode string
	aftSaleRemark  string
	details        *[]event.ZAfterSaleOrderDetail
	exDetails      *[]event.ZAfterSaleOrderExDetail
}

func (d *testAfterSaleOrderSyncData) GetRtid() string                            { return d.rtid }
func (d *testAfterSaleOrderSyncData) GetTid() string                             { return d.tid }
func (d *testAfterSaleOrderSyncData) GetTotal() float64                          { return d.total }
func (d *testAfterSaleOrderSyncData) GetPrivilege() float64                      { return d.privilege }
func (d *testAfterSaleOrderSyncData) GetPostFee() float64                        { return d.postFee }
func (d *testAfterSaleOrderSyncData) GetCreated() string                         { return d.created }
func (d *testAfterSaleOrderSyncData) GetAftSaleType() string                     { return d.aftSaleType }
func (d *testAfterSaleOrderSyncData) GetReasonCode() string                      { return d.reasonCode }
func (d *testAfterSaleOrderSyncData) GetLogistBillCode() string                  { return d.logistBillCode }
func (d *testAfterSaleOrderSyncData) GetAftSaleRemark() string                   { return d.aftSaleRemark }
func (d *testAfterSaleOrderSyncData) GetDetails() *[]event.ZAfterSaleOrderDetail { return d.details }
func (d *testAfterSaleOrderSyncData) GetExDetails() *[]event.ZAfterSaleOrderExDetail {
	return d.exDetails
}

// ---- mock ZAfterSaleOrderSyncReqEvent ----

type testAfterSaleOrderSyncReq struct {
	*testEvent
	*testBase
	orders *[]event.ZAfterSaleOrderSyncReqData
}

func (r *testAfterSaleOrderSyncReq) GetOrders() *[]event.ZAfterSaleOrderSyncReqData { return r.orders }

// ---- mock OrderStatusGetReqDataEvent ----

type testOrderStatusGetReqData struct {
	orderCodes *[]string
}

func (d *testOrderStatusGetReqData) GetOrderCodes() *[]string { return d.orderCodes }

// ---- tests ----

func TestSyncOrderStatus(t *testing.T) {
	ev := new(testEvent)
	ev.cache_ = cache.NewCache()
	var cfg IntegrationConfig
	ev.cfg = &cfg
	bs := new(testBase)

	orders := []event.ZOrderStatusSyncReqData{
		&testOrderStatusSyncData{tid: "test001", status: "Payed", refundStatus: "Normal"},
	}
	req := &testOrderStatusSyncReq{
		testEvent: ev,
		testBase:  bs,
		orders:    &orders,
	}

	service := New()
	resp, err := service.OrderSyncer().SyncOrderStatus(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.GetCode() != 0 {
		t.Errorf("expected code 0, got %d", resp.GetCode())
	}
	t.Logf("message=%s", resp.GetMessage())
}

func TestSyncAfterSaleOrders(t *testing.T) {
	ev := new(testEvent)
	ev.cache_ = cache.NewCache()
	var cfg IntegrationConfig
	ev.cfg = &cfg
	bs := new(testBase)

	details := []event.ZAfterSaleOrderDetail{
		&testAfterSaleDetail{oid: "detail001", backQty: 1, backTotal: 100.0},
	}
	orders := []event.ZAfterSaleOrderSyncReqData{
		&testAfterSaleOrderSyncData{
			rtid:        "rt001",
			tid:         "test001",
			total:       100.0,
			privilege:   0,
			postFee:     0,
			created:     "2025-01-01 12:00:00",
			aftSaleType: "JustRefund",
			reasonCode:  "01",
			details:     &details,
		},
	}
	req := &testAfterSaleOrderSyncReq{
		testEvent: ev,
		testBase:  bs,
		orders:    &orders,
	}

	service := New()
	resp, err := service.OrderSyncer().SyncAfterSaleOrders(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.GetCode() != 0 {
		t.Errorf("expected code 0, got %d", resp.GetCode())
	}
	t.Logf("message=%s", resp.GetMessage())
}

func TestGetOrderStatus(t *testing.T) {
	ev := new(testEvent)
	ev.cache_ = cache.NewCache()
	var cfg IntegrationConfig
	ev.cfg = &cfg

	codes := []string{"test001"}
	req := event.ZOrderStatusGetReqEvent{
		Event:                      ev,
		OrderStatusGetReqDataEvent: &testOrderStatusGetReqData{orderCodes: &codes},
	}

	service := New()
	resp, err := service.OrderSyncer().GetOrderStatus(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.GetCode() != 0 {
		t.Errorf("expected code 0, got %d", resp.GetCode())
	}
	t.Logf("message=%s, orderstatus count=%d", resp.GetMessage(), len(resp.GetOrderStatus()))
}
