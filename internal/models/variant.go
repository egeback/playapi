package models

//Variant is
type Variant struct {
	ID               string      `json:"-" groups:"api`
	URL              string      `json:"url" groups:"api`
	PlatformSpecific interface{} `json:"platform_specific" groups:"api"`
}
