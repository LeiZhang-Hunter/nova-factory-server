package craftRouteServiceImpl

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
	"nova-factory-server/app/business/craft/craftRouteDao"
	"nova-factory-server/app/business/craft/craftRouteModels"
	v1 "nova-factory-server/app/business/craft/craftRouteModels/api/v1"
	"nova-factory-server/app/business/craft/craftRouteService"
	"nova-factory-server/app/constant/control"
	craft2 "nova-factory-server/app/constant/craft"
	"nova-factory-server/app/utils/uuid"
	"strconv"
)

type CraftRouteServiceImpl struct {
	dao             craftRouteDao.ICraftRouteDao
	processDao      craftRouteDao.IProcessDao
	routeProcessDao craftRouteDao.IRouteProcessDao
	contextDao      craftRouteDao.IProcessContextDao
	bomDao          craftRouteDao.ISysProRouteProductBomDao
	productDao      craftRouteDao.ISysProRouteProductDao
	routeConfigDao  craftRouteDao.ISysCraftRouteConfigDao
}

func NewCraftRouteServiceImpl(dao craftRouteDao.ICraftRouteDao,
	processDao craftRouteDao.IProcessDao,
	routeProcessDao craftRouteDao.IRouteProcessDao,
	contextDao craftRouteDao.IProcessContextDao,
	bomDao craftRouteDao.ISysProRouteProductBomDao,
	productDao craftRouteDao.ISysProRouteProductDao,
	routeConfigDao craftRouteDao.ISysCraftRouteConfigDao) craftRouteService.ICraftRouteService {
	return &CraftRouteServiceImpl{
		dao:             dao,
		processDao:      processDao,
		routeProcessDao: routeProcessDao,
		contextDao:      contextDao,
		bomDao:          bomDao,
		productDao:      productDao,
		routeConfigDao:  routeConfigDao,
	}
}

func (craft *CraftRouteServiceImpl) AddCraftRoute(c *gin.Context, route *craftRouteModels.SysCraftRouteRequest) (*craftRouteModels.SysCraftRoute, error) {
	return craft.dao.AddCraftRoute(c, route)
}

func (craft *CraftRouteServiceImpl) UpdateCraftRoute(c *gin.Context, route *craftRouteModels.SysCraftRouteRequest) (*craftRouteModels.SysCraftRoute, error) {
	return craft.dao.UpdateCraftRoute(c, route)
}

func (craft *CraftRouteServiceImpl) RemoveCraftRoute(c *gin.Context, ids []int64) error {
	return craft.dao.RemoveCraftRoute(c, ids)
}

func (craft *CraftRouteServiceImpl) SelectCraftRoute(c *gin.Context, req *craftRouteModels.SysCraftRouteListReq) (*craftRouteModels.SysCraftRouteListData, error) {
	return craft.dao.SelectCraftRoute(c, req)
}

func (craft *CraftRouteServiceImpl) DetailCraftRoute(c *gin.Context, req *craftRouteModels.SysCraftRouteDetailRequest) (*craftRouteModels.SysCraftRouteConfig, error) {
	// 读取详情
	info, err := craft.routeConfigDao.GetById(uint64(req.RouteID))
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (craft *CraftRouteServiceImpl) SaveCraftRoute(c *gin.Context, topo *craftRouteModels.ProcessTopo) (*craftRouteModels.SysCraftRouteConfig, error) {
	if topo == nil {
		return nil, errors.New("参数错误")
	}
	if topo.Route == nil {
		return nil, errors.New("工艺路线参数错误")
	}
	info, err := craft.dao.GetById(c, topo.Route.RouteID)
	if err != nil {
		zap.L().Error("读取工艺路线失败", zap.Error(err))
		return nil, err
	}
	if info == nil {
		return nil, errors.New("工艺路线不存在")
	}

	if len(topo.Edges) == 0 {
		return nil, errors.New("请设置工艺路线")
	}

	if len(topo.Nodes) == 0 {
		return nil, errors.New("请设置工艺节点")
	}

	// 工序id集合
	var processIdDataMap map[string]*craftRouteModels.ProcessData = make(map[string]*craftRouteModels.ProcessData)
	// 工序编号
	var nodeIds map[string]string = make(map[string]string)
	var processNodesMap map[string]*craftRouteModels.ProcessData = make(map[string]*craftRouteModels.ProcessData)
	// 开始节点数量
	var startCount uint32 = 0
	// 开始节点对应的target
	var beginTargetData *craftRouteModels.ProcessData

	// 遍历所有节点
	for _, node := range topo.Nodes {
		if node == nil {
			return nil, errors.New("工艺节点不能为空")
		}

		_, ok := craft2.NodeMap[node.Type]
		if !ok {
			return nil, errors.New(fmt.Sprintf("工艺节点类型%s不存在", node.Type))
		}

		if node.Data == nil {
			return nil, errors.New("工序节点数据不存在")
		}

		if node.Data.Config == nil {
			return nil, errors.New(fmt.Sprintf("节点%s数据不存在", node.Id))
		}

		if node.Id == craft2.START_NAME {
			startCount++
		}

		if node.Type == craft2.NODE_PROCESS_TYPE {
			var data craftRouteModels.ProcessData
			err := mapstructure.Decode(node.Data.Config, &data)
			if err != nil {
				zap.L().Error("decode error ", zap.Error(err))
				continue
			}
			if data.ProcessId == "" {
				continue
			}
			processIdDataMap[data.ProcessId] = &data
			processNodesMap[node.Id] = &data
		}

		nodeIds[node.Id] = node.Id
	}

	//遍历所有边
	for _, edge := range topo.Edges {
		if edge == nil {
			return nil, errors.New("工艺路线边不能为空")
		}
		_, ok := craft2.EdgeMap[edge.Type]
		if !ok {
			return nil, errors.New(fmt.Sprintf("工艺路线%s不存在", edge.Type))
		}

		if edge.Source == edge.Target {
			return nil, errors.New(fmt.Sprintf("开始节点%s和结束节点%s不能相同", edge.Source, edge.Target))
		}

		_, ok = nodeIds[edge.Source]
		if !ok {
			if edge.Source != craft2.START_NAME {
				return nil, errors.New(fmt.Sprintf("工艺路线 source: %s 不存在", edge.Source))
			}
		}

		_, ok = nodeIds[edge.Target]
		if !ok {
			return nil, errors.New(fmt.Sprintf("工艺路线 target: %s 不存在", edge.Target))
		}

		if edge.Source == craft2.START_NAME {
			value, ok := processNodesMap[edge.Target]
			if !ok {
				return nil, errors.New(fmt.Sprintf("工艺节点%s不存在", edge.Target))
			}
			beginTargetData = value
		}
	}

	if startCount == 0 {
		return nil, errors.New("开始节点不存在")
	}

	if startCount > 1 {
		return nil, errors.New("只能有一个开始节点")
	}

	if len(processIdDataMap) == 0 {
		return nil, errors.New("工序节点数据不存在")
	}

	var processList []int64 = make([]int64, 0)
	var processCheckMap = make(map[int64]bool)
	for _, process := range processIdDataMap {
		parseInt, err := strconv.ParseInt(process.ProcessId, 10, 64)
		if err != nil {
			zap.L().Error("parse int error", zap.Error(err))
			return nil, err
		}
		processList = append(processList, parseInt)
	}

	// 读取工序列表
	processes, err := craft.processDao.GetByIds(c, processList)
	if err != nil {
		return nil, err
	}

	if len(processList) != len(processes) {

		for _, process := range processes {
			processCheckMap[process.ProcessID] = true
		}

		for _, processId := range processList {
			_, ok := processCheckMap[processId]
			if !ok {
				return nil, errors.New(fmt.Sprintf("工序id %d 不存在", processId))
			}
		}
		return nil, errors.New("工序节点数据不存在")
	}

	// 读取工序内容
	processContexts, err := craft.contextDao.GetByProcessIds(c, processList)
	if err != nil {
		return nil, err
	}

	// 组装工序配置
	content, err := craft.loadV1ProcessTopo(c, processes, processContexts, info, beginTargetData)
	if err != nil {
		return nil, err
	}

	if content == nil {
		return nil, errors.New("生成工艺路线配置失败")
	}

	data, err := craft.routeConfigDao.Save(c, uint64(topo.Route.RouteID), topo, content)
	if err != nil {
		return nil, err
	}
	return data, err
}

func (craft *CraftRouteServiceImpl) loadV1ProcessTopo(c *gin.Context,
	processes []*craftRouteModels.SysProProcess,
	processContexts []*craftRouteModels.SysProProcessContent,
	routerConfig *craftRouteModels.SysCraftRoute,
	beginTargetData *craftRouteModels.ProcessData) ([]byte, error) {
	router := v1.NewRouter()
	router.Id = uint64(routerConfig.RouteID)
	router.Name = routerConfig.RouteName
	router.Md5 = ""

	beginNextProcessId, err := strconv.ParseUint(beginTargetData.ProcessId, 10, 64)
	if err != nil {
		return nil, err
	}

	router.Begin = &v1.Begin{
		NextProcessId: beginNextProcessId,
	}

	contextList := make(map[uint64][]v1.ProcessContext, 0)

	for _, process := range processes {
		contextList[uint64(process.ProcessID)] = make([]v1.ProcessContext, 0)
	}

	// 组装工序内容
	for _, processContext := range processContexts {
		_, ok := contextList[processContext.ProcessID]
		if !ok {
			contextList[processContext.ProcessID] = make([]v1.ProcessContext, 0)
		}

		var controlRule *v1.ControlRules
		if processContext.Extension != "" {
			var rule craftRouteModels.ControlRule
			err := json.Unmarshal([]byte(processContext.Extension), &rule)
			if err != nil {
				zap.L().Error("解析触发规则失败", zap.Error(err))
				return nil, err
			}
			controlRule = v1.NewControlRules()

			if processContext.ControlType == string(control.Pid) {
				if rule.PidRules == nil {
					continue
				}
				controlRule.PidRules.DeviceId = rule.PidRules.DeviceId
				controlRule.PidRules.DataId = rule.PidRules.DataId
				controlRule.PidRules.ActualSignal = rule.PidRules.ActualSignal
				controlRule.PidRules.Proportional = rule.PidRules.Proportional
				controlRule.PidRules.Integral = rule.PidRules.Integral
				controlRule.PidRules.Derivative = rule.PidRules.Derivative

				for _, action := range rule.PidRules.Actions {
					controlRule.PidRules.Actions = append(controlRule.PidRules.Actions, &v1.DeviceAction{
						DeviceId:    action.DeviceId,
						DataId:      action.DataId,
						Value:       action.Value,
						DataFormat:  action.DataFormat,
						ControlMode: action.ControlMode,
						Condition:   action.Condition,
						Interval:    action.Interval,
					})
				}
			}

			if processContext.ControlType == string(control.Threshold) {
				for _, v := range rule.TriggerRules.Cases {
					for _, caseValue := range v.Conditions {
						var info v1.DeviceRuleInfo
						info.DataId = caseValue.DataId
						info.DeviceId = caseValue.DeviceId
						controlRule.TriggerRules.Rule.DataId = append(controlRule.TriggerRules.Rule.DataId, info)
					}
				}
				combinedRule := ""
				first := true
				for caseKey, caseRule := range rule.TriggerRules.Cases {
					if !first {
						combinedRule += " " + caseRule.Connector + " "
					} else {
						first = false
						combinedRule += "("
					}

					firstCondition := true
					conditionRule := ""
					for k, condition := range caseRule.Conditions {
						if !firstCondition {
							conditionRule += " " + condition.Connector + " "
						} else {
							firstCondition = false
							conditionRule += "("
						}
						conditionRule += condition.Rule
						if k == len(caseRule.Conditions)-1 {
							conditionRule += ")"
						}
					}

					combinedRule += conditionRule
					if len(rule.TriggerRules.Cases)-1 == caseKey {
						combinedRule += ")"
					}
					continue
				}
				controlRule.TriggerRules.Rule.Rule = combinedRule

				for _, action := range rule.TriggerRules.Actions {
					controlRule.TriggerRules.Actions = append(controlRule.TriggerRules.Actions, &v1.DeviceAction{
						DeviceId:    action.DeviceId,
						DataId:      action.DataId,
						Value:       action.Value,
						DataFormat:  action.DataFormat,
						ControlMode: action.ControlMode,
						Condition:   action.Condition,
						Interval:    action.Interval,
					})
				}
			}

		}

		contextList[processContext.ProcessID] = append(contextList[processContext.ProcessID], v1.ProcessContext{
			ContentID:      processContext.ContentID,
			ProcessID:      processContext.ProcessID,
			ControlName:    processContext.ControlName,
			ControllerType: processContext.ControlType,
			ControlRules:   controlRule,
		})
	}

	// 组装工序
	for _, process := range processes {
		processContext, ok := contextList[uint64(process.ProcessID)]
		if !ok {
			return nil, errors.New(fmt.Sprintf("工序不存在 %d", process.ProcessID))
		}
		router.Processes = append(router.Processes, &v1.Process{
			ProcessId: uint64(process.ProcessID),
			Name:      process.ProcessName,
			Context:   processContext,
		})
	}

	content, err := json.Marshal(router)
	if err != nil {
		return nil, err
	}

	md5 := uuid.MakeMd5(content)
	router.Md5 = md5
	content, err = json.Marshal(router)
	if err != nil {
		return nil, err
	}
	return content, nil
}
