package aiDataSetDaoImpl

import (
	"nova-factory-server/app/business/ai/agent/aiDataSetDao"
	"nova-factory-server/app/business/ai/agent/aiDataSetModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IAiModelProviderDaoImpl struct {
	db            *gorm.DB
	providerTable string
	llmDao        aiDataSetDao.IAiLLMDao
}

func NewIAiModelProviderDaoImpl(db *gorm.DB, llmDao aiDataSetDao.IAiLLMDao) aiDataSetDao.IAiModelProviderDao {
	return &IAiModelProviderDaoImpl{
		db:            db,
		providerTable: "ai_model_provider",
		llmDao:        llmDao,
	}
}

func (i *IAiModelProviderDaoImpl) ListWithLLM(c *gin.Context, req *aiDataSetModels.SysAiModelProviderListReq) (*aiDataSetModels.SysAiModelProviderListData, error) {
	db := i.db.Table(i.providerTable).Where("state = ?", commonStatus.NORMAL)
	if req != nil && req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req != nil && req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}
	size := 20
	page := int64(1)
	if req != nil && req.Size > 0 {
		size = int(req.Size)
	}
	if req != nil && req.Page > 0 {
		page = req.Page
	}
	offset := int((page - 1) * int64(size))
	db = baizeContext.GetGormDataScope(c, db)

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, err
	}
	rows := make([]*aiDataSetModels.SysAiModelProvider, 0)
	if err := db.Order("rank DESC,id DESC").Offset(offset).Limit(size).Find(&rows).Error; err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return &aiDataSetModels.SysAiModelProviderListData{
			Rows:  rows,
			Total: total,
		}, nil
	}

	names := make([]string, 0, len(rows))
	for _, item := range rows {
		names = append(names, item.Name)
	}
	llms, err := i.llmDao.ListByFIDs(c, names)
	if err != nil {
		return nil, err
	}
	group := make(map[string][]*aiDataSetModels.SysAiLLM)
	for _, llm := range llms {
		group[llm.Fid] = append(group[llm.Fid], llm)
	}
	for _, item := range rows {
		item.LLMs = group[item.Name]
	}
	return &aiDataSetModels.SysAiModelProviderListData{
		Rows:  rows,
		Total: total,
	}, nil
}
