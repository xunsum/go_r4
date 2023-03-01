package models

type Collection struct {
	Colid string `json:"colid"`
	Uid   string `json:"uid"`
	Vid   string `json:"vid"`
}

type FullCollection struct {
	Name  string
	Colid string
	Uid   string
	Vid   string
}

func (Collection) TableName() string {
	return "collections"
}
