package kdniao

import (
	"errors"

	"github.com/ttlv/kdniao"
	"github.com/ttlv/kdniao/sdk"
	"go.uber.org/zap"
	"nova-factory-server/app/business/shop/logistics/client/api"
)

// Adapter 快递鸟客户端适配器，实现 api.ExpressClient 接口
type Adapter struct {
	api sdk.ApiExpressQuery
}

// NewAdapter 从 Config 创建快递鸟适配器
func NewAdapter(cfg *Config) (api.ExpressClient, error) {
	if cfg == nil || cfg.EBusinessID == "" || cfg.AppKey == "" {
		return &Adapter{}, errors.New("cfg is nil")
	}
	config := kdniao.NewKdniaoConfig(cfg.EBusinessID, cfg.AppKey)
	logger := kdniao.NewKdniaoLogger()
	return &Adapter{
		api: sdk.NewExpressQuery(config, logger),
	}, nil
}

// Query 即时查询物流轨迹，返回统一模型
func (a *Adapter) Query(shipperCode, logisticCode string) (api.ExpressQueryResult, error) {
	if a.api == (sdk.ApiExpressQuery{}) {
		return &errorResult{reason: "快递鸟未配置，请在 config.yaml 中设置 kdniao.e_business_id 和 kdniao.app_key"}, nil
	}

	req := a.api.GetRequest(logisticCode)
	req.ShipperCode = shipperCode
	kdResp, err := a.api.GetResponse(req)
	if err != nil {
		zap.L().Error("get response error", zap.Error(err))
		return nil, err
	}

	if !kdResp.IsSuccess() {
		reason := kdResp.Reason
		if reason == "" {
			reason = "未知错误"
		}
		return &errorResult{reason: reason}, nil
	}

	return &queryResultWrapper{kdResp}, nil
}
