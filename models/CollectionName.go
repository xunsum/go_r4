package models

type CollectionName struct {
	Colid string `json:"colid"`
	Name  string `json:"name"`
	Uid   string `json:"uid"`
}

func (CollectionName) TableName() string {
	return "collection_names"
}
