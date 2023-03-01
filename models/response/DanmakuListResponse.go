package response

import "go_round4/models"

type DanmakuListResponse struct {
	Data   []models.Danmaku `json:"data"`
	Error  string           `json:"error"`
	Msg    string           `json:"msg"`
	Status int64            `json:"status"`
}
