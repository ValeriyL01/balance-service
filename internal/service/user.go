package service

import (
	"errors"

	"github.com/ValeriyL01/balance-service/internal/database"
	"github.com/ValeriyL01/balance-service/internal/models"
	"github.com/ValeriyL01/balance-service/internal/utils"
)

type UserService struct {
	userDB *database.UserDB
}

func NewUserService(userDB *database.UserDB) *UserService {
	return &UserService{userDB: userDB}
}
func (s *UserService) Register(registerData models.RegisterRequest) error {

	existingUser, _ := s.userDB.GetUserByUsername(registerData.Username)
	if existingUser != nil {
		return errors.New("имя пользователя уже существует")
	}

	hashedPassword, err := utils.HashPassword(registerData.Password)
	if err != nil {
		return err
	}

	user := &models.User{
		Username:     registerData.Username,
		Email:        registerData.Email,
		PasswordHash: hashedPassword,
	}

	err = s.userDB.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}
func (s *UserService) Login(loginData models.LoginRequest) error {

	user, err := s.userDB.GetUserByUsername(loginData.Username)
	if err != nil {
		return errors.New("неверный юзернейм")
	}

	if !utils.CheckPasswordHash(loginData.Password, user.PasswordHash) {
		return errors.New("неверный пароль")
	}

	return nil
}
