//go:build !iot

package iot

import "github.com/google/wire"

var ProviderSet = wire.NewSet()
