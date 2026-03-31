package aidatasetserviceimpl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/business/ai/agent/aidatasetservice"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type IAssistantServiceImpl struct {
	config RagFlowConfig
	client *http.Client
}

func NewIAssistantServiceImpl(client *http.Client) aidatasetservice.IAssistantService {
	var config RagFlowConfig
	err := viper.UnmarshalKey("dataSet", &config)
	if err != nil {
		panic(err)
	}
	return &IAssistantServiceImpl{
		config: config,
		client: client,
	}
}

func (i *IAssistantServiceImpl) AddAssistant(c *gin.Context, req *aidatasetmodels.CreateAssistantRequest) (*aidatasetmodels.CreateAssistantResponse, error) {
	content, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/chats",
		i.config.Host), bytes.NewBuffer(content))
	if err != nil {
		zap.L().Error("添加助理失败", zap.Error(err))
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+i.config.ApiKey)
	resp, err := i.client.Do(request)
	if err != nil {
		zap.L().Error("添加助理失败", zap.Error(err))
		return nil, errors.New("添加助理失败")
	}
	if resp == nil {
		return nil, errors.New("添加助理失败")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("添加助理失败", zap.Error(err))
		return nil, errors.New("添加助理失败")
	}
	var response aidatasetmodels.CreateAssistantResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errors.New("添加助理失败")
	}
	if response.Code != 0 {
		zap.L().Error(fmt.Sprintf("添加助理失败: %s", response.Message))
		return nil, errors.New("添加助理失败")
	}
	return &response, nil
}
func (i *IAssistantServiceImpl) UpdateAssistant(c *gin.Context, req *aidatasetmodels.UpdateAssistantRequest) (*aidatasetmodels.UpdateAssistantResponse, error) {
	chatId := req.ChatId
	req.ChatId = ""
	content, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/chats/%s",
		i.config.Host, chatId), bytes.NewBuffer(content))
	if err != nil {
		zap.L().Error("添加助理失败", zap.Error(err))
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+i.config.ApiKey)
	resp, err := i.client.Do(request)
	if err != nil {
		zap.L().Error("添加助理失败", zap.Error(err))
		return nil, errors.New("添加助理失败")
	}
	if resp == nil {
		return nil, errors.New("添加助理失败")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("添加助理失败", zap.Error(err))
		return nil, errors.New("添加助理失败")
	}
	var response aidatasetmodels.UpdateAssistantResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errors.New("添加助理失败")
	}
	if response.Code != 0 {
		zap.L().Error(fmt.Sprintf("添加助理失败: %s", response.Message))
		return nil, errors.New("添加助理失败")
	}
	return &response, nil
}

func (i *IAssistantServiceImpl) DeleteAssistant(c *gin.Context, ids []string) error {
	var req aidatasetmodels.DeleteAssistantRequest
	req.Ids = ids
	content, err := json.Marshal(req)
	if err != nil {
		zap.L().Error("删除助理失败", zap.Error(err))
		return errors.New("删除助理失败")
	}
	request, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/chats",
		i.config.Host), bytes.NewBuffer(content))
	if err != nil {
		zap.L().Error("删除助理失败", zap.Error(err))
		return errors.New("删除助理失败")
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+i.config.ApiKey)
	resp, err := i.client.Do(request)
	if err != nil {
		zap.L().Error("删除助理失败", zap.Error(err))
		return errors.New("删除助理失败")
	}
	if resp == nil {
		return errors.New("删除助理失败")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("删除助理失败", zap.Error(err))
		return errors.New("删除助理失败")
	}
	var response aidatasetmodels.DeleteAssistantResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return errors.New("删除助理失败")
	}
	if response.Code != 0 {
		zap.L().Error(fmt.Sprintf("删除助理失败: %s", response.Message))
		return errors.New("删除助理失败")
	}
	return nil
}

func (i *IAssistantServiceImpl) ListAssistant(c *gin.Context, req *aidatasetmodels.GetAssistantRequest) (*aidatasetmodels.GetAssistantResponse, error) {
	params := url.Values{}
	if req.Id != "" {
		params.Add("id", req.Id)
	}
	if req.Name != "" {
		params.Add("name", req.Name)
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

	request, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/chats?%s",
		i.config.Host, params.Encode()), nil)
	if err != nil {
		zap.L().Error("读取助理列表失败", zap.Error(err))
		return &aidatasetmodels.GetAssistantResponse{}, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+i.config.ApiKey)
	resp, err := i.client.Do(request)
	if err != nil {
		zap.L().Error("读取助理列表失败", zap.Error(err))
		return &aidatasetmodels.GetAssistantResponse{}, errors.New("读取助理列表失败")
	}
	if resp == nil {
		return &aidatasetmodels.GetAssistantResponse{}, errors.New("读取助理列表失败")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("读取助理列表失败", zap.Error(err))
		return &aidatasetmodels.GetAssistantResponse{}, errors.New("读取助理列表失败")
	}
	var response aidatasetmodels.GetAssistantResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return &aidatasetmodels.GetAssistantResponse{}, errors.New("读取助理列表失败")
	}
	if response.Code != 0 {
		zap.L().Error(fmt.Sprintf("读取助理列表失败: %s", response.Message))
		return &aidatasetmodels.GetAssistantResponse{}, errors.New("读取助理列表失败")
	}
	return &response, nil
}
