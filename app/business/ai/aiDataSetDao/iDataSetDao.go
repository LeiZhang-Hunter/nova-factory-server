package aiDataSetDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/aiDataSetModels"
)

type IDataSetDao interface {
	GetByName(c *gin.Context, name string) (*aiDataSetModels.SysDataset, error)
	GetById(c *gin.Context, id int64) (*aiDataSetModels.SysDataset, error)
	Create(c *gin.Context, dataset *aiDataSetModels.DataSetCreateResponse) (*aiDataSetModels.SysDataset, error)
	UpdateData(c *gin.Context, id int64, dataset *aiDataSetModels.DataSetData) (*aiDataSetModels.SysDataset, error)
	Update(c *gin.Context, datasetId int64, request *aiDataSetModels.UpdateDataSetRequest) (*aiDataSetModels.SysDataset, error)
	SelectByList(c *gin.Context, request *aiDataSetModels.DatasetListReq) (*aiDataSetModels.SysDatasetListData, error)
	DeleteByIds(c *gin.Context, ids []int64) error
}
