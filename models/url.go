package models

type Url struct {
	Id int `json:"-"`
	Code  string `json:"code"`
	ToUrl string `json:"url"`
}

func (this Url) TableName() string {
	return "urls"
}
