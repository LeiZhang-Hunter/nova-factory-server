package qcIndexService

import (
	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/qcIndex/qcIndexApi"

	"github.com/gin-gonic/gin"
)

type QcIndexService interface {
	List(c *gin.Context, req *qcIndexApi.QcIndexListReq) (*qcIndexApi.QcIndexListRes, error)
	Create(c *gin.Context, req *qcIndexApi.QcIndexCreateReq) (*baize.EmptyResponse, error)
	Update(c *gin.Context, req *qcIndexApi.QcIndexUpdateReq) (*baize.EmptyResponse, error)
	Delete(c *gin.Context, indexIds []int64) error
	GetById(c *gin.Context, indexId int64) (*qcIndexApi.QcIndexData, error)
}
