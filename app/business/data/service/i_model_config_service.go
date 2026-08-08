package service

import (
	"nova-factory-server/app/business/data/models/dto"

	"github.com/gin-gonic/gin"
)

// IModelConfigService 数据平台模型配置服务接口。
type IModelConfigService interface {
	Get(ctx *gin.Context) (*dto.ModelConfigResponse, error)
	Update(ctx *gin.Context, req *dto.UpdateModelConfigRequest) (*dto.ModelConfigResponse, error)
}
