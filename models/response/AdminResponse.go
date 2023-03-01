package response

import "go_round4/models"

type AdminResponse struct {
	Data   models.Admin `json:"data"`
	Error  string       `json:"error"`
	Msg    string       `json:"msg"`
	Status int64        `json:"status"`
}
