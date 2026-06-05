package models

type CompanyInfoVo struct {
	CompanyName       string `json:"companyName" db:"company_name"`
	CompanyDetail     string `json:"companyDetail" db:"company_detail"`
	ContactPhone      string `json:"contactPhone" db:"contact_phone"`
	ContactPerson     string `json:"contactPerson" db:"contact_person"`
	Email             string `json:"email" db:"email"`
	Address           string `json:"address" db:"address"`
	LogoUrl           string `json:"logoUrl" db:"logo_url"`
	BusinessLicenseNo string `json:"businessLicenseNo" db:"business_license_no"`
	Remark            string `json:"remark" db:"remark"`
}
