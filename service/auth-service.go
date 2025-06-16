package service

import (
	"context"
	"errors"
	"fmt"
	"gcw/dto"
	"gcw/entity"
	"gcw/helper"
	"gcw/repository"
	"os"

	"google.golang.org/api/idtoken"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepository repository.UserRepository
	googleClientId string
}

func NewAuthService(ur repository.UserRepository) *AuthService {
	return &AuthService{
		userRepository: ur,
		googleClientId: os.Getenv("GOOGLE_CLIENT_ID"),
	}
}

func (s *AuthService) GetUserByGoogleIdToken(idToken string) (*entity.User, error) {
	user := &entity.User{}
	ctx := context.Background()

	payload, err := idtoken.Validate(ctx, idToken, s.googleClientId)
	if err != nil {
		return nil, err
	}

	email, ok := payload.Claims["email"].(string)
	if !ok {
		return nil, errors.New("email not found in google id token")
	}

	picture, ok := payload.Claims["picture"].(string)
	if !ok {
		return nil, errors.New("picture not found in google id token")
	}

	err = s.userRepository.FindByEmail(email, user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			user.Email = email
			user.ProfilePicture = picture

			name, ok := payload.Claims["name"].(string)
			if ok {
				user.Name = name
			} else {
				user.Name = email
			}

			err = s.userRepository.Create(user)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return user, nil
}

func (s *AuthService) GetUserById(id uint64) (*entity.User, error) {
	user := &entity.User{}
	err := s.userRepository.FindById(id, user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (s *AuthService) FindByEmail(email string) (*entity.User, error) {
	user := &entity.User{}
	err := s.userRepository.FindByEmail(email, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) Registration(data *dto.RegisterDTO) (*entity.User, error) {
	user := &entity.User{}
	err := s.userRepository.FindByEmail(data.Email, user)
	if err != nil && err.Error() == "record not found" {
		dataUser, err := s.registerAccount(data)
		if err != nil {
			return nil, err
		}
		return dataUser, nil
	} else if err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("user already registered")
}

func (s *AuthService) LoginService(data *dto.LoginDTO) (*entity.User, error) {
	user := &entity.User{}
	err := s.userRepository.FindByEmail(data.Email, user)
	if err != nil {
		return nil, fmt.Errorf("email or password is wrong")
	}

	ok := helper.CheckPasswordHash(data.Password, user.Password)
	if !ok {
		return nil, fmt.Errorf("email or password is wrong")
	}

	return user, nil
}

func (s *AuthService) registerAccount(data *dto.RegisterDTO) (*entity.User, error) {
	password, err := helper.HashPassword(data.Password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Email:    data.Email,
		Password: password,
		Name:     data.Name,
		Role:     "user", // Default role
	}

	if err := s.userRepository.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}
