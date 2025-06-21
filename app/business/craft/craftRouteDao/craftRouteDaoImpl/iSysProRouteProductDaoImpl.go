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

type ISysProRouteProductDaoImpl struct {
	db        *gorm.DB
	tableName string
}

func NewISysProRouteProductDaoImpl(db *gorm.DB) craftRouteDao.ISysProRouteProductDao {
	return &ISysProRouteProductDaoImpl{
		db:        db,
		tableName: "sys_pro_route_product",
	}
}

func (product *ISysProRouteProductDaoImpl) Add(c *gin.Context, info *craftRouteModels.SysProRouteSetProduct) (*craftRouteModels.SysProRouteProduct, error) {
	data := craftRouteModels.NewSysProRouteProduct(info)
	data.RecordID = snowflake.GenID()
	data.DeptID = baizeContext.GetDeptId(c)
	data.SetCreateBy(baizeContext.GetUserId(c))
	ret := product.db.Table(product.tableName).Create(data)
	return data, ret.Error
}

func (product *ISysProRouteProductDaoImpl) Update(c *gin.Context, info *craftRouteModels.SysProRouteSetProduct) (*craftRouteModels.SysProRouteProduct, error) {
	data := craftRouteModels.NewSysProRouteProduct(info)
	data.SetUpdateBy(baizeContext.GetUserId(c))
	ret := product.db.Table(product.tableName).Where("record_id = ?", data.RecordID).Updates(data)
	return data, ret.Error
}

func (product *ISysProRouteProductDaoImpl) Remove(c *gin.Context, ids []string) error {
	ret := product.db.Table(product.tableName).Where("record_id in (?)", ids).Update("state", commonStatus.DELETE)
	return ret.Error
}

func (product *ISysProRouteProductDaoImpl) List(c *gin.Context, req *craftRouteModels.SysProRouteProductReq) (*craftRouteModels.SysProRouteProductList, error) {
	db := product.db.Table(product.tableName).Table(product.tableName)

	size := 0
	if req.ItemCode != "" {
		db = db.Where("item_code = ?", req.ItemCode)
	}
	if req.ItemName != "" {
		db = db.Where("item_name LIKE ?", "%"+req.ItemName+"%")
	}
	if req == nil || req.Size <= 0 {
		size = 20
	} else {
		size = int(req.Size)
	}
	offset := 0
	if req == nil || req.Page <= 0 {
		req.Page = 1
	} else {
		offset = int((req.Page - 1) * req.Size)
	}
	db = db.Where("state", commonStatus.NORMAL)
	db = baizeContext.GetGormDataScope(c, db)
	var dto []*craftRouteModels.SysProRouteProduct

	var total int64
	ret := db.Count(&total)
	if ret.Error != nil {
		return &craftRouteModels.SysProRouteProductList{
			Rows:  make([]*craftRouteModels.SysProRouteProduct, 0),
			Total: 0,
		}, ret.Error
	}

	ret = db.Offset(offset).Order("create_time desc").Limit(size).Find(&dto)
	if ret.Error != nil {
		return &craftRouteModels.SysProRouteProductList{
			Rows:  make([]*craftRouteModels.SysProRouteProduct, 0),
			Total: 0,
		}, ret.Error
	}
	return &craftRouteModels.SysProRouteProductList{
		Rows:  dto,
		Total: total,
	}, nil
}

func (s *ISysProRouteProductDaoImpl) GetByRouteId(c *gin.Context, routeID int64) ([]*craftRouteModels.SysProRouteSetProduct, error) {
	var boms []*craftRouteModels.SysProRouteSetProduct
	ret := s.db.Table(s.tableName).Where("route_id = ?", routeID).Where("state = ?", commonStatus.NORMAL).Find(&boms)
	return boms, ret.Error
}
