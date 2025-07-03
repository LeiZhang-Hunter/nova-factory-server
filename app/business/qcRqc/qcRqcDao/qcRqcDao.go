package qcRqcDao

import (
	"context"
	"nova-factory-server/app/business/qcRqc/qcRqcApi"
	"nova-factory-server/app/business/qcRqc/qcRqcModels"
)

// IQcRqcDao 退料检验单数据访问接口
type IQcRqcDao interface {
	// SelectQcRqcList 查询退料检验单列表
	SelectQcRqcList(ctx context.Context, req *qcRqcApi.QcRqcQueryReq) ([]*qcRqcModels.QcRqc, int64, error)

	// SelectQcRqcById 根据ID查询退料检验单
	SelectQcRqcById(ctx context.Context, id int64) (*qcRqcModels.QcRqc, error)

	// InsertQcRqc 新增退料检验单
	InsertQcRqc(ctx context.Context, rqc *qcRqcModels.QcRqc) error

	// UpdateQcRqc 更新退料检验单
	UpdateQcRqc(ctx context.Context, rqc *qcRqcModels.QcRqc) error

	// DeleteQcRqcByIds 批量删除退料检验单
	DeleteQcRqcByIds(ctx context.Context, ids []int64) error
}
