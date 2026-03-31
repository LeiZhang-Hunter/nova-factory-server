package cameraController

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewCameraController,
	wire.Struct(new(CameraController), "*"))

type CameraController struct {
	Camera *Camera
}
