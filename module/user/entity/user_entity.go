package user_entity

type UserDTO struct {
	ID       string
	Email    string
	Name     string
	Password string
	Address  string
}

type (
	UserServiceInterface interface {
		Register(req UserDTO) error
		Login(req UserDTO) (UserDTO error)
		GetAllUser() ([]UserDTO, error)
	}

	UserRepositoryInterface interface {
		InsertUser(userDTO UserDTO) error
		CheckUser(userDTO UserDTO) (UserDTO error)
		GetAllUser() ([]UserDTO, error)
	}
)
