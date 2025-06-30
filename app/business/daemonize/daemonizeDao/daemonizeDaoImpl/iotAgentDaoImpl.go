package daemonizeDaoImpl

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/v2/database/gdb"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"nova-factory-server/app/business/daemonize/daemonizeDao"
	"nova-factory-server/app/business/daemonize/daemonizeModels"
	"nova-factory-server/app/constant/agent"
	"nova-factory-server/app/constant/commonStatus"
	redisKey "nova-factory-server/app/constant/redis"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/baizeContext"
	"strconv"
	"time"
)

type IotAgentDaoImpl struct {
	tableName string
	db        *gorm.DB
	cache     cache.Cache
}

func NewIotAgentDaoImpl(db *gorm.DB, cache cache.Cache) daemonizeDao.IotAgentDao {
	return &IotAgentDaoImpl{
		db:        db,
		tableName: "sys_iot_agent",
		cache:     cache,
	}
}

// GetByObjectId 根据objectId获取agent信息
func (s *IotAgentDaoImpl) GetByObjectId(ctx context.Context, objectId uint64) (agent *daemonizeModels.SysIotAgent, err error) {
	var data *daemonizeModels.SysIotAgent
	ret := s.db.Table(s.tableName).Where("object_id = ?", objectId).Where("state = ?", commonStatus.NORMAL).First(&data)
	if ret.Error != nil {
		if ret.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return data, ret.Error
	}
	return data, nil
}

// GetAgentStateByLastHeartbeatTime 根据上次心跳时间获取agent状态
func (s *IotAgentDaoImpl) GetAgentStateByLastHeartbeatTime(lastHeartbeatTime *gtime.Time) (state int) {
	duration := time.Now().Unix() - lastHeartbeatTime.Unix()
	// 心跳时间小于十分钟都认为是在线
	if duration < int64(agent.CHECK_ONLINE_DURATION) {
		return agent.AgentStateOnline
	} else {
		return agent.AgentStateOffline
	}
}

// GetDistinctVersionList 获取版本号列表
func (s *IotAgentDaoImpl) GetDistinctVersionList(ctx context.Context) (versionList []string, err error) {
	doAgentList := make([]*daemonizeModels.SysIotAgent, 0)
	ret := s.db.Table(s.tableName).Where("state = ?", commonStatus.NORMAL).Select("version").Distinct().Find(&doAgentList)
	if ret.Error != nil {
		zap.L().Error("query agent version list error, err: %v", zap.Error(err))
		return nil, ret.Error
	}
	versionList = make([]string, 0, len(doAgentList))
	for _, doAgent := range doAgentList {
		versionList = append(versionList, doAgent.Version)
	}
	return versionList, nil
}

// Create 创建agent
func (s *IotAgentDaoImpl) Create(ctx context.Context, doAgent *daemonizeModels.SysIotAgent) (*daemonizeModels.SysIotAgent, error) {
	ret := s.db.Table(s.tableName).Create(doAgent)
	if ret.Error != nil {
		zap.L().Error("create agent error", zap.Error(ret.Error))
		return nil, ret.Error
	}
	return doAgent, nil
}

// UpdateHeartBeat 更新心跳包
func (s *IotAgentDaoImpl) UpdateHeartBeat(ctx context.Context, data *daemonizeModels.SysIotAgent) error {
	if data.ObjectID == 0 {
		return nil
	}
	s.cache.Set(ctx, fmt.Sprintf("%s%d", redisKey.AGENT_HEADETBEAT_CACHE, data.ObjectID), fmt.Sprintf("%d", gtime.Now().Unix()),
		300*time.Second)
	return nil
}

// Update 更细agent信息
func (s *IotAgentDaoImpl) Update(ctx context.Context, data *daemonizeModels.SysIotAgent) (*daemonizeModels.SysIotAgent, error) {
	model := s.db.Table(s.tableName)

	model = model.Where("object_id", data.ObjectID)
	ret := model.Updates(data)
	if ret.Error != nil {
		zap.L().Error("agent[%+v] update db error: %v", zap.Error(ret.Error))
		return nil, ret.Error
	}
	return data, ret.Error
}

// DeleteByObjectId 根据objectId删除agent
func (s *IotAgentDaoImpl) DeleteByObjectId(ctx context.Context, objectId uint64) error {
	ret := s.db.Table(s.tableName).Where("object_id = ?", objectId).Update("state", commonStatus.DELETE)
	if ret.Error != nil {
		zap.L().Error("agent[%+v] delete db error: %v", zap.Error(ret.Error))
		return ret.Error
	}
	return nil
}

// DeleteByObjectIdList 根据objectId删除agent
func (s *IotAgentDaoImpl) DeleteByObjectIdList(ctx context.Context, objectIdList []uint64) error {
	ret := s.db.Table(s.tableName).Where("object_id  in (?)", objectIdList).Update("state", commonStatus.DELETE)
	if ret.Error != nil {
		zap.L().Error("agent[%v] delete db error: %v", zap.Error(ret.Error))
		return ret.Error
	}
	return nil
}

// UpdateConfig 更新agent配置
func (s *IotAgentDaoImpl) UpdateConfig(ctx context.Context, configId uint64, objectIdList []uint64) (err error) {
	model := s.db.Table(s.tableName)
	if len(objectIdList) > 0 {
		model = model.Where("object_id in (?)", objectIdList)
	} else {
		return errors.New("object id is empty")
	}
	ret := model.Update("config_id", configId)
	if ret.Error != nil {
		zap.L().Error("agent[%v] update config[%v] error: %v", zap.Error(ret.Error))
		return ret.Error
	}
	return
}

// UpdateLastConfig 更新agent配置
func (s *IotAgentDaoImpl) UpdateLastConfig(ctx context.Context, configId uint64, objectIdList []uint64) (err error) {
	model := s.db.Table(s.tableName)
	if len(objectIdList) > 0 {
		model = model.Where("object_id in (?)", objectIdList)
	} else {
		return errors.New("object id is empty")
	}
	ret := model.Update("last_config_id", configId)
	if ret.Error != nil {
		zap.L().Error("agent[%v] update config[%v] error: %v", zap.Error(ret.Error))
		return ret.Error
	}
	return
}

// UpdateOperateStateByObjectIds 更新agent配置
func (s *IotAgentDaoImpl) UpdateOperateStateByObjectIds(ctx context.Context, objectId []uint64, operateState uint8) (err error) {
	model := s.db.Table(s.tableName)
	ret := model.Where("object_id in (?)", objectId).Updates(gdb.Map{
		"operate_state": operateState,
		"Operate_time":  gtime.Now(),
	})
	if ret.Error != nil {
		zap.L().Error("agent update state[%v] error: %v", zap.Error(ret.Error))
		return ret.Error
	}
	return nil
}

// getOffLineAgentId 读取在线的agent
func (s *IotAgentDaoImpl) getOnLineAgentId(ctx context.Context, cid string, page int, size int) ([]string, error) {

	return []string{}, nil
}

// GetAgentList 获取agent列表
func (s *IotAgentDaoImpl) GetAgentList(ctx *gin.Context, req *daemonizeModels.SysIotAgentListReq) (*daemonizeModels.SysIotAgentListData, error) {
	agentList := make([]*daemonizeModels.SysIotAgent, 0)
	model := s.db.Table(s.tableName)
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

	var total int64
	model = model.Where("state", commonStatus.NORMAL)
	ret := model.Count(&total)
	if ret.Error != nil {
		return &daemonizeModels.SysIotAgentListData{
			Rows:  make([]*daemonizeModels.SysIotAgent, 0),
			Total: 0,
		}, ret.Error
	}

	model = baizeContext.GetGormDataScope(ctx, model)
	ret = model.Offset(offset).Order("create_time desc").Limit(size).Find(&agentList)
	if ret.Error != nil {
		zap.L().Error("agent list query db error: %v", zap.Error(ret.Error))
		return &daemonizeModels.SysIotAgentListData{
			Rows:  make([]*daemonizeModels.SysIotAgent, 0),
			Total: 0,
		}, ret.Error
	}

	var keys []string = make([]string, 0)
	for _, data := range agentList {
		keys = append(keys, fmt.Sprintf("%s%d", redisKey.AGENT_HEADETBEAT_CACHE, data.ObjectID))
	}

	var heartBeatMap map[string]int64 = make(map[string]int64)
	slice := s.cache.MGet(ctx, keys)
	if slice != nil {
		values := slice.Val()
		for index, key := range keys {
			if values[index] == nil {
				heartBeatMap[key] = 0
			} else {
				v := values[index].(string)
				parseInt, err := strconv.ParseInt(v, 10, 64)
				if err != nil {
					zap.L().Error("agent get heart beat map error: %v", zap.Error(err))
					continue
				}
				heartBeatMap[key] = parseInt
			}

		}
	}
	offLineCount := 0
	onLineCount := 0
	for k, data := range agentList {
		key := fmt.Sprintf("%s%d", redisKey.AGENT_HEADETBEAT_CACHE, data.ObjectID)
		heartBeat, ok := heartBeatMap[key]
		if !ok {
			continue
		}
		duration := time.Now().Unix() - heartBeat
		if duration < int64(agent.CHECK_ONLINE_DURATION) {
			onLineCount++
			agentList[k].Active = agent.ONLINE
		} else {
			offLineCount++
			agentList[k].Active = agent.OFFLINE
		}
		heartBeatTime := time.Unix(heartBeat, 0)
		agentList[k].LastHeartbeatTime = &heartBeatTime
	}

	return &daemonizeModels.SysIotAgentListData{
		Rows:         agentList,
		Total:        total,
		OffLineCount: int64(offLineCount),
		OnLineCount:  int64(onLineCount),
	}, nil
}

func (i *IotAgentDaoImpl) Remove(c *gin.Context, ids []string) error {
	ret := i.db.Table(i.tableName).Where("object_id in (?)", ids).Update("state", commonStatus.DELETE)
	return ret.Error
}
