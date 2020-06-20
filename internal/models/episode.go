package models

// Episode is
type Episode struct {
	ID              string    `json:"id" groups:"api"`
	Name            string    `json:"name" groups:"api"`
	SvtID           string    `json:"svtId" groups:"api"`
	VideoSvtID      string    `json:"videoSvtId" groups:"api"`
	Slug            string    `json:"slug" groups:"api"`
	LongDescription string    `json:"longDescription" groups:"api"`
	ImageURL        string    `json:"imageUrl" groups:"api"`
	URL             string    `json:"url" groups:"api"`
	Duration        float64   `json:"duration" groups:"api`
	Number          string    `json:"number" groups:"api"`
	ValidFrom       string    `json:"validFrom" groups:"api"`
	ValidTo         string    `json:"validTo" groups:"api"`
	Varants         []Variant `json:"variants" groups:"api"`
	//Data interface
}
