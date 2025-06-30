package qcIndexServiceImpl

import (
	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/qcIndex/qcIndexApi"
	"nova-factory-server/app/business/qcIndex/qcIndexDao"
	"nova-factory-server/app/business/qcIndex/qcIndexService"

	"github.com/gin-gonic/gin"
)

type QcIndexServiceImpl struct {
	dao qcIndexDao.QcIndexDao
}

func NewQcIndexServiceImpl(dao qcIndexDao.QcIndexDao) qcIndexService.QcIndexService {
	return &QcIndexServiceImpl{
		dao: dao,
	}
}

// List 获取检测项列表
func (s *QcIndexServiceImpl) List(c *gin.Context, req *qcIndexApi.QcIndexListReq) (*qcIndexApi.QcIndexListRes, error) {
	total, data, err := s.dao.List(c, req)
	if err != nil {
		return nil, err
	}
	// 结构体从新赋值
	var list []*qcIndexApi.QcIndexData
	for i := range data {
		list = append(list, &qcIndexApi.QcIndexData{
			IndexId:      int64(data[i].ID),
			IndexCode:    data[i].IndexCode,
			IndexName:    data[i].IndexName,
			IndexType:    data[i].IndexType,
			QcTool:       data[i].QcTool,
			QcResultType: data[i].QcResultType,
			QcResultSpc:  data[i].QcResultSpc,
			Remark:       data[i].Remark,
			Attr1:        data[i].Attr1,
			Attr2:        data[i].Attr2,
			Attr3:        data[i].Attr3,
			Attr4:        data[i].Attr4,
			CreateBy:     data[i].CreateBy,
			CreateTime:   data[i].CreatedAt,
			UpdateBy:     data[i].UpdateBy,
			UpdateTime:   data[i].UpdatedAt,
		})
	}
	return &qcIndexApi.QcIndexListRes{
		Rows:  list,
		Total: total,
	}, nil
}

// Create 创建检测项
func (s *QcIndexServiceImpl) Create(c *gin.Context, req *qcIndexApi.QcIndexCreateReq) (*baize.EmptyResponse, error) {
	return s.dao.Create(c, req)
}

// Update 修改检测项
func (s *QcIndexServiceImpl) Update(c *gin.Context, req *qcIndexApi.QcIndexUpdateReq) (*baize.EmptyResponse, error) {
	return s.dao.Update(c, req)
}

// Delete 删除检测项
func (s *QcIndexServiceImpl) Delete(c *gin.Context, indexIds []int64) error {
	return s.dao.Delete(c, indexIds)
}

// GetById 根据ID获取检测项
func (s *QcIndexServiceImpl) GetById(c *gin.Context, indexId int64) (*qcIndexApi.QcIndexData, error) {
	index, err := s.dao.GetById(c, indexId)
	if err != nil {
		return nil, err
	}
	return &qcIndexApi.QcIndexData{
		IndexId:      int64(index.ID),
		IndexCode:    index.IndexCode,
		IndexName:    index.IndexName,
		IndexType:    index.IndexType,
		QcTool:       index.QcTool,
		QcResultType: index.QcResultType,
		QcResultSpc:  index.QcResultSpc,
		Remark:       index.Remark,
		Attr1:        index.Attr1,
		Attr2:        index.Attr2,
		Attr3:        index.Attr3,
		Attr4:        index.Attr4,
		CreateBy:     index.CreateBy,
		CreateTime:   index.CreatedAt,
		UpdateBy:     index.UpdateBy,
		UpdateTime:   index.UpdatedAt,
	}, nil
}
