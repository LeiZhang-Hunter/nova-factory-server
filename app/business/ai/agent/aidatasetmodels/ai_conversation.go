package aidatasetmodels

import "nova-factory-server/app/baize"

type AiConversation struct {
	ID             int64  `json:"id" gorm:"column:id"`
	Name           string `json:"name" gorm:"column:name"`
	Message        string `json:"message" gorm:"column:message"`
	LLMProviderID  string `json:"llmProviderId" gorm:"column:llm_provider_id"`
	LLMModelID     string `json:"llmModelId" gorm:"column:llm_model_id"`
	EnableThinking int32  `json:"enableThinking" gorm:"column:enable_thinking"`
	ChatMode       string `json:"chatMode" gorm:"column:chat_mode"`
	DeptID         int64  `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

type SetAiConversation struct {
	Name           string `json:"name" binding:"required"`
	Message        string `json:"message"`
	LLMProviderID  string `json:"llmProviderId"`
	LLMModelID     string `json:"llmModelId"`
	EnableThinking int32  `json:"enableThinking"`
	ChatMode       string `json:"chatMode"`
}

type AiConversationQuery struct {
	Name string `form:"name"`
	baize.BaseEntityDQL
}

type AiConversationListData struct {
	Rows  []*AiConversation `json:"rows"`
	Total int64             `json:"total"`
}

// SendMessageInput input for sending a message
type SendMessageInput struct {
	ConversationID int64  `json:"conversation_id,string"`
	Content        string `json:"content"`
	TabID          string `json:"tab_id"`
}

// StopGenerationInput input for stopping message generation
type StopGenerationInput struct {
	ConversationID int64  `json:"conversation_id,string"`
	TabID          string `json:"tab_id"`
}
