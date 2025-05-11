package aiDataSetDao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/aiDataSetModels"
)

type IDataSetDocumentDao interface {
	Create(c *gin.Context, id int64, response *aiDataSetModels.UploadDocumentResponse) ([]*aiDataSetModels.SysDatasetDocument, error)
	GetById(c *gin.Context, documentId int64) (*aiDataSetModels.SysDatasetDocument, error)
	Update(c *gin.Context, id int64, response *aiDataSetModels.PutDocumentRequest) (*aiDataSetModels.SysDatasetDocument, error)
	SelectByList(c *gin.Context, req *aiDataSetModels.ListDocumentRequest) (*aiDataSetModels.ListDocumentData, error)
	GetByIds(c *gin.Context, documentId []string) ([]*aiDataSetModels.SysDatasetDocument, error)
	RemoveByIds(c *gin.Context, documentId []string) error
	GetByDocumentUuids(c *gin.Context, documentUuids []string) ([]*aiDataSetModels.SysDatasetDocument, error)
}
