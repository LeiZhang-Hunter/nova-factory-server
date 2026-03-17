//go:build !ai

package ai

import (
	"nova-factory-server/app/routes"

	"github.com/google/wire"
)

func NewGinEngine(
	app *routes.App) *AI {
	return &AI{}
}

var ProviderSet = wire.NewSet(NewGinEngine)
