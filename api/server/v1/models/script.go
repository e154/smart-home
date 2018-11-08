package models

import (
	"time"
)

type ScriptLang string

const (
	ScriptLangLua        = ScriptLang("lua")
	ScriptLangCoffee     = ScriptLang("coffeescript")
	ScriptLangJavascript = ScriptLang("javascript")
)

type NewScript struct {
	Lang        ScriptLang `json:"lang"`
	Name        string     `json:"name"`
	Source      string     `json:"source"`
	Description string     `json:"description"`
}

type UpdateScript struct {
	Lang        ScriptLang `json:"lang"`
	Name        string     `json:"name"`
	Source      string     `json:"source"`
	Description string     `json:"description"`
}

type Script struct {
	Id          int64      `json:"id"`
	Lang        ScriptLang `json:"lang"`
	Name        string     `json:"name"`
	Source      string     `json:"source"`
	Description string     `json:"description"`
	Compiled    string     `json:"compiled"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Scripts []*Script

type ResponseScript struct {
	Code ResponseType `json:"code"`
	Data struct {
		Script *Script `json:"script"`
	} `json:"data"`
}

type ResponseScriptList struct {
	Code ResponseType `json:"code"`
	Data struct {
		Items  []*Script `json:"items"`
		Limit  int64     `json:"limit"`
		Offset int64     `json:"offset"`
		Total  int64     `json:"total"`
	} `json:"data"`
}
