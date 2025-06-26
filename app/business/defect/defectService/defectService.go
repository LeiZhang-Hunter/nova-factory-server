package defectService

import (
	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/defect/defectApi"

	"github.com/gin-gonic/gin"
)

type DefectService interface {
	List(c *gin.Context, route *defectApi.DefectListReq) (*defectApi.DefectListRes, error)
	Create(c *gin.Context, req *defectApi.DefectCreateReq) (*baize.EmptyResponse, error)
	Update(c *gin.Context, req *defectApi.DefectUpdateReq) (*baize.EmptyResponse, error)
	Delete(c *gin.Context, defectIds []int64) error
	GetById(c *gin.Context, defectId int64) (*defectApi.DefectData, error)
}
