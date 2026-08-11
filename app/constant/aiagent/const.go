package aiagent

// CORE sub agent类型
var (
	CORE = "core"
)

// agent字典
var (
	SubAgentType = "sub_agent_type"
	CoreSubAgent = "core_sub_agent"
)

// node type
var (
	Master = "master"
	Sub    = "sub"
)

type ConfigUpdate string

var ConfigPublishType ConfigUpdate = "publish"
var ConfigInitType ConfigUpdate = "init"
var ConfigRemoveType ConfigUpdate = "remove"

var AgentType string = "agent"

// DataAgentType 数据平台（data 模块）使用的智能体类型。
var DataAgentType = "data"
