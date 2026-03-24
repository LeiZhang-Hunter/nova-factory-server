package aiDataSetDaoImpl

import (
	"nova-factory-server/app/business/ai/agent/aiDataSetDao"
	"nova-factory-server/app/business/ai/agent/aiDataSetModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IAiLLMSettingDaoImpl struct {
	db    *gorm.DB
	table string
}

func NewIAiLLMSettingDaoImpl(db *gorm.DB) aiDataSetDao.IAiLLMSettingDao {
	return &IAiLLMSettingDaoImpl{
		db:    db,
		table: "ai_llm_setting",
	}
}

func (i *IAiLLMSettingDaoImpl) Set(c *gin.Context, req *aiDataSetModels.SetSysAiLLMSetting) (*aiDataSetModels.SysAiLLMSetting, error) {
	data := &aiDataSetModels.SysAiLLMSetting{
		ID:        req.ID,
		Name:      req.Name,
		PublicKey: req.PublicKey,
		LlmID:     req.LlmID,
		EmbdID:    req.EmbdID,
		AsrID:     req.AsrID,
		Img2txtID: req.Img2txtID,
		RerankID:  req.RerankID,
		TtsID:     req.TtsID,
		ParserIDs: req.ParserIDs,
		Credit:    req.Credit,
		Status:    req.Status,
	}
	if data.ID == 0 {
		data.ID = snowflake.GenID()
		data.DeptID = baizeContext.GetDeptId(c)
		data.State = commonStatus.NORMAL
		data.SetCreateBy(baizeContext.GetUserId(c))
		if err := i.db.Table(i.table).Create(data).Error; err != nil {
			return nil, err
		}
		return data, nil
	}
	data.SetUpdateBy(baizeContext.GetUserId(c))
	if err := i.db.Table(i.table).Where("id = ?", data.ID).Where("state = ?", commonStatus.NORMAL).
		Select("name", "public_key", "llm_id", "embd_id", "asr_id", "img2txt_id", "rerank_id", "tts_id", "parser_ids", "credit", "status", "update_by", "update_time").
		Updates(data).Error; err != nil {
		return nil, err
	}
	var ret aiDataSetModels.SysAiLLMSetting
	if err := i.db.Table(i.table).Where("id = ?", data.ID).First(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (i *IAiLLMSettingDaoImpl) Get(c *gin.Context, req *aiDataSetModels.GetSysAiLLMSettingReq) (*aiDataSetModels.SysAiLLMSetting, error) {
	var ret aiDataSetModels.SysAiLLMSetting
	db := i.db.Table(i.table).Where("state = ?", commonStatus.NORMAL)
	if req != nil && req.ID > 0 {
		db = db.Where("id = ?", req.ID)
	} else {
		db = db.Where("dept_id = ?", baizeContext.GetDeptId(c))
	}
	if err := db.Order("update_time desc").First(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}
