package controller

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewPipelineRuleController,
	NewServiceConnectionController,
	NewCollectorController,
	NewPipelineRuleAgentController,
	NewCollectorRuleController,
	wire.Struct(new(DataControllers), "*"),
)

type DataControllers struct {
	PipelineRuleController      *PipelineRuleController
	ServiceConnectionController *ServiceConnectionController
	CollectorController         *CollectorController
	PipelineRuleAgentController *PipelineRuleAgentController
	CollectorRuleController     *CollectorRuleController
}
