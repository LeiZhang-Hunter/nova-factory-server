package daemonizeServiceImpl

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/alert/alertDao"
	"nova-factory-server/app/business/alert/alertModels"
	"nova-factory-server/app/business/asset/device/deviceDao"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/daemonize/daemonizeService"
	device2 "nova-factory-server/app/constant/device"
	"nova-factory-server/app/utils/gateway/v1/api"
	logalertIntercept "nova-factory-server/app/utils/gateway/v1/config/app/intercept/logalert"
	"nova-factory-server/app/utils/gateway/v1/config/app/sink/alertwebhook"
	"nova-factory-server/app/utils/gateway/v1/config/app/sink/metric_exporter"
	"nova-factory-server/app/utils/gateway/v1/config/app/source/bhps7"
	"nova-factory-server/app/utils/gateway/v1/config/cfg"
	"nova-factory-server/app/utils/gateway/v1/config/interceptor"
	"nova-factory-server/app/utils/gateway/v1/config/logalert"
	"nova-factory-server/app/utils/gateway/v1/config/pipeline"
	"nova-factory-server/app/utils/gateway/v1/config/sink"
	source2 "nova-factory-server/app/utils/gateway/v1/config/source"
	"nova-factory-server/app/utils/snowflake"
	"strconv"
	"time"
)

type iGatewayConfigServiceImpl struct {
	deviceDao       deviceDao.IDeviceDao
	templateDao     deviceDao.IDeviceTemplateDao
	templateDataDao deviceDao.ISysModbusDeviceConfigDataDao
	ruleDao         alertDao.AlertRuleDao
	sinkDao         alertDao.AlertSinkTemplateDao
}

func NewIGatewayConfigServiceImpl(
	deviceDao deviceDao.IDeviceDao,
	templateDao deviceDao.IDeviceTemplateDao,
	templateDataDao deviceDao.ISysModbusDeviceConfigDataDao,
	ruleDao alertDao.AlertRuleDao,
	sinkDao alertDao.AlertSinkTemplateDao) daemonizeService.IGatewayConfigService {
	return &iGatewayConfigServiceImpl{
		deviceDao:       deviceDao,
		templateDao:     templateDao,
		templateDataDao: templateDataDao,
		ruleDao:         ruleDao,
		sinkDao:         sinkDao,
	}
}

func OfBhps7Device(vo *deviceModels.DeviceVO, data []*deviceModels.SysModbusDeviceConfigData) *bhps7.Device {
	var device bhps7.Device
	device.DeviceId = strconv.FormatUint(vo.DeviceId, 10)
	if vo.Name != nil {
		device.Name = *vo.Name
	} else {
		device.Name = strconv.FormatUint(vo.DeviceId, 10)
	}
	device.Template = make([]bhps7.DataTypeConfig, 0)
	for _, v := range data {
		device.Template = append(device.Template, bhps7.DataTypeConfig{
			Name:         v.Name,
			Protocol:     "",
			DataId:       uint64(v.DeviceConfigID),
			Annotation:   v.Name,
			DataFormat:   api.DataValueType(v.DataFormat),
			Unit:         v.Unit,
			TemplateId:   uint64(v.TemplateID),
			Position:     uint16(v.Register),
			FunctionCode: uint16(v.FunctionCode),
			Sort:         api.ByteOrder(v.Sort),
		})
	}

	return &device
}

// Generate 渲染Agent配置
func (i *iGatewayConfigServiceImpl) Generate(c *gin.Context, gatewayId int64) (*pipeline.PipelineConfig, error) {
	// 生成配置
	devices, err := i.deviceDao.GetLocalByGateWayId(c, gatewayId)
	if err != nil {
		return nil, err
	}

	var deviceAddressMap map[string][]*deviceModels.DeviceVO = make(map[string][]*deviceModels.DeviceVO)
	var templateIdMap map[uint64]uint64 = make(map[uint64]uint64)
	var templateIds []uint64 = make([]uint64, 0)
	var templatesMap map[uint64][]*deviceModels.SysModbusDeviceConfigData = make(map[uint64][]*deviceModels.SysModbusDeviceConfigData)
	piplines := pipeline.NewPipelineConfig()
	pipelinesConfig := pipeline.NewConfig()

	// 组装设备
	for _, device := range devices {
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
		source := source2.Config{
			Enabled: &enabled,
			Name:    fmt.Sprintf("gateway-%s", addr),
			Type:    "bhp_s7_gateway",
		}

		var config bhps7.Config
		config.Address = addr
		config.Devices = make([]bhps7.Device, 0)

		for _, d := range devicesData {
			templateDatas, ok := templatesMap[d.DeviceProtocolId]
			if !ok {
				continue
			}
			bhps7Device := OfBhps7Device(d, templateDatas)
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

		var exportConfig metric_exporter.Config
		exportConfig.Address = "localhost:6000"
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
	schedulerConfig.Sources = append(schedulerConfig.Sources, &source)

	// 读取模板下的所有模板数据
	piplines.Pipelines = append(piplines.Pipelines, *pipelinesConfig)
	piplines.Pipelines = append(piplines.Pipelines, *schedulerConfig)

	// 通过网关id读取告警配置信息
	alertConfig := pipeline.NewConfig()
	alertRule, err := i.ruleDao.GetOnlineByGatewayId(c, uint64(gatewayId))
	if err != nil {
		zap.L().Error("GetOnlineByGatewayId() failed", zap.Error(err))
		return nil, errors.New("告警规则不存在")
	}
	if alertRule == nil {
		alertRule = &alertModels.SysAlert{}
	}
	sinkInfo, err := i.sinkDao.GetById(c, uint64(alertRule.TemplateID))
	if err != nil {
		return nil, errors.New("告警模板不存在")
	}
	notifierSourceConfig := source2.Config{
		Enabled: &scheduleEnabled,
		Name:    "notifier",
		Type:    "notifier",
	}
	alertConfig.Sources = append(schedulerConfig.Sources, &notifierSourceConfig)

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
	return piplines, nil
}
