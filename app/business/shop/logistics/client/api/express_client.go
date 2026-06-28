package api

// ExpressClient 物流轨迹查询客户端接口（厂商无关）
// 每种第三方快递服务（快递鸟、快递100等）各自实现此接口
type ExpressClient interface {
	// Query 即时查询物流轨迹，返回统一模型
	Query(shipperCode, logisticCode string) (ExpressQueryResult, error)
}
