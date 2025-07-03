package qcIqcDao

import (
	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/qcIqc/qcIqcApi"
	"nova-factory-server/app/business/qcIqc/qcIqcModel"

	"github.com/gin-gonic/gin"
)

// IQcIqcDao 来料检验单数据访问接口
type IQcIqcDao interface {
	// List 查询来料检验单列表
	List(c *gin.Context, req *qcIqcApi.QcIqcQueryReq) (int64, []*qcIqcModel.QcIqc, error)

	// Create 新增来料检验单
	Create(c *gin.Context, req *qcIqcApi.QcIqcCreateReq) (*baize.EmptyResponse, error)

	// Update 修改来料检验单
	Update(c *gin.Context, req *qcIqcApi.QcIqcUpdateReq) (*baize.EmptyResponse, error)

	// Delete 批量删除来料检验单
	Delete(c *gin.Context, iqcIds []int64) error

	// GetById 根据ID查询来料检验单
	GetById(c *gin.Context, iqcId int64) (*qcIqcModel.QcIqc, error)
}
