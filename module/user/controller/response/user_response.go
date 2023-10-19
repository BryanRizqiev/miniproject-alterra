package user_response

type (
	StandartResponse struct {
		Message string `json:"message"`
	}
	LoginResponse struct {
		Message string `json:"message"`
		Token   string `json:"token"`
	}
)
