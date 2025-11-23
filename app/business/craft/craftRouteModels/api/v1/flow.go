package v1

type DeviceInfo struct {
	Address string `json:"address"`
}

type DeviceAction struct {
	Address     string      `json:"address"`
	DeviceId    string      `json:"deviceId"`
	DataId      string      `json:"dataId"`
	Value       interface{} `json:"value"`
	DataFormat  string      `json:"dataFormat"`
	ControlMode string      `json:"control_mode"`
	Condition   string      `json:"condition"`
	Interval    string      `json:"interval"`
}

type DeviceRuleInfo struct {
	DeviceId string `json:"deviceId"`
	DataId   string `json:"dataId"`
}

type DeviceRule struct {
	DataId []DeviceRuleInfo `json:"dataId"`
	Rule   string           `json:"rule"`
}

// DeviceTriggerRule 设备触发规则,阈值规则
type DeviceTriggerRule struct {
	Name       string          `json:"name"`
	Rule       *DeviceRule     `json:"rule"`
	Prompt     string          `json:"prompt"`
	DeviceInfo *DeviceInfo     `json:"deviceInfo"`
	Actions    []*DeviceAction `json:"actions"`
}

// PidRules 规则
type PidRules struct {
	Proportional int             `json:"proportional"`
	Integral     int             `json:"integral"`
	Derivative   int             `json:"derivative"`
	ActualSignal int             `json:"actualSignal"`
	DeviceId     string          `json:"device_id"`
	DataId       string          `json:"data_id"`
	Actions      []*DeviceAction `json:"actions"`
}

// ControlRules 控制算法
type ControlRules struct {
	TriggerRules *DeviceTriggerRule `json:"trigger_rules"`
	PidRules     *PidRules          `json:"pid_rules"`
}

// ProcessContext 工序内容
type ProcessContext struct {
	ContentID      uint64        `gorm:"column:content_id;primaryKey;autoIncrement:true;comment:内容ID" json:"content_id,string"` // 内容ID
	ProcessID      uint64        `gorm:"column:process_id;not null;comment:工序ID" json:"process_id,string" binding:"required"`   // 工序ID
	ControlName    string        `json:"control_name"`
	ControllerType string        `json:"controller_type"`
	ControlRules   *ControlRules `json:"control_rules"`
}

type Process struct {
	Name      string
	ProcessId uint64
	Context   []ProcessContext
}

type ProcessEdge struct {
	Source      uint64 `json:"source,string"`
	Target      uint64 `json:"target,string"`
	ProcessCode string `gorm:"column:process_code;comment:工序编码" json:"process_code"`   // 工序编码
	ProcessName string `gorm:"column:process_name;comment:工序名称" json:"process_name"`   // 工序名称
	KeyFlag     string `gorm:"column:key_flag;default:N;comment:关键工序" json:"key_flag"` // 关键工序
	IsCheck     string `gorm:"column:is_check;default:N;comment:是否检验" json:"is_check"` // 是否检验
}

type Begin struct {
	NextProcessId uint64 `json:"next_process_id,string"`
}

// Router 工艺路线
type Router struct {
	Processes []*Process
	Edge      []*ProcessEdge
	Begin     *Begin
	Name      string
	Id        uint64
	Md5       string
}

func NewRouter() *Router {
	return &Router{
		Processes: make([]*Process, 0),
		Edge:      make([]*ProcessEdge, 0),
		Begin:     &Begin{},
		Name:      "",
		Id:        0,
		Md5:       "",
	}
}

func NweDeviceTriggerRule() *DeviceTriggerRule {
	return &DeviceTriggerRule{
		Rule:    &DeviceRule{},
		Actions: make([]*DeviceAction, 0),
	}
}

func NewControlRules() *ControlRules {
	return &ControlRules{
		TriggerRules: NweDeviceTriggerRule(),
		PidRules: &PidRules{
			Actions: make([]*DeviceAction, 0),
		},
	}
}
