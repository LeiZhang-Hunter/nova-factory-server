package aiDataSetDaoImpl

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"nova-factory-server/app/business/ai/aiDataSetDao"
	"nova-factory-server/app/business/ai/aiDataSetModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
	"time"
)

type DataSetDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewDataSetDaoImpl(db *gorm.DB) aiDataSetDao.IDataSetDao {
	return &DataSetDaoImpl{
		db:        db,
		tableName: "sys_dataset",
	}
}

func (d *DataSetDaoImpl) GetById(c *gin.Context, id int64) (*aiDataSetModels.SysDataset, error) {
	var info *aiDataSetModels.SysDataset
	ret := d.db.Table(d.tableName).Where("dataset_id = ?", id).Where("state = ?", commonStatus.NORMAL).First(&info)
	if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return info, ret.Error
}

func (d *DataSetDaoImpl) GetByName(c *gin.Context, name string) (*aiDataSetModels.SysDataset, error) {
	var info *aiDataSetModels.SysDataset
	ret := d.db.Table(d.tableName).Where("dataset_name = ?", name).Where("state = ?", commonStatus.NORMAL).First(&info)
	if errors.Is(ret.Error, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return info, ret.Error
}

func (d *DataSetDaoImpl) Create(c *gin.Context, dataset *aiDataSetModels.DataSetCreateResponse) (*aiDataSetModels.SysDataset, error) {
	if dataset.Data.Name == "" {
		return &aiDataSetModels.SysDataset{}, errors.New("数据集名字不能为空")
	}

	createDate, err := time.Parse(time.RFC1123, dataset.Data.CreateDate)
	if err != nil {
		zap.L().Debug("time.Parse(time.RFC1123, dataset.Data.CreateDate) error", zap.Error(err))
		createDate = time.Now()
	}
	updateDate, err := time.Parse(time.RFC1123, dataset.Data.UpdateDate)
	if err != nil {
		zap.L().Debug("time.Parse(time.RFC1123, dataset.Data.UpdateDate) error", zap.Error(err))
		updateDate = time.Now()
	}

	var parseConfig []byte
	if dataset.Data.ParserConfig != nil {
		parseConfig, err = json.Marshal(dataset.Data.ParserConfig)
		if err != nil {
			return &aiDataSetModels.SysDataset{}, errors.New(err.Error())
		}
	}

	var data aiDataSetModels.SysDataset = aiDataSetModels.SysDataset{
		DatasetID:                     snowflake.GenID(),
		DatasetAvatar:                 dataset.Data.Avatar,
		DatasetChunkMethod:            dataset.Data.ChunkMethod,
		DatasetCreateDate:             createDate,
		DatasetCreateTime:             dataset.Data.CreateTime,
		DatasetCreatedBy:              dataset.Data.CreatedBy,
		DatasetDescription:            dataset.Data.Description,
		DatasetDocumentCount:          int64(dataset.Data.DocumentCount),
		DatasetEmbeddingModel:         dataset.Data.EmbeddingModel,
		DatasetUUID:                   dataset.Data.Id,
		DatasetLanguage:               dataset.Data.Language,
		DatasetName:                   dataset.Data.Name,
		DatasetPagerank:               dataset.Data.Pagerank,
		DatasetParserConfig:           string(parseConfig),
		DatasetPermission:             dataset.Data.Permission,
		DatasetSimilarityThreshold:    dataset.Data.SimilarityThreshold,
		DatasetTokenNum:               int64(dataset.Data.TokenNum),
		DatasetUpdateDate:             updateDate,
		DatasetUpdateTime:             dataset.Data.UpdateTime,
		DatasetVectorSimilarityWeight: dataset.Data.VectorSimilarityWeight,
		DeptID:                        baizeContext.GetDeptId(c),
	}
	data.SetCreateBy(baizeContext.GetUserId(c))
	ret := d.db.Table(d.tableName).Create(&data)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return &data, nil
}

func (d *DataSetDaoImpl) Update(c *gin.Context, datasetId int64, request *aiDataSetModels.UpdateDataSetRequest) (*aiDataSetModels.SysDataset, error) {
	if request.Name == "" {
		return &aiDataSetModels.SysDataset{}, errors.New("数据集名字不能为空")
	}

	var info *aiDataSetModels.SysDataset
	ret := d.db.Table(d.tableName).First(&info)
	if ret.Error != nil {
		return nil, ret.Error
	}
	if info == nil {
		return nil, nil
	}
	var parseConfig []byte
	var err error
	if request.ParserConfig != nil {
		parseConfig, err = json.Marshal(request.ParserConfig)
		if err != nil {
			return &aiDataSetModels.SysDataset{}, errors.New(err.Error())
		}
	}
	info.DatasetName = request.Name
	info.DatasetEmbeddingModel = request.EmbeddingModel
	info.DatasetChunkMethod = request.ChunkMethod
	info.DatasetParserConfig = string(parseConfig)
	info.SetUpdateBy(baizeContext.GetUserId(c))
	ret = d.db.Table(d.tableName).Where("dataset_id = ?", datasetId).Updates(info)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return info, nil
}

func (d *DataSetDaoImpl) SelectByList(c *gin.Context, request *aiDataSetModels.DatasetListReq) (*aiDataSetModels.SysDatasetListData, error) {
	db := d.db.Table(d.tableName)

	if request != nil && request.Name != "" {
		db = db.Where("dataset_name LIKE ?", "%"+request.Name+"%")
	}
	size := 0
	if request == nil || request.Size <= 0 {
		size = 20
	} else {
		size = int(request.Size)
	}
	offset := 0
	if request == nil || request.Page <= 0 {
		request.Page = 1
	} else {
		offset = int((request.Page - 1) * request.Size)
	}
	db = db.Where("state", commonStatus.NORMAL)
	db = baizeContext.GetGormDataScope(c, db)
	var dto []*aiDataSetModels.SysDataset

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &aiDataSetModels.SysDatasetListData{
			Rows:  make([]*aiDataSetModels.SysDataset, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &aiDataSetModels.SysDatasetListData{
			Rows:  make([]*aiDataSetModels.SysDataset, 0),
			Total: 0,
		}, ret.Error
	}
	return &aiDataSetModels.SysDatasetListData{
		Rows:  dto,
		Total: total,
	}, nil
}

func (d *DataSetDaoImpl) DeleteByIds(c *gin.Context, ids []int64) error {
	if len(ids) == 0 {
		return errors.New("请选中数据")
	}
	ret := d.db.Table(d.tableName).Where("dataset_id in (?)", ids).Delete(&aiDataSetModels.SysDataset{})
	return ret.Error
}

func (d *DataSetDaoImpl) UpdateData(c *gin.Context, id int64, dataset *aiDataSetModels.DataSetData) (*aiDataSetModels.SysDataset, error) {
	if dataset.Name == "" {
		return &aiDataSetModels.SysDataset{}, errors.New("数据集名字不能为空")
	}

	createDate, err := time.Parse(time.RFC1123, dataset.CreateDate)
	if err != nil {
		zap.L().Debug("time.Parse(time.RFC1123, dataset.CreateDate) error", zap.Error(err))
		createDate = time.Now()
	}
	updateDate, err := time.Parse(time.RFC1123, dataset.UpdateDate)
	if err != nil {
		zap.L().Debug("time.Parse(time.RFC1123, dataset.UpdateDate) error", zap.Error(err))
		updateDate = time.Now()
	}

	var data aiDataSetModels.SysDataset = aiDataSetModels.SysDataset{
		DatasetID:                     id,
		DatasetAvatar:                 dataset.Avatar,
		DatasetChunkMethod:            dataset.ChunkMethod,
		DatasetCreateDate:             createDate,
		DatasetCreateTime:             dataset.CreateTime,
		DatasetCreatedBy:              dataset.CreatedBy,
		DatasetDescription:            dataset.Description,
		DatasetDocumentCount:          int64(dataset.DocumentCount),
		DatasetEmbeddingModel:         dataset.EmbeddingModel,
		DatasetUUID:                   dataset.Id,
		DatasetLanguage:               dataset.Language,
		DatasetName:                   dataset.Name,
		DatasetPagerank:               dataset.Pagerank,
		DatasetParserConfig:           "",
		DatasetPermission:             dataset.Permission,
		DatasetSimilarityThreshold:    dataset.SimilarityThreshold,
		DatasetTokenNum:               int64(dataset.TokenNum),
		DatasetUpdateDate:             updateDate,
		DatasetUpdateTime:             dataset.UpdateTime,
		DatasetVectorSimilarityWeight: dataset.VectorSimilarityWeight,
		DeptID:                        baizeContext.GetDeptId(c),
	}
	data.SetUpdateBy(baizeContext.GetUserId(c))
	ret := d.db.Table(d.tableName).Where("dataset_id = ?", id).Updates(&data)
	if ret.Error != nil {
		return nil, ret.Error
	}
	return &data, nil
}
