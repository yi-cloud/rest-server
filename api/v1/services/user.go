package services

import (
	"github.com/yi-cloud/rest-server/models"
	"github.com/yi-cloud/rest-server/pkg/db"
)

type UserService struct {
	userRepo *models.UserRepository
}

func NewUserService(userRepo *models.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetAllUsers() ([]models.User, error) {
	ret, err := s.userRepo.FindAll()
	if err == nil && ret == nil {
		ret = []models.User{}
	}
	return ret, err
}

func (s *UserService) GetUser(name, mobilePhone, nickName string) (*models.User, error) {
	user, err := s.userRepo.FindOne(mobilePhone)
	if db.IsRecordNotFound(err) {
		return s.userRepo.Create(name, mobilePhone, nickName, "")
	}
	return user, err
}

func (s *UserService) GetUserByMobile(mobilePhone string) (*models.User, error) {
	return s.userRepo.FindOne(mobilePhone)
}
