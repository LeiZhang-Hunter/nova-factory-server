package replace

type Config struct {
	Old string `yaml:"old,omitempty"`
	New string `yaml:"new,omitempty"`
}
