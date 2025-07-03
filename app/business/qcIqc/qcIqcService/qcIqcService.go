package qcIqcService

import (
	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/qcIqc/qcIqcApi"
	"nova-factory-server/app/business/qcIqc/qcIqcModel"

	"github.com/gin-gonic/gin"
)

// IQcIqcService 来料检验单服务接口
type IQcIqcService interface {
	// List 获取来料检验单列表
	List(c *gin.Context, req *qcIqcApi.QcIqcQueryReq) (*qcIqcApi.QcIqcListRes, error)

	// Create 创建来料检验单
	Create(c *gin.Context, req *qcIqcApi.QcIqcCreateReq) (*baize.EmptyResponse, error)

	// Update 更新来料检验单
	Update(c *gin.Context, req *qcIqcApi.QcIqcUpdateReq) (*baize.EmptyResponse, error)

	// Delete 批量删除来料检验单
	Delete(c *gin.Context, iqcIds []int64) error

	// GetById 根据ID获取来料检验单
	GetById(c *gin.Context, iqcId int64) (*qcIqcModel.QcIqc, error)
}
