package aiDataSetController

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/aiDataSetModels"
	"nova-factory-server/app/utils/baizeContext"
)

// SessionCreate 创建助理会话
// @Summary 创建助理会话
// @Description 创建助理会话
// @Tags 工业智能体/会话管理
// @Param  object body aiDataSetModels.CreateSessionsRequest true "使用 chat 助手创建会话"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /ai/dataset/session/create [post]
func (d *Dataset) SessionCreate(c *gin.Context) {
	req := new(aiDataSetModels.CreateSessionsRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	response, err := d.chartsService.SessionCreate(c, req)
	if err != nil {
		baizeContext.Waring(c, "助手创建会话失败")
		return
	}
	baizeContext.SuccessData(c, response.Data)
}

// SessionUpdate 更新助理会话
// @Summary 更新助理会话
// @Description 更新助理会话
// @Tags 工业智能体/会话管理
// @Param  object body aiDataSetModels.UpdateSessionsRequest true "更新助理会话"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /ai/dataset/session/update [post]
func (d *Dataset) SessionUpdate(c *gin.Context) {
	req := new(aiDataSetModels.UpdateSessionsRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	_, err = d.chartsService.SessionUpdate(c, req)
	if err != nil {
		baizeContext.Waring(c, "更新会话失败")
		return
	}
	baizeContext.Success(c)
}

// SessionList 助理会话列表
// @Summary 助理会话列表
// @Description 助理会话列表
// @Tags 工业智能体/会话管理
// @Param  object query aiDataSetModels.ListSessionRequest true "助理会话列表"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /ai/dataset/session/list [get]
func (d *Dataset) SessionList(c *gin.Context) {
	req := new(aiDataSetModels.ListSessionRequest)
	err := c.ShouldBind(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	if req.ChatId == "" {
		baizeContext.SuccessData(c, &aiDataSetModels.ListSessionResponse{})
		return
	}
	list, err := d.chartsService.SessionList(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, list.Data)
}

// SessionRemove 删除会话
// @Summary 删除会话
// @Description 删除会话
// @Tags 工业智能体/会话管理
// @Param  object body aiDataSetModels.DeleteSessionRequest true "删除会话"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /ai/dataset/session/remove [delete]
func (d *Dataset) SessionRemove(c *gin.Context) {
	req := new(aiDataSetModels.DeleteSessionRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	_, err = d.chartsService.SessionRemove(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}

// ChartsCompletions 与聊天助手交谈
// @Summary 与聊天助手交谈
// @Description 与聊天助手交谈
// @Tags 工业智能体/会话管理
// @Param  object body aiDataSetModels.ChatsCompletionsRequest true "与聊天助手交谈"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /ai/dataset/session/charts/completions [post]
func (d *Dataset) ChartsCompletions(c *gin.Context) {
	/**
	{
	    "chat_id": "68afa28e2c1b11f0be8e0242ac1b0006",
	    "question": "你是谁11",
	    "session_id": "66b449bd73b74ab490eb8df6011382ae",
	    "stream": false
	}
	*/
	req := new(aiDataSetModels.ChatsCompletionsRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	completions, err := d.chartsService.ChatsCompletions(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	if req.Stream == false {
		baizeContext.SuccessData(c, completions.Data)
	}
	return
}

// AgentSessionCreate 创建agent会话
// @Summary 创建agent会话
// @Description 创建agent会话
// @Tags 工业智能体/会话管理
// @Param  object body aiDataSetModels.SessionAgentCreate true "创建agent会话"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /ai/dataset/session/agents/create [post]
func (d *Dataset) AgentSessionCreate(c *gin.Context) {
	req := new(aiDataSetModels.SessionAgentCreate)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	response, err := d.chartsService.AgentSessionCreate(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, response)
	return
}

// AgentCompletions Agent聊天
// @Summary Agent聊天
// @Description Agent聊天
// @Tags 工业智能体/会话管理
// @Param  object body aiDataSetModels.AgentsCompletionsRequest true "Agent聊天"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /ai/dataset/session/agents/completions [post]
func (d *Dataset) AgentCompletions(c *gin.Context) {
	req := new(aiDataSetModels.AgentsCompletionsRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	response, err := d.chartsService.AgentsCompletions(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	if req.Stream {
		baizeContext.SuccessData(c, response)
	}
}

// AgentsSessionList Agent会话列表
// @Summary Agent会话列表
// @Description Agent会话列表
// @Tags 工业智能体/会话管理
// @Param  object query aiDataSetModels.ListAgentSessionsRequest true "Agent会话列表请求参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /ai/dataset/session/agents/list [get]
func (d *Dataset) AgentsSessionList(c *gin.Context) {
	req := new(aiDataSetModels.ListAgentSessionsRequest)
	err := c.ShouldBindQuery(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	resp, err := d.chartsService.AgentSessionList(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, resp.Data)
}

// AgentsSessionRemove 删除Agent会话
// @Summary 删除Agent会话
// @Description 删除Agent会话
// @Tags 工业智能体/会话管理
// @Param  object body aiDataSetModels.RemoveAgentSessionsRequest true "删除Agent会话"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /ai/dataset/session/agents/delete [delete]
func (d *Dataset) AgentsSessionRemove(c *gin.Context) {
	req := new(aiDataSetModels.RemoveAgentSessionsRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	_, err = d.chartsService.AgentSessionRemove(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}

// ConversationRelatedQuestions 相关提问
// @Summary 相关提问
// @Description 相关提问
// @Tags 工业智能体/会话管理
// @Param  object body aiDataSetModels.ConversationRelatedQuestionsRequest true "相关提问"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /ai/dataset/session/conversation/related_questions [post]
func (d *Dataset) ConversationRelatedQuestions(c *gin.Context) {
	req := new(aiDataSetModels.ConversationRelatedQuestionsRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	questions, err := d.chartsService.ConversationRelatedQuestions(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, questions.Data)
}

// Ask 智能问答
// @Summary 智能问答
// @Description 智能问答
// @Tags 工业智能体/会话管理
// @Param  object body aiDataSetModels.AskRequest true "相关提问"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /ai/dataset/session/ask [post]
func (d *Dataset) Ask(c *gin.Context) {
	req := new(aiDataSetModels.AskRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	err = d.chartsService.Ask(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
}

// AgentList Agen列表
// @Summary Agen列表
// @Description Agen列表
// @Tags 工业智能体/会话管理
// @Param  object query aiDataSetModels.ListAgentSessionsRequest true "Agen话列表请求参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /ai/dataset/agents/list [get]
func (d *Dataset) AgentList(c *gin.Context) {
	req := new(aiDataSetModels.ListAgentRequest)
	err := c.ShouldBindQuery(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	resp, err := d.chartsService.AgentList(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, resp.Data)
}
