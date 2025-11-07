package prediction

import "time"

type Config struct {
	Parallelism   uint16            `yaml:"parallelism,omitempty" default:"1" validate:"required,gte=1,lte=100"`
	Model         string            //预测模型
	TimeWindow    time.Duration     `yaml:"time_window"`    //预测时间窗口
	PredictLength uint64            `yaml:"predict_length"` //预测长度
	Param         map[string]string `json:"param"`          //预测参数
}
