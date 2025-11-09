package scheduler

import "time"

// TaskConfig 读取任务配置
type TaskConfig struct {
	Enabled  bool          `yaml:"enabled" default:"false"`
	Version  string        `yaml:"version" default:"v1"`
	Host     string        `yaml:"host"`
	PollTime time.Duration `yaml:"poll_time" default:"5m"`
	Limit    uint64        `yaml:"limit" default:"50"`
}

// Config 配置
type Config struct {
	//Routers map[string]*Router
	Task TaskConfig `yaml:"task,omitempty"`
}
