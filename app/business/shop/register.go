//go:build shop
// +build shop

package shop

import (
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	GinProviderSet,
)
