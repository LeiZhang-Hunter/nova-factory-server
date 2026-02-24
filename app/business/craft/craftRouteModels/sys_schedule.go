package craftRouteModels

// ScheduleReq 调度请求
type ScheduleReq struct {
	GatewayId int64  `json:"gateway_id,string"`
	UserName  string `json:"user_name"`
	Password  string `json:"password"`
	Size      int64  `form:"pageSize" default:"10000"` //数量
}
