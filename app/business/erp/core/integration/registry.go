package integration

import (
	"errors"
	"nova-factory-server/app/business/erp/core/integration/api"
	_ "nova-factory-server/app/business/erp/core/integration/grasp"
	_ "nova-factory-server/app/business/erp/core/integration/kingdee"
	"strings"
	"sync"
)

type Registry struct {
	delegate *api.Registry
}

func newRegistry() *Registry {
	return &Registry{
		delegate: api.GetRegistry(),
	}
}

var (
	registryOnce sync.Once
	registryIns  *Registry
)

func GetRegistry() *Registry {
	registryOnce.Do(func() {
		registryIns = newRegistry()
	})
	return registryIns
}

func (r *Registry) Register(kind api.Kind, factory api.Factory) error {
	return r.delegate.Register(kind, factory)
}

func (r *Registry) Create(kind api.Kind) (api.Client, error) {
	return r.delegate.Create(kind)
}

func (r *Registry) CreateByType(tp string) (api.Client, error) {
	kind := resolveKind(tp)
	if kind == "" {
		return nil, errors.New("未支持的集成类型: " + tp)
	}
	if c, err := r.delegate.Create(kind); err == nil {
		return c, nil
	}
	return r.delegate.Create(api.Kind(strings.TrimSpace(tp)))
}

func Create(kind api.Kind) (api.Client, error) {
	return GetRegistry().Create(kind)
}

func CreateByType(tp string) (api.Client, error) {
	return GetRegistry().CreateByType(tp)
}

func resolveKind(tp string) api.Kind {
	t := strings.TrimSpace(tp)
	switch {
	case t == "":
		return ""
	case strings.EqualFold(t, string(api.KindGuanJiaPo)):
		return api.KindGuanJiaPo
	case strings.Contains(strings.ToLower(t), "gjp"):
		return api.KindGuanJiaPo
	case strings.Contains(t, "管家婆"):
		return api.KindGuanJiaPo
	case strings.EqualFold(t, string(api.KindKingdee)):
		return api.KindKingdee
	case strings.Contains(t, "金蝶"):
		return api.KindKingdee
	default:
		return api.Kind(t)
	}
}
