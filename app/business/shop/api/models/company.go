package models

// CompanyInfo stores the admin-managed company information used by the mini app.
type CompanyInfo struct {
	CompanyName   string `json:"companyName" gorm:"column:company_name"`
	CompanyDetail string `json:"companyDetail" gorm:"column:company_detail"`
	ContactPhone  string `json:"contactPhone" gorm:"column:contact_phone"`
	Email         string `json:"email" gorm:"column:email"`
	Address       string `json:"address" gorm:"column:address"`
	LogoUrl       string `json:"logoUrl" gorm:"column:logo_url"`
}

// CompanyInfoResp is the public company information returned to the mini app.
type CompanyInfoResp struct {
	CompanyName   string `json:"companyName"`
	CompanyDetail string `json:"companyDetail"`
	ContactPhone  string `json:"contactPhone"`
	Email         string `json:"email"`
	Address       string `json:"address"`
	LogoUrl       string `json:"logoUrl"`
}
