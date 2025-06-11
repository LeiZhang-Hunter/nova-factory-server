package pipeline

import (
	"nova-factory-server/app/utils/gateway/v1/config/interceptor"
	"nova-factory-server/app/utils/gateway/v1/config/queue"
	"nova-factory-server/app/utils/gateway/v1/config/sink"
	"nova-factory-server/app/utils/gateway/v1/config/source"
	"time"
)

type PipelineConfig struct {
	Pipelines []Config `yaml:"pipelines" validate:"dive,required"`
}

func NewPipelineConfig() *PipelineConfig {
	return &PipelineConfig{
		Pipelines: make([]Config, 0),
	}
}

type Config struct {
	Name             string        `yaml:"name,omitempty" validate:"required"`
	CleanDataTimeout time.Duration `yaml:"cleanDataTimeout,omitempty" default:"5s"`

	Queue        *queue.Config         `yaml:"queue,omitempty" validate:"dive,required"`
	Interceptors []*interceptor.Config `yaml:"interceptors,omitempty"`
	Sources      []*source.Config      `yaml:"sources,omitempty" validate:"dive,required"`
	Sink         *sink.Config          `yaml:"sink,omitempty" validate:"dive,required"`
}

func NewConfig() *Config {
	return &Config{
		Interceptors: make([]*interceptor.Config, 0),
		Sources:      make([]*source.Config, 0),
	}
}
