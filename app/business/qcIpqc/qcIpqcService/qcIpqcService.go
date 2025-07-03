package qcIpqcService

import (
	"context"
	"nova-factory-server/app/business/qcIpqc/qcIpqcApi"
)

type IQcIpqcService interface {
	GetQcIpqcList(ctx context.Context, req *qcIpqcApi.QcIpqcQueryReq) ([]*qcIpqcApi.QcIpqc, int64, error)
	GetQcIpqcById(ctx context.Context, id int64) (*qcIpqcApi.QcIpqc, error)
	CreateQcIpqc(ctx context.Context, req *qcIpqcApi.QcIpqcCreateReq) error
	UpdateQcIpqc(ctx context.Context, req *qcIpqcApi.QcIpqcUpdateReq) error
	DeleteQcIpqcByIds(ctx context.Context, ids []int64) error
}
