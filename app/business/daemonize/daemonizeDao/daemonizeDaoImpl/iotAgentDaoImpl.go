package daemonizeDaoImpl

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"nova-factory-server/app/business/daemonize/daemonizeDao"
	"nova-factory-server/app/business/daemonize/daemonizeModels"
	"nova-factory-server/app/constant/agent"
	"nova-factory-server/app/constant/commonStatus"
	redisKey "nova-factory-server/app/constant/redis"
	"nova-factory-server/app/datasource/cache"
	"nova-factory-server/app/utils/baizeContext"
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
func (s *IotAgentDaoImpl) Create(ctx context.Context, doAgent *daemonizeModels.SysIotAgent) (data *daemonizeModels.SysIotAgent, err error) {
	ret := s.db.Table(s.tableName).Create(&daemonizeModels.SysIotAgent{
		ObjectID:   doAgent.ObjectID,
		Name:       doAgent.Name,
		Version:    doAgent.Version,
		ConfigUUID: doAgent.ConfigUUID,
		Ipv4:       doAgent.Ipv4,
		Ipv6:       doAgent.Ipv6,
		Username:   doAgent.Username,
		Password:   doAgent.Password,
		State:      doAgent.State,
	})
	if ret.Error != nil {
		zap.L().Error("create agent error", zap.Error(err))
		return nil, ret.Error
	}
	return doAgent, nil
}

// UpdateHeartBeat 更新心跳包
func (s *IotAgentDaoImpl) UpdateHeartBeat(ctx context.Context, data *daemonizeModels.SysIotAgent) error {
	ret := s.cache.ZAdd(ctx, fmt.Sprintf("%s%s", redisKey.AGENT_HEADETBEAT_CACHE, ""), redis.Z{
		Score:  float64(gtime.Now().Unix()),
		Member: data.ObjectID,
	})
	if ret.Err() != nil {
		return ret.Err()
	}
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
func (s *IotAgentDaoImpl) UpdateConfig(ctx context.Context, configUuid string, objectIdList []uint64) (err error) {
	model := s.db.Table(s.tableName)
	if len(objectIdList) > 0 {
		model = model.Where("object_id in (?)", objectIdList)
	}
	ret := model.Where("config_uuid != ?", configUuid).Update("config_uuid", configUuid)
	if ret.Error != nil {
		zap.L().Error("agent[%v] update config[%v] error: %v", zap.Error(ret.Error))
		return ret.Error
	}
	return
}

// UpdateOperateState 更新agent配置
//func (s *IotAgentDaoImpl) UpdateOperateState(ctx context.Context, objectId uint64, operateState uint8) (err error) {
//	ret := s.db.Table(s.tableName).Where("object_id = ?", objectId).Updates(gdb.Map{
//		"operate_state": operateState,
//		"operate_time":  gtime.Now(),
//	})
//	if ret.Error != nil {
//		return ret.Error
//	}
//	return nil
//}

// GetByObjectIds 根据objectIds获取agent信息
//func (s *IotAgentDaoImpl) GetByObjectIds(ctx context.Context, objectIds []uint64) (agent []daemonizeModels.SysIotAgent, err error) {
//	var agents *[]daemonizeModels.SysIotAgent
//	//ret := s.db.Table(s.tableName).Where("object in (?)", objectIds).
//	//	Where("state = ?", commonStatus.NORMAL).Find(&agents)
//	//if agents == nil {
//	//	return make([]daemonizeModels.SysIotAgent, 0), nil
//	//}
//	//for k, agentInfo := range *agents {
//	//	score, err := s.getLastHeartBeatTime(ctx, ctx.Value(common.Cid).(string), agentInfo.ObjectId)
//	//	if err != nil {
//	//		zap.L().Error("get agent heart beat error: %s", zap.Error(err))
//	//	}
//	//	(*agents)[k].LastHeartbeatTime = gtime.NewFromTimeStamp(score)
//	//}
//	//if ret.Error != nil {
//	//	zap.L().Error("agent id[%v] query db error: %v", zap.Error(ret.Error))
//	//	return nil, ret.Error
//	//}
//	return *agents, nil
//}

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
	//offlineTime := time.Now().Unix() - int64(agent.CHECK_ONLINE_DURATION)
	//key := fmt.Sprintf("%s%s", redis.AGENT_HEADETBEAT_CACHE, cid)
	//cmd := s.cache.ZRangeByScore(ctx, key, &redis2.ZRangeBy{
	//	Min: (page-1)*size,
	//})
	//value, err := g.Redis().Do(ctx, "ZRANGEBYSCORE", key, fmt.Sprintf("(%d", offlineTime), "+inf", "LIMIT", (page-1)*size, page*size)
	//if err != nil {
	//	return nil, err
	//}
	//fmt.Println(value)
	//return value.Strings(), nil
	return []string{}, nil
}

// // getOffLineAgentId 读取离线的agent
//
//	func (s *IotAgentDaoImpl) getOffLineAgentId(ctx context.Context, cid string, page int, size int) ([]string, error) {
//		offlineTime := time.Now().Unix() - int64(common.CHECK_ONLINE_DURATION)
//		key := fmt.Sprintf("%s%s", common.AGENT_HEADETBEAT_CACHE, cid)
//		value, err := g.Redis().Do(ctx, "ZRANGEBYSCORE", key, fmt.Sprintf("(%d", offlineTime), "+inf", "LIMIT", (page-1)*size, page*size)
//		if err != nil {
//			return nil, err
//		}
//		fmt.Println(value)
//		return value.Strings(), nil
//	}
//
// // getOffLineAgentId 读取离线的agent
//
//	func (s *IotAgentDaoImpl) getLastHeartBeatTime(ctx context.Context, cid string, object uint64) (int64, error) {
//		key := fmt.Sprintf("%s%s", common.AGENT_HEADETBEAT_CACHE, cid)
//		score, err := g.Redis().ZScore(ctx, key, object)
//		if err != nil {
//			return 0, err
//		}
//		return int64(score), nil
//	}
//

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
	ret := model.Count(&total)
	if ret.Error != nil {
		return &daemonizeModels.SysIotAgentListData{
			Rows:  make([]*daemonizeModels.SysIotAgent, 0),
			Total: 0,
		}, ret.Error
	}

	model = model.Where("state", commonStatus.NORMAL)
	model = baizeContext.GetGormDataScope(ctx, model)
	ret = model.Offset(offset).Order("create_time desc").Limit(size).Find(&agentList)
	if ret.Error != nil {
		zap.L().Error("agent list query db error: %v", zap.Error(ret.Error))
		return &daemonizeModels.SysIotAgentListData{
			Rows:  make([]*daemonizeModels.SysIotAgent, 0),
			Total: 0,
		}, ret.Error
	}

	return &daemonizeModels.SysIotAgentListData{
		Rows:  agentList,
		Total: total,
	}, nil
}

func (i *IotAgentDaoImpl) Remove(c *gin.Context, ids []string) error {
	ret := i.db.Table(i.tableName).Where("object_id in (?)", ids).Update("state", commonStatus.DELETE)
	return ret.Error
}
