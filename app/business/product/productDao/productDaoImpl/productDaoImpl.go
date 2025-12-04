package productDaoImpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewSysProductLaboratoryDao)
