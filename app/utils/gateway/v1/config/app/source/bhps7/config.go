package bhps7

import (
	"nova-factory-server/app/utils/gateway/v1/api"
	"time"
)

type DataTypeConfig struct {
	TemplateId   uint64 `yaml:"templateId"`
	Name         string
	Protocol     string
	Annotation   string
	DataFormat   api.DataValueType `yaml:"dataFormat"`
	Unit         string            `yaml:"unit"`
	Position     uint16
	FunctionCode uint16 `yaml:"functionCode"`
	deviceId     string
	Sort         api.ByteOrder `yaml:"sort"`
}

type Device struct {
	DeviceId string           `yaml:"deviceId"`
	Name     string           `yaml:"name"`
	Template []DataTypeConfig `yaml:"template,omitempty"`
}

type Config struct {
	Enabled     *bool         `yaml:"enabled,omitempty"`
	Address     string        `yaml:"address,omitempty"`
	SlaveId     byte          `yaml:"slaveId,omitempty"`
	CollectTime time.Duration `yaml:"collectTime,omitempty" default:"5s"` // default 2 * read.timeout
	Devices     []Device      `yaml:"devices,omitempty"`
	Length      uint16        `yaml:"length,omitempty"`
}
