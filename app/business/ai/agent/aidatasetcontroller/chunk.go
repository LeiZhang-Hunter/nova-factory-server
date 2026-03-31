package aidatasetcontroller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/utils/baizeContext"
)

// ChunkList chunk列表
// @Summary chunk列表
// @Description chunk列表
// @Tags 工业智能体/分块管理
// @Param  object query aidatasetmodels.ChunkListReq true "设备分组列表请求参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /ai/dataset/chunk/list [get]
func (d *Dataset) ChunkList(c *gin.Context) {
	req := new(aidatasetmodels.ChunkListReq)
	err := c.ShouldBind(req)
	if err != nil {
		zap.L().Error("解析错误", zap.Error(err))
		baizeContext.SuccessData(c, &aidatasetmodels.ChunkListResponse{})
		return
	}
	documents, err := d.iChunkService.ChunkList(c, req)
	if err != nil {
		zap.L().Error("chunk列表错误", zap.Error(err))
		baizeContext.SuccessData(c, &aidatasetmodels.ChunkListResponse{})
		return
	}
	baizeContext.SuccessData(c, documents.Data)
}

// AddChunk 添加chunk
// @Summary 添加chunk
// @Description 添加chunk
// @Tags 工业智能体/分块管理
// @Param  object body aidatasetmodels.AddChunkReq true "添加chunk参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /ai/dataset/chunk/add [post]
func (d *Dataset) AddChunk(c *gin.Context) {
	req := new(aidatasetmodels.AddChunkReq)
	err := c.ShouldBindJSON(req)
	if err != nil {
		zap.L().Error("解析错误", zap.Error(err))
		baizeContext.SuccessData(c, &aidatasetmodels.ChunkListResponse{})
		return
	}
	documents, err := d.iChunkService.AddChunk(c, req)
	if err != nil {
		zap.L().Error("chunk列表错误", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, documents)
}

// RemoveChunk 移除chunk
// @Summary 移除chunk
// @Description 移除chunk
// @Tags 工业智能体/分块管理
// @Param  object body aidatasetmodels.RemoveChunkReq true "移除chunk"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "移除chunk成功"
// @Router /ai/dataset/chunk/remove [delete]
func (d *Dataset) RemoveChunk(c *gin.Context) {
	req := new(aidatasetmodels.RemoveChunkReq)
	err := c.ShouldBindJSON(req)
	if err != nil {
		zap.L().Error("解析错误", zap.Error(err))
		baizeContext.SuccessData(c, &aidatasetmodels.ChunkListResponse{})
		return
	}
	err = d.iChunkService.RemoveChunk(c, req)
	if err != nil {
		zap.L().Error("删除chunk", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}

// UpdateChunk 更新chunk
// @Summary 更新chunk
// @Description 更新chunk
// @Tags 工业智能体/分块管理
// @Param  object body aidatasetmodels.UpdateChunkReq true "请求参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /ai/dataset/chunk/update [put]
func (d *Dataset) UpdateChunk(c *gin.Context) {
	req := new(aidatasetmodels.UpdateChunkReq)
	err := c.ShouldBindJSON(req)
	if err != nil {
		zap.L().Error("解析错误", zap.Error(err))
		baizeContext.SuccessData(c, &aidatasetmodels.ChunkListResponse{})
		return
	}
	_, err = d.iChunkService.UpdateChunk(c, req)
	if err != nil {
		zap.L().Error("更新chunk失败", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.Success(c)
}

// RetrievalChunk 检索chunk
// @Summary 检索chunk
// @Description 检索chunk
// @Tags 工业智能体/分块管理
// @Param  object body aidatasetmodels.RetrievalListReq true "检索chunk请求参数"
// @Produce application/json
// @Success 200 {object}  response.ResponseData "获取成功"
// @Router /ai/dataset/chunk/retrieval [post]
func (d *Dataset) RetrievalChunk(c *gin.Context) {
	req := new(aidatasetmodels.RetrievalListReq)
	err := c.ShouldBindJSON(req)
	if err != nil {
		zap.L().Error("解析错误", zap.Error(err))
		baizeContext.Waring(c, "解析错误")
		return
	}
	response, err := d.iChunkService.RetrievalChunk(c, req)
	if err != nil {
		zap.L().Error("检索chunk失败", zap.Error(err))
		baizeContext.Waring(c, err.Error())
		return
	}
	baizeContext.SuccessData(c, response.Data)
}
