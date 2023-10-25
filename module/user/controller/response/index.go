package user_reqres

type (
	UserPresentataion struct {
		Id        string `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		DOB       string `json:"message"`
		Address   string `json:"address"`
		Phone     string `json:"phone"`
		Photo     string `json:"photo"`
		Role      string `json:"role"`
		CreatedAt string `json:"created_at"`
	}
)

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

	GetRequestingUserRes struct {
		Message string              `json:"message"`
		Data    []UserPresentataion `json:"data"`
	}
)
