package models

import "nova-factory-server/app/business/shop/logistics/client/api"

// TrackingQueryRequest 即时查询物流轨迹请求
type TrackingQueryRequest struct {
	Outsid      string `json:"outsid" binding:"required"`      // 物流单号
	CompanyCode string `json:"companyCode" binding:"required"` // 物流公司编码（对接快递鸟 ShipperCode）
}

// TrackingTraceNode 物流轨迹节点
type TrackingTraceNode struct {
	AcceptTime    string `json:"acceptTime"`    // 发生时间
	AcceptStation string `json:"acceptStation"` // 轨迹描述
	Location      string `json:"location"`      // 所在城市
	Action        string `json:"action"`        // 动作
}

// TrackingQueryResponse 即时查询物流轨迹响应
type TrackingQueryResponse struct {
	Outsid      string               `json:"outsid"`      // 物流单号
	CompanyCode string               `json:"companyCode"` // 物流公司编码
	CompanyName string               `json:"companyName"` // 物流公司名称
	State       string               `json:"state"`       // 物流状态：0无轨迹 1已揽收 2在途中 3签收 4问题件
	StateDesc   string               `json:"stateDesc"`   // 状态描述
	IsSigned    bool                 `json:"isSigned"`    // 是否已签收
	FromCache   bool                 `json:"fromCache"`   // 是否来自缓存
	Location    string               `json:"location"`    // 所在城市
	StateName   string               `json:"stateName"`
	Traces      []*TrackingTraceNode `json:"traces"` // 轨迹节点列表
}

func ConvertResult(result api.ExpressQueryResult) *TrackingQueryResponse {
	traceList := result.Traces()
	traces := make([]*TrackingTraceNode, 0, len(traceList))
	for _, t := range traceList {
		traces = append(traces, &TrackingTraceNode{
			AcceptTime:    t.AcceptTime(),
			AcceptStation: t.AcceptStation(),
			Location:      t.Location(),
			Action:        t.Action(),
		})
	}

	state := result.State()
	return &TrackingQueryResponse{
		Outsid:      result.LogisticCode(),
		CompanyCode: result.ShipperCode(),
		State:       state,
		StateDesc:   result.GetStateName(),
		Traces:      traces,
		Location:    result.Location(),
		StateName:   result.GetStateName(),
	}
}

// TrackingRecordSet 物流轨迹记录保存参数
type TrackingRecordSet struct {
	Outsid      string `json:"outsid"`
	CompanyCode string `json:"companyCode"`
	TraceJSON   string `json:"traceJson"`
	SignedTime  string `json:"signedTime"`
	OriginInfo  string `json:"originInfo"`
	DestInfo    string `json:"destInfo"`
}
