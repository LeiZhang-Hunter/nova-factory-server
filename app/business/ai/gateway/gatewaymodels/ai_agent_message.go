package gatewaymodels

import "nova-factory-server/app/baize"

type AIAgentMessage struct {
	ID              int64  `json:"id,string" gorm:"column:id"`
	ConversationID  int64  `json:"conversationId,string" gorm:"column:conversation_id"`
	Role            string `json:"role" gorm:"column:role"`
	Content         string `json:"content" gorm:"column:content"`
	ProviderID      string `json:"providerId" gorm:"column:provider_id"`
	ModelID         string `json:"modelId" gorm:"column:model_id"`
	Status          string `json:"status" gorm:"column:status"`
	Error           string `json:"error" gorm:"column:error"`
	InputTokens     int64  `json:"inputTokens" gorm:"column:input_tokens"`
	OutputTokens    int64  `json:"outputTokens" gorm:"column:output_tokens"`
	FinishReason    string `json:"finishReason" gorm:"column:finish_reason"`
	ToolCalls       string `json:"toolCalls" gorm:"column:tool_calls"`
	ToolCallID      string `json:"toolCallId" gorm:"column:tool_call_id"`
	ToolCallName    string `json:"toolCallName" gorm:"column:tool_call_name"`
	ThinkingContent string `json:"thinkingContent" gorm:"column:thinking_content"`
	Segments        string `json:"segments" gorm:"column:segments"`
	ImagesJSON      string `json:"imagesJson" gorm:"column:images_json"`
	DeptID          int64  `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

type AIAgentMessageUpsert struct {
	ID              int64  `json:"id,string"`
	ConversationID  int64  `json:"conversationId,string"`
	Role            string `json:"role"`
	Content         string `json:"content"`
	ProviderID      string `json:"providerId"`
	ModelID         string `json:"modelId"`
	Status          string `json:"status"`
	Error           string `json:"error"`
	InputTokens     int64  `json:"inputTokens"`
	OutputTokens    int64  `json:"outputTokens"`
	FinishReason    string `json:"finishReason"`
	ToolCalls       string `json:"toolCalls"`
	ToolCallID      string `json:"toolCallId"`
	ToolCallName    string `json:"toolCallName"`
	ThinkingContent string `json:"thinkingContent"`
	Segments        string `json:"segments"`
	ImagesJSON      string `json:"imagesJson"`
}

type AIAgentMessageQuery struct {
	ConversationID  int64   `form:"conversationId"`
	ConversationIDs []int64 `form:"-"`
	Role            string  `form:"role"`
	ProviderID      string  `form:"providerId"`
	ModelID         string  `form:"modelId"`
	Status          string  `form:"status"`
	Page            int64   `form:"page"`
	Size            int64   `form:"size"`
}

type AIAgentMessageListData struct {
	Rows  []*AIAgentMessage `json:"rows"`
	Total int64             `json:"total"`
}
