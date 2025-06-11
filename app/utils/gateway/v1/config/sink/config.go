package sink

import (
	"nova-factory-server/app/utils/gateway/v1/config/cfg"
	"nova-factory-server/app/utils/gateway/v1/config/concurrency"
	"nova-factory-server/app/utils/gateway/v1/config/sink/codec"
)

type Config struct {
	Enabled     *bool              `yaml:"enabled,omitempty"`
	Name        string             `yaml:"name,omitempty"`
	Type        string             `yaml:"type,omitempty" validate:"required"`
	Properties  cfg.CommonCfg      `yaml:",inline"`
	Parallelism int                `yaml:"parallelism,omitempty" default:"1" validate:"required,gte=1,lte=100"`
	Codec       codec.Config       `yaml:"codec,omitempty" validate:"dive"`
	Concurrency concurrency.Config `yaml:"concurrency,omitempty"`
}
