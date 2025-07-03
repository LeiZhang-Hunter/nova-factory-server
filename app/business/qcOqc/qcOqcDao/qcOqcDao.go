package qcOqcDao

import (
	"context"
	"nova-factory-server/app/business/qcOqc/qcOqcApi"
	"nova-factory-server/app/business/qcOqc/qcOqcModels"
)

// IQcOqcDao 出货检验单数据访问接口
type IQcOqcDao interface {
	// SelectQcOqcList 查询出货检验单列表
	SelectQcOqcList(ctx context.Context, req *qcOqcApi.QcOqcQueryReq) ([]*qcOqcModels.QcOqc, int64, error)

	// SelectQcOqcById 根据ID查询出货检验单
	SelectQcOqcById(ctx context.Context, id int64) (*qcOqcModels.QcOqc, error)

	// InsertQcOqc 新增出货检验单
	InsertQcOqc(ctx context.Context, oqc *qcOqcModels.QcOqc) error

	// UpdateQcOqc 更新出货检验单
	UpdateQcOqc(ctx context.Context, oqc *qcOqcModels.QcOqc) error

	// DeleteQcOqcByIds 批量删除出货检验单
	DeleteQcOqcByIds(ctx context.Context, ids []int64) error
}
