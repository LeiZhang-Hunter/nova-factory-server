package aiDataSetService

import (
	"nova-factory-server/app/business/ai/agent/aiDataSetModels"

	"github.com/gin-gonic/gin"
)

type IAiLLMSettingService interface {
	Set(c *gin.Context, req *aiDataSetModels.SetSysAiLLMSetting) (*aiDataSetModels.SysAiLLMSetting, error)
	Get(c *gin.Context, req *aiDataSetModels.GetSysAiLLMSettingReq) (*aiDataSetModels.SysAiLLMSetting, error)
}
