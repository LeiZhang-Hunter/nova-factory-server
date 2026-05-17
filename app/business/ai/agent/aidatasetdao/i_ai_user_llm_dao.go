package aidatasetdao

import (
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"

	"github.com/gin-gonic/gin"
)

type IAiUserLLMDao interface {
	Set(c *gin.Context, req *aidatasetmodels.SetSysUserLLM) (*aidatasetmodels.SysUserLLM, error)
	Get(c *gin.Context, req *aidatasetmodels.GetSysUserLLMReq) ([]*aidatasetmodels.SysUserLLM, error)
	Remove(c *gin.Context, req *aidatasetmodels.GetSysUserLLMReq) error
	// GetByFidAndLlm 读取用户的llm配置信息
	GetByFidAndLlm(factory string, llm string) (*aidatasetmodels.SysUserLLM, error)
}
