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
	"mime/multipart"
	"net/http"
	"net/url"
	"nova-factory-server/app/business/ai/aiDataSetDao"
	"nova-factory-server/app/business/ai/aiDataSetModels"
	"nova-factory-server/app/business/ai/aiDataSetService"
	"nova-factory-server/app/utils/file"
	"strconv"
	"time"
)

type IDataSetDocumentServiceImpl struct {
	client              *http.Client
	config              RagFlowConfig
	iDatasetDao         aiDataSetDao.IDataSetDao
	iDatasetDocumentDao aiDataSetDao.IDataSetDocumentDao
}

func NewIDataSetDocumentServiceImpl(client *http.Client, iDatasetDao aiDataSetDao.IDataSetDao, iDatasetDocumentDao aiDataSetDao.IDataSetDocumentDao) aiDataSetService.IDataSetDocumentService {
	var config RagFlowConfig
	err := viper.UnmarshalKey("dataSet", &config)
	if err != nil {
		panic(err)
	}
	return &IDataSetDocumentServiceImpl{
		client:              client,
		config:              config,
		iDatasetDao:         iDatasetDao,
		iDatasetDocumentDao: iDatasetDocumentDao,
	}
}

func (i *IDataSetDocumentServiceImpl) UploadFile(c *gin.Context, datasetId int64) ([]*aiDataSetModels.SysDatasetDocument, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	// 创建一个缓冲区来存储表单数据
	//var requestBody bytes.Buffer
	multiPartform, _ := c.MultipartForm()
	err := file.CreatFormFiles(&b, multiPartform, w)
	if err != nil {
		zap.L().Error("上传文档失败", zap.Error(err))
		return make([]*aiDataSetModels.SysDatasetDocument, 0), errors.New("上传文档失败")
	}
	datasetInfo, err := i.iDatasetDao.GetById(c, datasetId)
	if err != nil {
		zap.L().Error("上传文档失败", zap.Error(err))
		return make([]*aiDataSetModels.SysDatasetDocument, 0), errors.New("上传文档失败")
	}
	if datasetInfo == nil {
		return make([]*aiDataSetModels.SysDatasetDocument, 0), errors.New("知识库不存在")
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/datasets/%s/documents", i.config.Host, datasetInfo.DatasetUUID), &b)
	if err != nil {
		zap.L().Error("上传文档失败", zap.Error(err))
		return make([]*aiDataSetModels.SysDatasetDocument, 0), errors.New("上传文档失败")
	}
	if req == nil {
		return make([]*aiDataSetModels.SysDatasetDocument, 0), errors.New("上传文档失败")
	}

	// 创建一个多部分写入器
	req.Header.Set("Authorization", "Bearer "+i.config.ApiKey)
	req.Header.Set("Content-Type", w.FormDataContentType())
	resp, err := i.client.Do(req)
	if err != nil {
		zap.L().Error("上传文档失败", zap.Error(err))
		return make([]*aiDataSetModels.SysDatasetDocument, 0), errors.New("上传文档失败")
	}
	if resp == nil {
		zap.L().Error("更新知识库失败", zap.Error(err))
		return nil, errors.New("更新知识库失败")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("创建知识库失败", zap.Error(err))
		return make([]*aiDataSetModels.SysDatasetDocument, 0), errors.New("上传文档失败")
	}
	defer resp.Body.Close()
	var response aiDataSetModels.UploadDocumentResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		zap.L().Error("创建知识库失败", zap.Error(err))
		return make([]*aiDataSetModels.SysDatasetDocument, 0), errors.New("上传文档失败")
	}

	if response.Code != 0 {
		zap.L().Error("上传文档失败:" + response.Message)
		return make([]*aiDataSetModels.SysDatasetDocument, 0), errors.New("上传文档失败")
	}

	return i.iDatasetDocumentDao.Create(c, datasetId, &response)
}

func (i *IDataSetDocumentServiceImpl) PutFile(c *gin.Context, documentId int64, request *aiDataSetModels.PutDocumentRequest) (*aiDataSetModels.SysDatasetDocument, error) {
	info, err := i.iDatasetDocumentDao.GetById(c, documentId)
	if err != nil {
		zap.L().Error("读取文档失败", zap.Error(err))
		return nil, errors.New("读取文档失败")
	}
	if info == nil {
		zap.L().Error("文档不存在", zap.Error(err))
		return nil, errors.New("文档不存在")
	}
	var apiRequest aiDataSetModels.PutDocumentRequest
	apiRequest.Name = request.Name
	apiRequest.ParserConfig = request.ParserConfig
	apiRequest.ChunkMethod = request.ChunkMethod
	apiRequest.MetaFields = request.MetaFields
	content, err := json.Marshal(apiRequest)
	if err != nil {
		zap.L().Error("读取文档失败", zap.Error(err))
		return nil, errors.New("读取文档失败")
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/datasets/%s/documents/%s",
		i.config.Host, info.DatasetDatasetUUID, info.DatasetDocumentUUID), bytes.NewBuffer(content))
	if err != nil {
		zap.L().Error("更新知识库失败", zap.Error(err))
		return nil, errors.New("更新知识库失败")
	}
	if req == nil {
		return nil, errors.New("更新知识库失败")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+i.config.ApiKey)
	resp, err := i.client.Do(req)
	if err != nil {
		zap.L().Error("更新知识库失败", zap.Error(err))
		return nil, errors.New("更新知识库失败")
	}
	if resp == nil {
		zap.L().Error("更新知识库失败", zap.Error(err))
		return nil, errors.New("更新知识库失败")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("更新知识库失败", zap.Error(err))
		return nil, errors.New("更新知识库失败")
	}
	defer resp.Body.Close()
	var response aiDataSetModels.PutDocumentResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errors.New("更新知识库失败")
	}

	if response.Code != 0 {
		zap.L().Error("更新知识库失败" + response.Message)
		return nil, errors.New("更新知识库失败")
	}
	return i.iDatasetDocumentDao.Update(c, documentId, request)

}

func (i *IDataSetDocumentServiceImpl) DownloadFile(c *gin.Context, documentId int64) (*aiDataSetModels.SysDatasetDocument, error) {

	info, err := i.iDatasetDocumentDao.GetById(c, documentId)
	if err != nil {
		zap.L().Error("下载文档失败", zap.Error(err))
		return nil, errors.New("下载文档失败")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/datasets/%s/documents/%s",
		i.config.Host, info.DatasetDatasetUUID, info.DatasetDocumentUUID), nil)
	if err != nil {
		zap.L().Error("下载文档失败", zap.Error(err))
		return nil, errors.New("下载文档失败")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+i.config.ApiKey)
	resp, err := i.client.Do(req)
	if err != nil {
		zap.L().Error("下载文档失败", zap.Error(err))
		return nil, errors.New("下载文档失败")
	}
	if resp == nil {
		zap.L().Error("下载文档失败", zap.Error(err))
		return nil, errors.New("下载文档失败")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("下载文档失败")
	}
	// 设置 HTTP 头部
	c.Writer.Header().Set("Transfer-Encoding", "chunked") // 告诉浏览器，分段的流式的输出数据
	c.Writer.Header().Set("Content-Disposition", "attachment; filename="+info.DatasetLocation)
	//c.Writer.Header().Set("Content-Type", "application/octet-stream") // 根据文件类型设置
	//c.Writer.Header().Set("Accept-Ranges", "bytes")
	c.Writer.Flush()

	// 循环读取响应体
	for {
		buffer := make([]byte, 1024)
		n, err := resp.Body.Read(buffer)
		if err != nil {
			if err != io.EOF {
				zap.L().Error("下载文档失败", zap.Error(err))
			} else {
				// 处理读取到的数据
				c.Writer.Write(buffer[:n])
				c.Writer.Flush()
			}
			break
		}

		// 处理读取到的数据
		c.Writer.Write(buffer[:n])
		c.Writer.Flush()
	}
	return nil, nil
}

func (i *IDataSetDocumentServiceImpl) ListDocument(c *gin.Context, datasetId int64, req *aiDataSetModels.ListDocumentRequest) (*aiDataSetModels.ListDocumentData, error) {
	if datasetId == 0 {
		return &aiDataSetModels.ListDocumentData{}, nil
	}
	datasetInfo, err := i.iDatasetDao.GetById(c, datasetId)
	if err != nil {
		zap.L().Error("上传文档失败", zap.Error(err))
		return &aiDataSetModels.ListDocumentData{}, errors.New("读取知识库错误")
	}
	if datasetInfo == nil {
		return &aiDataSetModels.ListDocumentData{}, errors.New("知识库不存在")
	}
	params := url.Values{}
	if req.DocumentId != "" {
		params.Add("document_id", req.DocumentId)
	}
	if req.Keywords != "" {
		params.Add("keywords", req.Keywords)
	}
	if req.DocumentName != "" {
		params.Add("document_name", req.DocumentName)
	}
	if req.Page > 0 {
		params.Add("page", strconv.FormatInt(req.Page, 10))
	}
	if req.Size > 0 {
		params.Add("page_size", strconv.FormatInt(req.Size, 10))
	}
	if req.OrderBy != "" {
		params.Add("orderby", req.OrderBy)
	}
	if req.IsAsc != "" {
		params.Add("orderby", req.IsAsc)
	}
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/datasets/%s/documents?%s",
		i.config.Host, datasetInfo.DatasetUUID, params.Encode()), nil)
	if err != nil {
		zap.L().Error("下载文档失败", zap.Error(err))
		return nil, errors.New("下载文档失败")
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+i.config.ApiKey)
	resp, err := i.client.Do(request)
	if err != nil {
		zap.L().Error("下载文档失败", zap.Error(err))
		return nil, errors.New("下载文档失败")
	}
	if resp == nil {
		zap.L().Error("下载文档失败", zap.Error(err))
		return nil, errors.New("下载文档失败")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("下载文档失败")
	}

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response aiDataSetModels.ListDocumentResponseData
	err = json.Unmarshal(all, &response)
	if err != nil {
		return nil, err
	}
	if response.Code != 0 {
		return nil, errors.New(response.Message)
	}

	var uuids []string = make([]string, 0)
	for _, v := range response.Data.Docs {
		uuids = append(uuids, v.Id)
	}

	// 读取数据库索引
	documents, err := i.iDatasetDocumentDao.GetByDocumentUuids(c, uuids)
	if err != nil {
		return nil, err
	}

	var documentMap map[string]*aiDataSetModels.SysDatasetDocument = make(map[string]*aiDataSetModels.SysDatasetDocument)
	for _, v := range documents {
		documentMap[v.DatasetDocumentUUID] = v
	}

	var list aiDataSetModels.ListDocumentData
	list.Rows = make([]*aiDataSetModels.SysDatasetDocumentData, 0)
	list.Total = int64(response.Data.Total)
	for _, v := range response.Data.Docs {
		documentData, ok := documentMap[v.Id]
		var documentId int64 = 0
		var databaseId int64 = 0
		var deptId int64 = 0
		var datasetDatasetUUID string
		var datasetLanguage string
		if ok {
			documentId = documentData.DocumentID
			databaseId = documentData.DatasetID
			datasetDatasetUUID = documentData.DatasetDatasetUUID
			datasetLanguage = documentData.DatasetLanguage
			deptId = documentData.DeptID
		}
		var document aiDataSetModels.SysDatasetDocumentData
		document.DocumentID = documentId
		document.DatasetID = databaseId
		document.DatasetChunkMethod = v.ChunkMethod
		document.DatasetCreatedBy = v.CreatedBy
		document.DatasetDocumentUUID = v.Id
		document.DatasetDatasetUUID = datasetDatasetUUID
		document.DatasetLanguage = datasetLanguage
		document.DatasetLocation = v.Location
		document.DatasetName = v.Name
		document.DatasetParserConfig = v.ParserConfig
		document.DatasetRun = v.Run
		document.DatasetSize = int64(v.Size)
		document.DatasetThumbnail = v.Thumbnail
		document.DatasetType = v.Type
		document.DeptID = deptId
		document.ChunkCount = v.ChunkCount
		document.Run = v.Run
		document.Progress = v.Progress

		layout := "Mon, 02 Jan 2006 15:04:05 MST" // 注意这里的布局要和输入字符串的格式相匹配
		// 使用 Parse 函数解析时间字符串
		t, err := time.Parse(layout, v.CreateDate)
		if err == nil {
			document.CreateTime = &t
		} else {
			now := time.Now()
			document.CreateTime = &now
		}
		// 使用 Parse 函数解析时间字符串
		t, err = time.Parse(layout, v.UpdateDate)
		if err == nil {
			document.UpdateTime = &t
		} else {
			now := time.Now()
			document.UpdateTime = &now
		}
		list.Rows = append(list.Rows, &document)
	}
	return &list, nil
}

func (i *IDataSetDocumentServiceImpl) RemoveDocument(c *gin.Context, request *aiDataSetModels.DeleteDocumentRequest) error {
	datasetInfo, err := i.iDatasetDao.GetById(c, request.DatasetId)
	if err != nil {
		zap.L().Error("删除文档失败", zap.Error(err))
		return errors.New("删除文档失败")
	}

	if datasetInfo == nil {
		return errors.New("知识库不存在")
	}

	documents, err := i.iDatasetDocumentDao.GetByIds(c, request.DocumentIds)
	if len(documents) == 0 {
		return errors.New("文档不存在")
	}

	var deleteRequest aiDataSetModels.DeleteDocumentApiRequest
	deleteRequest.Ids = make([]string, 0)
	for _, v := range documents {
		deleteRequest.Ids = append(deleteRequest.Ids, v.DatasetDocumentUUID)
	}
	content, err := json.Marshal(deleteRequest)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/datasets/%s/documents",
		i.config.Host, datasetInfo.DatasetUUID), bytes.NewBuffer(content))
	if err != nil {
		zap.L().Error("删除文档失败", zap.Error(err))
		return errors.New("删除文档失败")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+i.config.ApiKey)
	resp, err := i.client.Do(req)
	if err != nil {
		zap.L().Error("删除文档失败", zap.Error(err))
		return errors.New("删除文档失败")
	}
	if resp == nil {
		zap.L().Error("删除文档失败", zap.Error(err))
		return errors.New("删除文档失败")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("删除文档失败")
	}

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response aiDataSetModels.DeleteDocumentApiResponse
	err = json.Unmarshal(all, &response)
	if err != nil {
		zap.L().Error("删除文档失败", zap.Error(err))
		return errors.New("删除文档失败")
	}

	if response.Code != 0 {
		zap.L().Error(fmt.Sprintf("删除文档失败, code: %d;message: %s", response.Code, response.Message))
		return errors.New("删除文档失败")
	}

	err = i.iDatasetDocumentDao.RemoveByIds(c, request.DocumentIds)
	if err != nil {
		return err
	}

	return nil
}

func (i *IDataSetDocumentServiceImpl) StartParse(c *gin.Context, request *aiDataSetModels.ParseDocumentApiRequest) error {
	datasetInfo, err := i.iDatasetDao.GetById(c, request.DatasetId)
	if err != nil {
		zap.L().Error("解析文档失败", zap.Error(err))
		return errors.New("解析文档失败")
	}

	if datasetInfo == nil {
		return errors.New("知识库不存在")
	}

	documents, err := i.iDatasetDocumentDao.GetByIds(c, request.DocumentIds)
	if len(documents) == 0 {
		return errors.New("文档不存在")
	}

	var parseRequest aiDataSetModels.ParseDocumentApiRequest
	parseRequest.DocumentIds = make([]string, 0)
	for _, v := range documents {
		parseRequest.DocumentIds = append(parseRequest.DocumentIds, v.DatasetDocumentUUID)
	}
	content, err := json.Marshal(parseRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/datasets/%s/chunks",
		i.config.Host, datasetInfo.DatasetUUID), bytes.NewBuffer(content))
	if err != nil {
		zap.L().Error("解析文档失败", zap.Error(err))
		return errors.New("解析文档失败")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+i.config.ApiKey)
	resp, err := i.client.Do(req)
	if err != nil {
		zap.L().Error("解析文档失败", zap.Error(err))
		return errors.New("解析文档失败")
	}
	if resp == nil {
		zap.L().Error("解析文档失败", zap.Error(err))
		return errors.New("解析文档失败")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("解析文档失败")
	}

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("解析文档失败", zap.Error(err))
		return errors.New("解析文档失败")
	}

	var response aiDataSetModels.ParseDocumentApiResponse
	err = json.Unmarshal(all, &response)
	if err != nil {
		zap.L().Error("解析文档失败", zap.Error(err))
		return errors.New("解析文档失败")
	}

	if response.Code != 0 {
		zap.L().Error(fmt.Sprintf("解析文档失败, code: %d;message: %s", response.Code, response.Message))
		return errors.New("解析文档失败")
	}

	return nil
}

func (i *IDataSetDocumentServiceImpl) StopParse(c *gin.Context, request *aiDataSetModels.ParseDocumentApiRequest) error {
	datasetInfo, err := i.iDatasetDao.GetById(c, request.DatasetId)
	if err != nil {
		zap.L().Error("停止解析文档失败", zap.Error(err))
		return errors.New("停止解析文档失败")
	}

	if datasetInfo == nil {
		return errors.New("知识库不存在")
	}

	documents, err := i.iDatasetDocumentDao.GetByIds(c, request.DocumentIds)
	if len(documents) == 0 {
		return errors.New("文档不存在")
	}

	var parseRequest aiDataSetModels.ParseDocumentApiRequest
	parseRequest.DocumentIds = make([]string, 0)
	for _, v := range documents {
		parseRequest.DocumentIds = append(parseRequest.DocumentIds, v.DatasetDocumentUUID)
	}
	content, err := json.Marshal(parseRequest)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/datasets/%s/chunks",
		i.config.Host, datasetInfo.DatasetUUID), bytes.NewBuffer(content))
	if err != nil {
		zap.L().Error("停止解析文档失败", zap.Error(err))
		return errors.New("停止解析文档失败")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+i.config.ApiKey)
	resp, err := i.client.Do(req)
	if err != nil {
		zap.L().Error("停止解析文档失败", zap.Error(err))
		return errors.New("停止解析文档失败")
	}
	if resp == nil {
		zap.L().Error("停止解析文档失败", zap.Error(err))
		return errors.New("停止解析文档失败")
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New("停止解析文档失败")
	}

	all, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("停止解析文档失败", zap.Error(err))
		return errors.New("停止解析文档失败")
	}

	var response aiDataSetModels.ParseDocumentApiResponse
	err = json.Unmarshal(all, &response)
	if err != nil {
		zap.L().Error("停止解析文档失败", zap.Error(err))
		return errors.New("停止解析文档失败")
	}

	if response.Code != 0 {
		zap.L().Error(fmt.Sprintf("停止解析文档失败, code: %d;message: %s", response.Code, response.Message))
		return errors.New("停止解析文档失败")
	}

	return nil
}
