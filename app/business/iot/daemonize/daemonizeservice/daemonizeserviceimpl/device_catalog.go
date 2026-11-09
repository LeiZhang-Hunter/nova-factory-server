package daemonizeserviceimpl

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"go.uber.org/zap"

	"nova-factory-server/app/business/iot/asset/device/devicemodels"
	catalogresponse "nova-factory-server/app/business/iot/daemonize/daemonizemodels/response"
)

// Catalog 根据网关设备快照生成完整配置目录。
func (s *iGatewayConfigServiceImpl) Catalog(ctx context.Context, id int64) (*catalogresponse.Catalog, error) {
	snapshot, err := s.deviceDao.GetGatewayCatalog(ctx, id)
	if err != nil {
		return nil, err
	}

	zap.L().Info("已读取网关设备配置快照", zap.Int64("gatewayId", id), zap.Int("devices", len(snapshot.Devices)), zap.Int("templates", len(snapshot.Templates)), zap.Int("points", len(snapshot.Points)))
	return buildDeviceCatalog(id, snapshot)
}

// buildDeviceCatalog 按设备协议转换快照，并计算目录版本。
func buildDeviceCatalog(id int64, snapshot *catalogresponse.CatalogSnapshot) (*catalogresponse.Catalog, error) {
	catalog := &catalogresponse.Catalog{GatewayID: strconv.FormatInt(id, 10), Devices: []catalogresponse.Device{}}
	templates := make(map[int64]*devicemodels.SysDeviceTemplate, len(snapshot.Templates))
	points := make(map[int64][]*devicemodels.SysModbusDeviceConfigData)
	for _, template := range snapshot.Templates {
		templates[template.TemplateID] = template
	}
	for _, point := range snapshot.Points {
		points[point.TemplateID] = append(points[point.TemplateID], point)
	}

	for _, device := range snapshot.Devices {
		converter, err := catalogConverterFor(device.ProtocolType)
		if err != nil {
			return nil, fmt.Errorf("device %d: %w", device.DeviceId, err)
		}

		template := templates[int64(device.DeviceProtocolId)]
		data := points[int64(device.DeviceProtocolId)]
		if template == nil || len(data) == 0 {
			return nil, fmt.Errorf("device %d template is incomplete", device.DeviceId)
		}

		item, err := converter.convert(device, template, data)
		if err != nil {
			return nil, fmt.Errorf("device %d template %d: %w", device.DeviceId, template.TemplateID, err)
		}
		catalog.Devices = append(catalog.Devices, item)
	}

	body, err := json.Marshal(catalog)
	if err != nil {
		return nil, err
	}
	if len(body) > 16<<20 {
		return nil, errors.New("catalog exceeds size limit")
	}

	catalog.Version = fmt.Sprintf("%x", sha256.Sum256(body))
	return catalog, nil
}
