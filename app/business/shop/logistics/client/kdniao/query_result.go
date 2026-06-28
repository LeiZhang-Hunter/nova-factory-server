package kdniao

import (
	kdniaoResp "github.com/ttlv/kdniao/response"
	"nova-factory-server/app/business/shop/logistics/client/api"
)

// queryResultWrapper 将快递鸟响应包装为 api.ExpressQueryResult
type queryResultWrapper struct {
	resp kdniaoResp.ExpressQueryResponse
}

func (w *queryResultWrapper) OrderCode() string      { return w.resp.OrderCode }
func (w *queryResultWrapper) ShipperCode() string    { return w.resp.ShipperCode }
func (w *queryResultWrapper) LogisticCode() string   { return w.resp.LogisticCode }
func (w *queryResultWrapper) Callback() string       { return w.resp.Callback }
func (w *queryResultWrapper) Success() bool          { return w.resp.Success }
func (w *queryResultWrapper) Reason() string         { return w.resp.Reason }
func (w *queryResultWrapper) State() string          { return w.resp.State }
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
