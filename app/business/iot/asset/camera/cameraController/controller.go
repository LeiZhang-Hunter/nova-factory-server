package cameraController

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewCameraController,
	wire.Struct(new(Controller), "*"))

type Controller struct {
	Camera *CameraController
}
