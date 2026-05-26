//go:build !ai

package impl

import (
	"errors"

	apiModels "nova-factory-server/app/business/shop/api/models"
	"nova-factory-server/app/business/shop/api/service"

	"github.com/gin-gonic/gin"
)

type IApiShopVoiceServiceNoSupportImpl struct{}

func NewIApiShopVoiceServiceImpl() service.IApiShopVoiceService {
	return &IApiShopVoiceServiceNoSupportImpl{}
}

func (s *IApiShopVoiceServiceNoSupportImpl) ProcessTurn(c *gin.Context, emitter service.IApiShopVoiceEmitter, req *apiModels.ShopVoiceSubmitReq) error {
	return errors.New("voice service requires ai build tag")
}
