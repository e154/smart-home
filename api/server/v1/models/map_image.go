package models

type MapImage struct {
	Id      int64  `json:"id,omitempty"`
	Image   *Image `json:"image"`
	ImageId int64  `json:"image_id"`
	Style   string `json:"style"`
}
