package kdniao

type config struct {
	SystemName  string `json:"systemName"`
	Credentials struct {
		EBusinessID string `json:"e_business_id"`
		AppKey      string `json:"app_key"`
	} `json:"credentials"`
}
