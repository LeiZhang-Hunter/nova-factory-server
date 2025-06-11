package concurrency

type Goroutine struct {
	InitThreshold    int     `yaml:"initThreshold,omitempty" default:"16" validate:"gte=1"`
	MaxGoroutine     int     `yaml:"maxGoroutine,omitempty" default:"30" validate:"gte=1"`
	UnstableTolerate int     `yaml:"unstableTolerate,omitempty" default:"3" validate:"gte=1"`
	ChannelLenOfCap  float64 `yaml:"channelLenOfCap,omitempty" default:"0.4" validate:"gt=0"`
}

type Rtt struct {
	BlockJudgeThreshold string  `yaml:"blockJudgeThreshold,omitempty" default:"120%"`
	NewRttWeigh         float64 `yaml:"newRttWeigh,omitempty" default:"0.5" validate:"gte=0,lte=1"`
}

type Ratio struct {
	Multi             int `yaml:"multi,omitempty" default:"2" validate:"gt=1"`
	Linear            int `yaml:"linear,omitempty" default:"2" validate:"gt=1"`
	LinearWhenBlocked int `yaml:"linearWhenBlocked,omitempty" default:"4" validate:"gt=1"`
}

type Duration struct {
	Unstable int `yaml:"unstable,omitempty" default:"15" validate:"gte=1"`
	Stable   int `yaml:"stable,omitempty" default:"30" validate:"gte=1"`
}

type Config struct {
	Enable    bool       `yaml:"enabled,omitempty"`
	Goroutine *Goroutine `yaml:"goroutine,omitempty"`
	Rtt       *Rtt       `yaml:"rtt,omitempty"`
	Ratio     *Ratio     `yaml:"ratio,omitempty"`
	Duration  *Duration  `yaml:"duration,omitempty"`
}
