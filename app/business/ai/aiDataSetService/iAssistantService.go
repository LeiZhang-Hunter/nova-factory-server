package aiDataSetService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/aiDataSetModels"
)

type IAssistantService interface {
	AddAssistant(c *gin.Context, req *aiDataSetModels.CreateAssistantRequest) (*aiDataSetModels.CreateAssistantResponse, error)
	UpdateAssistant(c *gin.Context, req *aiDataSetModels.UpdateAssistantRequest) (*aiDataSetModels.UpdateAssistantResponse, error)
	ListAssistant(c *gin.Context, req *aiDataSetModels.GetAssistantRequest) (*aiDataSetModels.GetAssistantResponse, error)
	DeleteAssistant(c *gin.Context, ids []string) error
}
