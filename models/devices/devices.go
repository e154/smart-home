package devices

type BaseResponse struct {
	Error     string  `json:"error"`
	ErrorCode string  `json:"error_code"`
	Time      float64 `json:"time"`
}
