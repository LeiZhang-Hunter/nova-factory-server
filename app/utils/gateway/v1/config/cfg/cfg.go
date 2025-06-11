package cfg

import "nova-factory-server/app/utils/yaml"

type CommonCfg map[string]interface{}

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
