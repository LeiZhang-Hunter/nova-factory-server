package aiDataSetModels

import "nova-factory-server/app/baize"

type Llm struct {
	FrequencyPenalty float64 `json:"frequency_penalty"`
	ModelName        string  `json:"model_name"`
	PresencePenalty  float64 `json:"presence_penalty"`
	Temperature      float64 `json:"temperature"`
	TopP             float64 `json:"top_p"`
}

type Prompt struct {
	EmptyResponse            string      `json:"empty_response"`
	KeywordsSimilarityWeight float64     `json:"keywords_similarity_weight"`
	Opener                   string      `json:"opener"`
	Prompt                   string      `json:"prompt"`
	RerankModel              string      `json:"rerank_model"`
	SimilarityThreshold      float64     `json:"similarity_threshold"`
	TopN                     int         `json:"top_n"`
	Variables                []Variables `json:"variables"`
}

// CreateAssistantRequest 创建聊天代理

type CreateAssistantRequest struct {
	DatasetIds   []string      `json:"dataset_ids"`
	Name         string        `json:"name"`
	Avatar       string        `json:"avatar,omitempty"`
	Llm          *Llm          `json:"llm,omitempty"`
	Prompt       *Prompt       `json:"prompt,omitempty"`
	Description  string        `json:"description,omitempty"`
	Icon         string        `json:"icon,omitempty"`
	TopN         int           `json:"top_n,omitempty"`
	TopK         int           `json:"top_k,omitempty"`
	PromptConfig *PromptConfig `json:"prompt_config,omitempty"`
}

type Variables struct {
	Key      string `json:"key"`
	Optional bool   `json:"optional"`
}

type PromptConfig struct {
	System          string `json:"system"`
	Prologue        string `json:"prologue"`
	Parameters      string `json:"parameters"`
	EmptyResponse   string `json:"empty_response"`
	Quote           string `json:"quote"`
	Tts             string `json:"tts"`
	RefineMultiturn string `json:"refine_multiturn"`
}

type CreateAssistantData struct {
	Avatar      string   `json:"avatar"`
	CreateDate  string   `json:"create_date"`
	CreateTime  int64    `json:"create_time"`
	DatasetIds  []string `json:"dataset_ids"`
	Description string   `json:"description"`
	DoRefer     string   `json:"do_refer"`
	Id          string   `json:"id"`
	Language    string   `json:"language"`
	Llm         *Llm     `json:"llm,omitempty"`
	Name        string   `json:"name"`
	Prompt      *Prompt  `json:"prompt,omitempty"`
	PromptType  string   `json:"prompt_type"`
	Status      string   `json:"status"`
	TenantId    string   `json:"tenant_id"`
	TopK        int      `json:"top_k"`
	UpdateDate  string   `json:"update_date"`
	UpdateTime  int64    `json:"update_time"`
}

type CreateAssistantResponse struct {
	Code    int                  `json:"code"`
	Data    *CreateAssistantData `json:"data"`
	Message string               `json:"message"`
}

// UpdateAssistantRequest 更新聊天代理
type UpdateAssistantRequest struct {
	ChatId     string   `json:"chat_id"`
	DatasetIds []string `json:"dataset_ids,omitempty"`
	Name       string   `json:"name"`
	Avatar     string   `json:"avatar,omitempty"`
	Llm        *Llm     `json:"llm,omitempty"`
	Prompt     *Prompt  `json:"prompt,omitempty"`
}

type UpdateAssistantResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// DeleteAssistantRequest 删除聊天代理
type DeleteAssistantRequest struct {
	Ids []string `json:"ids"`
}

type DeleteAssistantResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// GetAssistantRequest 读取助理列表
type GetAssistantRequest struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	baize.BaseEntityDQL
}

type GetAssistantResponse struct {
	Code int `json:"code"`
	Data []struct {
		Avatar      string   `json:"avatar"`
		CreateDate  string   `json:"create_date"`
		CreateTime  int64    `json:"create_time"`
		Description string   `json:"description"`
		DoRefer     string   `json:"do_refer"`
		Id          string   `json:"id"`
		DatasetIds  []string `json:"dataset_ids"`
		Language    string   `json:"language"`
		Llm         *Llm     `json:"llm"`
		Name        string   `json:"name"`
		Prompt      *Prompt  `json:"prompt"`
		PromptType  string   `json:"prompt_type"`
		Status      string   `json:"status"`
		TenantId    string   `json:"tenant_id"`
		TopK        int      `json:"top_k"`
		UpdateDate  string   `json:"update_date"`
		UpdateTime  int64    `json:"update_time"`
	} `json:"data"`
	Message string `json:"message"`
}
