package aidatasetdao

import (
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"

	"github.com/gin-gonic/gin"
)

type IAiLLMSettingDao interface {
	Set(c *gin.Context, req *aidatasetmodels.SetSysAiLLMSetting) (*aidatasetmodels.SysAiLLMSetting, error)
	Get(c *gin.Context, req *aidatasetmodels.GetSysAiLLMSettingReq) (*aidatasetmodels.SysAiLLMSetting, error)
}
