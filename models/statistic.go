package models

type Statistic struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Value       int32  `json:"value"`
	Diff        int32  `json:"diff"`
}
