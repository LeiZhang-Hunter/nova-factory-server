package aiDataSetController

import (
	"context"
	"time"

	"nova-factory-server/app/business/ai/agent/aiDataSetModels"
	"nova-factory-server/app/utils/baizeContext"
	models2 "nova-factory-server/app/utils/llm/models"

	"charm.land/fantasy"
	"github.com/gin-gonic/gin"
)

func (d *Dataset) ModelRequest(c *gin.Context) {
	req := new(aiDataSetModels.ModelRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		baizeContext.ParameterError(c)
		return
	}
	timeout := req.TimeoutSeconds
	if timeout <= 0 {
		timeout = 60
	}
	ctx, cancel := context.WithTimeout(c, time.Duration(timeout)*time.Second)
	defer cancel()
	providerResult, err := models2.CreateProvider(ctx, &models2.ProviderConfig{
		ModelString:    req.Model,
		SystemPrompt:   req.SystemPrompt,
		ProviderAPIKey: req.ProviderAPIKey,
		ProviderURL:    req.ProviderURL,
		MaxTokens:      req.MaxTokens,
		Temperature:    req.Temperature,
		TopP:           req.TopP,
		TopK:           req.TopK,
	})
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	if providerResult.Closer != nil {
		defer providerResult.Closer.Close()
	}
	prompt := make(fantasy.Prompt, 0, 2)
	if req.SystemPrompt != "" {
		prompt = append(prompt, fantasy.NewSystemMessage(req.SystemPrompt))
	}
	prompt = append(prompt, fantasy.NewUserMessage(req.Prompt))
	resp, err := providerResult.Model.Generate(ctx, fantasy.Call{
		Prompt: prompt,
	})
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	result := &aiDataSetModels.ModelResponse{
		Provider:         providerResult.Model.Provider(),
		Model:            providerResult.Model.Model(),
		Content:          resp.Content.Text(),
		PromptTokens:     int(resp.Usage.InputTokens),
		CompletionTokens: int(resp.Usage.OutputTokens),
		TotalTokens:      int(resp.Usage.TotalTokens),
	}
	baizeContext.SuccessData(c, result)
}
