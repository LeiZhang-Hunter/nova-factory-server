package defectServiceImpl

import (
	"nova-factory-server/app/baize"
	"nova-factory-server/app/business/defect/defectApi"
	"nova-factory-server/app/business/defect/defectDao"
	"nova-factory-server/app/business/defect/defectService"

	"github.com/gin-gonic/gin"
)

type DefectServiceImpl struct {
	dao defectDao.DefectDao
}

func NewDefectServiceImpl(dao defectDao.DefectDao) defectService.DefectService {
	return &DefectServiceImpl{
		dao: dao,
	}
}

// List 获取缺陷列表
func (s *DefectServiceImpl) List(c *gin.Context, req *defectApi.DefectListReq) (*defectApi.DefectListRes, error) {
	total, data, err := s.dao.List(c, req)
	if err != nil {
		return nil, err
	}
	// 结构体从新赋值
	var list []*defectApi.DefectData
	for i := range data {
		list = append(list, &defectApi.DefectData{
			Id:          int64(data[i].ID),
			DefectCode:  data[i].DefectCode,
			DefectName:  data[i].DefectName,
			IndexType:   data[i].IndexType,
			DefectLevel: data[i].DefectLevel,
			DeptId:      data[i].DeptId,
			State:       data[i].State,
			CreateBy:    data[i].CreateBy,
			UpdateBy:    data[i].UpdateBy,
			Remark:      data[i].Remark,
			Attr1:       data[i].Attr1,
			Attr2:       data[i].Attr2,
			Attr3:       data[i].Attr3,
			Attr4:       data[i].Attr4,
		})
	}
	return &defectApi.DefectListRes{
		Rows:  list,
		Total: total,
	}, nil
}

// Create 创建缺陷
func (s *DefectServiceImpl) Create(c *gin.Context, req *defectApi.DefectCreateReq) (*baize.EmptyResponse, error) {
	return s.dao.Create(c, req)
}

// Update 修改缺陷
func (s *DefectServiceImpl) Update(c *gin.Context, req *defectApi.DefectUpdateReq) (*baize.EmptyResponse, error) {
	return s.dao.Update(c, req)
}

// Delete 删除缺陷
func (s *DefectServiceImpl) Delete(c *gin.Context, defectIds []int64) error {
	return s.dao.Delete(c, defectIds)
}

// GetById 根据ID获取缺陷
func (s *DefectServiceImpl) GetById(c *gin.Context, defectId int64) (*defectApi.DefectData, error) {
	defect, err := s.dao.GetById(c, defectId)
	if err != nil {
		return nil, err
	}
	return &defectApi.DefectData{
		Id:          int64(defect.ID),
		DefectCode:  defect.DefectCode,
		DefectName:  defect.DefectName,
		DefectLevel: defect.DefectLevel,
		IndexType:   defect.IndexType,
		Attr1:       defect.Attr1,
		Attr2:       defect.Attr2,
		Attr3:       defect.Attr3,
		Attr4:       defect.Attr4,
		Remark:      defect.Remark,
	}, nil
}
