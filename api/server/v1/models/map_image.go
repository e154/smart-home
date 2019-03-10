package models

type MapImage struct {
	Id      int64  `json:"id"`
	ImageId int64  `json:"image_id"`
	Style   string `json:"style"`
}
