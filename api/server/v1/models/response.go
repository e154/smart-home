package models

type ResponseSuccess struct {
	Code string   `json:"code"`
	Data struct{} `json:"data"`
}

type NewObjectSuccess struct {
	Code ResponseType `json:"code"`
	Data struct {
		Id int64 `json:"id"`
	} `json:"data"`
}