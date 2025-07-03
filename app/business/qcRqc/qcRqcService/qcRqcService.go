package qcRqcService

import (
	"context"
	"nova-factory-server/app/business/qcRqc/qcRqcApi"
)

// IQcRqcService 退料检验单服务接口
type IQcRqcService interface {
	// GetQcRqcList 获取退料检验单列表
	GetQcRqcList(ctx context.Context, req *qcRqcApi.QcRqcQueryReq) ([]*qcRqcApi.QcRqc, int64, error)

	// GetQcRqcById 根据ID获取退料检验单
	GetQcRqcById(ctx context.Context, rqcId int64) (*qcRqcApi.QcRqc, error)

	// CreateQcRqc 创建退料检验单
	CreateQcRqc(ctx context.Context, req *qcRqcApi.QcRqcCreateReq) error

	// UpdateQcRqc 更新退料检验单
	UpdateQcRqc(ctx context.Context, req *qcRqcApi.QcRqcUpdateReq) error

	// DeleteQcRqcByIds 批量删除退料检验单
	DeleteQcRqcByIds(ctx context.Context, rqcIds []int64) error
}
