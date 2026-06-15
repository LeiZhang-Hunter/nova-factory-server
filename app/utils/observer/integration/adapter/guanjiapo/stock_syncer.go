package guanjiapo

import (
	"context"
	"encoding/json"
	"nova-factory-server/app/utils/observer/integration/adapter/guanjiapo/model"
	"nova-factory-server/app/utils/observer/integration/api"
	"nova-factory-server/app/utils/observer/integration/event"
	"nova-factory-server/app/utils/observer/integration/result"
	"strings"

	"gopkg.in/errgo.v2/errors"
)

type stockSyncer struct {
	tokenURL string
	mode     string
}

func newStockSyncer(tokenURL string, mode string) api.StockSearcher {
	return &stockSyncer{tokenURL: tokenURL, mode: mode}
}

// SearchStocks 查询库存数据（emall.stock.get）。
func (s *stockSyncer) SearchStocks(ctx context.Context, req event.ZStockGetReqEvent) (result.StockGetResponse, error) {
	if req.GetPage() < 1 {
		return nil, errors.New("page不能小于1")
	}
	if req.GetPageSize() < 1 {
		return nil, errors.New("pagesize不能小于1")
	}
	snapshot, err := parseSnapshot(req.Config())
	if err != nil {
		return nil, err
	}
	token, err := resolveAccessToken(ctx, snapshot, req.GetCache())
	if err != nil {
		return nil, err
	}
	body := map[string]any{
		"page":     req.GetPage(),
		"pagesize": req.GetPageSize(),
	}
	if req.GetSkuCode() != nil {
		body["skucode"] = *req.GetSkuCode()
	}
	if req.GetGoodsCode() != nil {
		body["goodscode"] = *req.GetGoodsCode()
	}
	if req.GetWhsCode() != nil {
		body["whscode"] = *req.GetWhsCode()
	}
	if req.GetIsContainWhs() != nil {
		body["iscontainwhs"] = *req.GetIsContainWhs()
	}
	respBytes, err := doSignedPost(ctx, s.tokenURL, snapshot, token, "emall.stock.get", body)
	if err != nil {
		return nil, err
	}
	ret := &model.StockGetResponse{}
	if err = json.Unmarshal(respBytes, ret); err != nil {
		return nil, errors.New("库存查询响应解析失败: " + strings.TrimSpace(string(respBytes)))
	}
	if ret.Code != 0 {
		return nil, errors.New(ret.Message)
	}
	return ret, nil
}
