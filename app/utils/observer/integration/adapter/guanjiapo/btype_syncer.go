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

type btypeSyncer struct {
	tokenURL string
	mode     string
}

func newBtypeSyncer(tokenURL string, mode string) api.BtypeSearcher {
	return &btypeSyncer{tokenURL: tokenURL, mode: mode}
}

// GetBtypes 查询往来单位数据（emall.btype.get）。
func (s *btypeSyncer) GetBtypes(ctx context.Context, req event.ZBtypeGetReqEvent) (result.BtypeGetResponse, error) {
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
	if req.GetBtypeCode() != nil {
		body["btypecode"] = *req.GetBtypeCode()
	}
	if req.GetBtypeName() != nil {
		body["btypename"] = *req.GetBtypeName()
	}
	if req.GetTel() != nil {
		body["tel"] = *req.GetTel()
	}
	respBytes, err := doSignedPost(ctx, s.tokenURL, snapshot, token, "emall.btype.get", body)
	if err != nil {
		return nil, err
	}
	ret := &model.BtypeGetResponse{}
	if err = json.Unmarshal(respBytes, ret); err != nil {
		return nil, errors.New("往来单位查询响应解析失败: " + strings.TrimSpace(string(respBytes)))
	}
	if ret.Code != 0 {
		return nil, errors.New(ret.Message)
	}
	return ret, nil
}
