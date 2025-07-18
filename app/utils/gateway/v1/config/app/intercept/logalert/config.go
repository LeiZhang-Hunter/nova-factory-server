package logalert

import (
	"nova-factory-server/app/utils/gateway/v1/config/interceptor"
)

type Config struct {
	interceptor.ExtensionConfig `yaml:",inline"`
	AlertId                     string                 `yaml:"alertId" json:"alert_id"`
	Matcher                     Matcher                `yaml:"matcher,omitempty" json:"matcher"`
	Additions                   map[string]interface{} `yaml:"additions,omitempty" json:"additions"`
	Ignore                      []DeviceMetric         `yaml:"ignore,omitempty" json:"ignore"`
	Advanced                    Advanced               `yaml:"advanced,omitempty" json:"advanced"`
	SendOnlyMatched             bool                   `yaml:"sendOnlyMatched,omitempty" json:"sendOnlyMatched"`
}

type DeviceMetric struct {
	DeviceId string `yaml:"deviceId" json:"deviceId"`
	DataId   string `yaml:"dataId" json:"dataId"`
}

type Matcher struct {
	Contains     []DeviceMetric `yaml:"contains,omitempty" json:"contains"`
	TargetHeader string         `yaml:"target,omitempty" json:"targetHeader"`
}

type Advanced struct {
	Enable    bool     `yaml:"enabled" json:"enabled"`
	Mode      []string `yaml:"mode,omitempty" json:"mode"`
	Duration  uint64   `yaml:"duration,omitempty" json:"duration"`
	MatchType string   `yaml:"matchType,omitempty" json:"matchType"`
	Rules     []Rule   `yaml:"rules,omitempty" json:"rules"`
}

type Rule struct {
	MatchType string  `yaml:"matchType,omitempty" json:"matchType"`
	Groups    []Group `yaml:"groups,omitempty" json:"groups"`
}

type Group struct {
	Key          string `yaml:"key,omitempty" json:"key"`
	Name         string `yaml:"name,omitempty" json:"name"`
	Operator     string `yaml:"operator,omitempty" json:"operator"`
	OperatorName string `yaml:"operatorName,omitempty" json:"operatorName"`
	Value        string `yaml:"value,omitempty" json:"value"`
}
