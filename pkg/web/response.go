package web

type response struct {
	Data any `json:"data"`
}

type errResponse struct {
	Status  int    `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Response(data any) response {
	return response{
		Data: data,
	}
}

func ErrResponse(status int, code string, message string) errResponse {
	return errResponse{
		Status:  status,
		Code:    code,
		Message: message,
	}
}
