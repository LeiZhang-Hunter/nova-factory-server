package alertController

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewAlert, NewAlertTemplate, NewAlertLog, NewAlertAction, NewAlertAiReason, wire.Struct(new(Controller), "*"))

type Controller struct {
	Alert         *Alert
	AlertTemplate *AlertTemplate
	AlertLog      *AlertLog
	AlertAction   *AlertAction
	AlertAiReason *AlertAiReason
}
