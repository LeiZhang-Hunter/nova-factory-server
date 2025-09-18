package aiDataSetServiceImpl

import "github.com/google/wire"

var ProviderSet = wire.NewSet(NewIDataSetServiceImpl, NewHttpClient, NewIDataSetDocumentServiceImpl,
	NewIChunkServiceImpl, NewIAssistantServiceImpl, NewIChartServiceImpl, NewIAiPredictionServiceImpl)
