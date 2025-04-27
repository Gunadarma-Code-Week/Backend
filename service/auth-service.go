package service

import (
	"context"
	"errors"
	"gcw/entity"
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
