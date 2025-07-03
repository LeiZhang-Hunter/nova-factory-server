package qcOqcService

import (
	"context"
	"nova-factory-server/app/business/qcOqc/qcOqcApi"
)

// IQcOqcService 出货检验单服务接口
type IQcOqcService interface {
	// GetQcOqcList 获取出货检验单列表
	GetQcOqcList(ctx context.Context, req *qcOqcApi.QcOqcQueryReq) ([]*qcOqcApi.QcOqc, int64, error)

	// GetQcOqcById 根据ID获取出货检验单
	GetQcOqcById(ctx context.Context, oqcId int64) (*qcOqcApi.QcOqc, error)

	// CreateQcOqc 创建出货检验单
	CreateQcOqc(ctx context.Context, req *qcOqcApi.QcOqcCreateReq) error

	// UpdateQcOqc 更新出货检验单
	UpdateQcOqc(ctx context.Context, req *qcOqcApi.QcOqcUpdateReq) error

	// DeleteQcOqc 删除出货检验单
	DeleteQcOqc(ctx context.Context, req *qcOqcApi.QcOqcDeleteReq) error
}
