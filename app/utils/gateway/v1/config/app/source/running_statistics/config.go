package running_statistics

type GroupRule struct {
	Key          string `yaml:"key"`
	Name         string `yaml:"name"`
	Operator     string `yaml:"operator"`
	OperatorName string `yaml:"operatorName"`
	Value        string `yaml:"value"`
}

type StatisticsRule struct {
	MatchType string      `yaml:"matchType"`
	Groups    []GroupRule `yaml:"groups"`
	RunStatus int         `yaml:"run_status"`
}

type Expression struct {
	Rules []StatisticsRule `yaml:"rules"`
}

type DeviceStatistics struct {
	DeviceId   string     `yaml:"deviceId"`
	Table      string     `yaml:"table"`
	Expression Expression `yaml:"expression"`
}

type Config struct {
	DeviceStatistics []*DeviceStatistics `yaml:"statistics"`
}
