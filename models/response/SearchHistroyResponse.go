package response

import (
	"go_round4/models"
)

type SearchHistoryListResponse struct {
	Status int                    `json:"status"`
	Data   []models.SearchHistory `json:"data"`
	Msg    string                 `json:"msg"`
	Error  string                 `json:"error"`
}

type SearchHistoryResponse struct {
	Status int                  `json:"status"`
	Data   models.SearchHistory `json:"data"`
	Msg    string               `json:"msg"`
	Error  string               `json:"error"`
}
