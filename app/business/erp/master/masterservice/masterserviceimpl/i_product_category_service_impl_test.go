package masterserviceimpl

import (
	"strings"
	"testing"

	"nova-factory-server/app/business/erp/erpbiz"
	"nova-factory-server/app/business/erp/master/mastermodels"

	"github.com/gin-gonic/gin"
)

func TestProductCategoryServiceImplCreateGeneratesCodeWhenEmpty(t *testing.T) {
	dao := &productCategoryServiceTestDAO{}
	service := NewProductCategoryService(dao)

	req := &mastermodels.ProductCategoryUpsert{
		Name: "分类A",
		Code: "   ",
	}

	data, err := service.Create(nil, req)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if dao.createReq == nil {
		t.Fatal("expected dao.Create to be called")
	}
	if dao.createReq.Code == "" {
		t.Fatal("expected code to be auto generated")
	}
	if !strings.HasPrefix(dao.createReq.Code, "PCAT") {
		t.Fatalf("expected generated code to start with PCAT, got %q", dao.createReq.Code)
	}
	if data == nil || data.Code == "" {
		t.Fatal("expected created data to include generated code")
	}
}

func TestProductCategoryServiceImplUpdateGeneratesCodeWhenEmpty(t *testing.T) {
	dao := &productCategoryServiceTestDAO{}
	service := NewProductCategoryService(dao)

	req := &mastermodels.ProductCategoryUpsert{
		ID:   1,
		Name: "分类B",
		Code: " ",
	}

	data, err := service.Update(nil, req)
	if err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if dao.updateReq == nil {
		t.Fatal("expected dao.Update to be called")
	}
	if dao.updateReq.Code == "" {
		t.Fatal("expected code to be auto generated")
	}
	if !strings.HasPrefix(dao.updateReq.Code, "PCAT") {
		t.Fatalf("expected generated code to start with PCAT, got %q", dao.updateReq.Code)
	}
	if data == nil || data.Code == "" {
		t.Fatal("expected updated data to include generated code")
	}
}

type productCategoryServiceTestDAO struct {
	createReq *mastermodels.ProductCategoryUpsert
	updateReq *mastermodels.ProductCategoryUpsert
}

func (d *productCategoryServiceTestDAO) Create(c *gin.Context, req *mastermodels.ProductCategoryUpsert) (*mastermodels.ProductCategory, error) {
	d.createReq = req
	return &mastermodels.ProductCategory{
		ID:   1,
		Name: req.Name,
		Code: req.Code,
	}, nil
}

func (d *productCategoryServiceTestDAO) Update(c *gin.Context, req *mastermodels.ProductCategoryUpsert) (*mastermodels.ProductCategory, error) {
	d.updateReq = req
	return &mastermodels.ProductCategory{
		ID:   req.ID,
		Name: req.Name,
		Code: req.Code,
	}, nil
}

func (d *productCategoryServiceTestDAO) DeleteByIDs(c *gin.Context, ids []int64) error {
	return nil
}

func (d *productCategoryServiceTestDAO) GetByID(c *gin.Context, id int64) (*mastermodels.ProductCategory, error) {
	return nil, nil
}

func (d *productCategoryServiceTestDAO) GetByColumn(c *gin.Context, column string, value any) (*mastermodels.ProductCategory, error) {
	return nil, nil
}

func (d *productCategoryServiceTestDAO) ListPage(c *gin.Context, req *mastermodels.ProductCategoryQuery) (*erpbiz.PageResult[mastermodels.ProductCategory], error) {
	return &erpbiz.PageResult[mastermodels.ProductCategory]{}, nil
}

func (d *productCategoryServiceTestDAO) List(c *gin.Context, req *mastermodels.ProductCategoryQuery) (*mastermodels.ProductCategoryListData, error) {
	return &mastermodels.ProductCategoryListData{}, nil
}
