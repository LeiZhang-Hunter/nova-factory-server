package service

import (
	"nova-factory-server/app/business/data/models/dto"

	"github.com/gin-gonic/gin"
)

// ICollectorRuleService 采集器规则下发服务接口
type ICollectorRuleService interface {
	Dispatch(c *gin.Context, collectorID, operator string, req *dto.BindCollectorRulesRequest) (*dto.DispatchLogResponse, error)
	PreviewProposed(c *gin.Context, collectorID string, req *dto.BindCollectorRulesRequest) (*dto.DispatchPreviewResponse, error)
	PreviewCurrent(c *gin.Context, collectorID string) (*dto.DispatchPreviewResponse, error)
	PullRules(c *gin.Context, deviceID, token, clientMD5 string) (*dto.PullRulesResponse, error)
	ListLogs(c *gin.Context, offset, limit int, collectorID string) ([]dto.DispatchLogResponse, int64, error)
}
