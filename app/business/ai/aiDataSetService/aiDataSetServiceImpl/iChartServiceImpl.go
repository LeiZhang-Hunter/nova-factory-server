package aiDataSetServiceImpl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/tmaxmax/go-sse"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"nova-factory-server/app/business/ai/aiDataSetModels"
	"nova-factory-server/app/business/ai/aiDataSetService"
	"strconv"
	"time"
)

type IChartServiceImpl struct {
	config RagFlowConfig
	client *http.Client
}

func NewIChartServiceImpl(client *http.Client) aiDataSetService.IChartService {
	var config RagFlowConfig
	err := viper.UnmarshalKey("dataSet", &config)
	if err != nil {
		panic(err)
	}
	return &IChartServiceImpl{
		config: config,
		client: client,
	}
}

func (i *IChartServiceImpl) SessionCreate(c *gin.Context, req *aiDataSetModels.CreateSessionsRequest) (*aiDataSetModels.CreateSessionsResponse, error) {
	content, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/chats/%s/sessions",
		i.config.Host, req.ChatId), bytes.NewBuffer(content))
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
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("添加助理失败", zap.Error(err))
		return nil, errors.New("添加助理失败")
	}
	var response aiDataSetModels.CreateSessionsResponse
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

func (i *IChartServiceImpl) SessionUpdate(c *gin.Context, req *aiDataSetModels.UpdateSessionsRequest) (*aiDataSetModels.UpdateSessionsResponse, error) {
	var data aiDataSetModels.UpdateSessionsApiRequest
	data.Name = req.Name
	data.UserId = req.UserId
	content, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/chats/%s/sessions/%s",
		i.config.Host, req.ChatId, req.SessionId), bytes.NewBuffer(content))
	if err != nil {
		zap.L().Error("更新助理失败", zap.Error(err))
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+i.config.ApiKey)
	resp, err := i.client.Do(request)
	if err != nil {
		zap.L().Error("更新助理失败", zap.Error(err))
		return nil, errors.New("更新助理失败")
	}
	defer resp.Body.Close()
	if resp == nil {
		return nil, errors.New("更新助理失败")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("更新助理失败", zap.Error(err))
		return nil, errors.New("更新助理失败")
	}
	var response aiDataSetModels.UpdateSessionsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errors.New("更新助理失败")
	}
	if response.Code != 0 {
		zap.L().Error(fmt.Sprintf("更新助理失败: %s", response.Message))
		return nil, errors.New("更新助理失败")
	}
	return &response, nil
}

func (i *IChartServiceImpl) SessionList(c *gin.Context, req *aiDataSetModels.ListSessionRequest) (*aiDataSetModels.ListSessionResponse, error) {
	params := url.Values{}
	if req.Id != "" {
		params.Add("id", req.Id)
	}
	if req.Name != "" {
		params.Add("name", req.Name)
	}
	if req.Page > 0 {
		params.Add("page", strconv.FormatInt(req.Page, 10))
	} else {
		params.Add("page", "1")
	}
	if req.Size > 0 {
		params.Add("page_size", strconv.FormatInt(req.Size, 10))
	} else {
		params.Add("page_size", "10")
	}

	request, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/chats/%s/sessions?%s",
		i.config.Host, req.ChatId, params.Encode()), nil)
	if err != nil {
		zap.L().Error("读取会话列表失败", zap.Error(err))
		return &aiDataSetModels.ListSessionResponse{}, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+i.config.ApiKey)
	resp, err := i.client.Do(request)
	if err != nil {
		zap.L().Error("读取会话列表失败", zap.Error(err))
		return &aiDataSetModels.ListSessionResponse{}, errors.New("读取会话列表失败")
	}
	if resp == nil {
		return &aiDataSetModels.ListSessionResponse{}, errors.New("读取会话列表失败")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("读取会话列表失败", zap.Error(err))
		return &aiDataSetModels.ListSessionResponse{}, errors.New("读取会话列表失败")
	}
	var response aiDataSetModels.ListSessionResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return &aiDataSetModels.ListSessionResponse{}, errors.New("读取会话列表失败")
	}
	if response.Code != 0 {
		zap.L().Error(fmt.Sprintf("读取会话列表失败: %s", response.Message))
		return &aiDataSetModels.ListSessionResponse{}, errors.New("读取会话列表失败")
	}
	for dk, d := range response.Data {
		for mk, m := range d.Messages {
			for rk, r := range m.Reference {
				response.Data[dk].Messages[mk].Reference[rk].DocUrl = fmt.Sprintf("%s/v1/document/get/%s", i.config.ImageUrl, r.DocumentId)
			}
		}
	}
	return &response, nil
}

func (i *IChartServiceImpl) SessionRemove(c *gin.Context, req *aiDataSetModels.DeleteSessionRequest) (*aiDataSetModels.DeleteSessionResponse, error) {
	var data aiDataSetModels.DeleteSessionApiRequest
	data.Ids = req.Ids
	content, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/chats/%s/sessions",
		i.config.Host, req.ChatId), bytes.NewBuffer(content))
	if err != nil {
		zap.L().Error("删除助理失败", zap.Error(err))
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+i.config.ApiKey)
	resp, err := i.client.Do(request)
	if err != nil {
		zap.L().Error("删除助理失败", zap.Error(err))
		return nil, errors.New("删除助理失败")
	}
	if resp == nil {
		return nil, errors.New("删除助理失败")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("删除助理失败", zap.Error(err))
		return nil, errors.New("删除助理失败")
	}
	var response aiDataSetModels.DeleteSessionResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errors.New("删除助理失败")
	}
	if response.Code != 0 {
		zap.L().Error(fmt.Sprintf("删除助理失败: %s", response.Message))
		return nil, errors.New("删除助理失败")
	}
	return &response, nil
}

func (i *IChartServiceImpl) AgentSessionList(c *gin.Context, req *aiDataSetModels.ListAgentSessionsRequest) (*aiDataSetModels.AgentSessionListResponse, error) {
	params := url.Values{}
	if req.Dsl != true {
		params.Add("dsl", "false")
	}
	if req.UserId != "" {
		params.Add("user_id", req.UserId)
	}
	if req.Page > 0 {
		params.Add("page", strconv.FormatInt(req.Page, 10))
	}
	if req.Size > 0 {
		params.Add("page_size", strconv.FormatInt(req.Size, 10))
	}
	if req.IsAsc != "" {
		params.Add("orderby", req.IsAsc)
	}

	request, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/agents/%s/sessions?%s",
		i.config.Host, req.AgentId, params.Encode()), nil)
	if err != nil {
		zap.L().Error("读取会话列表失败", zap.Error(err))
		return &aiDataSetModels.AgentSessionListResponse{}, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+i.config.ApiKey)
	resp, err := i.client.Do(request)
	if err != nil {
		zap.L().Error("读取会话列表失败", zap.Error(err))
		return &aiDataSetModels.AgentSessionListResponse{}, errors.New("读取会话列表失败")
	}
	if resp == nil {
		return &aiDataSetModels.AgentSessionListResponse{}, errors.New("读取会话列表失败")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("读取会话列表失败", zap.Error(err))
		return &aiDataSetModels.AgentSessionListResponse{}, errors.New("读取会话列表失败")
	}
	var response aiDataSetModels.AgentSessionListResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return &aiDataSetModels.AgentSessionListResponse{}, errors.New("读取会话列表失败")
	}
	if response.Code != 0 {
		zap.L().Error(fmt.Sprintf("读取会话列表失败: %s", response.Message))
		return &aiDataSetModels.AgentSessionListResponse{}, errors.New("读取会话列表失败")
	}
	return &response, nil
}

func (i *IChartServiceImpl) AgentSessionRemove(c *gin.Context, req *aiDataSetModels.RemoveAgentSessionsRequest) (*aiDataSetModels.RemoveAgentSessionsResponse, error) {
	var data aiDataSetModels.DeleteSessionApiRequest
	data.Ids = req.Ids
	content, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/agents/%s/sessions",
		i.config.Host, req.AgentId), bytes.NewBuffer(content))
	if err != nil {
		zap.L().Error("删除助理会话失败", zap.Error(err))
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+i.config.ApiKey)
	resp, err := i.client.Do(request)
	if err != nil {
		zap.L().Error("删除助理会话失败", zap.Error(err))
		return nil, errors.New("删除助理会话失败")
	}
	if resp == nil {
		return nil, errors.New("删除助理会话失败")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("删除助理会话失败", zap.Error(err))
		return nil, errors.New("删除助理会话失败")
	}
	var response aiDataSetModels.RemoveAgentSessionsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errors.New("删除助理会话失败")
	}
	if response.Code != 0 {
		zap.L().Error(fmt.Sprintf("删除助理会话失败: %s", response.Message))
		return nil, errors.New("删除助理失败")
	}
	return &response, nil
}

func (i *IChartServiceImpl) ConversationRelatedQuestions(c *gin.Context, req *aiDataSetModels.ConversationRelatedQuestionsRequest) (*aiDataSetModels.ConversationRelatedQuestionsResponse, error) {
	content, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/sessions/related_questions",
		i.config.Host), bytes.NewBuffer(content))
	if err != nil {
		zap.L().Error("提问失败", zap.Error(err))
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+i.config.ApiKey)
	resp, err := i.client.Do(request)
	if err != nil {
		zap.L().Error("提问失败", zap.Error(err))
		return nil, errors.New("提问失败")
	}
	defer resp.Body.Close()
	if resp == nil {
		return nil, errors.New("提问失败")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("提问失败", zap.Error(err))
		return nil, errors.New("提问失败")
	}
	var response aiDataSetModels.ConversationRelatedQuestionsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errors.New("提问失败")
	}
	if response.Code != 0 {
		zap.L().Error(fmt.Sprintf("提问失败: %s", response.Message))
		return nil, errors.New("提问失败")
	}
	return &response, nil
}

func (i *IChartServiceImpl) AgentList(c *gin.Context, req *aiDataSetModels.ListAgentRequest) (*aiDataSetModels.ListAgentResponse, error) {
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
	if req.IsAsc != "" {
		params.Add("orderby", req.IsAsc)
	}

	request, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/agents?%s",
		i.config.Host, params.Encode()), nil)
	if err != nil {
		zap.L().Error("读取会话列表失败", zap.Error(err))
		return &aiDataSetModels.ListAgentResponse{}, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+i.config.ApiKey)
	request.Header.Set("Accept", "text/event-stream")
	request.Header.Set("Cache-Control", "no-cache")
	resp, err := i.client.Do(request)
	if err != nil {
		zap.L().Error("读取会话列表失败", zap.Error(err))
		return &aiDataSetModels.ListAgentResponse{}, errors.New("读取会话列表失败")
	}
	if resp == nil {
		return &aiDataSetModels.ListAgentResponse{}, errors.New("读取会话列表失败")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("读取会话列表失败", zap.Error(err))
		return &aiDataSetModels.ListAgentResponse{}, errors.New("读取会话列表失败")
	}
	var response aiDataSetModels.ListAgentResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return &aiDataSetModels.ListAgentResponse{}, errors.New("读取会话列表失败")
	}
	if response.Code != 0 {
		zap.L().Error(fmt.Sprintf("读取会话列表失败: %s", response.Message))
		return &aiDataSetModels.ListAgentResponse{}, errors.New("读取会话列表失败")
	}
	return &response, nil
}

func (i *IChartServiceImpl) ChatsCompletions(c *gin.Context, req *aiDataSetModels.ChatsCompletionsRequest) (*aiDataSetModels.ChatsCompletionsResponse, error) {
	var data aiDataSetModels.ChatsCompletionsApiRequest
	data.Question = req.Question
	data.Stream = req.Stream
	data.SessionId = req.SessionId
	data.UserId = req.UserId
	content, err := json.Marshal(data)
	if err != nil {
		zap.L().Error("与聊天助手交谈失败", zap.Error(err))
		return nil, err
	}
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/chats/%s/completions",
		i.config.Host, req.ChatId), bytes.NewBuffer(content))
	if err != nil {
		zap.L().Error("与聊天助手交谈失败", zap.Error(err))
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+i.config.ApiKey)
	var client *http.Client = i.client
	if req.Stream {
		transport := &http.Transport{
			MaxIdleConns:        100,              // 整个客户端的最大空闲连接数
			MaxIdleConnsPerHost: 60,               // 每个主机的最大空闲连接数
			MaxConnsPerHost:     50,               // 每个主机的最大连接数
			IdleConnTimeout:     60 * time.Second, // 空闲连接的超时时间

			DisableKeepAlives: false, // 不禁用连接保持活动
			ForceAttemptHTTP2: true,  // 尝试使用HTTP/2
		}
		// 创建HTTP客户端并设置自定义的Transport
		client = &http.Client{
			Transport: transport,
			Timeout:   300 * time.Second, // 请求的总超时时间
		}
	}
	resp, err := client.Do(request)
	if err != nil {
		zap.L().Error("与聊天助手交谈失败", zap.Error(err))
		return nil, errors.New("与聊天助手交谈失败")
	}
	if resp == nil {
		return nil, errors.New("与聊天助手交谈失败")
	}
	defer resp.Body.Close()

	var response aiDataSetModels.ChatsCompletionsResponse
	if !req.Stream {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			zap.L().Error("与聊天助手交谈失败", zap.Error(err))
			return nil, errors.New("与聊天助手交谈失败")
		}
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, errors.New("与聊天助手交谈失败")
		}
		if response.Code != 0 {
			zap.L().Error(fmt.Sprintf("与聊天助手交谈失败: %s", response.Message))
			return nil, errors.New("与聊天助手交谈失败")
		}
	} else {
		// 声明数据格式为event stream
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		// 禁用nginx缓存,防止nginx会缓存数据导致数据流是一段一段的
		c.Writer.Header().Set("X-Accel-Buffering", "no")
		w := c.Writer
		flusher, _ := w.(http.Flusher)
		flusher.Flush()
		for ev, err := range sse.Read(resp.Body, nil) {
			if err != nil {
				fmt.Fprintf(w, "event: error\n")
				fmt.Fprintf(w, "data: %s\n\n", err.Error())
				// handle read error
				break // can end the loop as Read stops on first error anyway
			}
			fmt.Fprintf(w, "data: "+ev.Data+"\n\n")
			// Do something with the events, parse the JSON or whatever.
			flusher.Flush()
		}
	}

	return &response, nil
}

func (i *IChartServiceImpl) AgentSessionCreate(c *gin.Context, req *aiDataSetModels.SessionAgentCreate) (*aiDataSetModels.SessionAgentResponse, error) {
	content, err := json.Marshal(req.Data)
	if err != nil {
		zap.L().Error("创建Agent会话失败", zap.Error(err))
		return nil, err
	}
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/agents/%s/sessions",
		i.config.Host, req.AgentId), bytes.NewBuffer(content))
	if err != nil {
		zap.L().Error("创建Agent会话失败", zap.Error(err))
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+i.config.ApiKey)
	resp, err := i.client.Do(request)
	if err != nil {
		zap.L().Error("创建Agent会话失败", zap.Error(err))
		return nil, errors.New("创建Agent会话失败")
	}
	if resp == nil {
		return nil, errors.New("创建Agent会话失败")
	}
	defer resp.Body.Close()

	var response aiDataSetModels.SessionAgentResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("创建Agent会话失败", zap.Error(err))
		return nil, errors.New("创建Agent会话失败")
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errors.New("创建Agent会话失败")
	}
	if response.Code != 0 {
		zap.L().Error(fmt.Sprintf("创建Agent会话失败: %s", response.Message))
		return nil, errors.New("创建Agent会话失败")
	}

	return &response, nil
}

func (i *IChartServiceImpl) AgentsCompletions(c *gin.Context, req *aiDataSetModels.AgentsCompletionsRequest) (*aiDataSetModels.AgentsCompletionsApiResponse, error) {
	var data aiDataSetModels.AgentsCompletionsApiRequest
	data.Question = req.Question
	data.Stream = req.Stream
	data.SessionId = req.SessionId
	data.UserId = req.UserId
	data.SyncDsl = req.SyncDsl
	content, err := json.Marshal(data)
	if err != nil {
		zap.L().Error("与聊天助手交谈失败", zap.Error(err))
		return nil, err
	}
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/agents/%s/completions",
		i.config.Host, req.AgentId), bytes.NewBuffer(content))
	if err != nil {
		zap.L().Error("与聊天助手交谈失败", zap.Error(err))
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+i.config.ApiKey)

	var client *http.Client = i.client
	if req.Stream {
		transport := &http.Transport{
			MaxIdleConns:        100,              // 整个客户端的最大空闲连接数
			MaxIdleConnsPerHost: 60,               // 每个主机的最大空闲连接数
			MaxConnsPerHost:     50,               // 每个主机的最大连接数
			IdleConnTimeout:     60 * time.Second, // 空闲连接的超时时间

			DisableKeepAlives: false, // 不禁用连接保持活动
			ForceAttemptHTTP2: true,  // 尝试使用HTTP/2
		}
		// 创建HTTP客户端并设置自定义的Transport
		client = &http.Client{
			Transport: transport,
			Timeout:   300 * time.Second, // 请求的总超时时间
		}
	}
	resp, err := client.Do(request)
	if err != nil {
		zap.L().Error("与聊天助手交谈失败", zap.Error(err))
		return nil, errors.New("与聊天助手交谈失败")
	}
	if resp == nil {
		return nil, errors.New("与聊天助手交谈失败")
	}
	defer resp.Body.Close()

	var response aiDataSetModels.AgentsCompletionsApiResponse
	if !req.Stream {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			zap.L().Error("与聊天助手交谈失败", zap.Error(err))
			return nil, errors.New("与聊天助手交谈失败")
		}
		err = json.Unmarshal(body, &response)
		if err != nil {
			return nil, errors.New("与聊天助手交谈失败")
		}
		if response.Code != 0 {
			zap.L().Error(fmt.Sprintf("与聊天助手交谈失败: %s", response.Message))
			return nil, errors.New("与聊天助手交谈失败")
		}
	} else {
		// 声明数据格式为event stream
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		// 禁用nginx缓存,防止nginx会缓存数据导致数据流是一段一段的
		c.Writer.Header().Set("X-Accel-Buffering", "no")
		w := c.Writer
		flusher, _ := w.(http.Flusher)
		flusher.Flush()
		for ev, err := range sse.Read(resp.Body, nil) {
			if err != nil {
				fmt.Fprintf(w, "event: error\n")
				fmt.Fprintf(w, "data: %s\n\n", err.Error())
				// handle read error
				break // can end the loop as Read stops on first error anyway
			}
			fmt.Fprintf(w, "data: "+ev.Data+"\n\n")
			// Do something with the events, parse the JSON or whatever.
			flusher.Flush()
		}

	}

	return &response, nil
}

func (i *IChartServiceImpl) Ask(c *gin.Context, req *aiDataSetModels.AskRequest) error {
	content, err := json.Marshal(req)
	if err != nil {
		return err
	}
	request, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/sessions/ask",
		i.config.Host), bytes.NewBuffer(content))
	if err != nil {
		zap.L().Error("回答失败", zap.Error(err))
		return err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+i.config.ApiKey)
	resp, err := i.client.Do(request)
	if err != nil {
		zap.L().Error("回答失败", zap.Error(err))
		return errors.New("回答失败")
	}
	if resp == nil {
		return errors.New("回答失败")
	}
	defer resp.Body.Close()

	w := c.Writer
	flusher, _ := w.(http.Flusher)

	var isSse bool = false
	for ev, err := range sse.Read(resp.Body, nil) {
		if !isSse {
			// 声明数据格式为event stream
			c.Writer.Header().Set("Content-Type", "text/event-stream")
			c.Writer.Header().Set("Cache-Control", "no-cache")
			c.Writer.Header().Set("Connection", "keep-alive")
			// 禁用nginx缓存,防止nginx会缓存数据导致数据流是一段一段的
			c.Writer.Header().Set("X-Accel-Buffering", "no")
			flusher.Flush()
		}
		isSse = true
		if err != nil {
			fmt.Fprintf(w, "event: error\n")
			fmt.Fprintf(w, "data: %s\n\n", err.Error())
			// handle read error
			break // can end the loop as Read stops on first error anyway
		}
		fmt.Fprintf(w, "data: "+ev.Data+"\n\n")
		// Do something with the events, parse the JSON or whatever.
		flusher.Flush()
	}

	if !isSse {
		return errors.New("智能回答失败")
	}
	return nil
}
