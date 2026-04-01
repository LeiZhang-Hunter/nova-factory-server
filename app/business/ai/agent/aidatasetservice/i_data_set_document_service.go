package aidatasetservice

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
)

type IDataSetDocumentService interface {
	UploadFile(c *gin.Context, datasetId int64) ([]*aidatasetmodels.SysDatasetDocument, error)
	PutFile(c *gin.Context, documentId int64, request *aidatasetmodels.PutDocumentRequest) (*aidatasetmodels.SysDatasetDocument, error)
	DownloadFile(c *gin.Context, documentId int64) (*aidatasetmodels.SysDatasetDocument, error)
	ListDocument(c *gin.Context, databaseId int64, req *aidatasetmodels.ListDocumentRequest) (*aidatasetmodels.ListDocumentData, error)
	RemoveDocument(c *gin.Context, request *aidatasetmodels.DeleteDocumentRequest) error
	StartParse(c *gin.Context, request *aidatasetmodels.ParseDocumentApiRequest) error
	StopParse(c *gin.Context, request *aidatasetmodels.ParseDocumentApiRequest) error
}
