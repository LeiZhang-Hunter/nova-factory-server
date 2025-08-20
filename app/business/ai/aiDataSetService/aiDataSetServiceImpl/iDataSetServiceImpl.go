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
	"nova-factory-server/app/business/ai/aiDataSetDao"
	"nova-factory-server/app/business/ai/aiDataSetModels"
	"nova-factory-server/app/business/ai/aiDataSetService"
	"nova-factory-server/app/business/system/systemDao"
	"nova-factory-server/app/business/system/systemModels"
	"nova-factory-server/app/utils/baizeContext"
)

type IDataSetServiceImpl struct {
	config      RagFlowConfig
	client      *http.Client
	iDataSetDao aiDataSetDao.IDataSetDao
	iDeptDao    systemDao.IDeptDao
	iUserDao    systemDao.IUserDao
}

func NewIDataSetServiceImpl(iDataSetDao aiDataSetDao.IDataSetDao, iDeptDao systemDao.IDeptDao, iUserDao systemDao.IUserDao, client *http.Client) aiDataSetService.IDataSetService {
	var config RagFlowConfig
	err := viper.UnmarshalKey("dataSet", &config)
	if err != nil {
		panic(err)
	}

	return &IDataSetServiceImpl{
		client:      client,
		config:      config,
		iUserDao:    iUserDao,
		iDataSetDao: iDataSetDao,
		iDeptDao:    iDeptDao,
	}
}

func (i *IDataSetServiceImpl) CreateDataSet(c *gin.Context, request *aiDataSetModels.DataSetRequest) (*aiDataSetModels.SysDataset, error) {
	deptId := baizeContext.GetDeptId(c)
	dept := i.iDeptDao.SelectDeptById(c, deptId)
	if dept == nil {
		return nil, errors.New("部门不存在")
	}
	var originName string = request.Name
	request.Name = dept.DeptName + "-" + request.Name

	info, err := i.iDataSetDao.GetByName(c, originName)
	if err != nil {
		zap.L().Error("创建知识库失败", zap.Error(err))
		return nil, errors.New("创建知识库失败")
	}
	if info != nil {
		return nil, errors.New("知识库已经存在")
	}
	content, err := json.Marshal(request)
	if err != nil {
		zap.L().Error("创建知识库失败", zap.Error(err))
		return nil, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/datasets", i.config.Host), bytes.NewBuffer(content))
	if err != nil {
		zap.L().Error("创建知识库失败", zap.Error(err))
		return nil, err
	}
	if req == nil {
		return nil, errors.New("创建知识库失败")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+i.config.ApiKey)
	resp, err := i.client.Do(req)
	if err != nil {
		zap.L().Error("创建知识库失败", zap.Error(err))
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("创建知识库失败", zap.Error(err))
		return nil, err
	}
	var response aiDataSetModels.DataSetCreateResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		zap.L().Error("创建知识库失败", zap.Error(err))
		return nil, err
	}

	if response.Code != 0 {
		return nil, errors.New(response.Message)
	}
	response.Data.Name = originName
	create, err := i.iDataSetDao.Create(c, &response)
	if err != nil {
		zap.L().Error("创建知识库失败", zap.Error(err))
		return nil, err
	}
	return create, nil
}

func (i *IDataSetServiceImpl) UpdateDataSet(c *gin.Context, request *aiDataSetModels.UpdateDataSetRequest, id int64) (*aiDataSetModels.SysDataset, error) {
	deptId := baizeContext.GetDeptId(c)
	dept := i.iDeptDao.SelectDeptById(c, deptId)
	if dept == nil {
		return nil, errors.New("部门不存在")
	}
	orginName := request.Name
	request.Name = dept.DeptName + "-" + request.Name

	info, err := i.iDataSetDao.GetById(c, id)
	if err != nil {
		zap.L().Error("创建知识库失败", zap.Error(err))
		return nil, errors.New("更新知识库失败")
	}
	if info == nil {
		return nil, errors.New("知识库不存在")
	}

	if request.ParserConfig != nil && request.ParserConfig.Graphrag != nil && !request.ParserConfig.Graphrag.UseGraphrag {
		request.ParserConfig.Graphrag = nil
	}
	if request.ParserConfig != nil && request.ParserConfig.Raptor != nil && !request.ParserConfig.Raptor.UseRaptor {
		request.ParserConfig.Raptor = nil
	}
	content, err := json.Marshal(request)
	if err != nil {
		zap.L().Error("更新知识库失败", zap.Error(err))
		return nil, err
	}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/datasets/"+info.DatasetUUID, i.config.Host), bytes.NewBuffer(content))
	if err != nil {
		zap.L().Error("更新知识库失败", zap.Error(err))
		return nil, err
	}
	if req == nil {
		return nil, errors.New("更新知识库失败")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+i.config.ApiKey)
	resp, err := i.client.Do(req)
	if err != nil {
		zap.L().Error("创建知识库失败", zap.Error(err))
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("创建知识库失败", zap.Error(err))
		return nil, err
	}
	var response aiDataSetModels.DataSetCreateResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		zap.L().Error("创建知识库失败", zap.Error(err))
		return nil, err
	}

	if response.Code != 0 {
		return nil, errors.New(response.Message)
	}
	request.Name = orginName
	create, err := i.iDataSetDao.Update(c, info.DatasetID, request)
	if err != nil {
		return nil, err
	}
	return create, nil
}

func (i *IDataSetServiceImpl) SelectDataSet(c *gin.Context, request *aiDataSetModels.DatasetListReq) (*aiDataSetModels.SysDatasetListData, error) {
	list, err := i.iDataSetDao.SelectByList(c, request)
	if err != nil {
		return nil, err
	}

	//  读取用户id集合
	userIdMap := make(map[int64]bool)
	for _, v := range list.Rows {
		if v.CreateBy > 0 {
			userIdMap[v.CreateBy] = true
		}

		if v.UpdateBy > 0 {
			userIdMap[v.UpdateBy] = true
		}
	}

	// 格式化服务id
	userIds := make([]int64, 0)
	for k, _ := range userIdMap {
		if k > 0 {
			userIds = append(userIds, k)
		}
	}
	users := i.iUserDao.SelectByUserIds(c, userIds)
	userVoMap := make(map[int64]*systemModels.SysUserDML)
	for _, v := range users {
		userVoMap[v.UserId] = v
	}

	for k, v := range list.Rows {
		var createUserName string
		var updateUserName string
		userVo, ok := userVoMap[v.CreateBy]
		if ok {
			createUserName = userVo.UserName
		}

		userVo, ok = userVoMap[v.UpdateBy]
		if ok {
			updateUserName = userVo.UserName
		}
		list.Rows[k].CreateUserName = createUserName
		list.Rows[k].UpdateUserName = updateUserName
	}

	//
	return list, nil
}

func (i *IDataSetServiceImpl) DeleteDataSet(c *gin.Context, id int64) error {
	info, err := i.iDataSetDao.GetById(c, id)
	if err != nil {
		zap.L().Error("删除知识库失败", zap.Error(err))
		return errors.New("删除知识库失败")
	}
	if info == nil {
		return errors.New("知识库不存在")
	}
	var deleteReq aiDataSetModels.SysDatasetDeleteReq
	deleteReq.Ids = []string{info.DatasetUUID}
	deleteReqContent, err := json.Marshal(deleteReq)
	if err != nil {
		zap.L().Error("删除知识库失败", zap.Error(err))
		return errors.New("删除知识库失败")
	}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/api/v1/datasets", i.config.Host), bytes.NewBuffer(deleteReqContent))
	if err != nil {
		zap.L().Error("删除知识库失败", zap.Error(err))
		return err
	}
	if req == nil {
		return errors.New("删除知识库失败")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+i.config.ApiKey)
	resp, err := i.client.Do(req)
	if err != nil {
		zap.L().Error("创建知识库失败", zap.Error(err))
		return errors.New("删除知识库失败")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("删除知识库失败", zap.Error(err))
		return errors.New("删除知识库失败")
	}
	var deleteResponse aiDataSetModels.DataSetCreateResponse
	err = json.Unmarshal(body, &deleteResponse)
	if err != nil {
		return errors.New("删除知识库失败")
	}
	if deleteResponse.Code != 0 {
		return errors.New(deleteResponse.Message)
	}
	err = i.iDataSetDao.DeleteByIds(c, []int64{
		info.DatasetID,
	})
	if err != nil {
		zap.L().Error("删除知识库失败", zap.Error(err))
		return errors.New("删除知识库失败")
	}
	return nil
}

func (i *IDataSetServiceImpl) SyncDataSet(c *gin.Context, id int64) error {
	info, err := i.iDataSetDao.GetById(c, id)
	if err != nil {
		zap.L().Error("同步知识库失败", zap.Error(err))
		return errors.New("同步知识库失败")
	}
	params := url.Values{}
	params.Add("id", info.DatasetUUID)

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/datasets?%s", i.config.Host, params.Encode()), nil)
	if err != nil {
		zap.L().Error("删除知识库失败", zap.Error(err))
		return err
	}
	if req == nil {
		return errors.New("删除知识库失败")
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+i.config.ApiKey)
	resp, err := i.client.Do(req)
	if err != nil {
		zap.L().Error("创建知识库失败", zap.Error(err))
		return errors.New("删除知识库失败")
	}

	if resp == nil {
		zap.L().Error("创建知识库失败 resp == nil")
		return errors.New("删除知识库失败")
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("删除知识库失败", zap.Error(err))
		return errors.New("删除知识库失败")
	}

	var list aiDataSetModels.DatasetListResponse
	err = json.Unmarshal(body, &list)
	if err != nil {
		zap.L().Error("同步知识库失败", zap.Error(err))
		return err
	}
	if list.Code != 0 {
		zap.L().Error(fmt.Sprintf("同步知识库失败: %s", list.Message))
		return errors.New(list.Message)
	}
	if len(list.Data) == 0 {
		return nil
	}
	i.iDataSetDao.UpdateData(c, id, &list.Data[0])
	return nil
}

func (i *IDataSetServiceImpl) GetInfoById(c *gin.Context, id int64) (*aiDataSetModels.SysDataset, error) {
	info, err := i.iDataSetDao.GetById(c, id)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, errors.New("知识库不存在")
	}
	if len(info.DatasetParserConfig) != 0 {
		var parseConfig aiDataSetModels.ParserConfig
		err := json.Unmarshal([]byte(info.DatasetParserConfig), &parseConfig)
		if err != nil {
			return nil, err
		}
		info.ParserConfig = &parseConfig
	}

	return info, nil
}
