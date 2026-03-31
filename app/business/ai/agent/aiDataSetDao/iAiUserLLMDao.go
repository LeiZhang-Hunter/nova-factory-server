package aiDataSetDao

import (
	"nova-factory-server/app/business/ai/agent/aiDataSetModels"

	"github.com/gin-gonic/gin"
)

type IAiUserLLMDao interface {
	Set(c *gin.Context, req *aiDataSetModels.SetSysUserLLM) (*aiDataSetModels.SysUserLLM, error)
	Get(c *gin.Context, req *aiDataSetModels.GetSysUserLLMReq) ([]*aiDataSetModels.SysUserLLM, error)
}
