package models

type Comment struct {
	Comid      string `json:"comid"`
	Content    string `json:"content"`
	Uid        string `json:"uid"`
	Vid        string `json:"vid"`
	Visibility string `json:"visibility"`
	ReplyTo    string `json:"reply_to"`
}

func (Comment) TableName() string {
	return "comments"
}
