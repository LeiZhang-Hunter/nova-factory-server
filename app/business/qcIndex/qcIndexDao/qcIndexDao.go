package qcIndexDao

import (
	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/qcIndex/qcIndexApi"
	"nova-factory-server/app/business/qcIndex/qcIndexModel"

	"github.com/gin-gonic/gin"
)

type QcIndexDao interface {
	List(c *gin.Context, req *qcIndexApi.QcIndexListReq) (int64, []*qcIndexModel.QcIndex, error)
	Create(c *gin.Context, req *qcIndexApi.QcIndexCreateReq) (*baize.EmptyResponse, error)
	Update(c *gin.Context, req *qcIndexApi.QcIndexUpdateReq) (*baize.EmptyResponse, error)
	Delete(c *gin.Context, indexIds []int64) error
	GetById(c *gin.Context, indexId int64) (*qcIndexModel.QcIndex, error)
}
