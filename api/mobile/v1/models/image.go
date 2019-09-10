package models

// swagger:model
type NewImage struct {
	Thumb    string `json:"thumb"`
	Image    string `json:"image"`
	MimeType string `json:"mime_type"`
	Title    string `json:"title"`
	Size     int64  `json:"size"`
	Name     string `json:"name"`
}

// swagger:model
type UpdateImage struct {
	Id       int64  `json:"id"`
	Thumb    string `json:"thumb"`
	Url      string `json:"url"`
	Image    string `json:"image"`
	MimeType string `json:"mime_type"`
	Title    string `json:"title"`
	Size     int64  `json:"size"`
	Name     string `json:"name"`
}

// swagger:model
type Image struct {
	Id        int64     `json:"id"`
	Thumb     string    `json:"thumb"`
	Url       string    `json:"url"`
	Image     string    `json:"image"`
	MimeType  string    `json:"mime_type"`
	Title     string    `json:"title"`
	Size      int64     `json:"size"`
	Name      string    `json:"name"`
}