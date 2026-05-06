package camera_exporter

type Config struct {
	Address string `yaml:"address,omitempty" default:"127.0.0.1:6000"`
}
