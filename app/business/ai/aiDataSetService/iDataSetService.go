package aiDataSetService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/aiDataSetModels"
)

type IDataSetService interface {
	CreateDataSet(c *gin.Context, request *aiDataSetModels.DataSetRequest) (*aiDataSetModels.SysDataset, error)
	UpdateDataSet(c *gin.Context, request *aiDataSetModels.UpdateDataSetRequest, id int64) (*aiDataSetModels.SysDataset, error)
	SelectDataSet(c *gin.Context, request *aiDataSetModels.DatasetListReq) (*aiDataSetModels.SysDatasetListData, error)
	DeleteDataSet(c *gin.Context, id int64) error
	SyncDataSet(c *gin.Context, id int64) error
}
