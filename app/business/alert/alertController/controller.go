package alertController

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewAlert, NewAlertTemplate, wire.Struct(new(Controller), "*"))

type Controller struct {
	Alert         *Alert
	AlertTemplate *AlertTemplate
}
