package cameraServiceImpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(
	NewCameraService)
