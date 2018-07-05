package models

type Url struct {
	Id int
	Code  string
	ToUrl string
}

func (this Url) TableName() string {
	return "urls"
}
