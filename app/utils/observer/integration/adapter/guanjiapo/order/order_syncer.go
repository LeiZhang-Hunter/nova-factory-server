package order

import (
	"context"
	"encoding/json"
	"strings"

	"nova-factory-server/app/utils/observer/integration/adapter/guanjiapo/client"
	"nova-factory-server/app/utils/observer/integration/adapter/guanjiapo/model"
	"nova-factory-server/app/utils/observer/integration/api"
	"nova-factory-server/app/utils/observer/integration/event"
	"nova-factory-server/app/utils/observer/integration/result"

	"gopkg.in/errgo.v2/errors"
)

type orderSyncer struct {
	tokenURL string
	mode     string
}

// New 创建管家婆订单同步能力实现。
func New(tokenURL string, mode string) api.OrderSyncer {
	return &orderSyncer{tokenURL: tokenURL, mode: mode}
}

// SyncOrders 同步订单至管家婆（emall.order.synchronize）。
func (s *orderSyncer) SyncOrders(ctx context.Context, req event.OrderEvent) (result.OrderSyncResponse, error) {
	if req == nil || len(req.GetOrders()) == 0 {
		return nil, errors.New("orders不能为空")
	}
	snapshot, err := client.ParseSnapshot(req.Config())
	if err != nil {
		return nil, err
	}
	token, err := client.ResolveAccessToken(ctx, snapshot, req.GetCache())
	if err != nil {
		return nil, err
	}
	body := toOrderSyncOrder(req)
	respBytes, err := client.DoSignedPost(ctx, s.tokenURL, snapshot, token, "emall.order.synchronize", body)
	if err != nil {
		return nil, err
	}
	ret := &model.OrderSyncResponse{}
	if err = json.Unmarshal(respBytes, ret); err != nil {
		return nil, err
	}
	if ret.Code != 0 {
		return nil, errors.New(ret.Message)
	}
	for _, o := range ret.Orders {
		if o.BillCode == "" {
			return nil, errors.New(o.Message)
		}
	}
	return ret, nil
}

// SyncOrderStatus 订单状态同步（emall.orderstatus.synchronize）。
func (s *orderSyncer) SyncOrderStatus(ctx context.Context, req event.ZOrderStatusSyncReqEvent) (result.OrderStatusSyncResponse, error) {
	if req == nil || req.GetOrders() == nil {
		return nil, errors.New("orders不能为空")
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
		"orders": req.GetOrders(),
	}
	respBytes, err := client.DoSignedPost(ctx, s.tokenURL, snapshot, token, "emall.orderstatus.synchronize", body)
	if err != nil {
		return nil, err
	}
	ret := &model.OrderStatusSyncResponse{}
	if err = json.Unmarshal(respBytes, ret); err != nil {
		return nil, errors.New("订单状态同步响应解析失败: " + strings.TrimSpace(string(respBytes)))
	}
	if ret.Code != 0 {
		return nil, errors.New(ret.Message)
	}
	return ret, nil
}

// SyncAfterSaleOrders 售后订单同步（emall.afterorder.synchronize）。
func (s *orderSyncer) SyncAfterSaleOrders(ctx context.Context, req event.ZAfterSaleOrderSyncReqEvent) (result.AfterSaleOrderSyncResponse, error) {
	if req == nil || req.GetOrders() == nil {
		return nil, errors.New("orders不能为空")
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
		"orders": toSyncAfterSaleOrders(req),
	}
	respBytes, err := client.DoSignedPost(ctx, s.tokenURL, snapshot, token, "emall.afterorder.synchronize", body)
	if err != nil {
		return nil, err
	}
	ret := &model.AfterSaleOrderSyncResponse{}
	if err = json.Unmarshal(respBytes, ret); err != nil {
		return nil, errors.New("售后订单同步响应解析失败: " + strings.TrimSpace(string(respBytes)))
	}
	if ret.Code != 0 {
		return nil, errors.New(ret.Message)
	}
	return ret, nil
}

// GetOrderStatus 查询订单状态（emall.orderstatus.get）。
func (s *orderSyncer) GetOrderStatus(ctx context.Context, req event.ZOrderStatusGetReqEvent) (result.OrderStatusGetResponse, error) {
	if req.GetOrderCodes() == nil {
		return nil, errors.New("ordercodes不能为空")
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
		"ordercodes": *req.GetOrderCodes(),
	}
	respBytes, err := client.DoSignedPost(ctx, s.tokenURL, snapshot, token, "emall.orderstatus.get", body)
	if err != nil {
		return nil, err
	}
	ret := &model.OrderStatusGetResponse{}
	if err = json.Unmarshal(respBytes, ret); err != nil {
		return nil, errors.New("订单状态查询响应解析失败: " + strings.TrimSpace(string(respBytes)))
	}
	if ret.Code != 0 {
		return nil, errors.New(ret.Message)
	}
	return ret, nil
}
