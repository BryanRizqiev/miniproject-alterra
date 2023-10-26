package global_response

type StandartResponse struct {
	Message string `json:"message"`
}

type StandartResponseWithData struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}
