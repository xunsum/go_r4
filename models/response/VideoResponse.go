package response

import "go_round4/models"

type VideoListResponse struct {
	Data   []models.Video `json:"data"`
	Error  string         `json:"error"`
	Msg    string         `json:"msg"`
	Status int64          `json:"status"`
}

type VideoResponse struct {
	Data   models.Video `json:"data"`
	Error  string       `json:"error"`
	Msg    string       `json:"msg"`
	Status int64        `json:"status"`
}
