package models

type Like struct {
	Uid string `json:"uid"`
	Vid string `json:"vid"`
}

func (Like) TableName() string {
	return "likes"
}
