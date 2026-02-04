package craft

type FLOW_TYPE string

type EdgeMapType map[string]FLOW_TYPE

type NodeType map[string]FLOW_TYPE

var EdgeMap EdgeMapType = map[string]FLOW_TYPE{
	"default":    "default",
	"straight":   "straight",
	"step":       "step",
	"smoothstep": "smoothstep",
	"bezier":     "bezier",
	"custom":     "custom",
}

const (
	START_NAME        = "start-node"
	NODE_PROCESS_TYPE = "process"
	NODE_INIT_TYPE    = "init"
	NODE_FINISH_TYPE  = "finish"
)

var NodeMap NodeType = map[string]FLOW_TYPE{
	"process": "process",
	"product": "product",
	"default": "default",
	"init":    "init",
	"finish":  "finish",
}
