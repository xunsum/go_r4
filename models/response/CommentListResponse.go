package response

import "go_round4/models"

type CommentListResponse struct {
	Data   []models.Comment `json:"data"`
	Error  string           `json:"error"`
	Msg    string           `json:"msg"`
	Status int64            `json:"status"`
}
