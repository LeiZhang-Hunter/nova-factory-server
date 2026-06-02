//go:build !ai

package impl

import (
	"context"
	"errors"

	"nova-factory-server/app/business/shop/api/service"
)

type IApiShopASRServiceNoSupportImpl struct{}

func NewIApiShopASRServiceImpl() service.IApiShopASRService {
	return &IApiShopASRServiceNoSupportImpl{}
}

func (s *IApiShopASRServiceNoSupportImpl) Recognize(ctx context.Context, audioBase64 string) (string, error) {
	return "", errors.New("asr service requires ai build tag")
}
