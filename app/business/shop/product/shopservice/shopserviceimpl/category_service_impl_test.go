package shopserviceimpl

import (
	"testing"

	"nova-factory-server/app/business/shop/product/shopmodels"

	"github.com/gin-gonic/gin"
)

func TestShopCategoryServiceImplCreateGeneratesCategoryCodeWhenEmpty(t *testing.T) {
	dao := &categoryServiceTestDAO{}
	service := &ShopCategoryServiceImpl{dao: dao}

	req := &shopmodels.CategoryUpsert{
		CategoryName: "分类A",
		CategoryCode: "   ",
	}

	_, err := service.Create(nil, req)
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if dao.createReq == nil {
		t.Fatal("expected dao.Create to be called")
	}
	if dao.createReq.CategoryCode == "" {
		t.Fatal("expected category code to be auto generated")
	}
	if len(dao.createReq.CategoryCode) <= len("CAT") || dao.createReq.CategoryCode[:3] != "CAT" {
		t.Fatalf("expected generated category code to start with CAT, got %q", dao.createReq.CategoryCode)
	}
}

type categoryServiceTestDAO struct {
	createReq *shopmodels.CategoryUpsert
}

func (d *categoryServiceTestDAO) Create(c *gin.Context, req *shopmodels.CategoryUpsert) (*shopmodels.Category, error) {
	d.createReq = req
	return &shopmodels.Category{
		ID:           1,
		CategoryName: req.CategoryName,
		CategoryCode: req.CategoryCode,
	}, nil
}

func (d *categoryServiceTestDAO) Update(c *gin.Context, req *shopmodels.CategoryUpsert) (*shopmodels.Category, error) {
	return nil, nil
}

func (d *categoryServiceTestDAO) DeleteByIDs(c *gin.Context, ids []int64) error {
	return nil
}

func (d *categoryServiceTestDAO) GetByID(c *gin.Context, id int64) (*shopmodels.Category, error) {
	return nil, nil
}

func (d *categoryServiceTestDAO) All(c *gin.Context) ([]*shopmodels.Category, error) {
	return nil, nil
}

func (d *categoryServiceTestDAO) ListByIDs(c *gin.Context, ids []int64) ([]*shopmodels.Category, error) {
	return nil, nil
}

func (d *categoryServiceTestDAO) List(c *gin.Context, req *shopmodels.CategoryQuery) (*shopmodels.CategoryListData, error) {
	return nil, nil
}
