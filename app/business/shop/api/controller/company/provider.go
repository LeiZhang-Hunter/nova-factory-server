package company

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewCompanyInfo, wire.Struct(new(Controller), "*"))

type Controller struct {
	CompanyInfo *CompanyInfo
}
