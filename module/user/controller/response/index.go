package user_reqres

type (
	StandartResponse struct {
		Message string `json:"message"`
	}
	LoginResponse struct {
		Message string `json:"message"`
		Token   string `json:"token"`
	}
)
