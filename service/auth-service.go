package service

import (
	"errors"
	"gcw/dto"
	"gcw/entity"
	"gcw/repository"
	"log"

	"github.com/mashingan/smapping"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authService struct {
	userRepository repository.UserRepository
}

type AuthService interface {
	Register(*dto.UserRequestDTO) (*entity.User, error)
	IsDuplicateUsername(string) bool
	FindByUsername(string) (*entity.User, error)
	VerifyPassword(string, string) bool
}

func NewAuthService(ur repository.UserRepository) AuthService {
	return &authService{
		userRepository: ur,
	}
}

func (s *authService) FindByUsername(username string) (*entity.User, error) {
	user := &entity.User{}
	err := s.userRepository.FindByUsername(username, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) IsDuplicateUsername(username string) bool {
	user := &entity.User{}
	err := s.userRepository.FindByUsername(username, user)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}

func (s *authService) Register(u *dto.UserRequestDTO) (*entity.User, error) {
	user := &entity.User{}
	if err := smapping.FillStruct(user, smapping.MapFields(u)); err != nil {
		return nil, err
	}
	user.Password = hashPassword([]byte(user.Password))
	if err := s.userRepository.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *authService) VerifyPassword(hashedPwd string, plainPassword string) bool {
	byteHash := []byte(hashedPwd)
	bytePwd := []byte(plainPassword)
	err := bcrypt.CompareHashAndPassword(byteHash, bytePwd)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func hashPassword(password []byte) string {
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		panic("Failed to hash a password")
	}
	return string(hash)
}
