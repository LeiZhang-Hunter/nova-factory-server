package aidatasetcontroller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io"
	"net/http"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/utils/baizeContext"
	"strings"
	"time"
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

// GetRagFlowDocumentPreview 获取 RagFlow 预览文件
// @Summary 获取 RagFlow 预览文件
// @Description 携带 RagFlow API Key 拉取预览资源并透传给客户端
// @Tags 工业智能体/文档管理
// @Param doc_id path string true "RagFlow 文档预览资源ID"
// @Produce application/octet-stream
// @Success 200 {file} file "预览文件"
// @Router /ai/dataset/document/get/{doc_id} [get]
func (d *Dataset) GetRagFlowDocumentPreview(c *gin.Context) {
	docID := strings.TrimSpace(c.Param("doc_id"))
	if docID == "" {
		baizeContext.ParameterError(c)
		return
	}
	apiKey := strings.TrimSpace(viper.GetString("dataSet.api_key"))
	if apiKey == "" {
		apiKey = strings.TrimSpace(viper.GetString("dataset.api_key"))
	}
	if apiKey == "" {
		baizeContext.Waring(c, "ragflow api key未配置")
		return
	}
	baseURL := strings.TrimSpace(viper.GetString("dataSet.image_url"))
	if baseURL == "" {
		baseURL = strings.TrimSpace(viper.GetString("dataset.image_url"))
	}
	if baseURL == "" {
		baseURL = strings.TrimSpace(viper.GetString("dataSet.host"))
	}
	if baseURL == "" {
		baseURL = strings.TrimSpace(viper.GetString("dataset.host"))
	}
	if baseURL == "" {
		baizeContext.Waring(c, "ragflow地址未配置")
		return
	}

	previewURL := fmt.Sprintf("%s/v1/document/get/%s", strings.TrimRight(baseURL, "/"), docID)
	req, err := http.NewRequestWithContext(c, http.MethodGet, previewURL, nil)
	if err != nil {
		zap.L().Error("create ragflow request failed", zap.Error(err))
		baizeContext.Waring(c, "请求失败")
		return
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)

	httpClient := &http.Client{Timeout: 120 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		zap.L().Error("ragflow request failed", zap.Error(err))
		baizeContext.Waring(c, "请求失败")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		msg := strings.TrimSpace(string(body))
		if msg == "" {
			msg = "预览文件获取失败"
		}
		baizeContext.Waring(c, msg)
		return
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	c.Writer.Header().Set("Content-Type", contentType)
	if cd := resp.Header.Get("Content-Disposition"); cd != "" {
		c.Writer.Header().Set("Content-Disposition", cd)
	}
	if cc := resp.Header.Get("Cache-Control"); cc != "" {
		c.Writer.Header().Set("Cache-Control", cc)
	}
	if lm := resp.Header.Get("Last-Modified"); lm != "" {
		c.Writer.Header().Set("Last-Modified", lm)
	}
	if etag := resp.Header.Get("ETag"); etag != "" {
		c.Writer.Header().Set("ETag", etag)
	}
	c.Status(http.StatusOK)
	_, _ = io.Copy(c.Writer, resp.Body)
}
