package aiDataSetServiceImpl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"nova-factory-server/app/business/ai/aiDataSetModels"
	"nova-factory-server/app/business/ai/aiDataSetService"
	"strconv"
)

type IChunkServiceImpl struct {
	config RagFlowConfig
	client *http.Client
}

func NewIChunkServiceImpl(client *http.Client) aiDataSetService.IChunkService {
	var config RagFlowConfig
	err := viper.UnmarshalKey("dataSet", &config)
	if err != nil {
		panic(err)
	}
	return &IChunkServiceImpl{
		config: config,
		client: client,
	}
}

func (i *IChunkServiceImpl) ChunkList(c *gin.Context, req *aiDataSetModels.ChunkListReq) (*aiDataSetModels.ChunkListResponse, error) {
	params := url.Values{}
	if req.Keywords != "" {
		params.Add("keywords", req.Keywords)
	}
	if req.Page > 0 {
		params.Add("page", strconv.FormatInt(req.Page, 10))
	}
	if req.Size > 0 {
		params.Add("page_size", strconv.FormatInt(req.Size, 10))
	}
	if req.Id != "" {
		params.Add("id", req.Id)
	}
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/datasets/%s/documents/%s/chunks?%s",
		i.config.Host, req.DatasetUuid, req.DocumentUuid, params.Encode()), nil)
	if err != nil {
		zap.L().Error("读取chunk列表失败", zap.Error(err))
		return &aiDataSetModels.ChunkListResponse{}, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+i.config.ApiKey)
	resp, err := i.client.Do(request)
	if err != nil {
		zap.L().Error("读取chunk列表失败", zap.Error(err))
		return &aiDataSetModels.ChunkListResponse{}, errors.New("读取CHUNK列表失败")
	}
	if resp == nil {
		return &aiDataSetModels.ChunkListResponse{}, errors.New("读取CHUNK列表失败")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("读取CHUNK列表失败", zap.Error(err))
		return &aiDataSetModels.ChunkListResponse{}, errors.New("读取CHUNK列表失败")
	}
	var response aiDataSetModels.ChunkListResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		zap.L().Error("读取chunk列表失败", zap.Error(err))
		return &aiDataSetModels.ChunkListResponse{}, errors.New("读取chunk列表失败")
	}
	if response.Code != 0 {
		zap.L().Error(fmt.Sprintf("读取CHUNK列表失败: %s", response.Message))
		return &aiDataSetModels.ChunkListResponse{}, errors.New("读取chunk列表失败")
	}
	return &response, nil
}

func (i *IChunkServiceImpl) AddChunk(c *gin.Context, req *aiDataSetModels.AddChunkReq) (*aiDataSetModels.AddChunkResponse, error) {
	if req.ImportantKeywords == nil {
		req.ImportantKeywords = []string{}
	}
	if req.Questions == nil {
		req.Questions = []string{}
	}
	content, err := json.Marshal(req)
	if err != nil {
		zap.L().Error("读取chunk列表失败", zap.Error(err))
		return &aiDataSetModels.AddChunkResponse{}, errors.New("添加chunk失败")
	}
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/datasets/%s/documents/%s/chunks",
		i.config.Host, req.DatasetUuid, req.DocumentUuid), bytes.NewBuffer(content))
	if err != nil {
		zap.L().Error("读取chunk列表失败", zap.Error(err))
		return &aiDataSetModels.AddChunkResponse{}, errors.New("添加chunk失败")
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+i.config.ApiKey)
	resp, err := i.client.Do(request)
	if err != nil {
		zap.L().Error("添加chunk失败", zap.Error(err))
		return &aiDataSetModels.AddChunkResponse{}, errors.New("添加chunk失败")
	}
	if resp == nil {
		return &aiDataSetModels.AddChunkResponse{}, errors.New("添加chunk失败")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("添加chunk失败", zap.Error(err))
		return &aiDataSetModels.AddChunkResponse{}, errors.New("添加chunk失败")
	}
	var response aiDataSetModels.AddChunkResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return &aiDataSetModels.AddChunkResponse{}, errors.New("添加chunk失败")
	}
	if response.Code != 0 {
		zap.L().Error(fmt.Sprintf("添加chunk失败: %s", response.Message))
		return &aiDataSetModels.AddChunkResponse{}, errors.New("添加chunk失败")
	}
	return &response, nil
}

func (i *IChunkServiceImpl) RemoveChunk(c *gin.Context, req *aiDataSetModels.RemoveChunkReq) error {
	if len(req.ChunkIds) == 0 {
		return errors.New("请选择分片")
	}

	content, err := json.Marshal(req)
	if err != nil {
		zap.L().Error("读取chunk列表失败", zap.Error(err))
		return errors.New("删除chunk失败")
	}
	request, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/datasets/%s/documents/%s/chunks",
		i.config.Host, req.DatasetUuid, req.DocumentUuid), bytes.NewBuffer(content))
	if err != nil {
		zap.L().Error("读取chunk列表失败", zap.Error(err))
		return errors.New("删除chunk失败")
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+i.config.ApiKey)
	resp, err := i.client.Do(request)
	if err != nil {
		zap.L().Error("删除chunk失败", zap.Error(err))
		return errors.New("删除chunk失败")
	}
	if resp == nil {
		return errors.New("删除chunk失败")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("删除chunk失败", zap.Error(err))
		return errors.New("删除chunk失败")
	}
	var response aiDataSetModels.AddChunkResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return errors.New("删除chunk失败")
	}
	if response.Code != 0 {
		zap.L().Error(fmt.Sprintf("删除chunk失败: %s", response.Message))
		return errors.New("删除chunk失败")
	}
	return nil
}

func (i *IChunkServiceImpl) UpdateChunk(c *gin.Context, req *aiDataSetModels.UpdateChunkReq) (*aiDataSetModels.UpdateChunkResponse, error) {
	if req.ChunkUuid == "" {
		return nil, errors.New("更新chunk失败")
	}

	if req.DatasetUuid == "" {
		return nil, errors.New("更新chunk失败")
	}

	if req.DocumentUuid == "" {
		return nil, errors.New("更新chunk失败")
	}

	if req.ImportantKeywords == nil {
		req.ImportantKeywords = []string{}
	}
	if req.Questions == nil {
		req.Questions = []string{}
	}

	content, err := json.Marshal(req)
	if err != nil {
		zap.L().Error("更新chunk失败", zap.Error(err))
		return nil, errors.New("更新chunk失败")
	}
	request, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/datasets/%s/documents/%s/chunks/%s",
		i.config.Host, req.DatasetUuid, req.DocumentUuid, req.ChunkUuid), bytes.NewBuffer(content))
	if err != nil {
		zap.L().Error("更新chunk失败", zap.Error(err))
		return nil, errors.New("更新chunk失败")
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+i.config.ApiKey)
	resp, err := i.client.Do(request)
	if err != nil {
		zap.L().Error("更新chunk失败", zap.Error(err))
		return nil, errors.New("更新chunk失败")
	}
	if resp == nil {
		return nil, errors.New("更新chunk失败")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("更新chunk失败", zap.Error(err))
		return nil, errors.New("更新chunk失败")
	}
	var response aiDataSetModels.UpdateChunkResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errors.New("更新chunk失败")
	}
	if response.Code != 0 {
		zap.L().Error(fmt.Sprintf("更新chunk失败: %s", response.Message))
		return nil, errors.New("更新chunk失败")
	}
	return &response, nil
}

func (i *IChunkServiceImpl) RetrievalChunk(c *gin.Context, req *aiDataSetModels.RetrievalListReq) (*aiDataSetModels.RetrievalApiListResponse, error) {
	var apiReq aiDataSetModels.RetrievalApiListReq
	apiReq.Of(req)
	content, err := json.Marshal(apiReq)
	if err != nil {
		zap.L().Error("检索chunk失败", zap.Error(err))
		return nil, errors.New("检索chunk失败")
	}
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/retrieval",
		i.config.Host), bytes.NewBuffer(content))
	if err != nil {
		zap.L().Error("检索chunk失败", zap.Error(err))
		return nil, errors.New("检索chunk失败")
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+i.config.ApiKey)
	resp, err := i.client.Do(request)
	if err != nil {
		zap.L().Error("检索chunk失败", zap.Error(err))
		return nil, errors.New("检索chunk失败")
	}
	if resp == nil {
		return nil, errors.New("检索chunk失败")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("检索chunk失败", zap.Error(err))
		return nil, errors.New("检索chunk失败")
	}
	var response aiDataSetModels.RetrievalApiListResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		zap.L().Error("检索chunk失败", zap.Error(err))
		return nil, errors.New("检索chunk失败")
	}
	if response.Code != 0 {
		zap.L().Error(fmt.Sprintf("检索chunk失败: %s", response.Message))
		return nil, errors.New("检索chunk失败")
	}
	return &response, nil
}
