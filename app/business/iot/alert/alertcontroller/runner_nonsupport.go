//go:build !ai
// +build !ai

package alertcontroller

import "nova-factory-server/app/business/iot/alert/alertmodels"

// Runner degrades to a no-op worker when AI support is disabled.
type Runner struct{}

var runner = &Runner{}

func NewRunner() *Runner {
	return runner
}

func GetAlertRunner() *Runner {
	return runner
}

func (r *Runner) Push(_ *alertmodels.AlertLogInfo) {}

func (r *Runner) Run() {}

func (r *Runner) Stop() {}
