package service

import (
	apiModels "nova-factory-server/app/business/shop/api/models"

	"github.com/gin-gonic/gin"
)

type IApiShopVoiceEmitter interface {
	SendEvent(event *apiModels.ShopVoiceServerEvent) error
	SendAudioChunk(chunk []byte) error
}

type IApiShopVoiceService interface {
	ProcessTurn(c *gin.Context, emitter IApiShopVoiceEmitter, req *apiModels.ShopVoiceSubmitReq) error
}
