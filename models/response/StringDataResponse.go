package response

// StringDataResponse 用于返回简单文字信息
type StringDataResponse struct {
	Status int    `json:"status"`
	Data   string `json:"data"`
	Msg    string `json:"msg"`
	Error  string `json:"error"`
}
