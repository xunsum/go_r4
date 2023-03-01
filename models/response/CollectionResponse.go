package response

import "go_round4/models"

type CollectionListResponse struct {
	Data   []models.FullCollection `json:"data"`
	Error  string                  `json:"error"`
	Msg    string                  `json:"msg"`
	Status int64                   `json:"status"`
}

type CollectionResponse struct {
	Data   models.FullCollection `json:"data"`
	Error  string                `json:"error"`
	Msg    string                `json:"msg"`
	Status int64                 `json:"status"`
}
