//go:build erp
// +build erp

package erp_api

import (
	qqdcontroller "nova-factory-server/app/business/erp_api/controller/qqd"
	"nova-factory-server/app/routes"

	"github.com/google/wire"
)

var GinProviderSet = wire.NewSet(NewGinEngine)

func NewGinEngine(app *routes.App, qqd *qqdcontroller.Controller) *ErpAPI {
	group := app.Engine.Group("/api/v1/erp-api")
	qqd.PublicRoutes(group)
	return &ErpAPI{}
}
