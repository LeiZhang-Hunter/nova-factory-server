package ipc

type Device struct {
	ID       string `yaml:"id,omitempty"`
	Name     string `yaml:"name,omitempty"`
	Address  string `yaml:"address,omitempty"`
	Port     uint16 `yaml:"port,omitempty"`
	Brand    string `yaml:"brand,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}
