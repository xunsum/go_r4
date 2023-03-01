package response

import "go_round4/models"

type UserListResponse struct {
	Data   []models.User `json:"data"`
	Error  string        `json:"error"`
	Msg    string        `json:"msg"`
	Status int64         `json:"status"`
}

type UserResponse struct {
	Data   models.User `json:"data"`
	Error  string      `json:"error"`
	Msg    string      `json:"msg"`
	Status int64       `json:"status"`
}
