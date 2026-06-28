package kdniao

import (
	kdniaoResp "github.com/ttlv/kdniao/response"
	"nova-factory-server/app/business/shop/logistics/client/api"
)

// queryResultWrapper 将快递鸟响应包装为 api.ExpressQueryResult
type queryResultWrapper struct {
	resp      kdniaoResp.ExpressQueryResponse
	StateName string
}

var stateMap = map[string]string{
	"0": "货物暂无轨迹信息",
	"1": "快递员已上门揽收快递",
	"2": "货物在途中运输",
	"3": "快件被签收",
	"4": "货物运输途中存在异常",
	"5": "快递被转寄到新地址",
}

func newQueryResultWrapper(resp *kdniaoResp.ExpressQueryResponse) *queryResultWrapper {
	q := &queryResultWrapper{resp: *resp}
	stateName, ok := stateMap[resp.State]
	if ok {
		q.StateName = stateName
	}
	return q
}

func (w *queryResultWrapper) OrderCode() string      { return w.resp.OrderCode }
func (w *queryResultWrapper) ShipperCode() string    { return w.resp.ShipperCode }
func (w *queryResultWrapper) LogisticCode() string   { return w.resp.LogisticCode }
func (w *queryResultWrapper) Callback() string       { return w.resp.Callback }
func (w *queryResultWrapper) Success() bool          { return w.resp.Success }
func (w *queryResultWrapper) Reason() string         { return w.resp.Reason }
func (w *queryResultWrapper) State() string          { return w.resp.State }
func (w *queryResultWrapper) GetStateName() string   { return w.StateName }
func (w *queryResultWrapper) StateEx() string        { return w.resp.StateEx }
func (w *queryResultWrapper) Location() string       { return w.resp.Location }
func (w *queryResultWrapper) Station() string        { return w.resp.Station }
func (w *queryResultWrapper) StationTel() string     { return w.resp.StationTel }
func (w *queryResultWrapper) StationAdd() string     { return w.resp.StationAdd }
func (w *queryResultWrapper) DeliveryMan() string    { return w.resp.DeliveryMan }
func (w *queryResultWrapper) DeliveryManTel() string { return w.resp.DeliveryManTel }
func (w *queryResultWrapper) NextCity() string       { return w.resp.NextCity }

func (w *queryResultWrapper) Traces() []api.ExpressTraces {
	traces := make([]api.ExpressTraces, 0, len(w.resp.Traces))
	for _, t := range w.resp.Traces {
		traces = append(traces, &traceWrapper{t})
	}
	return traces
}

// traceWrapper 将快递鸟轨迹节点包装为 api.ExpressTraces
type traceWrapper struct {
	trace kdniaoResp.Traces
}

func (t *traceWrapper) AcceptTime() string    { return t.trace.AcceptTime }
func (t *traceWrapper) AcceptStation() string { return t.trace.AcceptStation }
func (t *traceWrapper) Location() string      { return t.trace.Location }
func (t *traceWrapper) Action() string        { return t.trace.Action }
func (t *traceWrapper) Remark() string        { return t.trace.Remark }

// errorResult 实现 api.ExpressQueryResult 的错误结果
type errorResult struct {
	reason string
}

func (e *errorResult) OrderCode() string           { return "" }
func (e *errorResult) ShipperCode() string         { return "" }
func (e *errorResult) LogisticCode() string        { return "" }
func (e *errorResult) Callback() string            { return "" }
func (e *errorResult) Success() bool               { return false }
func (e *errorResult) Reason() string              { return e.reason }
func (e *errorResult) State() string               { return "" }
func (e *errorResult) StateEx() string             { return "" }
func (e *errorResult) Location() string            { return "" }
func (e *errorResult) Station() string             { return "" }
func (e *errorResult) StationTel() string          { return "" }
func (e *errorResult) StationAdd() string          { return "" }
func (e *errorResult) DeliveryMan() string         { return "" }
func (e *errorResult) DeliveryManTel() string      { return "" }
func (e *errorResult) NextCity() string            { return "" }
func (e *errorResult) Traces() []api.ExpressTraces { return nil }
func (e *errorResult) GetStateName() string {
	return ""
}
