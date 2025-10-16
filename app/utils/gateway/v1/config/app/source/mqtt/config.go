package mqtt

import "nova-factory-server/app/utils/gateway/v1/api"

type DataTypeConfig struct {
	TemplateId uint64 `yaml:"templateId"`
	DataId     uint64 `yaml:"dataId"`
	Name       string
	Expression string
	DataFormat api.DataValueType `yaml:"dataFormat,omitempty"`
	Unit       string            `yaml:"unit,omitempty"`
}

type Device struct {
	DeviceId string           `yaml:"deviceId"`
	Name     string           `yaml:"name"`
	Topic    string           `yaml:"topic"`
	Template []DataTypeConfig `yaml:"template,omitempty"`
}

type Config struct {
	Address  string   `yaml:"address,omitempty"`
	ClientId string   `yaml:"client_id,omitempty"`
	Username string   `yaml:"username,omitempty"`
	Password string   `yaml:"password,omitempty"`
	Devices  []Device `yaml:"devices,omitempty"`
}
