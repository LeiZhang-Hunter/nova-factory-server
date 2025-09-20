package aiDataSetDaoImpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewDataSetDaoImpl, NewIDataSetDocumentDaoImpl, NewIAiPredictionListDaoImpl, NewIAiPredictionExceptionDaoImpl)
