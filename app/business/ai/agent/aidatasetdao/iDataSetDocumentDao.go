package aidatasetdao

import (
	"github.com/gin-gonic/gin"
	"nova-factory-server/app/business/ai/agent/aidatasetmodels"
)

type IDataSetDocumentDao interface {
	Create(c *gin.Context, id int64, response *aidatasetmodels.UploadDocumentResponse) ([]*aidatasetmodels.SysDatasetDocument, error)
	GetById(c *gin.Context, documentId int64) (*aidatasetmodels.SysDatasetDocument, error)
	Update(c *gin.Context, id int64, response *aidatasetmodels.PutDocumentRequest) (*aidatasetmodels.SysDatasetDocument, error)
	SelectByList(c *gin.Context, req *aidatasetmodels.ListDocumentRequest) (*aidatasetmodels.ListDocumentData, error)
	GetByIds(c *gin.Context, documentId []string) ([]*aidatasetmodels.SysDatasetDocument, error)
	RemoveByIds(c *gin.Context, documentId []string) error
	GetByDocumentUuids(c *gin.Context, documentUuids []string) ([]*aidatasetmodels.SysDatasetDocument, error)
}
