package user_request

type RegisterRequest struct {
	Name     string `json:"name" form:"name" validate:"required,min=3"`
	Email    string `json:"email" form:"email" validate:"required,email,min=3"`
	Password string `json:"password" form:"password" validate:"required,min=3"`
	Address  string `json:"address" form:"address"`
	DOB      string `json:"dob" form:"dob"`
	Phone    string `json:"phone" form:"phone"`
}

type UpdateRequest struct {
	Name    string `json:"name" form:"name" validate:"required,min=3"`
	DOB     string `json:"dob" form:"dob"`
	Address string `json:"address" form:"address"`
	Phone   string `json:"phone" form:"phone" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email,min=3"`
	Password string `json:"password" form:"password" validate:"required,min=3"`
}

type ApproveVerificationReq struct {
	Role string `json:"role" form:"role" validate:"required,min=4,oneof=user verified_user"`
}
