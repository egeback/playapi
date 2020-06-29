package models

//Variant struct that contains information regarding differnt versions of the same media
type Variant struct {
	ID               string      `json:"-" groups:"api"`
	URL              string      `json:"url" groups:"api"`
	PlatformSpecific interface{} `json:"platform_specific" groups:"api"`
}
