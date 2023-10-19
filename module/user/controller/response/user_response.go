package user_response

type (
	StandartResponse struct {
		Message string `json:"message"`
	}
	LoginResponse struct {
		Message string `json:"message"`
		Token   string `json:"token"`
	}
	VerifyResponse struct {
		Message string `json:"message"`
		UserID  string `json:"user_id"`
		Email   string `json:"email"`
	}
)
