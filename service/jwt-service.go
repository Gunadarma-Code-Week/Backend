package service

import (
	"fmt"
	"gcw/dto"
	"gcw/entity"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/mashingan/smapping"
)

type jwtCustomClaim struct {
	dto.UserResponseDTO
	jwt.StandardClaims
}
type JwtService struct {
	secretKey  string
	refreshKey string
	issuer     string
}

func NewJwtService() *JwtService {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey == "" {
		secretKey = "12345"
	}

	refreshSecretKey := os.Getenv("JWT_REFRESH_SECRET")
	if refreshSecretKey == "" {
		refreshSecretKey = "12345refresh"
	}

	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		issuer = "gcw"
	}

	return &JwtService{
		secretKey:  secretKey,
		refreshKey: refreshSecretKey,
		issuer:     issuer,
	}
}

func (j *JwtService) GenerateToken(user *entity.User) string {
	userResponseDTO := dto.UserResponseDTO{}
	smapping.FillStruct(&userResponseDTO, smapping.MapFields(user))

	claims := &jwtCustomClaim{
		UserResponseDTO: userResponseDTO,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 3).Unix(),
			Issuer:    j.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		panic(err)
	}
	return signedToken
}

func (j *JwtService) GenerateRefreshToken(user *entity.User) string {
	userResponseDTO := dto.UserResponseDTO{}
	smapping.FillStruct(&userResponseDTO, smapping.MapFields(user))

	refreshClaims := &jwtCustomClaim{
		UserResponseDTO: userResponseDTO,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(0, 0, 3).Unix(),
			Issuer:    j.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedToken, err := token.SignedString([]byte(j.refreshKey))
	if err != nil {
		panic(err)
	}
	return signedToken
}

func (j *JwtService) validateToken(token string, isRefresh bool) (*jwt.Token, error) {
	return jwt.Parse(token, func(t_ *jwt.Token) (interface{}, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t_.Header["alg"])
		}
		var secretKey string
		if isRefresh {
			secretKey = j.refreshKey
		} else {
			secretKey = j.secretKey
		}
		return []byte(secretKey), nil
	})
}

func (j *JwtService) GetClaimsByToken(token string) (jwt.MapClaims, error) {
	aToken, err := j.validateToken(token, false)
	if err != nil {
		return nil, err
	}
	claims := aToken.Claims.(jwt.MapClaims)
	return claims, nil
}

func (j *JwtService) GetClaimsByRefreshToken(token string) (jwt.MapClaims, error) {
	aToken, err := j.validateToken(token, true)
	if err != nil {
		return nil, err
	}
	claims := aToken.Claims.(jwt.MapClaims)
	return claims, nil
}
