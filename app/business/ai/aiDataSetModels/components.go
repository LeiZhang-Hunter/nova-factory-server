package aiDataSetModels

type AgentComponent struct {
	Downstream []interface{} `json:"downstream"`
	Obj        struct {
		ComponentName string                 `json:"component_name"`
		Inputs        []interface{}          `json:"inputs"`
		Output        interface{}            `json:"output"`
		Params        map[string]interface{} `json:"params"`
	} `json:"obj"`
	Upstream []interface{} `json:"upstream"`
}

type AgentGraph struct {
	Edges []interface{} `json:"edges"`
	Nodes []struct {
		Data struct {
			Label string `json:"label"`
			Name  string `json:"name"`
			Form  struct {
			} `json:"form,omitempty"`
		} `json:"data"`
		Dragging bool   `json:"dragging"`
		Height   int    `json:"height"`
		Id       string `json:"id"`
		Position struct {
			X float64 `json:"x"`
			Y float64 `json:"y"`
		} `json:"position"`
		PositionAbsolute struct {
			X float64 `json:"x"`
			Y float64 `json:"y"`
		} `json:"positionAbsolute"`
		Selected       bool   `json:"selected"`
		SourcePosition string `json:"sourcePosition"`
		TargetPosition string `json:"targetPosition"`
		Type           string `json:"type"`
		Width          int    `json:"width"`
	} `json:"nodes"`
}
