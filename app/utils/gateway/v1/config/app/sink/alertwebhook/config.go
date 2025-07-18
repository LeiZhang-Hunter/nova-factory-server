package alertwebhook

import "nova-factory-server/app/utils/gateway/v1/config/logalert"

type Config struct {
	Addr                 string `yaml:"addr,omitempty"`
	logalert.AlertConfig `yaml:",inline"`
}
