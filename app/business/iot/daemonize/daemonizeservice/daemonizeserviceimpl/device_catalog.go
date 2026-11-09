package daemonizeserviceimpl

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"nova-factory-server/app/business/iot/asset/device/devicemodels"
	catalogresponse "nova-factory-server/app/business/iot/daemonize/daemonizemodels/response"
)

// Catalog converts a consistent asset snapshot into passive device configuration.
func (s *iGatewayConfigServiceImpl) Catalog(ctx context.Context, id int64) (*catalogresponse.Catalog, error) {
	snapshot, err := s.deviceDao.GetGatewayCatalog(ctx, id)
	if err != nil {
		return nil, err
	}
	catalog := &catalogresponse.Catalog{
		GatewayID: strconv.FormatInt(id, 10),
		Devices:   []catalogresponse.Device{},
	}
	templates := map[int64]*devicemodels.SysDeviceTemplate{}
	for _, t := range snapshot.Templates {
		templates[t.TemplateID] = t
	}
	points := map[int64][]catalogresponse.Point{}
	for _, p := range snapshot.Points {
		if p.Register == nil || p.FunctionCode < 1 || p.FunctionCode > 4 {
			return nil, fmt.Errorf("template %d has unsupported measurement configuration", p.TemplateID)
		}
		points[p.TemplateID] = append(points[p.TemplateID], catalogresponse.Point{
			ID:           strconv.FormatInt(p.DeviceConfigID, 10),
			Name:         p.Name,
			Type:         p.Type,
			Register:     *p.Register,
			FunctionCode: int8(p.FunctionCode),
			Sort:         p.Sort,
			Unit:         p.Unit,
			Precision:    p.Precision,
			Expression:   p.Expression,
			DataFormat:   p.DataFormat,
		})
	}
	for _, d := range snapshot.Devices {
		if d.ProtocolType != "modbus-tcp" && d.ProtocolType != "modbus_tcp" {
			return nil, fmt.Errorf("device %d protocol is not supported by catalog v1", d.DeviceId)
		}
		t := templates[int64(d.DeviceProtocolId)]
		if t == nil || len(points[t.TemplateID]) == 0 {
			return nil, fmt.Errorf("device %d template is incomplete", d.DeviceId)
		}
		var ext devicemodels.ExtensionInfo
		if json.Unmarshal([]byte(d.Extension), &ext) != nil || len(ext.LocalInfo) != 1 {
			return nil, fmt.Errorf("device %d requires exactly one local connection", d.DeviceId)
		}
		connection := ext.LocalInfo[0]
		if connection.Slave <= 0 || connection.Quantity == 0 || connection.Address == "" {
			return nil, fmt.Errorf("device %d local connection is incomplete", d.DeviceId)
		}
		config, _ := json.Marshal(map[string]any{
			"ipAddress":    connection.Address,
			"slaveAddress": connection.Slave,
			"maxRegisters": connection.Quantity,
		})
		name, number := "", ""
		if d.Name != nil {
			name = *d.Name
		}
		if d.Number != nil {
			number = *d.Number
		}
		catalog.Devices = append(catalog.Devices, catalogresponse.Device{
			ID:           strconv.FormatUint(d.DeviceId, 10),
			Name:         name,
			Number:       number,
			Protocol:     "modbus_tcp",
			Config:       config,
			TemplateID:   strconv.FormatInt(t.TemplateID, 10),
			TemplateName: t.Name,
			Points:       points[t.TemplateID],
		})
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
