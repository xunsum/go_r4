package models

type Video struct {
	Description string `json:"description"`
	Length      string `json:"length"` // 长度，秒为单位
	LikeCount   int    `json:"like_count"`
	Title       string `json:"title"`
	Type        string `json:"type"`        // 鬼畜 - 0 动漫 - 1 生活 - 2 美食 - 3
	Uid         string `json:"uid"`         // 上传者id
	UploadTime  string `json:"upload_time"` // unix time
	Vid         string `json:"vid"`
	Visibility  int    `json:"visibility"`
	Views       int    `json:"views"`
}

func (Video) TableName() string {
	return "videos"
}
