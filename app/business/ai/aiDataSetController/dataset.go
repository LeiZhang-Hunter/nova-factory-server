package aiDataSetController

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/ai/aiDataSetModels"
	"nova-factory-server/app/business/ai/aiDataSetService"
	"nova-factory-server/app/middlewares"
	"nova-factory-server/app/utils/baizeContext"
)

type Dataset struct {
	aiService               aiDataSetService.IDataSetService
	iDataSetDocumentService aiDataSetService.IDataSetDocumentService
	iChunkService           aiDataSetService.IChunkService
	assistantService        aiDataSetService.IAssistantService
	chartsService           aiDataSetService.IChartService
}

func NewDataset(aiService aiDataSetService.IDataSetService,
	iDataSetDocumentService aiDataSetService.IDataSetDocumentService,
	iChunkService aiDataSetService.IChunkService, assistantService aiDataSetService.IAssistantService,
	chartsService aiDataSetService.IChartService) *Dataset {
	return &Dataset{
		aiService:               aiService,
		iDataSetDocumentService: iDataSetDocumentService,
		iChunkService:           iChunkService,
		assistantService:        assistantService,
		chartsService:           chartsService,
	}
}

func (d *Dataset) PrivateRoutes(router *gin.RouterGroup) {
	ai := router.Group("/ai/dataset")
	ai.GET("/list", middlewares.HasPermission("ai:dataset"), d.GetAiDataSet)                          // 知识库列表
	ai.GET("/info", middlewares.HasPermission("ai:dataset:info"), d.GetAiDataGetInfo)                 // 知识库列表
	ai.GET("/config", middlewares.HasPermission("ai:dataset:info"), d.GetAiDataGetConfig)             // 知识库列表
	ai.POST("/create", middlewares.HasPermission("ai:dataset:create"), d.CreateAiDataSet)             // 创建知识库
	ai.PUT("/update/:dataset_id", middlewares.HasPermission("ai:dataset:update"), d.UpdateAiDataSet)  // 更新知识库
	ai.DELETE("/remove/:dataset_id", middlewares.HasPermission("ai:dataset:remove"), d.RemoveDataSet) // 删除知识库集合

	// 知识库文档
	ai.POST("/upload/document/:datasetId", middlewares.HasPermission("ai:dataset:upload:document"),
		d.UploadDocument) // 上传知识库文档
	ai.PUT("/update/document/:document_id", middlewares.HasPermission("ai:dataset:update:document"),
		d.PutDocument) // 更新知识库文档
	ai.GET("/download/document/:document_id", middlewares.HasPermission("ai:dataset:update:document"),
		d.DownloadDocument) // 更新知识库文档
	ai.GET("/list/document", middlewares.HasPermission("ai:dataset:update:document"),
		d.ListDocument) // 更新知识库文档
	ai.DELETE("/remove/document", middlewares.HasPermission("ai:dataset:remove:document"), d.RemoveDocument) //删除文档列表
	ai.POST("/start/document", middlewares.HasPermission("ai:dataset:start:document"), d.StartParseDocument) //解析文档
	ai.POST("/stop/document", middlewares.HasPermission("ai:dataset:stop:document"), d.StopParseDocument)    //停止解析文档

	// chunk 管理
	ai.GET("/chunk/list", middlewares.HasPermission("ai:dataset:chunk"), d.ChunkList)
	ai.POST("/chunk/add", middlewares.HasPermission("ai:dataset:chunk:add"), d.AddChunk)
	ai.DELETE("/chunk/remove", middlewares.HasPermission("ai:dataset:chunk:remove"), d.RemoveChunk)
	ai.PUT("/chunk/update", middlewares.HasPermission("ai:dataset:chunk:update"), d.UpdateChunk)
	ai.POST("/chunk/retrieval", middlewares.HasPermission("ai:dataset:chunk:retrieval"), d.RetrievalChunk)

	// assistant 助理代理
	ai.GET("/assistant/list", middlewares.HasPermission("ai:dataset:assistant"), d.ListAssistant)
	ai.PUT("/assistant/update", middlewares.HasPermission("ai:dataset:assistant:update"), d.UpdateAssistant)
	ai.POST("/assistant/create", middlewares.HasPermission("ai:dataset:chats:assistant:create"), d.CreateAssistant)
	ai.DELETE("/assistant/remove/:assistantIds", middlewares.HasPermission("ai:dataset:chats:assistant:remove"), d.RemoveAssistant)

	// chats 聊天
	// 使用chat助手创建会话
	ai.POST("/session/create", middlewares.HasPermission("ai:dataset:session:create"), d.SessionCreate)
	// 更新聊天助手会话
	ai.POST("/session/update", middlewares.HasPermission("ai:dataset:session:update"), d.SessionUpdate)
	// 列出聊天助手的会话
	ai.GET("/session/list", middlewares.HasPermission("ai:dataset:session:list"), d.SessionList)
	// 删除聊天助手的会话
	ai.DELETE("/session/remove", middlewares.HasPermission("ai:dataset:session:remove"), d.SessionRemove)
	// 与聊天助手交谈
	ai.POST("/session/charts/completions", middlewares.HasPermission("ai:dataset:session:charts:completions"), d.ChartsCompletions)
	// 使用 agent 创建会话
	ai.GET("/session/agents/create", middlewares.HasPermission("ai:dataset:session:agents:create"), d.AgentSessionCreate)
	// 与代理结合
	ai.GET("/session/agents/completions", middlewares.HasPermission("ai:dataset:session:agents:completions"), d.AgentCompletions)
	// list 代理会话
	ai.GET("/session/agents/list", middlewares.HasPermission("ai:dataset:session:agents:list"), d.AgentsSessionList)
	// 删除代理会话
	ai.DELETE("/session/agents/delete", middlewares.HasPermission("ai:dataset:session:agents:delete"), d.AgentsSessionRemove)
	// 相关提问
	ai.POST("/session/conversation/related_questions", middlewares.HasPermission("ai:dataset:session:conversation:related_questions"), d.ConversationRelatedQuestions)
	// 智能问答
	ai.POST("/session/ask", middlewares.HasPermission("ai:dataset:session:ask"), d.Ask)
	// 清单代理
	ai.GET("/agents/list", middlewares.HasPermission("ai:dataset:agents:list"), d.AgentList)
}

// GetAiDataSet 读取数据集列表
// @Summary 读取数据集列表
// @Description 读取数据集列表
// @Tags 工业智能体/知识库管理
// @Param  object query aiDataSetModels.DatasetListReq true "设备分组列表请求参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /ai/dataset/list [get]
func (d *Dataset) GetAiDataSet(c *gin.Context) {
	req := new(aiDataSetModels.DatasetListReq)
	err := c.ShouldBind(req)
	if err != nil {
		zap.L().Error("解析错误", zap.Error(err))
		baizeContext.SuccessData(c, &aiDataSetModels.DatasetListReq{})
		return
	}

	list, err := d.aiService.SelectDataSet(c, req)
	if err != nil {
		return
	}
	baizeContext.SuccessData(c, list)
}

// CreateAiDataSet 添加数据集
// @Summary 添加数据集
// @Description 添加数据集
// @Tags 工业智能体/知识库管理
// @Param  object body aiDataSetModels.DataSetRequest true "设备分组参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /ai/dataset/create [post]
func (d *Dataset) CreateAiDataSet(c *gin.Context) {
	req := new(aiDataSetModels.DataSetRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	set, err := d.aiService.CreateDataSet(c, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, set)
}

// RemoveDataSet 移除数据集
// @Summary 移除数据集
// @Description 移除数据集
// @Tags 工业智能体/知识库管理
// @Param  dataset_id path int64 true "dataset_id"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /ai/dataset/remove [delete]
func (d *Dataset) RemoveDataSet(c *gin.Context) {
	datasetId := baizeContext.ParamInt64(c, "dataset_id")
	if datasetId == 0 {
		baizeContext.Waring(c, "请输入知识库id")
		return
	}

	err := d.aiService.DeleteDataSet(c, datasetId)
	if err != nil {
		baizeContext.Waring(c, "删除失败")
		return
	}
	baizeContext.Success(c)
}

// UpdateAiDataSet 更新数据集
// @Summary 更新数据集
// @Description 更新数据集
// @Tags 工业智能体/知识库管理
// @Param  dataset_id path string true "dataset_id"
// @Param  object body aiDataSetModels.UpdateDataSetRequest true "设备分组参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /ai/dataset/update/{dataset_id} [put]
func (d *Dataset) UpdateAiDataSet(c *gin.Context) {
	req := new(aiDataSetModels.UpdateDataSetRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	datasetId := baizeContext.ParamInt64(c, "dataset_id")
	if datasetId == 0 {
		baizeContext.Waring(c, "请输入知识库id")
		return
	}

	set, err := d.aiService.UpdateDataSet(c, req, datasetId)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, set)
}

// GetAiDataGetInfo 读取数据集
// @Summary 读取数据集
// @Description 读取数据集
// @Tags 工业智能体/知识库管理
// @Param  object query aiDataSetModels.GetDatasetInfoReq true "读取数据集"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /ai/dataset/info [get]
func (d *Dataset) GetAiDataGetInfo(c *gin.Context) {
	req := new(aiDataSetModels.GetDatasetInfoReq)
	err := c.ShouldBind(req)
	if err != nil {
		zap.L().Error("解析错误", zap.Error(err))
		baizeContext.SuccessData(c, &aiDataSetModels.DatasetListReq{})
		return
	}

	list, err := d.aiService.GetInfoById(c, req.Id)
	if err != nil {
		return
	}
	baizeContext.SuccessData(c, list)
}

// GetAiDataGetConfig 读取知识库配置
// @Summary 读取知识库配置
// @Description 读取知识库配置
// @Tags 工业智能体/知识库管理
// @Produce application/json
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /ai/dataset/config [get]
func (d *Dataset) GetAiDataGetConfig(c *gin.Context) {
	baizeContext.SuccessData(c, aiDataSetModels.NewDataSetConfig())
}
