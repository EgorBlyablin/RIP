package services

import (
	"errors"
	"rip/internal/app/ds"
	"rip/internal/app/repositories"

	"gorm.io/gorm"
)

var (
	ErrorUserHasNoAccess = errors.New("user has no access")
)

type UsersService struct {
	usersRepository *repositories.UsersRepository
}

func NewUsersService(turbinesAppDB *gorm.DB) *UsersService {
	return &UsersService{
		usersRepository: repositories.NewUsersRepository(turbinesAppDB),
	}
}

func (service *UsersService) RegisterUser(user ds.CreateUser) (ds.User, error) {
	return service.usersRepository.CreateUser(user)
}

func (service *UsersService) GetUser(userId uint) (ds.User, error) {
	return service.usersRepository.GetUserByID(uint(userId))
}

func (service *UsersService) UpdateUser(userId uint, user ds.UpdateUser) (ds.User, error) {
	return service.usersRepository.UpdateUser(userId, user)
}

func (service *UsersService) CheckIsModerator(userId uint) (bool, error) {
	user, err := service.usersRepository.GetUserByID(userId)
	if err != nil {
		return false, err
	}
	return user.IsModerator, nil
}
