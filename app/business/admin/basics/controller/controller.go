package controller

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewCompanyInfo, wire.Struct(new(Basics), "*"))

type Basics struct {
	CompanyInfo *CompanyInfo
}
