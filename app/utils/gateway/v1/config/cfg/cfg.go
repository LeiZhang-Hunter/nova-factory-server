package cfg

import (
	"github.com/go-playground/validator/v10"
	"nova-factory-server/app/utils/yaml"
)

type CommonCfg map[string]interface{}

type Validator interface {
	Validate() error
}

func NewCommonCfg() CommonCfg {
	return make(map[string]interface{})
}

func Pack(config interface{}) (CommonCfg, error) {
	if config == nil {
		return nil, nil
	}

	out, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}

	ret := NewCommonCfg()
	err = yaml.Unmarshal(out, &ret)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// UnpackFromCommonCfg create an Unpack struct from CommonCfg
func UnpackFromCommonCfg(c CommonCfg, config interface{}) *UnPack {
	if c == nil {
		c = NewCommonCfg()
	}

	out, err := yaml.Marshal(c)
	return NewUnpack(out, config, err).unpack()
}

type UnPack struct {
	content []byte
	config  interface{}
	err     error
}

func NewUnpack(raw []byte, config interface{}, err error) *UnPack {
	return &UnPack{
		content: raw,
		config:  config,
		err:     err,
	}
}

func (u *UnPack) unpack() *UnPack {
	if u.err != nil {
		return u
	}
	if u.content == nil {
		return u
	}

	err := yaml.UnmarshalWithPrettyError(u.content, u.config)
	u.err = err
	return u
}

// validate will affect by `validate` tag, and call Validate() function in Config.
func validate(config interface{}) error {
	if config == nil {
		return nil
	}

	validate := validator.New()
	err := validate.Struct(config)
	if err != nil {
		return err
	}

	if cfg, ok := config.(Validator); ok {
		return cfg.Validate()
	}
	return nil
}

func (u *UnPack) Validate() *UnPack {
	if u.err != nil {
		return u
	}

	if err := validate(u.config); err != nil {
		u.err = err
	}

	return u
}

// Do return the error
func (u *UnPack) Do() error {
	return u.err
}
