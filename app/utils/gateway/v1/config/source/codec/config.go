package codec

import (
	"nova-factory-server/app/utils/gateway/v1/config/cfg"
)

type Config struct {
	Type          string `yaml:"type,omitempty" default:"json"`
	cfg.CommonCfg `yaml:",inline"`
}
