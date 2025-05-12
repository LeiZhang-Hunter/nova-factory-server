package aiDataSetDaoImpl

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/ai/aiDataSetDao"
	"nova-factory-server/app/business/ai/aiDataSetModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
)

type IDataSetDocumentDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewIDataSetDocumentDaoImpl(db *gorm.DB) aiDataSetDao.IDataSetDocumentDao {
	return &IDataSetDocumentDaoImpl{
		tableName: "sys_dataset_document",
		db:        db,
	}
}

func (i *IDataSetDocumentDaoImpl) Create(c *gin.Context, datasetId int64, response *aiDataSetModels.UploadDocumentResponse) ([]*aiDataSetModels.SysDatasetDocument, error) {
	var documents []*aiDataSetModels.SysDatasetDocument = make([]*aiDataSetModels.SysDatasetDocument, 0)
	for _, v := range response.Data {
		var document aiDataSetModels.SysDatasetDocument
		document.DocumentID = snowflake.GenID()
		document.DatasetID = datasetId
		document.DatasetChunkMethod = v.ChunkMethod
		document.DatasetCreatedBy = v.CreatedBy
		document.DatasetDocumentUUID = v.Id
		document.DatasetDatasetUUID = v.DatasetId
		document.DatasetLanguage = ""
		document.DatasetLocation = v.Location
		document.DatasetName = v.Name
		content, err := json.Marshal(v.ParserConfig)
		if err != nil {
			document.DatasetParserConfig = ""
		} else {
			document.DatasetParserConfig = string(content)
		}
		document.DatasetRun = v.Run
		document.DatasetSize = int64(v.Size)
		document.DatasetThumbnail = v.Thumbnail
		document.DatasetType = v.Type
		document.DeptID = baizeContext.GetDeptId(c)
		document.SetCreateBy(baizeContext.GetUserId(c))
		ret := i.db.Table(i.tableName).Create(&document)
		documents = append(documents, &document)
		if ret.Error != nil {
			return documents, ret.Error
		}
	}
	return documents, nil
}

func (i *IDataSetDocumentDaoImpl) GetById(c *gin.Context, documentId int64) (*aiDataSetModels.SysDatasetDocument, error) {
	var info *aiDataSetModels.SysDatasetDocument
	ret := i.db.Table(i.tableName).Where("document_id = ?", documentId).Where("state = ?", commonStatus.NORMAL).First(&info)
	if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return info, ret.Error
}

func (i *IDataSetDocumentDaoImpl) Update(c *gin.Context, id int64, request *aiDataSetModels.PutDocumentRequest) (*aiDataSetModels.SysDatasetDocument, error) {
	info, err := i.GetById(c, id)
	if err != nil {
		return nil, err
	}
	info.DatasetName = request.Name
	info.DatasetChunkMethod = request.ChunkMethod
	if request.ParserConfig != nil {
		config, err := json.Marshal(request.ParserConfig)
		if err != nil {
			return nil, err
		}
		info.DatasetParserConfig = string(config)
	}
	if len(request.MetaFields) != 0 {
		//config, err := json.Marshal(request.ParserConfig)
		//if err != nil {
		//	return nil, err
		//}
		//info.DatasetParserConfig = string(config)
	}

	ret := i.db.Table(i.tableName).Where("document_id = ?", id).Updates(info)
	return info, ret.Error
}

func (d *IDataSetDocumentDaoImpl) SelectByList(c *gin.Context, request *aiDataSetModels.ListDocumentRequest) (*aiDataSetModels.ListDocumentData, error) {
	//db := d.db.Table(d.tableName)
	//
	//if request != nil && request.Keywords != "" {
	//	db = db.Where("dataset_name LIKE ?", "%"+request.Keywords+"%")
	//}
	//size := 0
	//if request == nil || request.Size <= 0 {
	//	size = 20
	//} else {
	//	size = int(request.Size)
	//}
	//offset := 0
	//if request == nil || request.Page <= 0 {
	//	request.Page = 1
	//} else {
	//	offset = int((request.Page - 1) * request.Size)
	//}
	//db = db.Where("state", commonStatus.NORMAL)
	//db = baizeContext.GetGormDataScope(c, db)
	//var dto []*aiDataSetModels.SysDatasetDocument
	//
	//var total int64
	//ret := db.Count(&total)
	//if ret.Error != nil {
	//	return &aiDataSetModels.ListDocumentData{
	//		Rows:  make([]*aiDataSetModels.SysDatasetDocument, 0),
	//		Total: 0,
	//	}, ret.Error
	//}
	//
	//ret = db.Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	//if ret.Error != nil {
	//	return &aiDataSetModels.ListDocumentData{
	//		Rows:  make([]*aiDataSetModels.SysDatasetDocument, 0),
	//		Total: 0,
	//	}, ret.Error
	//}
	//return &aiDataSetModels.ListDocumentData{
	//	Rows:  dto,
	//	Total: total,
	//}, nil
	return nil, nil
}

func (d *IDataSetDocumentDaoImpl) GetByIds(c *gin.Context, documentId []string) ([]*aiDataSetModels.SysDatasetDocument, error) {
	var info []*aiDataSetModels.SysDatasetDocument
	ret := d.db.Table(d.tableName).Where("document_id in (?)", documentId).Where("state = ?", commonStatus.NORMAL).First(&info)
	if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return info, ret.Error
}

func (d *IDataSetDocumentDaoImpl) RemoveByIds(c *gin.Context, documentId []string) error {
	ret := d.db.Table(d.tableName).Where("document_id in (?)", documentId).Update("state", commonStatus.DELETE)
	if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
		return nil
	}
	return ret.Error
}

func (d *IDataSetDocumentDaoImpl) GetByDocumentUuids(c *gin.Context, documentUuids []string) ([]*aiDataSetModels.SysDatasetDocument, error) {
	var info []*aiDataSetModels.SysDatasetDocument
	ret := d.db.Table(d.tableName).Debug().Where("dataset_document_uuid in (?)", documentUuids).Where("state = ?", commonStatus.NORMAL).Find(&info)
	if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return info, ret.Error
}
