//go:build ai

package agent

import (
	"strings"

	apiModels "nova-factory-server/app/business/shop/api/models"
	shopService "nova-factory-server/app/business/shop/api/service"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ASR struct {
	service shopService.IApiShopASRService
}

func NewASR(service shopService.IApiShopASRService) *ASR {
	return &ASR{service: service}
}

func (asr *ASR) PrivateRoutes(router *gin.RouterGroup) {
	group := router.Group("/api/v1/app/shop/agent/conversations")
	group.POST("/asr", asr.Recognize)
}

func (asr *ASR) Recognize(c *gin.Context) {
	req := new(apiModels.ShopASRSubmitReq)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if strings.TrimSpace(req.AudioBase64) == "" {
		baizeContext.Waring(c, "audio_base64不能为空")
		return
	}

	zap.L().Info("[DEBUG-ASR-UPLOAD] asr post request",
		zap.String("file_name", req.FileName),
		zap.Int("audio_base64_len", len(req.AudioBase64)),
	)

	text, err := asr.service.Recognize(c.Request.Context(), req.AudioBase64)
	if err != nil {
		zap.L().Error("asr recognize failed",
			zap.String("file_name", req.FileName),
			zap.Error(err),
		)
		baizeContext.Waring(c, err.Error())
		return
	}

	baizeContext.SuccessData(c, gin.H{
		"text": text,
	})
}
