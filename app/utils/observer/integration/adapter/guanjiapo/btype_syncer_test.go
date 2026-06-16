package guanjiapo

import (
	"context"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/observer/integration/event"
	"testing"
)

// ---- mock BtypeGetReqDataEvent ----

type testBtypeGetReqData struct {
	page      int64
	pageSize  int64
	btypeCode *string
	btypeName *string
	tel       *string
}

func (d *testBtypeGetReqData) GetPage() int64        { return d.page }
func (d *testBtypeGetReqData) GetPageSize() int64    { return d.pageSize }
func (d *testBtypeGetReqData) GetBtypeCode() *string { return d.btypeCode }
func (d *testBtypeGetReqData) GetBtypeName() *string { return d.btypeName }
func (d *testBtypeGetReqData) GetTel() *string       { return d.tel }

// ---- tests ----

func TestGetBtypes(t *testing.T) {
	ev := new(testEvent)
	ev.cache_ = cache.NewCache()
	var cfg IntegrationConfig
	ev.cfg = &cfg

	req := event.ZBtypeGetReqEvent{
		Event:                ev,
		BtypeGetReqDataEvent: &testBtypeGetReqData{page: 1, pageSize: 10},
	}
	service := New()
	resp, err := service.BtypeSearcher().GetBtypes(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.GetCode() != 0 {
		t.Errorf("expected code 0, got %d", resp.GetCode())
	}
	t.Logf("total=%d, datas count=%d", resp.GetTotal(), len(resp.GetDatas()))
}
