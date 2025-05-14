package aiDataSetController

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/ai/aiDataSetModels"
	"nova-factory-server/app/constant/sessionStatus"
	"nova-factory-server/app/utils/baizeContext"
)

// UploadDocument 上传文档
// @Summary 上传文档
// @Description 上传文档
// @Tags 工业智能体/文档管理
// @Security BearerAuth
// @Accept multipart/form-data
// @Param  datasetId path int64 true "datasetId"
// @Param file formData file true "files"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /ai/dataset/upload/document/{datasetId} [post]
func (d *Dataset) UploadDocument(c *gin.Context) {
	// 解析multipart表单
	form, err := c.MultipartForm()
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}
	files := form.File["file"] // 获取名为"files"的所有文件
	if len(files) == 0 {
		baizeContext.Waring(c, "请选择要上传的文件")
		return
	}
	datasetId := baizeContext.ParamInt64(c, "datasetId")

	documents, err := d.iDataSetDocumentService.UploadFile(c, datasetId)
	if err != nil {
		baizeContext.SuccessData(c, documents)
		c.Set(sessionStatus.MsgKey, err.Error())
		return
	}
	d.aiService.SyncDataSet(c, datasetId)
	baizeContext.SuccessData(c, documents)
}

// PutDocument 更新文档
// @Summary 更新文档
// @Description 更新文档
// @Tags 工业智能体/文档管理
// @Security BearerAuth
// @Param  document_id path int64 true "document_id"
// @Param  object body aiDataSetModels.PutDocumentRequest true "更新文档请求"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /ai/dataset/update/document/{document_id} [put]
func (d *Dataset) PutDocument(c *gin.Context) {
	req := new(aiDataSetModels.PutDocumentRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}

	documentId := baizeContext.ParamInt64(c, "document_id")
	if documentId == 0 {
		baizeContext.Waring(c, "请输入文档id")
		return
	}

	documents, err := d.iDataSetDocumentService.PutFile(c, documentId, req)
	if err != nil {
		baizeContext.Waring(c, err.Error())
		c.Set(sessionStatus.MsgKey, err.Error())
		return
	}
	baizeContext.SuccessData(c, documents)
}

// DownloadDocument 下载文档
// @Summary 下载文档
// @Description 下载文档
// @Tags 工业智能体/文档管理
// @Security BearerAuth
// @Param  document_id path int64 true "document_id"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /ai/dataset/download/document/{document_id} [get]
func (d *Dataset) DownloadDocument(c *gin.Context) {
	documentId := baizeContext.ParamInt64(c, "document_id")
	d.iDataSetDocumentService.DownloadFile(c, documentId)
}

// ListDocument 文档列表
// @Summary 文档列表
// @Description 文档列表
// @Tags 工业智能体/文档管理
// @Security BearerAuth
// @Param  object query aiDataSetModels.ListDocumentRequest true "文档列表参数"
// @Success 200 {object}  response.ResponseData "设置分组成功"
// @Router /ai/dataset/list/document [get]
func (d *Dataset) ListDocument(c *gin.Context) {
	req := new(aiDataSetModels.ListDocumentRequest)
	err := c.ShouldBindQuery(req)
	if err != nil {
		zap.L().Error("解析错误", zap.Error(err))
		baizeContext.SuccessData(c, &aiDataSetModels.ListDocumentData{})
		return
	}
	documents, err := d.iDataSetDocumentService.ListDocument(c, req.DatasetId, req)
	if err != nil {
		zap.L().Error("文档列表错误", zap.Error(err))
		baizeContext.SuccessData(c, &aiDataSetModels.ListDocumentData{})
		return
	}
	baizeContext.SuccessData(c, documents)
}

// RemoveDocument 删除文档
// @Summary 删除文档
// @Description 删除文档
// @Tags 工业智能体/文档管理
// @Security BearerAuth
// @Param  object body aiDataSetModels.DeleteDocumentRequest true "删除文档请求参数"
// @Success 200 {object}  response.ResponseData "删除文档成功"
// @Router /ai/dataset/remove/document [delete]
func (d *Dataset) RemoveDocument(c *gin.Context) {

	req := new(aiDataSetModels.DeleteDocumentRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}

	if req.DatasetId == 0 {
		baizeContext.Waring(c, "请选择知识库")
		return
	}

	if len(req.DocumentIds) == 0 {
		baizeContext.Waring(c, "请选择文档")
		return
	}

	err = d.iDataSetDocumentService.RemoveDocument(c, req)
	if err != nil {
		zap.L().Error("文档列表错误", zap.Error(err))
		baizeContext.Waring(c, "删除失败")
		return
	}

	// 同步文档库
	d.aiService.SyncDataSet(c, req.DatasetId)
	baizeContext.Success(c)
}

// StartParseDocument 解析文档
// @Summary 解析文档
// @Description 解析文档
// @Tags 工业智能体/文档管理
// @Security BearerAuth
// @Param  object body aiDataSetModels.ParseDocumentApiRequest true "解析文档请求参数"
// @Success 200 {object}  response.ResponseData "解析文档成功"
// @Router /ai/dataset/start/document [post]
func (d *Dataset) StartParseDocument(c *gin.Context) {

	req := new(aiDataSetModels.ParseDocumentApiRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}

	if len(req.DocumentIds) == 0 {
		baizeContext.Waring(c, "请选择文档")
		return
	}

	err = d.iDataSetDocumentService.StartParse(c, req)
	if err != nil {
		zap.L().Error("文档列表错误", zap.Error(err))
		baizeContext.Waring(c, "解析文档失败")
		return
	}

	baizeContext.Success(c)
}

// StopParseDocument 停止解析文档
// @Summary 停止解析文档
// @Description 停止解析文档
// @Tags 工业智能体/文档管理
// @Security BearerAuth
// @Param  object body aiDataSetModels.ParseDocumentApiRequest true "停止解析文档参数"
// @Success 200 {object}  response.ResponseData "停止解析文档成功"
// @Router /ai/dataset/stop/document [post]
func (d *Dataset) StopParseDocument(c *gin.Context) {

	req := new(aiDataSetModels.ParseDocumentApiRequest)
	err := c.ShouldBindJSON(req)
	if err != nil {
		baizeContext.ParameterError(c)
		return
	}

	if req.DatasetId == 0 {
		baizeContext.Waring(c, "请选择知识库")
		return
	}

	if len(req.DocumentIds) == 0 {
		baizeContext.Waring(c, "请选择文档")
		return
	}

	err = d.iDataSetDocumentService.StopParse(c, req)
	if err != nil {
		zap.L().Error("文档列表错误", zap.Error(err))
		baizeContext.Waring(c, "停止解析失败")
		return
	}

	// 同步文档库
	d.aiService.SyncDataSet(c, req.DatasetId)
	baizeContext.Success(c)
}
