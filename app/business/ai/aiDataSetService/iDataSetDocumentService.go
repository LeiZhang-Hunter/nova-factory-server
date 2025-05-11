package aiDataSetService

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/aiDataSetModels"
)

type IDataSetDocumentService interface {
	UploadFile(c *gin.Context, datasetId int64) ([]*aiDataSetModels.SysDatasetDocument, error)
	PutFile(c *gin.Context, documentId int64, request *aiDataSetModels.PutDocumentRequest) (*aiDataSetModels.SysDatasetDocument, error)
	DownloadFile(c *gin.Context, documentId int64) (*aiDataSetModels.SysDatasetDocument, error)
	ListDocument(c *gin.Context, databaseId int64, req *aiDataSetModels.ListDocumentRequest) (*aiDataSetModels.ListDocumentData, error)
	RemoveDocument(c *gin.Context, request *aiDataSetModels.DeleteDocumentRequest) error
	StartParse(c *gin.Context, request *aiDataSetModels.ParseDocumentApiRequest) error
	StopParse(c *gin.Context, request *aiDataSetModels.ParseDocumentApiRequest) error
}
