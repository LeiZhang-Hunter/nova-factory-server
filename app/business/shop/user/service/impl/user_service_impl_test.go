package impl

import (
	"nova-factory-server/app/business/shop/user/models"
	"testing"

	"nova-factory-server/app/utils/bCryptPasswordEncoder"

	"github.com/gin-gonic/gin"
)

func TestShopUserServiceCreateFillDefault(t *testing.T) {
	dao := &mockShopUserDao{}
	service := &ShopUserServiceImpl{dao: dao}
	req := &models.UserUpsert{
		Username: "demo",
		UserType: 1,
		Password: "123456",
	}

	_, err := service.Create(&gin.Context{}, req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if req.UserID == "" {
		t.Fatal("expected auto generated user id")
	}
	if req.Status == nil || *req.Status {
		t.Fatal("expected default disabled status")
	}
	if !bCryptPasswordEncoder.CheckPasswordHash("123456", req.Password) {
		t.Fatal("expected password hashed")
	}
}

func TestShopUserServiceUpdateKeepPasswordWhenEmpty(t *testing.T) {
	dao := &mockShopUserDao{
		current: &models.User{
			ID:       1,
			Password: "hashed-password",
		},
	}
	service := &ShopUserServiceImpl{dao: dao}
	req := &models.UserUpsert{
		ID:       1,
		Username: "demo",
		UserType: 2,
	}

	_, err := service.Update(&gin.Context{}, req)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if req.Password != "hashed-password" {
		t.Fatalf("expected original password to be kept, got %s", req.Password)
	}
}

type mockShopUserDao struct {
	current *models.User
}

func (m *mockShopUserDao) Create(c *gin.Context, req *models.UserUpsert) (*models.User, error) {
	return &models.User{UserID: req.UserID, Status: req.Status}, nil
}

func (m *mockShopUserDao) Update(c *gin.Context, req *models.UserUpsert) (*models.User, error) {
	return &models.User{ID: req.ID, UserID: req.UserID, Status: req.Status}, nil
}

func (m *mockShopUserDao) DeleteByIDs(c *gin.Context, ids []int64) error {
	return nil
}

func (m *mockShopUserDao) GetByID(c *gin.Context, id int64) (*models.User, error) {
	return m.current, nil
}

func (m *mockShopUserDao) List(c *gin.Context, req *models.UserQuery) (*models.UserListData, error) {
	return nil, nil
}
