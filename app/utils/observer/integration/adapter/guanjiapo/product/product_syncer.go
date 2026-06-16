package product

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"nova-factory-server/app/utils/observer/integration/adapter/guanjiapo/client"
	"nova-factory-server/app/utils/observer/integration/adapter/guanjiapo/model"
	"nova-factory-server/app/utils/observer/integration/api"
	"nova-factory-server/app/utils/observer/integration/event"
	"nova-factory-server/app/utils/observer/integration/result"
)

const (
	methodProductGet    = "emall.goods.get"
	methodProductUpdate = "emall.goods.update"
)

type productSyncer struct {
	tokenURL string
	mode     string
}

// New 创建管家婆商品查询能力实现。
func New(tokenURL string, mode string) api.Product {
	return &productSyncer{
		tokenURL: tokenURL,
		mode:     mode,
	}
}

// SearchProducts 调用管家婆 emall.goods.get 查询商品数据。
func (s *productSyncer) SearchProducts(ctx context.Context, req event.ZProductGetReqEvent) (result.GoodsGetResponse, error) {
	if req.GetPage() < 1 {
		return nil, errors.New("page不能小于1")
	}
	if req.GetPageSize() < 1 {
		return nil, errors.New("pagesize不能小于1")
	}
	snapshot, err := client.ParseSnapshot(req.Config())
	if err != nil {
		return nil, err
	}
	token, err := client.ResolveAccessToken(ctx, snapshot, req.GetCache())
	if err != nil {
		return nil, err
	}
	body := map[string]any{
		"page":     req.GetPage(),
		"pagesize": req.GetPageSize(),
	}
	if req.GetReturnType() != nil {
		body["returntype"] = req.GetReturnType()
	}
	if req.GetGoodsName() != nil {
		body["goodsname"] = req.GetGoodsName()
	}
	if req.GetGoodsCode() != nil {
		body["goodscode"] = req.GetGoodsCode()
	}
	respBytes, err := client.DoSignedPost(ctx, s.tokenURL, snapshot, token, methodProductGet, body)
	if err != nil {
		return nil, err
	}
	ret := &model.ProductSearchResponse{}
	if uerr := json.Unmarshal(respBytes, ret); uerr != nil {
		return nil, errors.New("商品查询响应解析失败: " + strings.TrimSpace(string(respBytes)))
	}
	if ret.RespCode != 0 {
		return nil, errors.New(ret.RespMessage)
	}
	return ret, nil
}

// UpdateProductRemark 调用管家婆 emall.goods.update 更新商品备注。
func (s *productSyncer) UpdateProductRemark(ctx context.Context, req event.ZProductUpdateReqEvent) (result.ProductRemarkUpdateResponse, error) {
	if req.GetItems() == nil {
		return nil, errors.New("items不能为空")
	}
	snapshot, err := client.ParseSnapshot(req.Config())
	if err != nil {
		return nil, err
	}
	token, err := client.ResolveAccessToken(ctx, snapshot, req.GetCache())
	if err != nil {
		return nil, err
	}
	items := make([]map[string]any, 0, len(*req.GetItems()))
	for _, item := range *req.GetItems() {
		items = append(items, map[string]any{
			"goodsid": item.GetGoodsID(),
			"remark":  item.GetRemark(),
		})
	}
	body := map[string]any{
		"items": items,
	}
	respBytes, err := client.DoSignedPost(ctx, s.tokenURL, snapshot, token, methodProductUpdate, body)
	if err != nil {
		return nil, err
	}
	ret := &model.ProductRemarkUpdateResponse{}
	if uerr := json.Unmarshal(respBytes, ret); uerr != nil {
		return nil, errors.New("商品备注更新响应解析失败: " + strings.TrimSpace(string(respBytes)))
	}
	if ret.Code != 0 {
		return nil, errors.New(ret.Message)
	}
	return ret, nil
}

// ProductRelationQuery 商品对应关系查询（emall.goodsrelation.get）。
func (s *productSyncer) ProductRelationQuery(ctx context.Context, req event.ZProductRelationQueryReqEvent) (result.ProductRelationQueryResponse, error) {
	if req.GetPage() < 1 {
		return nil, errors.New("page不能小于1")
	}
	if req.GetPageSize() < 1 {
		return nil, errors.New("pagesize不能小于1")
	}
	snapshot, err := client.ParseSnapshot(req.Config())
	if err != nil {
		return nil, err
	}
	token, err := client.ResolveAccessToken(ctx, snapshot, req.GetCache())
	if err != nil {
		return nil, err
	}
	body := map[string]any{
		"page":     req.GetPage(),
		"pagesize": req.GetPageSize(),
	}
	if req.GetGoodsCode() != nil {
		body["goodscode"] = *req.GetGoodsCode()
	}
	if req.GetGoodsName() != nil {
		body["goodsname"] = *req.GetGoodsName()
	}
	respBytes, err := client.DoSignedPost(ctx, s.tokenURL, snapshot, token, "emall.goodsrelation.get", body)
	if err != nil {
		return nil, err
	}
	ret := &model.ProductRelationQueryResponse{}
	if uerr := json.Unmarshal(respBytes, ret); uerr != nil {
		return nil, errors.New("商品对应关系查询响应解析失败: " + strings.TrimSpace(string(respBytes)))
	}
	if ret.Code != 0 {
		return nil, errors.New(ret.Message)
	}
	return ret, nil
}
