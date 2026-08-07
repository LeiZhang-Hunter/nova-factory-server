package dao

import (
	"context"

	"nova-factory-server/app/business/data/models/entity"
)

type IPipelineRuleDAO interface {
	Create(context.Context, *entity.PipelineRule) error
	GetByID(context.Context, string) (*entity.PipelineRule, error)
	GetByName(context.Context, string, string) (*entity.PipelineRule, error)
	List(context.Context, int, int, string, string, string) ([]entity.PipelineRule, int64, error)
	Update(context.Context, *entity.PipelineRule) error
	Delete(context.Context, string) error
	CountConnectionReferences(context.Context, string) (int64, error)
}

type IServiceConnectionDAO interface {
	Create(context.Context, *entity.ServiceConnection) error
	GetByID(context.Context, string) (*entity.ServiceConnection, error)
	GetByName(context.Context, string, string) (*entity.ServiceConnection, error)
	List(context.Context, int, int, string, string, string) ([]entity.ServiceConnection, int64, error)
	Update(context.Context, *entity.ServiceConnection) error
	Delete(context.Context, string) error
}
