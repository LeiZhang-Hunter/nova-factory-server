package aiDataSetServiceImpl

import (
	"errors"
	"nova-factory-server/app/business/ai/agent/aiDataSetDao"
	"nova-factory-server/app/business/ai/agent/aiDataSetModels"
	"nova-factory-server/app/business/ai/agent/aiDataSetService"
	"strings"

	"github.com/gin-gonic/gin"
)

type IAiLLMSettingServiceImpl struct {
	dao    aiDataSetDao.IAiLLMSettingDao
	llmDao aiDataSetDao.IAiLLMDao
}

func NewIAiLLMSettingServiceImpl(dao aiDataSetDao.IAiLLMSettingDao, llmDao aiDataSetDao.IAiLLMDao) aiDataSetService.IAiLLMSettingService {
	return &IAiLLMSettingServiceImpl{
		dao:    dao,
		llmDao: llmDao,
	}
}

func (i *IAiLLMSettingServiceImpl) Set(c *gin.Context, req *aiDataSetModels.SetSysAiLLMSetting) (*aiDataSetModels.SysAiLLMSetting, error) {
	if strings.TrimSpace(req.PublicKey) == "" {
		return nil, errors.New("api key不能为空")
	}
	if strings.TrimSpace(req.LlmID) == "" || strings.TrimSpace(req.EmbdID) == "" ||
		strings.TrimSpace(req.AsrID) == "" || strings.TrimSpace(req.Img2txtID) == "" ||
		strings.TrimSpace(req.RerankID) == "" || strings.TrimSpace(req.ParserIDs) == "" {
		return nil, errors.New("模型配置不能为空")
	}
	ids := collectModelIDs(req)
	existed, err := i.llmDao.ListExistingLLMNames(c, ids)
	if err != nil {
		return nil, err
	}
	existMap := make(map[string]bool, len(existed))
	for _, v := range existed {
		existMap[v] = true
	}
	for _, id := range ids {
		if !existMap[id] {
			return nil, errors.New("模型不存在: " + id)
		}
	}
	return i.dao.Set(c, req)
}

func (i *IAiLLMSettingServiceImpl) Get(c *gin.Context, req *aiDataSetModels.GetSysAiLLMSettingReq) (*aiDataSetModels.SysAiLLMSetting, error) {
	return i.dao.Get(c, req)
}

func collectModelIDs(req *aiDataSetModels.SetSysAiLLMSetting) []string {
	ids := []string{req.LlmID, req.EmbdID, req.AsrID, req.Img2txtID, req.RerankID, req.TtsID}
	if req.ParserIDs != "" {
		ids = append(ids, strings.Split(req.ParserIDs, ",")...)
	}
	m := make(map[string]bool)
	ret := make([]string, 0, len(ids))
	for _, id := range ids {
		id = strings.TrimSpace(id)
		if id == "" || m[id] {
			continue
		}
		m[id] = true
		ret = append(ret, id)
	}
	return ret
}
