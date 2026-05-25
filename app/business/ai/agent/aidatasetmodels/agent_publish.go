package aidatasetmodels

type AgentPublishReq struct {
	GatewayId int64  `json:"gateway_id"`
	AgentId   string `json:"agentId" binding:"required"`
	Action    string `json:"action"`
	TimeoutMS int64  `json:"timeoutMs"`
}

type AgentPublishRes struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
}
