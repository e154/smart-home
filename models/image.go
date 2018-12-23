package models

import "time"

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
