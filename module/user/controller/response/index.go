package user_reqres

type (
	UserPresentataion struct {
		Id        string
		Name      string
		Email     string
		DOB       string
		Address   string
		Phone     string
		Photo     string
		Role      string
		CreatedAt string
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
		Message string
		Data    []UserPresentataion
	}
)
