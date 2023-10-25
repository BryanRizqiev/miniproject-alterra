package user_reqres

type (
	UserPresentataion struct {
		Id              string `json:"id"`
		Name            string `json:"name"`
		Email           string `json:"email"`
		DOB             string `json:"dob"`
		Address         string `json:"address"`
		Phone           string `json:"phone"`
		Photo           string `json:"photo"`
		Role            string `json:"role"`
		VerifiedEmailAt string `json:"verified_email_at"`
		RequestVerified string `json:"request_verified"`
		CreatedAt       string `json:"created_at"`
	}
)

type (
	LoginResponse struct {
		Message string `json:"message"`
		Token   string `json:"token"`
	}

	GetRequestingUserRes struct {
		Message string              `json:"message"`
		Data    []UserPresentataion `json:"data"`
	}
)
