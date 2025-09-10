package user

import "github.com/andredubov/rocket-factory/iam/internal/service"

type usersService struct {
	usersRepository service.UsersRepository
	passwordHasher  service.PasswordHasher
}

func NewUsersService(
	usersRepository service.UsersRepository,
	passwordHasher service.PasswordHasher,
) *usersService {
	return &usersService{
		usersRepository: usersRepository,
		passwordHasher:  passwordHasher,
	}
}
