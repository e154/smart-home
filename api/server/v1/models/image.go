package models

import "time"

type NewImage struct {
	Thumb    string `json:"thumb"`
	Image    string `json:"image"`
	MimeType string `json:"mime_type"`
	Title    string `json:"title"`
	Size     int64  `json:"size"`
	Name     string `json:"name"`
}

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

type Image struct {
	Id        int64     `json:"id"`
	Thumb     string    `json:"thumb"`
	Url       string    `json:"url"`
	Image     string    `json:"image"`
	MimeType  string    `json:"mime_type"`
	Title     string    `json:"title"`
	Size      int64     `json:"size"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type ImageListModel struct {
	Items []Image `json:"items"`
	Meta  struct {
		Limit        int `json:"limit"`
		Offset       int `json:"offset"`
		ObjectsCount int `json:"objects_count"`
	}
}