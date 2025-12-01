package craftRouteDaoImpl

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"nova-factory-server/app/business/craft/craftRouteDao"
	"nova-factory-server/app/business/craft/craftRouteModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
)

type ProcessContextDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewProcessContextDaoImpl(db *gorm.DB) craftRouteDao.IProcessContextDao {
	return &ProcessContextDaoImpl{
		db:        db,
		tableName: "sys_pro_process_content",
	}
}

func (p *ProcessContextDaoImpl) Add(c *gin.Context, processContext *craftRouteModels.SysProSetProcessContent) (*craftRouteModels.SysProProcessContent, error) {
	context := craftRouteModels.NewSysProProcessContent(processContext)
	context.ContentID = uint64(snowflake.GenID())
	context.SetCreateBy(baizeContext.GetUserId(c))
	if processContext.ControlRules != nil {
		content, err := json.Marshal(processContext.ControlRules)
		if err != nil {
			zap.L().Error("json marshal error", zap.Error(err))
		}
		context.Extension = string(content)
	}
	ret := p.db.Table(p.tableName).Create(context)
	return context, ret.Error
}

func (p *ProcessContextDaoImpl) Update(c *gin.Context, processContext *craftRouteModels.SysProSetProcessContent) (*craftRouteModels.SysProProcessContent, error) {
	context := craftRouteModels.NewSysProProcessContent(processContext)
	context.SetUpdateBy(baizeContext.GetUserId(c))
	content, err := json.Marshal(&processContext.ControlRules)
	if err != nil {
		zap.L().Error("json marshal error", zap.Error(err))
	}
	context.Extension = string(content)
	ret := p.db.Table(p.tableName).Where("content_id = ?", processContext.ContentID).Updates(context)
	return context, ret.Error
}

func (p *ProcessContextDaoImpl) Remove(c *gin.Context, ids []string) error {
	ret := p.db.Table(p.tableName).Where("content_id in (?)", ids).Update("state", commonStatus.DELETE)
	return ret.Error
}

func (p *ProcessContextDaoImpl) List(c *gin.Context, req *craftRouteModels.SysProProcessContextListReq) (*craftRouteModels.SysProProcessContextListData, error) {
	db := p.db.Table(p.tableName).Debug()

	db = db.Where("process_id = ?", req.ProcessID)

	size := 0
	if req == nil || req.Size <= 0 {
		size = 20
	} else {
		size = int(req.Size)
	}
	offset := 0
	if req == nil || req.Page <= 0 {
		req.Page = 1
	} else {
		offset = int((req.Page - 1) * req.Size)
	}
	db = db.Where("state", commonStatus.NORMAL)
	db = baizeContext.GetGormDataScope(c, db)
	var dto []*craftRouteModels.SysProProcessContent

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &craftRouteModels.SysProProcessContextListData{
			Rows:  make([]*craftRouteModels.SysProProcessContent, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &craftRouteModels.SysProProcessContextListData{
			Rows:  make([]*craftRouteModels.SysProProcessContent, 0),
			Total: 0,
		}, ret.Error
	}
	for k, _ := range dto {
		if dto[k].Extension != "" {
			var triggerRule craftRouteModels.ControlRule
			err := json.Unmarshal([]byte(dto[k].Extension), &triggerRule)
			if err != nil {
				zap.L().Error("json unmarshal error", zap.Error(err))
			}
			if triggerRule.TriggerRules == nil {
				triggerRule.TriggerRules = &craftRouteModels.TriggerRules{
					Actions: make([]craftRouteModels.ControllerAction, 0),
					Cases:   make([]craftRouteModels.TriggerCase, 0),
				}
			}

			if triggerRule.PidRules == nil {
				triggerRule.PidRules = &craftRouteModels.PidRules{
					Actions: make([]craftRouteModels.ControllerAction, 0),
				}
			}

			if triggerRule.PredictRules == nil {
				triggerRule.PredictRules = &craftRouteModels.PredictRules{
					Actions: make([]craftRouteModels.ControllerAction, 0),
					Cases:   make([]craftRouteModels.TriggerCase, 0),
				}
			}
			dto[k].ControlRules = &triggerRule
		}
	}
	return &craftRouteModels.SysProProcessContextListData{
		Rows:  dto,
		Total: total,
	}, nil
}

func (p *ProcessContextDaoImpl) GetByProcessIds(c *gin.Context, ids []int64) ([]*craftRouteModels.SysProProcessContent, error) {
	var data []*craftRouteModels.SysProProcessContent
	ret := p.db.Table(p.tableName).Where("process_id in (?)", ids).Find(&data)
	return data, ret.Error
}
