package models

type Danmaku struct {
	Content  string `json:"content"`
	Did      string `json:"did"`
	SentTime int    `json:"sent_time"` // 在视频中的哪个位置，秒为单位
	Uid      string `json:"uid"`
	Vid      string `json:"vid"`
}

func (Danmaku) TableName() string {
	return "danmakus"
}
