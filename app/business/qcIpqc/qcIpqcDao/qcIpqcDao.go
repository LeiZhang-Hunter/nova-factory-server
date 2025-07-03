package qcIpqcDao

import (
	"context"
	"nova-factory-server/app/business/qcIpqc/qcIpqcApi"
	"nova-factory-server/app/business/qcIpqc/qcIpqcModels"
)

type IQcIpqcDao interface {
	SelectQcIpqcList(ctx context.Context, req *qcIpqcApi.QcIpqcQueryReq) ([]*qcIpqcModels.QcIpqc, int64, error)
	SelectQcIpqcById(ctx context.Context, id int64) (*qcIpqcModels.QcIpqc, error)
	InsertQcIpqc(ctx context.Context, ipqc *qcIpqcModels.QcIpqc) error
	UpdateQcIpqc(ctx context.Context, ipqc *qcIpqcModels.QcIpqc) error
	DeleteQcIpqcByIds(ctx context.Context, ids []int64) error
}
