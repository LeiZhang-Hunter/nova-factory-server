package aiDataSetDaoImpl

import (
	"errors"
	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/ai/agent/aiDataSetDao"
	"nova-factory-server/app/business/ai/agent/aiDataSetModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/snowflake"
	"time"

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

	rows := make([]*aiDataSetModels.SysAiModelProvider, 0)
	if err := db.Order("`rank` DESC,id DESC").Find(&rows).Error; err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return &aiDataSetModels.SysAiModelProviderListData{
			Rows:  rows,
			Total: 0,
		}, nil
	}

	names := make([]string, 0, len(rows))
	for _, item := range rows {
		names = append(names, item.Name)
	}

	return &aiDataSetModels.SysAiModelProviderListData{
		Rows:  rows,
		Total: 0,
	}, nil
}

func (i *IAiModelProviderDaoImpl) UpsertFactoryProvider(tx *gorm.DB, item *aiDataSetModels.FactoryProviderUpsert, status int32, rank int32, now time.Time) (int64, error) {
	var row aiDataSetModels.AiModelProviderEntity
	kt := time.Now()
	err := tx.Table(i.providerTable).Where("name = ?", item.Name).Where("state = 0").Take(&row).Error
	if err == nil {
		update := aiDataSetModels.AiModelProviderEntity{
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
		update := aiDataSetModels.AiModelProviderEntity{
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
	create := aiDataSetModels.AiModelProviderEntity{
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
