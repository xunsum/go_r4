package models

type SearchHistory struct {
	Description string `json:"description"`
	Length      string `json:"length"`
	LikeCounts  string `json:"like_counts"`
	Name        string `json:"name"`
	SearchTime  int64  `json:"search_time"`
	Title       string `json:"title"`
	Type        int    `json:"type"`
	Uid         string `json:"uid"`
	UploadTime  string `json:"uploadTime"`
	ViewCounts  string `json:"view_counts"`
}

func (SearchHistory) TableName() string {
	return "search_histories"
}
