package order

import (
	"encoding/json"
	"testing"

	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/observer/integration/config"
	"nova-factory-server/app/utils/observer/integration/event"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type testOrderEvent struct {
	orders []event.OrderData
}

func (e *testOrderEvent) GetOrders() []event.OrderData { return e.orders }
func (e *testOrderEvent) Metadata() map[string]any     { return map[string]any{} }
func (e *testOrderEvent) Ptr() any                     { return e }
func (e *testOrderEvent) Config() config.Config        { return nil }
func (e *testOrderEvent) Action() event.EventType      { return "" }
func (e *testOrderEvent) GetCache() cache.Cache        { return nil }
func (e *testOrderEvent) GetCallback() event.Callback  { return nil }
func (e *testOrderEvent) GetDB() *gorm.DB              { return nil }
func (e *testOrderEvent) GetTransaction() bool         { return false }
func (e *testOrderEvent) GetCtx() *gin.Context         { return nil }

type testOrderData struct {
	tid               string
	weight            float64
	size              float64
	buyerNick         string
	buyerMessage      string
	sellerMemo        string
	total             float64
	postFee           float64
	receiverName      string
	receiverState     string
	receiverStateName string
	receiverCity      string
	receiverDistrict  string
	receiverAddress   string
	receiverPhone     string
	receiverMobile    string
	created           string
	status            string
	orderType         string
	invoiceName       string
	sellerFlag        string
	payTime           string
	bTypeCode         string
	details           []event.GoodsDetail
	accounts          []event.Account
}

func (d *testOrderData) Metadata() map[string]any        { return map[string]any{} }
func (d *testOrderData) Ptr() any                        { return d }
func (d *testOrderData) GetOrderNo() string              { return d.tid }
func (d *testOrderData) GetUserId() uint64               { return 0 }
func (d *testOrderData) GetWeight() float64              { return d.weight }
func (d *testOrderData) GetSize() float64                { return d.size }
func (d *testOrderData) GetBuyerNick() string            { return d.buyerNick }
func (d *testOrderData) GetBuyerMessage() string         { return d.buyerMessage }
func (d *testOrderData) GetSellerMemo() string           { return d.sellerMemo }
func (d *testOrderData) GetTotalAmount() float64         { return d.total }
func (d *testOrderData) GetPrivilege() float64           { return 0 }
func (d *testOrderData) GetPostFee() float64             { return d.postFee }
func (d *testOrderData) GetReceiverName() string         { return d.receiverName }
func (d *testOrderData) GetReceiverState() string        { return d.receiverState }
func (d *testOrderData) GetReceiverStateName() string    { return d.receiverStateName }
func (d *testOrderData) GetReceiverCity() string         { return d.receiverCity }
func (d *testOrderData) GetReceiverCityName() string     { return "" }
func (d *testOrderData) GetReceiverDistrict() string     { return d.receiverDistrict }
func (d *testOrderData) GetReceiverDistrictName() string { return "" }
func (d *testOrderData) GetReceiverAddress() string      { return d.receiverAddress }
func (d *testOrderData) GetReceiverPhone() string        { return d.receiverPhone }
func (d *testOrderData) GetReceiverMobile() string       { return d.receiverMobile }
func (d *testOrderData) GetCreated() string              { return d.created }
func (d *testOrderData) GetType() string                 { return d.orderType }
func (d *testOrderData) GetStatus() string               { return d.status }
func (d *testOrderData) GetInvoiceName() string          { return d.invoiceName }
func (d *testOrderData) GetSellerFlag() string           { return d.sellerFlag }
func (d *testOrderData) GetPayTime() string              { return d.payTime }
func (d *testOrderData) GetLogIstBTypeCode() string      { return "" }
func (d *testOrderData) GetLogIstBillCode() string       { return "" }
func (d *testOrderData) GetBTypeCode() string            { return d.bTypeCode }
func (d *testOrderData) GetDetails() []event.GoodsDetail { return d.details }
func (d *testOrderData) GetAccounts() []event.Account    { return d.accounts }
func (d *testOrderData) GetTransactionId() string        { return "" }
func (d *testOrderData) GetNotifyRaw() string            { return "" }
func (d *testOrderData) GetMchId() string                { return "" }
func (d *testOrderData) GetAppid() string                { return "" }
func (d *testOrderData) GetPayerOpenid() string          { return "" }
func (d *testOrderData) GetPayChannel() int              { return 0 }

type testGoodsDetail struct {
	oid      string
	numiid   int64
	skuid    int64
	num      float64
	outeriid string
	payment  float64
	picpath  string
	weight   float64
	size     float64
	unitid   int64
	unitqty  float64
}

func (d *testGoodsDetail) GetOid() string            { return d.oid }
func (d *testGoodsDetail) GetBarcode() string        { return "" }
func (d *testGoodsDetail) GetEshopGoodsId() string   { return "" }
func (d *testGoodsDetail) GetOuterIid() string       { return d.outeriid }
func (d *testGoodsDetail) GetEshopGoodsName() string { return "" }
func (d *testGoodsDetail) GetEshopSkuId() string     { return "" }
func (d *testGoodsDetail) GetEshopSkuName() string   { return "" }
func (d *testGoodsDetail) GetNumIid() int64          { return d.numiid }
func (d *testGoodsDetail) GetSkuId() int64           { return d.skuid }
func (d *testGoodsDetail) GetNum() float64           { return d.num }
func (d *testGoodsDetail) GetPayment() float64       { return d.payment }
func (d *testGoodsDetail) GetPicPath() string        { return d.picpath }
func (d *testGoodsDetail) GetWeight() float64        { return d.weight }
func (d *testGoodsDetail) GetSize() float64          { return d.size }
func (d *testGoodsDetail) GetUniTid() int64          { return d.unitid }
func (d *testGoodsDetail) GetUnitQty() float64       { return d.unitqty }
func (d *testGoodsDetail) GetRawData() string        { return "" }

type testAccount struct {
	financeCode string
	total       float64
}

func (a *testAccount) GetFinanceCode() string { return a.financeCode }
func (a *testAccount) GetTotal() float64      { return a.total }
func (a *testAccount) GetRawData() string     { return "" }

func TestToOrderSyncOrderJSON(t *testing.T) {
	req := &testOrderEvent{
		orders: []event.OrderData{
			&testOrderData{
				tid:               "T202606270001",
				weight:            1.25,
				size:              0.5,
				buyerNick:         "buyer-a",
				buyerMessage:      "请尽快发货",
				sellerMemo:        "优先处理",
				total:             128.8,
				postFee:           6,
				receiverName:      "张三",
				receiverState:     "浙江省",
				receiverStateName: "浙江省名称",
				receiverCity:      "杭州市",
				receiverDistrict:  "西湖区",
				receiverAddress:   "文三路 1 号",
				receiverPhone:     "0571-88888888",
				receiverMobile:    "13800000000",
				created:           "2026-06-27 10:00:00",
				status:            "Payed",
				orderType:         "NoCod",
				invoiceName:       "测试公司",
				sellerFlag:        "1",
				payTime:           "2026-06-27 10:05:00",
				bTypeCode:         "CUST001",
				details: []event.GoodsDetail{
					&testGoodsDetail{
						oid:      "D001",
						numiid:   1001,
						skuid:    2001,
						num:      2,
						outeriid: "SKU-001",
						payment:  120,
						picpath:  "https://example.com/a.png",
						weight:   1.2,
						size:     0.4,
						unitid:   1,
						unitqty:  2,
					},
				},
				accounts: []event.Account{
					&testAccount{financeCode: "WX", total: 128.8},
				},
			},
		},
	}

	converted := toOrderSyncOrder(req)
	gotBytes, err := json.Marshal(converted)
	if err != nil {
		t.Fatalf("marshal converted order: %v", err)
	}
	got := string(gotBytes)
	t.Logf("converted json: %s", got)

	const want = `{"orders":[{"tid":"T202606270001","weight":1.25,"size":0.5,"buyernick":"buyer-a","buyermessage":"请尽快发货","sellermemo":"优先处理","total":128.8,"postfee":6,"receivername":"张三","receiverstate":"浙江省","receivercity":"杭州市","receiverdistrict":"西湖区","receiveraddress":"文三路 1 号","receiverphone":"0571-88888888","receivermobile":"浙江省名称","receiverzip":"","created":"2026-06-27 10:00:00","status":"Payed","type":"NoCod","invoicename":"测试公司","sellerflag":"1","paytime":"2026-06-27 10:05:00","btypecode":"CUST001","getbtypeaddress":0,"details":[{"oid":"D001","numiid":1001,"skuid":2001,"num":2,"outeriid":"SKU-001","payment":120,"picpath":"https://example.com/a.png","weight":1.2,"size":0.4,"unitid":1,"unitqty":2}],"accounts":[{"financeCode":"WX","total":128.8}]}]}`
	if got != want {
		t.Fatalf("unexpected converted json\nwant: %s\n got: %s", want, got)
	}
}
