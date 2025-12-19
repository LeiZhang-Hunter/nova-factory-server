package daemonizeServiceImpl

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"nova-factory-server/app/business/ai/aiDataSetDao"
	"nova-factory-server/app/business/alert/alertDao"
	"nova-factory-server/app/business/asset/device/deviceDao"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/daemonize/daemonizeService"
	"nova-factory-server/app/business/system/systemDao"
	"nova-factory-server/app/business/system/systemModels"
	device2 "nova-factory-server/app/constant/device"
	"nova-factory-server/app/constant/gateway"
	iotdb2 "nova-factory-server/app/constant/iotdb"
	logalertIntercept "nova-factory-server/app/utils/gateway/v1/config/app/intercept/logalert"
	"nova-factory-server/app/utils/gateway/v1/config/app/sink/alertwebhook"
	"nova-factory-server/app/utils/gateway/v1/config/app/sink/metric_exporter"
	"nova-factory-server/app/utils/gateway/v1/config/app/sink/time_data_exporter"
	"nova-factory-server/app/utils/gateway/v1/config/app/source/bhps7"
	"nova-factory-server/app/utils/gateway/v1/config/app/source/mqtt"
	"nova-factory-server/app/utils/gateway/v1/config/app/source/prediction"
	"nova-factory-server/app/utils/gateway/v1/config/app/source/running_statistics"
	"nova-factory-server/app/utils/gateway/v1/config/app/source/scheduler"
	"nova-factory-server/app/utils/gateway/v1/config/cfg"
	"nova-factory-server/app/utils/gateway/v1/config/interceptor"
	"nova-factory-server/app/utils/gateway/v1/config/logalert"
	"nova-factory-server/app/utils/gateway/v1/config/pipeline"
	"nova-factory-server/app/utils/gateway/v1/config/render"
	"nova-factory-server/app/utils/gateway/v1/config/sink"
	source2 "nova-factory-server/app/utils/gateway/v1/config/source"
	"nova-factory-server/app/utils/snowflake"
	"strconv"
	"time"
)

type iGatewayConfigServiceImpl struct {
	deviceDao         deviceDao.IDeviceDao
	templateDao       deviceDao.IDeviceTemplateDao
	templateDataDao   deviceDao.ISysModbusDeviceConfigDataDao
	ruleDao           alertDao.AlertRuleDao
	sinkDao           alertDao.AlertSinkTemplateDao
	electricConfigDao systemDao.IDeviceElectricDao
	dictDataDao       systemDao.IDictDataDao
	predictionDao     aiDataSetDao.IAiPredictionControlDao
}

func NewIGatewayConfigServiceImpl(
	deviceDao deviceDao.IDeviceDao,
	templateDao deviceDao.IDeviceTemplateDao,
	templateDataDao deviceDao.ISysModbusDeviceConfigDataDao,
	ruleDao alertDao.AlertRuleDao,
	sinkDao alertDao.AlertSinkTemplateDao, electricConfigDao systemDao.IDeviceElectricDao,
	dictDataDao systemDao.IDictDataDao, predictionDao aiDataSetDao.IAiPredictionControlDao) daemonizeService.IGatewayConfigService {
	return &iGatewayConfigServiceImpl{
		deviceDao:         deviceDao,
		templateDao:       templateDao,
		templateDataDao:   templateDataDao,
		ruleDao:           ruleDao,
		sinkDao:           sinkDao,
		electricConfigDao: electricConfigDao,
		dictDataDao:       dictDataDao,
		predictionDao:     predictionDao,
	}
}

// Generate 渲染Agent配置
func (i *iGatewayConfigServiceImpl) Generate(c *gin.Context, gatewayId int64) (*pipeline.PipelineConfig, error) {
	// 生成配置
	devices, err := i.deviceDao.GetLocalByGateWayId(c, gatewayId)
	if err != nil {
		return nil, err
	}

	// 查找sink地址
	addresses := i.dictDataDao.SelectDictDataByType(c, gateway.GRPC_SINK_ADDRESS)
	if len(addresses) == 0 {
		return nil, errors.New("网关数据写入地址不能是空")
	}
	sink_address := addresses[0].DictValue
	serverHostInfo := i.dictDataDao.SelectDictDataByType(c, gateway.SERVER_HOST)
	if len(serverHostInfo) == 0 {
		return nil, errors.New("网关数据写入地址不能是空")
	}
	server_address := serverHostInfo[0].DictValue

	var deviceAddressMap map[string][]*deviceModels.DeviceVO = make(map[string][]*deviceModels.DeviceVO)
	var templateIdMap map[uint64]uint64 = make(map[uint64]uint64)
	var templateIds []uint64 = make([]uint64, 0)
	var templatesMap map[uint64][]*deviceModels.SysModbusDeviceConfigData = make(map[uint64][]*deviceModels.SysModbusDeviceConfigData)
	piplines := pipeline.NewPipelineConfig()
	pipelinesConfig := pipeline.NewConfig()

	// 组装设备
	for _, device := range devices {

		if device.ExtensionInfo == nil {
			continue
		}

		if len(device.ExtensionInfo.LocalInfo) == 0 && len(device.ExtensionInfo.LocalMqttInfo) == 0 {
			continue
		}

		if device.ProtocolType == device2.MQTT {
			address := device.ExtensionInfo.LocalMqttInfo[0].Address
			if address == "" {
				continue
			}
			if deviceAddressMap[address] == nil {
				deviceAddressMap[address] = make([]*deviceModels.DeviceVO, 0)
			}
			deviceAddressMap[address] = append(deviceAddressMap[address], device)
			if device.DeviceProtocolId != 0 {
				templateIdMap[device.DeviceProtocolId] = device.DeviceProtocolId
			}
		} else if device.ProtocolType == device2.MODBUS_TCP {
			address := device.ExtensionInfo.LocalInfo[0].Address
			if address == "" {
				continue
			}
			if deviceAddressMap[address] == nil {
				deviceAddressMap[address] = make([]*deviceModels.DeviceVO, 0)
			}
			deviceAddressMap[address] = append(deviceAddressMap[address], device)
			if device.DeviceProtocolId != 0 {
				templateIdMap[device.DeviceProtocolId] = device.DeviceProtocolId
			}
		}
	}

	if len(templateIdMap) == 0 {
		return nil, fmt.Errorf("设备未绑定模板")
	}

	for _, templateId := range templateIdMap {
		templateIds = append(templateIds, templateId)
	}

	// 读取设备模板
	templates, err := i.templateDao.GetByIds(c, templateIds)
	if err != nil {
		return nil, err
	}

	for _, templateValue := range templates {
		templatesMap[uint64(templateValue.TemplateID)] = make([]*deviceModels.SysModbusDeviceConfigData, 0)
	}

	// 读取设备模板数据
	templatesData, err := i.templateDataDao.GetByTemplateIds(c, templateIds)
	if err != nil {
		return nil, err
	}

	for _, data := range templatesData {
		templatesMap[uint64(data.TemplateID)] = append(templatesMap[uint64(data.TemplateID)], data)
	}

	//构建source
	for addr, devicesData := range deviceAddressMap {
		enabled := true

		if len(devicesData) == 0 {
			continue
		}

		if devicesData[0].ExtensionInfo == nil {
			continue
		}

		if len(devicesData[0].ExtensionInfo.LocalMqttInfo) == 0 && len(devicesData[0].ExtensionInfo.LocalInfo) == 0 {
			continue
		}

		deviceType := devicesData[0].ProtocolType

		if deviceType == device2.MODBUS_TCP {
			source := source2.Config{
				Enabled: &enabled,
				Name:    fmt.Sprintf("gateway-%s", addr),
				Type:    "bhp_s7_gateway",
			}

			var config bhps7.Config
			config.Address = addr
			if len(devicesData) > 0 && len(devicesData[0].ExtensionInfo.LocalInfo) > 0 {
				config.Quantity = devicesData[0].ExtensionInfo.LocalInfo[0].Quantity
			}
			config.Devices = make([]bhps7.Device, 0)

			for _, d := range devicesData {
				templateDatas, ok := templatesMap[d.DeviceProtocolId]
				if !ok {
					continue
				}
				bhps7Device := render.OfBhps7Device(d, templateDatas)
				if bhps7Device == nil {
					continue
				}
				config.Devices = append(config.Devices, *bhps7Device)
			}
			pack, err := cfg.Pack(config)
			if err != nil {
				zap.L().Error("cfg.Pack() failed", zap.Error(err))
				continue
			}
			source.Properties = pack
			pipelinesConfig.Sources = append(pipelinesConfig.Sources, &source)
		} else if deviceType == device2.MQTT {
			username := devicesData[0].ExtensionInfo.LocalMqttInfo[0].Username
			password := devicesData[0].ExtensionInfo.LocalMqttInfo[0].Password
			clientId := devicesData[0].ExtensionInfo.LocalMqttInfo[0].ClientId
			source := source2.Config{
				Enabled: &enabled,
				Name:    fmt.Sprintf("gateway-%s", addr),
				Type:    "mqtt",
			}

			var config mqtt.Config
			config.Address = addr
			config.Username = username
			config.ClientId = clientId
			config.Password = password
			config.Devices = make([]mqtt.Device, 0)

			for _, d := range devicesData {
				templateDatas, ok := templatesMap[d.DeviceProtocolId]
				if !ok {
					continue
				}
				mqttDevice := render.OfMqttDevice(d, templateDatas)
				if mqttDevice == nil {
					continue
				}
				mqttDevice.Topic = d.ExtensionInfo.LocalMqttInfo[0].Topic
				config.Devices = append(config.Devices, *mqttDevice)
			}
			pack, err := cfg.Pack(config)
			if err != nil {
				zap.L().Error("cfg.Pack() failed", zap.Error(err))
				continue
			}
			source.Properties = pack
			pipelinesConfig.Sources = append(pipelinesConfig.Sources, &source)
		}

		var exportConfig metric_exporter.Config
		exportConfig.Address = sink_address
		packContent, err := cfg.Pack(exportConfig)
		if err != nil {
			zap.L().Error("cfg.Pack() failed", zap.Error(err))
			continue
		}
		pipelinesConfig.Sink = &sink.Config{
			Enabled:     &enabled,
			Name:        fmt.Sprintf("sink-%d", snowflake.GenID()),
			Type:        "metric_exporter",
			Properties:  packContent,
			Parallelism: 1,
		}
	}
	pipelinesConfig.Name = fmt.Sprintf("gateway-pipeline-%d", gatewayId)

	schedulerConfig := pipeline.NewConfig()
	schedulerConfig.Name = "scheduler"
	var scheduleEnabled bool = true
	source := source2.Config{
		Enabled: &scheduleEnabled,
		Name:    "scheduler",
		Type:    "scheduler",
	}
	var taskConfig scheduler.Config
	taskConfig.Task.Enabled = true
	taskConfig.Task.Version = "v1"
	taskConfig.Task.Host = server_address
	taskConfig.Task.Limit = 50
	taskConfig.Task.PollTime = 5 * time.Minute
	packContent, err := cfg.Pack(taskConfig)
	source.Properties = packContent
	schedulerConfig.Sources = append(schedulerConfig.Sources, &source)

	// 读取模板下的所有模板数据
	piplines.Pipelines = append(piplines.Pipelines, *pipelinesConfig)
	piplines.Pipelines = append(piplines.Pipelines, *schedulerConfig)

	// ====================================== 告警配置 =====================================================
	// 通过网关id读取告警配置信息
	alertConfig := pipeline.NewConfig()
	alertRule, err := i.ruleDao.GetOnlineByGatewayId(c, uint64(gatewayId))
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			zap.L().Error("GetOnlineByGatewayId() failed", zap.Error(err))
			return nil, errors.New("读取告警配置失败")
		} else {
			alertRule = nil
		}

	}

	if alertRule != nil {
		sinkInfo, err := i.sinkDao.GetById(c, uint64(alertRule.TemplateID))
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("告警模板不存在")
			}
		}

		if sinkInfo != nil {
			notifierSourceConfig := source2.Config{
				Enabled: &scheduleEnabled,
				Name:    "notifier",
				Type:    "notifier",
			}
			alertConfig.Sources = append(alertConfig.Sources, &notifierSourceConfig)
			alertConfig.Name = fmt.Sprintf("gateway-notifier-%d", gatewayId)
			interceptorConfig := interceptor.Config{
				Enabled: &scheduleEnabled,
				Name:    "alert",
				Type:    "logAlert",
			}
			var logAlertConfig logalertIntercept.Config
			var advanced *logalertIntercept.Advanced
			var ignore []logalertIntercept.DeviceMetric
			var matcher *logalertIntercept.Matcher
			var additions map[string]interface{}
			if alertRule.Advanced != "" {
				err = json.Unmarshal([]byte(alertRule.Advanced), &advanced)
				if err != nil {
					zap.L().Error("json.Unmarshal failed", zap.Error(err))
				}
			}
			if alertRule.Ignore != "" {
				err = json.Unmarshal([]byte(alertRule.Ignore), &ignore)
				if err != nil {
					zap.L().Error("json.Unmarshal failed", zap.Error(err))
				}
			}
			if alertRule.Matcher != "" {
				err = json.Unmarshal([]byte(alertRule.Matcher), &matcher)
				if err != nil {
					zap.L().Error("json.Unmarshal failed", zap.Error(err))
				}
			}
			if alertRule.Additions != "" {
				err = json.Unmarshal([]byte(alertRule.Additions), &additions)
				if err != nil {
					zap.L().Error("json.Unmarshal failed", zap.Error(err))
				}
			}
			logAlertConfig.AlertId = strconv.FormatInt(alertRule.ID, 10)
			logAlertConfig.Ignore = ignore
			logAlertConfig.Matcher = *matcher
			logAlertConfig.Additions = additions
			logAlertConfig.SendOnlyMatched = true
			logAlertConfig.Advanced = *advanced
			logAlertConfig.Advanced.Enable = true
			pack, err := cfg.Pack(logAlertConfig)
			if err != nil {
				zap.L().Error("cfg.Pack() failed", zap.Error(err))
			}
			interceptorConfig.Properties = pack
			alertConfig.Interceptors = append(alertConfig.Interceptors, &interceptorConfig)

			//  告警输出
			var headers map[string]string = make(map[string]string)
			if len(sinkInfo.Headers) != 0 {
				err := json.Unmarshal([]byte(sinkInfo.Headers), &headers)
				if err != nil {
					zap.L().Error("json.Unmarshal failed", zap.Error(err))
				}
			}
			alertWebhookConfig := alertwebhook.Config{
				Addr: sinkInfo.Addr,
				AlertConfig: logalert.AlertConfig{
					Template:              sinkInfo.Template,
					Timeout:               time.Duration(sinkInfo.Timeout) * time.Second,
					Headers:               headers,
					Method:                sinkInfo.Method,
					SendLogAlertAtOnce:    true,
					SendNoDataAlertAtOnce: true,
				},
			}
			alertWebhookConfigPacked, err := cfg.Pack(alertWebhookConfig)
			alertConfig.Sink = &sink.Config{
				Enabled:     &scheduleEnabled,
				Name:        sinkInfo.Name,
				Type:        "alertWebhook",
				Parallelism: 1,
				Properties:  alertWebhookConfigPacked,
			}
			piplines.Pipelines = append(piplines.Pipelines, *alertConfig)
		}
	}

	// 渲染用电量聚合统计算法
	//  查询所有电流配置
	// ====================================== 电流配置 =====================================================
	electricConfig := pipeline.NewConfig()
	electricConfig.Name = "electric_statistics"
	all, err := i.electricConfigDao.All(c)
	if err != nil {
		return nil, err
	}
	var runningStatisticsConfig running_statistics.Config
	runningStatisticsConfig.DeviceStatistics = make([]*running_statistics.DeviceStatistics, 0)
	if len(all) != 0 {
		for _, config := range all {
			var statistics running_statistics.DeviceStatistics
			if config.DeviceID == 0 {
				continue
			}
			var ex systemModels.Expression
			err := json.Unmarshal([]byte(config.Expression), &ex)
			if err != nil {
				zap.L().Error("json error", zap.Error(err))
				continue
			}

			table := iotdb2.MakeRunDeviceTemplateName(config.DeviceID)
			statistics.Table = table
			statistics.DeviceId = fmt.Sprintf("%d", config.DeviceID)
			statistics.Expression = running_statistics.Expression{
				Rules: make([]running_statistics.StatisticsRule, 0),
			}
			for _, exRule := range ex.Rules {
				var rule running_statistics.StatisticsRule
				rule.RunStatus = int(exRule.RunStatus)
				rule.MatchType = exRule.MatchType
				rule.Groups = make([]running_statistics.GroupRule, 0)
				for _, group := range exRule.Groups {
					rule.Groups = append(rule.Groups, running_statistics.GroupRule{
						Key:          group.Key,
						Name:         group.Name,
						Operator:     group.Operator,
						OperatorName: group.OperatorName,
						Value:        group.Value,
					})
				}
				statistics.Expression.Rules = append(statistics.Expression.Rules, rule)
			}

			runningStatisticsConfig.DeviceStatistics = append(runningStatisticsConfig.DeviceStatistics, &statistics)
		}
	}
	electricConfigSource := source2.Config{
		Enabled: &scheduleEnabled,
		Name:    "running_statistics",
		Type:    "running_statistics",
	}
	pack, err := cfg.Pack(runningStatisticsConfig)
	if err != nil {
		zap.L().Error("cfg.Pack() failed", zap.Error(err))
		return nil, err
	}
	electricConfigSource.Properties = pack
	electricConfig.Sources = append(electricConfig.Sources, &electricConfigSource)
	// 安装sink
	exportConfig := time_data_exporter.Config{
		Address: sink_address,
	}
	packContent, err = cfg.Pack(exportConfig)
	if err != nil {
		zap.L().Error("cfg.Pack() failed", zap.Error(err))
		return nil, err
	}
	electricConfig.Sink = &sink.Config{
		Enabled:     &scheduleEnabled,
		Name:        fmt.Sprintf("sink-%d", snowflake.GenID()),
		Type:        "time_data_exporter",
		Properties:  packContent,
		Parallelism: 1,
	}
	piplines.Pipelines = append(piplines.Pipelines, *electricConfig)

	predictionInfo, err := i.predictionDao.Find(c)
	if err != nil {
		zap.L().Error("i.predictionDao.Find() failed", zap.Error(err))
		return nil, err
	}
	if predictionInfo != nil {
		//  查询预测配置
		predictionPipeline := pipeline.NewConfig()
		predictionPipeline.Name = "prediction"
		predictionPipelineSource := source2.Config{
			Enabled: &scheduleEnabled,
			Name:    "prediction",
			Type:    "prediction",
		}
		var predictionConfig prediction.Config
		predictionConfig.Parallelism = uint16(predictionInfo.Parallelism)
		predictionConfig.Model = predictionInfo.Model
		predictionConfig.TimeWindow = time.Duration(predictionInfo.Interval)
		predictionConfig.PredictLength = uint64(predictionInfo.PredictLength)
		//predictionConfig.
		pack, err = cfg.Pack(&predictionConfig)
		if err != nil {
			zap.L().Error("cfg.Pack() failed", zap.Error(err))
			return nil, err
		}
		predictionPipelineSource.Properties = pack

		predictionPipeline.Sources = append(predictionPipeline.Sources, &predictionPipelineSource)
		piplines.Pipelines = append(piplines.Pipelines, *predictionPipeline)

	}
	return piplines, nil
}
