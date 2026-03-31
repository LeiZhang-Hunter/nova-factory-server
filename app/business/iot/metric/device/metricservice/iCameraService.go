package metricservice

import (
	"context"
	v1 "github.com/novawatcher-io/nova-factory-payload/camera/v1"
)

func main() {

}

type ICameraService interface {
	Report(ctx context.Context, req *v1.CameraData) error
}
