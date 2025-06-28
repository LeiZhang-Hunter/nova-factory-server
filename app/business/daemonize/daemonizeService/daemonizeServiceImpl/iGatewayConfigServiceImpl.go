package daemonizeServiceImpl

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nova-factory-server/app/business/asset/device/deviceDao"
	"nova-factory-server/app/business/asset/device/deviceModels"
	"nova-factory-server/app/business/daemonize/daemonizeService"
	"nova-factory-server/app/utils/gateway/v1/api"
	"nova-factory-server/app/utils/gateway/v1/config/app/sink/metric_exporter"
	"nova-factory-server/app/utils/gateway/v1/config/app/source/bhps7"
	"nova-factory-server/app/utils/gateway/v1/config/cfg"
	"nova-factory-server/app/utils/gateway/v1/config/pipeline"
	"nova-factory-server/app/utils/gateway/v1/config/sink"
	source2 "nova-factory-server/app/utils/gateway/v1/config/source"
	"nova-factory-server/app/utils/snowflake"
	"strconv"
)

type iGatewayConfigServiceImpl struct {
	deviceDao       deviceDao.IDeviceDao
	templateDao     deviceDao.IDeviceTemplateDao
	templateDataDao deviceDao.ISysModbusDeviceConfigDataDao
}

func NewIGatewayConfigServiceImpl(
	deviceDao deviceDao.IDeviceDao,
	templateDao deviceDao.IDeviceTemplateDao,
	templateDataDao deviceDao.ISysModbusDeviceConfigDataDao) daemonizeService.IGatewayConfigService {
	return &iGatewayConfigServiceImpl{
		deviceDao:       deviceDao,
		templateDao:     templateDao,
		templateDataDao: templateDataDao,
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
	for _, device := range devices {
		if len(device.ExtensionInfo.LocalInfo) == 0 {
			continue
		}
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
	return piplines, nil
}
