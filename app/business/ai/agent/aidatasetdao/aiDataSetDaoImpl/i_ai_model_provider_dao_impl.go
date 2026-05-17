package aiDataSetDaoImpl

import (
	"errors"
	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/ai/agent/aidatasetdao"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/snowflake"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IAiModelProviderDaoImpl struct {
	db            *gorm.DB
	providerTable string
	llmDao        aidatasetdao.IAiLLMDao
}

func NewIAiModelProviderDaoImpl(db *gorm.DB, llmDao aidatasetdao.IAiLLMDao) aidatasetdao.IAiModelProviderDao {
	return &IAiModelProviderDaoImpl{
		db:            db,
		providerTable: "ai_model_provider",
		llmDao:        llmDao,
	}
}

func (i *IAiModelProviderDaoImpl) ListWithLLM(c *gin.Context, req *aidatasetmodels.SysAiModelProviderListReq) (*aidatasetmodels.SysAiModelProviderListData, error) {
	db := i.db.Table(i.providerTable).Where("state = ?", commonStatus.NORMAL)
	if req != nil && req.Name != "" {
		db = db.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req != nil && req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	rows := make([]*aidatasetmodels.SysAiModelProvider, 0)
	if err := db.Order("`rank` DESC,id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return &aidatasetmodels.SysAiModelProviderListData{
			Rows:  rows,
			Total: 0,
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
	attachLLMs(rows, llms)

	return &aidatasetmodels.SysAiModelProviderListData{
		Rows:  rows,
		Total: int64(len(rows)),
	}, nil
}

func (i *IAiModelProviderDaoImpl) ListEmbeddingWithLLM(c *gin.Context, req *aidatasetmodels.SysAiModelProviderListReq) (*aidatasetmodels.SysAiModelProviderListData, error) {
	data, err := i.ListWithLLM(c, req)
	if err != nil {
		return nil, err
	}
	if data == nil || len(data.Rows) == 0 {
		return data, nil
	}
	filtered := filterEmbeddingProviders(data.Rows)
	return &aidatasetmodels.SysAiModelProviderListData{
		Rows:  filtered,
		Total: int64(len(filtered)),
	}, nil
}

func attachLLMs(rows []*aidatasetmodels.SysAiModelProvider, llms []*aidatasetmodels.SysAiLLM) {
	llmMap := make(map[string][]*aidatasetmodels.SysAiLLM, len(rows))
	for _, llm := range llms {
		if llm == nil {
			continue
		}
		llmMap[llm.Fid] = append(llmMap[llm.Fid], llm)
	}
	for _, row := range rows {
		if row == nil {
			continue
		}
		row.LLMs = llmMap[row.Name]
		if row.LLMs == nil {
			row.LLMs = make([]*aidatasetmodels.SysAiLLM, 0)
		}
	}
}

func filterEmbeddingProviders(rows []*aidatasetmodels.SysAiModelProvider) []*aidatasetmodels.SysAiModelProvider {
	filtered := make([]*aidatasetmodels.SysAiModelProvider, 0, len(rows))
	for _, provider := range rows {
		if provider == nil {
			continue
		}
		llms := make([]*aidatasetmodels.SysAiLLM, 0, len(provider.LLMs))
		for _, llm := range provider.LLMs {
			if llm == nil {
				continue
			}
			if strings.EqualFold(strings.TrimSpace(llm.ModelType), "embedding") {
				llms = append(llms, llm)
			}
		}
		if len(llms) == 0 {
			continue
		}
		provider.LLMs = llms
		filtered = append(filtered, provider)
	}
	return filtered
}

func (i *IAiModelProviderDaoImpl) UpsertFactoryProvider(tx *gorm.DB, item *aidatasetmodels.FactoryProviderUpsert, status int32, rank int32, now time.Time) (int64, error) {
	var row aidatasetmodels.AiModelProviderEntity
	kt := time.Now()
	err := tx.Table(i.providerTable).Where("name = ?", item.Name).Where("state = 0").Take(&row).Error
	if err == nil {
		update := aidatasetmodels.AiModelProviderEntity{
			Logo:   item.Logo,
			Tags:   item.Tags,
			Status: status,
			Rank:   rank,
			BaseEntity: baize.BaseEntity{
				CreateTime: &kt,
				UpdateTime: &kt,
			},
		}
		return row.ID, tx.Table(i.providerTable).Where("id = ?", row.ID).
			Select("logo", "tags", "status", "rank", "update_time_db").Updates(&update).Error
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	err = tx.Table(i.providerTable).Where("name = ?", item.Name).Take(&row).Error
	if err == nil {
		update := aidatasetmodels.AiModelProviderEntity{
			Name:   item.Name,
			Logo:   item.Logo,
			Tags:   item.Tags,
			Status: status,
			Rank:   rank,
			State:  0,
			BaseEntity: baize.BaseEntity{
				CreateTime: &kt,
				UpdateTime: &kt,
			},
		}
		return row.ID, tx.Table(i.providerTable).Where("id = ?", row.ID).
			Select("name", "logo", "tags", "status", "rank", "state", "update_time_db").Updates(&update).Error
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}
	id := snowflake.GenID()
	create := aidatasetmodels.AiModelProviderEntity{
		ID:     id,
		Name:   item.Name,
		Logo:   item.Logo,
		Tags:   item.Tags,
		Status: status,
		Rank:   rank,
		State:  0,
		BaseEntity: baize.BaseEntity{
			CreateTime: &kt,
			UpdateTime: &kt,
		},
	}
	ret := tx.Table(i.providerTable).Create(&create)
	return id, ret.Error
}
