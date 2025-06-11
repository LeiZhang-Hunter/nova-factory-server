package codec

import (
	"nova-factory-server/app/utils/gateway/v1/config/cfg"
)

type Config struct {
	Type          string `yaml:"type,omitempty" default:"json"`
	PrintEvents   bool   `yaml:"printEvents,omitempty"`
	cfg.CommonCfg `yaml:",inline"`
}
