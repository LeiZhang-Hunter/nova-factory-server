package source

import (
	"nova-factory-server/app/utils/gateway/v1/config/cfg"
	"nova-factory-server/app/utils/gateway/v1/config/source/codec"
)

type Config struct {
	Enabled         *bool                  `yaml:"enabled,omitempty"`
	Name            string                 `yaml:"name,omitempty"`
	Type            string                 `yaml:"type,omitempty" validate:"required"`
	Properties      cfg.CommonCfg          `yaml:",inline"`
	FieldsUnderRoot bool                   `yaml:"fieldsUnderRoot,omitempty" default:"false"`
	FieldsUnderKey  string                 `yaml:"fieldsUnderKey,omitempty" default:"fields"`
	Fields          map[string]interface{} `yaml:"fields,omitempty"`
	FieldsFromEnv   map[string]string      `yaml:"fieldsFromEnv,omitempty"`
	FieldsFromPath  map[string]string      `yaml:"fieldsFromPath,omitempty"`
	Codec           *codec.Config          `yaml:"codec,omitempty"`

	TimestampKey      string `yaml:"timestampKey,omitempty"`
	TimestampLocation string `yaml:"timestampLocation,omitempty"`
	TimestampLayout   string `yaml:"timestampLayout,omitempty"`
	BodyKey           string `yaml:"bodyKey,omitempty"`
}
