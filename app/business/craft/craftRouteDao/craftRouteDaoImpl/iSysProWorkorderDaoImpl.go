package craftRouteDaoImpl

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nova-factory-server/app/business/craft/craftRouteDao"
	"nova-factory-server/app/business/craft/craftRouteModels"
	"nova-factory-server/app/constant/commonStatus"
	"nova-factory-server/app/utils/baizeContext"
	"nova-factory-server/app/utils/snowflake"
)

type WorkOrderDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewWorkOrderDaoImpl(db *gorm.DB) craftRouteDao.ISysProWorkorderDao {
	return &WorkOrderDaoImpl{
		tableName: "sys_pro_workorder",
		db:        db,
	}
}

func (w *WorkOrderDaoImpl) Add(c *gin.Context, info *craftRouteModels.SysSetProWorkorder) (*craftRouteModels.SysProWorkorder, error) {
	data := craftRouteModels.OfSysSetProWorkorder(info)
	data.WorkorderID = snowflake.GenID()
	data.DeptID = baizeContext.GetDeptId(c)
	data.SetCreateBy(baizeContext.GetUserId(c))
	ret := w.db.Table(w.tableName).Create(data)
	return data, ret.Error
}

func (w *WorkOrderDaoImpl) Update(c *gin.Context, info *craftRouteModels.SysSetProWorkorder) (*craftRouteModels.SysProWorkorder, error) {
	data := craftRouteModels.OfSysSetProWorkorder(info)
	data.SetUpdateBy(baizeContext.GetUserId(c))
	ret := w.db.Table(w.tableName).Where("workorder_id = ?", data.WorkorderID).Updates(data)
	return data, ret.Error
}

func (w *WorkOrderDaoImpl) Remove(c *gin.Context, ids []string) error {
	ret := w.db.Table(w.tableName).Where("workorder_id in (?)", ids).Update("state", commonStatus.DELETE)
	return ret.Error
}

func (w *WorkOrderDaoImpl) List(c *gin.Context, req *craftRouteModels.SysProWorkorderReq) (*craftRouteModels.SysProWorkorderList, error) {
	db := w.db.Table(w.tableName)

	size := 0
	if req.WorkorderCode != "" {
		db = db.Where("workorder_code = ?", req.WorkorderCode)
	}
	if req.WorkorderName != "" {
		db = db.Where("workorder_name LIKE ?", "%"+req.WorkorderName+"%")
	}
	if req.SourceCode != "" {
		db = db.Where("source_code LIKE ?", req.SourceCode)
	}
	if req.ProductCode != "" {
		db = db.Where("product_code = ?", req.ProductCode)
	}
	if req.ProductName != "" {
		db = db.Where("product_name LIKE ?", "%"+req.ProductName+"%")
	}
	if req.ClientCode != "" {
		db = db.Where("client_code = ?", req.ClientCode)
	}
	if req.ClientName != "" {
		db = db.Where("client_name = ?", req.ClientName)
	}
	if req.WorkorderType != "" {
		db = db.Where("workorder_type = ?", req.WorkorderType)
	}
	if req.RequestDate != "" {
		db = db.Where("request_date = ", req.RequestDate)
	}

	offset := 0
	if req == nil || req.Page <= 0 {
		req.Page = 1
	} else {
		offset = int((req.Page - 1) * req.Size)
	}

	db = db.Where("state", commonStatus.NORMAL)
	db = baizeContext.GetGormDataScope(c, db)
	var dto []*craftRouteModels.SysProWorkorder

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &craftRouteModels.SysProWorkorderList{
			Rows:  make([]*craftRouteModels.SysProWorkorder, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &craftRouteModels.SysProWorkorderList{
			Rows:  make([]*craftRouteModels.SysProWorkorder, 0),
			Total: 0,
		}, ret.Error
	}
	return &craftRouteModels.SysProWorkorderList{
		Rows:  dto,
		Total: total,
	}, nil
}
