package models

// swagger:model
type Image struct {
	Id    int64  `json:"id"`
	Thumb string `json:"thumb,omitempty"`
	Url   string `json:"url"`
}
