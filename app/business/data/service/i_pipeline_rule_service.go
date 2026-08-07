package service

import (
	"encoding/json"

	"nova-factory-server/app/business/data/models/dto"

	"github.com/gin-gonic/gin"
)

type IPipelineRuleService interface {
	Create(*gin.Context, *dto.CreatePipelineRuleRequest) (*dto.PipelineRuleResponse, error)
	Get(*gin.Context, string) (*dto.PipelineRuleResponse, error)
	List(*gin.Context, int, int, string, string, string) ([]dto.PipelineRuleResponse, int64, error)
	Update(*gin.Context, string, *dto.UpdatePipelineRuleRequest) (*dto.PipelineRuleResponse, error)
	Delete(*gin.Context, string) error
	ValidateConfig(*gin.Context, json.RawMessage, string) (json.RawMessage, error)
}

type IServiceConnectionService interface {
	Create(*gin.Context, *dto.CreateServiceConnectionRequest) (*dto.ServiceConnectionResponse, error)
	Get(*gin.Context, string) (*dto.ServiceConnectionResponse, error)
	List(*gin.Context, int, int, string, string, string) ([]dto.ServiceConnectionResponse, int64, error)
	Update(*gin.Context, string, *dto.UpdateServiceConnectionRequest) (*dto.ServiceConnectionResponse, error)
	Delete(*gin.Context, string) error
}
