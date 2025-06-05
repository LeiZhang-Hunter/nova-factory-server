package daemonizeDao

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

func NewIotAgentDaoImpl(db *gorm.DB, cache cache.Cache) *IotAgentDaoImpl {
	return &IotAgentDaoImpl{
		db:        db,
		tableName: "sys_iot_agent",
		cache:     cache,
	}
}

// GetByObjectId 根据objectId获取agent信息
func (s *IotAgentDaoImpl) GetByObjectId(ctx context.Context, objectId uint64) (agent *daemonizeModels.SysIotAgent, err error) {
	var data *daemonizeModels.SysIotAgent
	ret := s.db.Table(s.tableName).Where("object_id = ?", objectId).Where("status = ?", commonStatus.NORMAL).First(&data)
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

// GetList 获取agent列表
func (s *IotAgentDaoImpl) GetList(ctx *gin.Context, req *daemonizeModels.SysIotAgentListReq, state int, pageNum, pageSize int) (data *daemonizeModels.SysIotAgentListData, err error) {
	agentList := make([]*daemonizeModels.SysIotAgent, 0)
	model := s.db.Table(s.tableName)
	if req.Name != "" {
		model = model.Where("name  LIKE ?", "%"+req.Name+"%")
	}
	if req.Version != "" {
		model = model.Where("version  = ?", req.Version)
	}
	if req.ConfigUUID != "" {
		model = model.Where("config_uuid = ?", req.ConfigUUID)
	}
	model = model.Where("state", commonStatus.NORMAL)
	model = baizeContext.GetGormDataScope(ctx, model)
	// 在线和离线状态根据心跳时间判断，超过一分钟没有心跳为离线
	//offlineTime := timestamppb.New(time.Now().Truncate(time.Minute))
	//if state == consts.AgentStateOnline {
	//	model = model.Where("last_heartbeat_time > ?", gconv.String(offlineTime))
	//} else if state == consts.AgentStateOffline {
	//	model = model.Where("last_heartbeat_time <= ?", gconv.String(offlineTime))
	//}
	size := 0
	if req.Size <= 0 {
		size = 20
	} else {
		size = int(req.Size)
	}
	offset := 0
	if req.Page <= 0 {
		req.Page = 1
	} else {
		offset = int((req.Page - 1) * req.Size)
	}
	offset = (pageNum - 1) * pageSize

	var total int64
	ret := model.Count(&total)
	if ret.Error != nil {
		return &daemonizeModels.SysIotAgentListData{
			Rows:  make([]*daemonizeModels.SysIotAgent, 0),
			Total: 0,
		}, ret.Error
	}

	ret = model.Offset(offset).Order("create_time desc").Limit(size).Find(&agentList)
	if ret.Error != nil {
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
func (s *IotAgentDaoImpl) Create(ctx context.Context, doAgent *daemonizeModels.SysIotAgent) (id int64, err error) {
	ret := s.db.Table(s.tableName).Create(&daemonizeModels.SysIotAgent{
		ObjectID:    doAgent.ObjectID,
		CompanyUUID: doAgent.CompanyUUID,
		Name:        doAgent.Name,
		Version:     doAgent.Version,
		ConfigUUID:  doAgent.ConfigUUID,
		Ipv4:        doAgent.Ipv4,
		Ipv6:        doAgent.Ipv6,
		State:       doAgent.State,
	})
	if ret != nil {
		g.Log().Errorf(ctx, "agent[%+v] insert to db error: %v", doAgent, err.Error())
		return doAgent.ID, ret.Error
	}
	return
}

// UpdateHeartBeat 更新心跳包
func (s *IotAgentDaoImpl) UpdateHeartBeat(ctx context.Context, data *daemonizeModels.SysIotAgent) error {
	ret := s.cache.ZAdd(ctx, fmt.Sprintf("%s%s", redisKey.AGENT_HEADETBEAT_CACHE, data.CompanyUUID), redis.Z{
		Score:  float64(gtime.Now().Unix()),
		Member: data.ObjectID,
	})
	if ret.Err() != nil {
		return ret.Err()
	}
	return nil
}

// Update 更细agent信息
func (s *IotAgentDaoImpl) Update(ctx context.Context, data *daemonizeModels.SysIotAgent) (err error) {
	model := s.db.Table(s.tableName)
	doAgent := &daemonizeModels.SysIotAgent{}

	if data.CompanyUUID != "" {
		doAgent.CompanyUUID = data.CompanyUUID
	}
	if data.Name != "" {
		doAgent.Name = data.Name
	}
	if data.Version != "" {
		doAgent.Version = data.Version
	}
	if data.ConfigUUID != "" {
		doAgent.ConfigUUID = data.ConfigUUID
	}
	//if data.LastHeartbeatTime !=  {
	//	doAgent.LastHeartbeatTime = time.New(agent.LastHeartbeatTime.AsTime())
	//}
	doAgent.State = data.State

	if data.ID > 0 {
		model = model.Where("id", data.ID)
	} else if data.ObjectID > 0 {
		model = model.Where("object_id", data.ObjectID)
	}
	ret := model.Update("last_heartbeat_time", data.LastHeartbeatTime)
	if ret.Error != nil {
		zap.L().Error("agent[%+v] update db error: %v", zap.Error(err))
		return ret.Error
	}
	return
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

// GetEntityList 获取agent列表
//func (s *IotAgentDaoImpl) GetEntityList(ctx context.Context, agent *daemonizeModels.SysIotAgent, state int, pageNum, pageSize int) (agentList []*entity.DoAgent, total int, err error) {
//	agentList = make([]*entity.DoAgent, 0)
//	model := dao.DoAgent.Ctx(ctx).OmitEmpty()
//	if agent.Name != "" {
//		model = model.WhereLike("name", "%"+agent.Name+"%")
//	}
//	if agent.Version != "" {
//		model = model.Where(entity.DoAgent{
//			Version: agent.Version,
//		})
//	}
//	if agent.ConfigUuid != "" {
//		model = model.Where(entity.DoAgent{
//			ConfigUuid: agent.ConfigUuid,
//		})
//	}
//	model = model.Where(entity.DoAgent{
//		CompanyUuid: ctx.Value(common.Cid).(string),
//	})
//
//	// 在线和离线状态根据心跳时间判断，超过一分钟没有心跳为离线,从redis有序集合里查找
//	if state == consts.AgentStateOnline {
//		ids, err := s.getOnLineAgentId(ctx, ctx.Value(common.Cid).(string), pageNum, pageSize)
//		if err != nil {
//			g.Log().Errorf(ctx, err.Error())
//		} else {
//			if len(ids) == 0 {
//				return []*entity.DoAgent{}, 0, nil
//			}
//			model = model.WhereIn(dao.DoAgent.Columns().ObjectId, ids)
//		}
//	} else if state == consts.AgentStateOffline {
//		ids, err := s.getOffLineAgentId(ctx, ctx.Value(common.Cid).(string), pageNum, pageSize)
//		if err != nil {
//			g.Log().Errorf(ctx, err.Error())
//		} else {
//			model = model.WhereIn(dao.DoAgent.Columns().ObjectId, ids)
//		}
//	}
//	if pageNum <= 0 {
//		pageNum = 1
//	}
//	if pageSize <= 0 || pageSize > 50 {
//		pageSize = 50
//	}
//	offset := (pageNum - 1) * pageSize
//	err = model.Limit(pageSize).Offset(offset).
//		ScanAndCount(&agentList, &total, true)
//	if err != nil {
//		g.Log().Errorf(ctx, "agent list query db error: %v", err.Error())
//		return nil, 0, gerror.NewCode(gcode.CodeDbOperationError, "agent list query db error")
//	}
//	for _, agentInfo := range agentList {
//		score, err := s.getLastHeartBeatTime(ctx, ctx.Value(common.Cid).(string), agentInfo.ObjectId)
//		if err != nil {
//			g.Log().Errorf(ctx, "get agent heart beat error: %s", err.Error())
//		}
//		agentInfo.LastHeartbeatTime = gtime.NewFromTimeStamp(score)
//	}
//	return
//}
