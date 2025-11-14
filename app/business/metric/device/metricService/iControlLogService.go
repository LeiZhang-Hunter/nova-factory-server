package metricService

import (
	"context"
	v1 "github.com/novawatcher-io/nova-factory-payload/metric/grpc/v1"
)

type IControlLogService interface {
	// Export 导入控制日志
	Export(ctx context.Context, request *v1.ExportControlLogRequest) error
}
