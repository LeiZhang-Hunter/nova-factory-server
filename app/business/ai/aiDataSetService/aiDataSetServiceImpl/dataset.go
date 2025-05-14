package aiDataSetServiceImpl

type RagFlowConfig struct {
	ApiKey   string `mapstructure:"api_key"`
	Host     string `mapstructure:"host"`
	ImageUrl string `mapstructure:"image_url"`
}
