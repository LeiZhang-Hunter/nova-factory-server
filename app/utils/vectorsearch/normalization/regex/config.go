package regex

type Config struct {
	Pattern     string `yaml:"pattern,omitempty"`
	Replacement string `yaml:"replacement,omitempty"`
}
