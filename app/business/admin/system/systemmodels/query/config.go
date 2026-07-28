package query

import "nova-factory-server/app/baize"

type SysConfigDQL struct {
	ConfigName string `form:"configName" db:"config_name"` //参数名称
	ConfigKey  string `form:"configKey" db:"config_key"`   //参数键名
	ConfigType string `form:"configType" db:"config_type"` //系统内置（Y是 N否）
	baize.BaseEntityDQL
}
