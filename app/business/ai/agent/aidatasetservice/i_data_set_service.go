package aidatasetservice

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
)

type IDataSetService interface {
	CreateDataSet(c *gin.Context, request *aidatasetmodels.DataSetRequest) (*aidatasetmodels.SysDataset, error)
	UpdateDataSet(c *gin.Context, request *aidatasetmodels.UpdateDataSetRequest, id int64) (*aidatasetmodels.SysDataset, error)
	SelectDataSet(c *gin.Context, request *aidatasetmodels.DatasetListReq) (*aidatasetmodels.SysDatasetListData, error)
	GetInfoById(c *gin.Context, id int64) (*aidatasetmodels.SysDataset, error)
	DeleteDataSet(c *gin.Context, id int64) error
	SyncDataSet(c *gin.Context, id int64) error
}
