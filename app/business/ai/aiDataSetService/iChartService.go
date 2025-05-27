package aiDataSetService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/aiDataSetModels"
)

type IChartService interface {
	SessionCreate(c *gin.Context, req *aiDataSetModels.CreateSessionsRequest) (*aiDataSetModels.CreateSessionsResponse, error)
	SessionUpdate(c *gin.Context, req *aiDataSetModels.UpdateSessionsRequest) (*aiDataSetModels.UpdateSessionsResponse, error)
	SessionList(c *gin.Context, req *aiDataSetModels.ListSessionRequest) (*aiDataSetModels.ListSessionResponse, error)
	SessionRemove(c *gin.Context, req *aiDataSetModels.DeleteSessionRequest) (*aiDataSetModels.DeleteSessionResponse, error)
	AgentSessionList(c *gin.Context, req *aiDataSetModels.ListAgentSessionsRequest) (*aiDataSetModels.AgentSessionListResponse, error)
	AgentSessionRemove(c *gin.Context, req *aiDataSetModels.RemoveAgentSessionsRequest) (*aiDataSetModels.RemoveAgentSessionsResponse, error)
	ConversationRelatedQuestions(c *gin.Context, req *aiDataSetModels.ConversationRelatedQuestionsRequest) (*aiDataSetModels.ConversationRelatedQuestionsResponse, error)
	Ask(c *gin.Context, req *aiDataSetModels.AskRequest) error
	AgentList(c *gin.Context, req *aiDataSetModels.ListAgentRequest) (*aiDataSetModels.ListAgentResponse, error)
	ChatsCompletions(c *gin.Context, req *aiDataSetModels.ChatsCompletionsRequest) (*aiDataSetModels.ChatsCompletionsResponse, error)
	AgentSessionCreate(c *gin.Context, req *aiDataSetModels.SessionAgentCreate) (*aiDataSetModels.SessionAgentResponse, error)
	AgentsCompletions(c *gin.Context, req *aiDataSetModels.AgentsCompletionsRequest) (*aiDataSetModels.AgentsCompletionsApiResponse, error)
}
