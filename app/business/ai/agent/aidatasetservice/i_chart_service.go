package aidatasetservice

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
)

type IChartService interface {
	SessionCreate(c *gin.Context, req *aidatasetmodels.CreateSessionsRequest) (*aidatasetmodels.CreateSessionsResponse, error)
	SessionUpdate(c *gin.Context, req *aidatasetmodels.UpdateSessionsRequest) (*aidatasetmodels.UpdateSessionsResponse, error)
	SessionList(c *gin.Context, req *aidatasetmodels.ListSessionRequest) (*aidatasetmodels.ListSessionResponse, error)
	SessionRemove(c *gin.Context, req *aidatasetmodels.DeleteSessionRequest) (*aidatasetmodels.DeleteSessionResponse, error)
	AgentSessionList(c *gin.Context, req *aidatasetmodels.ListAgentSessionsRequest) (*aidatasetmodels.AgentSessionListResponse, error)
	AgentSessionRemove(c *gin.Context, req *aidatasetmodels.RemoveAgentSessionsRequest) (*aidatasetmodels.RemoveAgentSessionsResponse, error)
	ConversationRelatedQuestions(c *gin.Context, req *aidatasetmodels.ConversationRelatedQuestionsRequest) (*aidatasetmodels.ConversationRelatedQuestionsResponse, error)
	Ask(c *gin.Context, req *aidatasetmodels.AskRequest) error
	AgentList(c *gin.Context, req *aidatasetmodels.ListAgentRequest) (*aidatasetmodels.ListAgentResponse, error)
	ChatsCompletions(c *gin.Context, req *aidatasetmodels.ChatsCompletionsRequest) (*aidatasetmodels.ChatsCompletionsResponse, error)
	AgentSessionCreate(c *gin.Context, req *aidatasetmodels.SessionAgentCreate) (*aidatasetmodels.SessionAgentResponse, error)
	AgentsCompletions(c *gin.Context, req *aidatasetmodels.AgentsCompletionsRequest) (*aidatasetmodels.AgentsCompletionsApiResponse, error)
}
