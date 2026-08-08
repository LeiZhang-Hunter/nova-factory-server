package impl

import (
	"fmt"

	"nova-factory-server/app/business/data/dao"
	"nova-factory-server/app/business/data/models/dto"
	"nova-factory-server/app/business/data/models/entity"
	"nova-factory-server/app/business/data/service"

	"github.com/gin-gonic/gin"
)

// IModelConfigServiceImpl 模型配置服务实现。
type IModelConfigServiceImpl struct {
	dao dao.IModelConfigDAO
}

// NewIModelConfigServiceImpl 创建模型配置服务。
func NewIModelConfigServiceImpl(modelConfigDao dao.IModelConfigDAO) service.IModelConfigService {
	return &IModelConfigServiceImpl{dao: modelConfigDao}
}

// Get 获取当前模型配置；尚未保存时返回默认值。
func (s *IModelConfigServiceImpl) Get(ctx *gin.Context) (*dto.ModelConfigResponse, error) {
	config, err := s.dao.Get(ctx)
	if err != nil {
		return nil, err
	}
	if config == nil {
		config = entity.DefaultModelConfig()
	}
	return modelConfigResponse(config), nil
}

// Update 校验并保存模型配置。
func (s *IModelConfigServiceImpl) Update(ctx *gin.Context, req *dto.UpdateModelConfigRequest) (*dto.ModelConfigResponse, error) {
	config, err := s.dao.Get(ctx)
	if err != nil {
		return nil, err
	}
	if config == nil {
		config = entity.DefaultModelConfig()
	}

	config.Provider = req.Provider
	config.Model = req.Model
	config.EnableTemperature = req.EnableTemperature
	config.EnableTopP = req.EnableTopP
	config.EnableMaxTokens = req.EnableMaxTokens
	config.MaxContextCount = req.MaxContextCount
	config.RetrievalTopK = req.RetrievalTopK

	config.Temperature = entity.DefaultModelTemperature
	if req.Temperature != nil {
		config.Temperature = *req.Temperature
	}
	config.TopP = entity.DefaultModelTopP
	if req.TopP != nil {
		config.TopP = *req.TopP
	}
	config.MaxTokens = entity.DefaultModelMaxTokens
	if req.MaxTokens != nil {
		config.MaxTokens = *req.MaxTokens
	}
	config.RetrievalMatchThreshold = entity.DefaultModelMatchThreshold
	if req.RetrievalMatchThreshold != nil {
		config.RetrievalMatchThreshold = *req.RetrievalMatchThreshold
	}

	if err := validateModelConfig(config); err != nil {
		return nil, err
	}

	if err := s.dao.Save(ctx, config); err != nil {
		return nil, err
	}
	return modelConfigResponse(config), nil
}

func validateModelConfig(config *entity.ModelConfig) error {
	if config.Temperature < 0 || config.Temperature > 2 {
		return fmt.Errorf("Temperature 必须在 0 到 2 之间")
	}
	if config.TopP < 0 || config.TopP > 1 {
		return fmt.Errorf("Top P 必须在 0 到 1 之间")
	}
	if config.MaxTokens < 0 || config.MaxTokens > 200000 {
		return fmt.Errorf("Max Tokens 必须在 0 到 200000 之间")
	}
	if config.MaxContextCount < 0 || config.MaxContextCount > 100 {
		return fmt.Errorf("上下文轮数必须在 0 到 100 之间")
	}
	if config.RetrievalTopK < 0 || config.RetrievalTopK > 100 {
		return fmt.Errorf("检索 TopK 必须在 0 到 100 之间")
	}
	if config.RetrievalMatchThreshold < 0 || config.RetrievalMatchThreshold > 1 {
		return fmt.Errorf("匹配阈值必须在 0 到 1 之间")
	}
	return nil
}

func modelConfigResponse(config *entity.ModelConfig) *dto.ModelConfigResponse {
	return &dto.ModelConfigResponse{
		ID:                      config.ID,
		Provider:                config.Provider,
		Model:                   config.Model,
		Temperature:             config.Temperature,
		EnableTemperature:       config.EnableTemperature,
		TopP:                    config.TopP,
		EnableTopP:              config.EnableTopP,
		MaxTokens:               config.MaxTokens,
		EnableMaxTokens:         config.EnableMaxTokens,
		MaxContextCount:         config.MaxContextCount,
		RetrievalTopK:           config.RetrievalTopK,
		RetrievalMatchThreshold: config.RetrievalMatchThreshold,
		CreatedAt:               config.CreatedAt,
		UpdatedAt:               config.UpdatedAt,
	}
}
