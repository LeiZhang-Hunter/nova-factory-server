package gatewaymodels

import "nova-factory-server/app/baize"

// AIAgentOrchestration 智能体编排配置。
type AIAgentOrchestration struct {
	ID         int64  `json:"id,string" gorm:"column:id"`
	AgentID    int64  `json:"agentId,string" gorm:"column:agent_id"`
	Content    string `json:"content" gorm:"column:content"`
	Config     string `json:"config" gorm:"column:config"`
	ConfigMd5  string `json:"-" gorm:"column:config_md5"`
	ContentMd5 string `json:"-" gorm:"column:content_md5"`
	DeptID     int64  `json:"deptId" gorm:"column:dept_id"`
	baize.BaseEntity
	State int32 `json:"state" gorm:"column:state"`
}

// AIAgentOrchestrationUpsert 智能体编排保存参数。
type AIAgentOrchestrationUpsert struct {
	AgentID int64  `json:"agentId,string"`
	Content string `json:"content"`
	Config  string `json:"-"`
}

type Node struct {
	Id         string `json:"id"`
	Type       string `json:"type"`
	Dimensions struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"dimensions"`
	ComputedPosition struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z int     `json:"z"`
	} `json:"computedPosition"`
	HandleBounds struct {
		Source []struct {
			Id       interface{} `json:"id"`
			Type     string      `json:"type"`
			NodeId   string      `json:"nodeId"`
			Position string      `json:"position"`
			X        float64     `json:"x"`
			Y        float64     `json:"y"`
			Width    int         `json:"width"`
			Height   int         `json:"height"`
		} `json:"source"`
		Target []struct {
			Id       interface{} `json:"id"`
			Type     string      `json:"type"`
			NodeId   string      `json:"nodeId"`
			Position string      `json:"position"`
			X        float64     `json:"x"`
			Y        float64     `json:"y"`
			Width    int         `json:"width"`
			Height   int         `json:"height"`
		} `json:"target"`
	} `json:"handleBounds"`
	Selected    bool `json:"selected"`
	Dragging    bool `json:"dragging"`
	Resizing    bool `json:"resizing"`
	Initialized bool `json:"initialized"`
	IsParent    bool `json:"isParent"`
	Position    struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"position"`
	Data struct {
		Title       string `json:"title"`
		Subtitle    string `json:"subtitle"`
		Icon        string `json:"icon"`
		Accent      string `json:"accent"`
		Description string `json:"description"`
		Metrics     []struct {
			Label string `json:"label"`
			Value string `json:"value"`
		} `json:"metrics"`
		Config          *GraphAISubAgentUpsert `json:"config"`
		SourceAgentId   string                 `json:"sourceAgentId,omitempty"`
		SourceAgentName string                 `json:"sourceAgentName,omitempty"`
	} `json:"data"`
	Events struct {
	} `json:"events"`
}

// AgentOrchestrationConfig 前端出来的编排数据结构
type AgentOrchestrationConfig struct {
	AgentId string `json:"agentId"`
	Nodes   []Node `json:"nodes"`
	Edges   []struct {
		Id     string `json:"id"`
		Type   string `json:"type"`
		Source string `json:"source"`
		Target string `json:"target"`
		Data   struct {
			Label string `json:"label"`
		} `json:"data"`
		Events struct {
		} `json:"events"`
		Label     string `json:"label"`
		Animated  bool   `json:"animated"`
		MarkerEnd struct {
			Type   string  `json:"type"`
			Color  string  `json:"color"`
			Width  float64 `json:"width"`
			Height float64 `json:"height"`
		} `json:"markerEnd"`
		Style struct {
			Stroke      string  `json:"stroke"`
			StrokeWidth float64 `json:"strokeWidth"`
		} `json:"style"`
		SourceNode struct {
			Id         string `json:"id"`
			Type       string `json:"type"`
			Dimensions struct {
				Width  float64 `json:"width"`
				Height float64 `json:"height"`
			} `json:"dimensions"`
			ComputedPosition struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
				Z float64 `json:"z"`
			} `json:"computedPosition"`
			HandleBounds struct {
				Source []struct {
					Id       interface{} `json:"id"`
					Type     string      `json:"type"`
					NodeId   string      `json:"nodeId"`
					Position string      `json:"position"`
					X        float64     `json:"x"`
					Y        float64     `json:"y"`
					Width    float64     `json:"width"`
					Height   float64     `json:"height"`
				} `json:"source"`
				Target interface{} `json:"target"`
			} `json:"handleBounds"`
			Selected    bool `json:"selected"`
			Dragging    bool `json:"dragging"`
			Resizing    bool `json:"resizing"`
			Initialized bool `json:"initialized"`
			IsParent    bool `json:"isParent"`
			Position    struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
			} `json:"position"`
			Data struct {
				Title       string `json:"title"`
				Subtitle    string `json:"subtitle"`
				Icon        string `json:"icon"`
				Accent      string `json:"accent"`
				Description string `json:"description"`
				Metrics     []struct {
					Label string `json:"label"`
					Value string `json:"value"`
				} `json:"metrics"`
				Config struct {
					Name           string   `json:"name"`
					Description    string   `json:"description"`
					Model          string   `json:"model"`
					Prompt         string   `json:"prompt"`
					Status         string   `json:"status"`
					Tags           []string `json:"tags"`
					ExecuteMode    string   `json:"executeMode"`
					RoutingPolicy  string   `json:"routingPolicy"`
					ExpectedOutput string   `json:"expectedOutput"`
				} `json:"config"`
			} `json:"data"`
			Events struct {
			} `json:"events"`
		} `json:"sourceNode"`
		TargetNode struct {
			Id         string `json:"id"`
			Type       string `json:"type"`
			Dimensions struct {
				Width  float64 `json:"width"`
				Height float64 `json:"height"`
			} `json:"dimensions"`
			ComputedPosition struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
				Z float64 `json:"z"`
			} `json:"computedPosition"`
			HandleBounds struct {
				Source []struct {
					Id       interface{} `json:"id"`
					Type     string      `json:"type"`
					NodeId   string      `json:"nodeId"`
					Position string      `json:"position"`
					X        float64     `json:"x"`
					Y        float64     `json:"y"`
					Width    float64     `json:"width"`
					Height   float64     `json:"height"`
				} `json:"source"`
				Target []struct {
					Id       interface{} `json:"id"`
					Type     string      `json:"type"`
					NodeId   string      `json:"nodeId"`
					Position string      `json:"position"`
					X        float64     `json:"x"`
					Y        float64     `json:"y"`
					Width    float64     `json:"width"`
					Height   float64     `json:"height"`
				} `json:"target"`
			} `json:"handleBounds"`
			Selected    bool `json:"selected"`
			Dragging    bool `json:"dragging"`
			Resizing    bool `json:"resizing"`
			Initialized bool `json:"initialized"`
			IsParent    bool `json:"isParent"`
			Position    struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
			} `json:"position"`
			Data struct {
				Title       string `json:"title"`
				Subtitle    string `json:"subtitle"`
				Icon        string `json:"icon"`
				Accent      string `json:"accent"`
				Description string `json:"description"`
				Metrics     []struct {
					Label string `json:"label"`
					Value string `json:"value"`
				} `json:"metrics"`
				Config struct {
					Name                   string        `json:"name"`
					Description            string        `json:"description"`
					Instruction            string        `json:"instruction"`
					McpEnabled             bool          `json:"mcpEnabled"`
					McpServerIds           []interface{} `json:"mcpServerIds"`
					McpServerEnabledIds    []interface{} `json:"mcpServerEnabledIds"`
					LocalToolEnabled       bool          `json:"localToolEnabled"`
					LocalTools             []interface{} `json:"localTools"`
					AllowMcpServerIdsTools struct {
					} `json:"allowMcpServerIdsTools"`
					SourceAgentId   string   `json:"sourceAgentId"`
					SourceAgentName string   `json:"sourceAgentName"`
					Enable          bool     `json:"enable"`
					Model           string   `json:"model"`
					Prompt          string   `json:"prompt"`
					Status          string   `json:"status"`
					Tags            []string `json:"tags"`
					Speciality      string   `json:"speciality"`
					ReturnField     string   `json:"returnField"`
					Sla             string   `json:"sla"`
					SubAgentType    string   `json:"subAgentType"`
					CoreSubAgent    string   `json:"coreSubAgent"`
				} `json:"config"`
				SourceAgentId   string `json:"sourceAgentId"`
				SourceAgentName string `json:"sourceAgentName"`
			} `json:"data"`
			Events struct {
			} `json:"events"`
		} `json:"targetNode"`
		SourceX float64 `json:"sourceX"`
		SourceY float64 `json:"sourceY"`
		TargetX float64 `json:"targetX"`
		TargetY float64 `json:"targetY"`
	} `json:"edges"`
	SavedAt int64 `json:"savedAt"`
}

// AgentLoadConfig 下发的配置
type AgentLoadConfig struct {
	Agent    *AIAgent            `json:"agent"`
	SubAgent []*AISubAgentUpsert `json:"subAgent"`
}
