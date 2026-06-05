package normalization

import (
	"fmt"
	"nova-factory-server/app/utils/vectorsearch/normalization/api"
)

// Pipeline 表示一组按顺序执行的归一化步骤。
type Pipeline struct {
	config api.Config
	steps  []api.Step
	err    error
}

// NewPipeline 创建新的归一化 pipeline。
func NewPipeline(config api.Config) *Pipeline {
	filtered := make([]api.Step, 0, len(config.Interceptors))
	var initErr error
	for _, stepConfig := range config.Interceptors {
		if stepConfig == nil {
			continue
		}
		if stepConfig.Enabled != nil && !*stepConfig.Enabled {
			continue
		}

		step, err := api.NewStep(stepConfig.Type)
		if err != nil {
			initErr = err
			break
		}

		currentConfig := *stepConfig
		if currentConfig.Name == "" {
			currentConfig.Name = step.Type()
		}
		if err := step.Init(currentConfig); err != nil {
			initErr = fmt.Errorf("init normalization step %q failed: %w", currentConfig.Type, err)
			break
		}
		filtered = append(filtered, step)
	}
	return &Pipeline{
		config: config,
		steps:  filtered,
		err:    initErr,
	}
}

// Normalize 执行全部步骤并返回结果。
func (p *Pipeline) Normalize(input string) (Result, error) {
	ctx := NewContext(input)
	if p == nil {
		return ctx.Result(), nil
	}
	if p.err != nil {
		return Result{}, p.err
	}
	for _, step := range p.steps {
		if step == nil {
			continue
		}
		if err := step.Apply(ctx); err != nil {
			return Result{}, fmt.Errorf("normalization step %q failed: %w", step.Name(), err)
		}
	}
	return ctx.Result(), nil
}
