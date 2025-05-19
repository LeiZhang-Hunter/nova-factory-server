package aiDataSetModels

import "nova-factory-server/app/baize"

// CreateSessionsRequest 使用 chat 助手创建会话
type CreateSessionsRequest struct {
	ChatId string `json:"chat_id" binding:"required"`
	Name   string `json:"name" binding:"required"`
	UserId string `json:"user_id"`
}

type CreateSessionsResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		ChatId     string `json:"chat_id"`
		CreateDate string `json:"create_date"`
		CreateTime int64  `json:"create_time"`
		Id         string `json:"id"`
		Messages   []struct {
			Content string `json:"content"`
			Role    string `json:"role"`
		} `json:"messages"`
		Name       string `json:"name"`
		UpdateDate string `json:"update_date"`
		UpdateTime int64  `json:"update_time"`
	} `json:"data"`
}

// UpdateSessionsRequest 使用 chat 助手更新会话
type UpdateSessionsRequest struct {
	ChatId    string `json:"chat_id" binding:"required"`
	SessionId string `json:"session_id" binding:"required"`
	Name      string `json:"name"`
	UserId    string `json:"user_id"`
}

type UpdateSessionsApiRequest struct {
	Name   string `json:"name"`
	UserId string `json:"user_id"`
}

type UpdateSessionsResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ListSessionRequest 列表请求
type ListSessionRequest struct {
	ChatId string `json:"chat_id" form:"chat_id" binding:"required"`
	Name   string `json:"name" form:"name"`
	Id     string `json:"id" form:"id"`
	UserId string `json:"user_id" form:"user_id"`
	baize.BaseEntityDQL
}

type ListSessionResponse struct {
	Code int `json:"code"`
	Data []struct {
		Chat       string `json:"chat"`
		CreateDate string `json:"create_date"`
		CreateTime int64  `json:"create_time"`
		Id         string `json:"id"`
		Messages   []struct {
			Content string `json:"content"`
			Role    string `json:"role"`
		} `json:"messages"`
		Name       string `json:"name"`
		UpdateDate string `json:"update_date"`
		UpdateTime int64  `json:"update_time"`
	} `json:"data"`
	Message string `json:"message"`
}

// DeleteSessionRequest 删除请求
type DeleteSessionRequest struct {
	ChatId string   `json:"chat_id" form:"chat_id" binding:"required"`
	Ids    []string `json:"ids" form:"ids"`
}

type DeleteSessionApiRequest struct {
	Ids []string `json:"ids"`
}

type DeleteSessionResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ListAgentSessionsRequest list代理会话
type ListAgentSessionsRequest struct {
	AgentId string `json:"agent_id" form:"agent_id" binding:"required"`
	UserId  string `json:"user_id" form:"user_id"`
	Dsl     bool   `json:"dsl" form:"dsl"`
	baize.BaseEntityDQL
}

type AgentSessionListResponse struct {
	Code int `json:"code"`
	Data []struct {
		AgentId string `json:"agent_id"`
		Dsl     struct {
			Answer     []interface{}              `json:"answer"`
			Components map[string]*AgentComponent `json:"components"`
			Graph      *AgentGraph                `json:"graph"`
			History    []interface{}              `json:"history"`
			Messages   []interface{}              `json:"messages"`
			Path       []interface{}              `json:"path"`
			Reference  []interface{}              `json:"reference"`
		} `json:"dsl"`
		Id      string `json:"id"`
		Message []struct {
			Content string `json:"content"`
			Role    string `json:"role"`
		} `json:"message"`
		Source string `json:"source"`
		UserId string `json:"user_id"`
	} `json:"data"`
	Message string `json:"message"`
}

// RemoveAgentSessionsRequest 删除Agent会话
type RemoveAgentSessionsRequest struct {
	Ids     []string `json:"ids" form:"ids"`
	AgentId string   `json:"agent_id" form:"agent_id" binding:"required"`
}

type RemoveAgentSessionsApiRequest struct {
	Ids []string `json:"ids"`
}

type RemoveAgentSessionsResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ConversationRelatedQuestionsRequest 相关提问
type ConversationRelatedQuestionsRequest struct {
	Question string `json:"question"`
}

type ConversationRelatedQuestionsResponse struct {
	Code    int      `json:"code"`
	Data    []string `json:"data"`
	Message string   `json:"message"`
}

// ListAgentRequest 代理列表
type ListAgentRequest struct {
	Id   string `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
	baize.BaseEntityDQL
}

type ListAgentResponse struct {
	Code int `json:"code"`
	Data []struct {
		Avatar      interface{} `json:"avatar"`
		CanvasType  interface{} `json:"canvas_type"`
		CreateDate  string      `json:"create_date"`
		CreateTime  int64       `json:"create_time"`
		Description interface{} `json:"description"`
		Dsl         struct {
			Answer     []interface{}              `json:"answer"`
			Components map[string]*AgentComponent `json:"components"`
			Graph      struct {
				Edges []interface{} `json:"edges"`
				Nodes []struct {
					Data struct {
						Label string `json:"label"`
						Name  string `json:"name"`
					} `json:"data"`
					Height   int    `json:"height"`
					Id       string `json:"id"`
					Position struct {
						X float64 `json:"x"`
						Y float64 `json:"y"`
					} `json:"position"`
					SourcePosition string `json:"sourcePosition"`
					TargetPosition string `json:"targetPosition"`
					Type           string `json:"type"`
					Width          int    `json:"width"`
				} `json:"nodes"`
			} `json:"graph"`
			History   []interface{} `json:"history"`
			Messages  []interface{} `json:"messages"`
			Path      []interface{} `json:"path"`
			Reference []interface{} `json:"reference"`
		} `json:"dsl"`
		Id         string `json:"id"`
		Title      string `json:"title"`
		UpdateDate string `json:"update_date"`
		UpdateTime int64  `json:"update_time"`
		UserId     string `json:"user_id"`
	} `json:"data"`
	Message string `json:"message"`
}

// ChatsCompletionsRequest 聊天交流
type ChatsCompletionsRequest struct {
	ChatId    string `json:"chat_id" form:"chat_id" binding:"required"`
	Question  string `json:"question"`
	Stream    bool   `json:"stream"`
	SessionId string `json:"session_id"`
	UserId    string `json:"user_id"`
}

type ChatsCompletionsApiRequest struct {
	Question  string `json:"question"`
	Stream    bool   `json:"stream"`
	SessionId string `json:"session_id"`
	UserId    string `json:"user_id"`
}

type ChatsCompletionsResponse struct {
	Code    int                   `json:"code"`
	Data    *ChatsCompletionsData `json:"data,omitempty"`
	Message string                `json:"message,omitempty"`
}

type CompletionsReference struct {
	Total  int `json:"total"`
	Chunks []struct {
		Id               string   `json:"id"`
		Content          string   `json:"content"`
		DocumentId       string   `json:"document_id"`
		DocumentName     string   `json:"document_name"`
		DatasetId        string   `json:"dataset_id"`
		ImageId          string   `json:"image_id"`
		Similarity       float64  `json:"similarity"`
		VectorSimilarity float64  `json:"vector_similarity"`
		TermSimilarity   float64  `json:"term_similarity"`
		Positions        []string `json:"positions"`
	} `json:"chunks"`
	DocAggs []struct {
		DocName string `json:"doc_name"`
		DocId   string `json:"doc_id"`
		Count   int    `json:"count"`
	} `json:"doc_aggs"`
}

type ChatsCompletionsData struct {
	Answer    string                `json:"answer"`
	Reference *CompletionsReference `json:"reference,omitempty"`
	Prompt    string                `json:"prompt"`
	Id        string                `json:"id"`
	SessionId string                `json:"session_id"`
}

// SessionAgentCreate 创建agent会话
type SessionAgentCreate struct {
	AgentId string                 `json:"agent_id"`
	Data    map[string]interface{} `json:"data"`
}

type SessionAgentCreateRequest struct {
	Data map[string]interface{} `json:"data"`
}

type SessionAgentResponse struct {
	Code int `json:"code"`
	Data struct {
		AgentId string `json:"agent_id"`
		Dsl     struct {
			Answer     []interface{}              `json:"answer"`
			Components map[string]*AgentComponent `json:"components"`
			EmbedId    string                     `json:"embed_id"`
			Graph      *AgentGraph                `json:"graph"`
			History    []interface{}              `json:"history"`
			Messages   []interface{}              `json:"messages"`
			Path       [][]string                 `json:"path"`
			Reference  []interface{}              `json:"reference"`
		} `json:"dsl"`
		Id      string `json:"id"`
		Message []struct {
			Content string `json:"content"`
			Role    string `json:"role"`
		} `json:"message"`
		Source string `json:"source"`
		UserId string `json:"user_id"`
	} `json:"data"`
	Message string `json:"message,omitempty"`
}

// AgentsCompletionsRequest 和Agent聊天
type AgentsCompletionsRequest struct {
	AgentId   string `json:"agent_id"`
	Question  string `json:"question"`
	Stream    bool   `json:"stream"`
	SessionId string `json:"session_id"`
	UserId    string `json:"user_id"`
	SyncDsl   bool   `json:"sync_dsl"`
}

type AgentsCompletionsApiRequest struct {
	Question  string `json:"question"`
	Stream    bool   `json:"stream"`
	SessionId string `json:"session_id"`
	UserId    string `json:"user_id"`
	SyncDsl   bool   `json:"sync_dsl"`
}

type AgentsCompletionsApiResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		SessionId string                `json:"session_id"`
		Answer    string                `json:"answer"`
		Reference *CompletionsReference `json:"reference,omitempty"`
		Param     []struct {
			Key      string `json:"key"`
			Name     string `json:"name"`
			Optional bool   `json:"optional"`
			Type     string `json:"type"`
			Value    string `json:"value,omitempty"`
		} `json:"param"`
	} `json:"data"`
}
