package aidatasetdao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
)

type IDataSetDao interface {
	GetByName(c *gin.Context, name string) (*aidatasetmodels.SysDataset, error)
	GetById(c *gin.Context, id int64) (*aidatasetmodels.SysDataset, error)
	Create(c *gin.Context, dataset *aidatasetmodels.DataSetCreateResponse) (*aidatasetmodels.SysDataset, error)
	UpdateData(c *gin.Context, id int64, dataset *aidatasetmodels.DataSetData) (*aidatasetmodels.SysDataset, error)
	Update(c *gin.Context, datasetId int64, request *aidatasetmodels.UpdateDataSetRequest) (*aidatasetmodels.SysDataset, error)
	SelectByList(c *gin.Context, request *aidatasetmodels.DatasetListReq) (*aidatasetmodels.SysDatasetListData, error)
	DeleteByIds(c *gin.Context, ids []int64) error
}
