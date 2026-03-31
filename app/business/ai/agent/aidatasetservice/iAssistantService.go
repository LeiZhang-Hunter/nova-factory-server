package aidatasetservice

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
)

type IAssistantService interface {
	AddAssistant(c *gin.Context, req *aidatasetmodels.CreateAssistantRequest) (*aidatasetmodels.CreateAssistantResponse, error)
	UpdateAssistant(c *gin.Context, req *aidatasetmodels.UpdateAssistantRequest) (*aidatasetmodels.UpdateAssistantResponse, error)
	ListAssistant(c *gin.Context, req *aidatasetmodels.GetAssistantRequest) (*aidatasetmodels.GetAssistantResponse, error)
	DeleteAssistant(c *gin.Context, ids []string) error
}
