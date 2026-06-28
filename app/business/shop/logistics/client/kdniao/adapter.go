package kdniao

import (
	"encoding/json"
	"errors"
	"github.com/ttlv/kdniao"
	"github.com/ttlv/kdniao/sdk"
	"go.uber.org/zap"
	"nova-factory-server/app/business/shop/logistics/client/api"
)

const Name = "kdniao"

// Adapter 快递鸟客户端适配器，实现 api.ExpressClient 接口
type Adapter struct {
	api sdk.ApiExpressQuery
}

// NewAdapter 从 Config 创建快递鸟适配器
func NewAdapter(apiConfig api.Config) (api.ExpressClient, error) {
	var cfg *config
	err := json.Unmarshal([]byte(apiConfig.GetData()), &cfg)
	if err != nil {
		zap.L().Error("json unmarshal failed", zap.String("config", apiConfig.GetData()), zap.Error(err))
		return nil, err
	}

	if cfg == nil || cfg.Credentials.EBusinessID == "" || cfg.Credentials.AppKey == "" {
		return &Adapter{}, errors.New("cfg is nil")
	}
	kdniaoCfg := kdniao.NewKdniaoConfig(cfg.Credentials.EBusinessID, cfg.Credentials.AppKey)
	logger := kdniao.NewKdniaoLogger()
	return &Adapter{
		api: sdk.NewExpressQuery(kdniaoCfg, logger),
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

	return newQueryResultWrapper(&kdResp), nil
}
