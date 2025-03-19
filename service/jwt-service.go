package service

import (
	"fmt"
	"gcw/entity"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type jwtCustomClaim struct {
	UserId uint64 `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`

	jwt.StandardClaims
}
type JwtService struct {
	secretKey string
	issuer    string
}

func NewJwtService() *JwtService {
	return &JwtService{
		secretKey: getSecretKey(),
		issuer:    "url-shortener",
	}
}

func getSecretKey() string {
	secretKey := os.Getenv("JWT_SECRET")
	if secretKey != "" {
		secretKey = "12345"
	}
	return secretKey
}

func (j *JwtService) GenerateToken(user *entity.User) string {
	claims := &jwtCustomClaim{
		Email:  user.Email,
		Role:   *user.Role,
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 3).Unix(),
			Issuer:    j.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (j *JwtService) GenerateRefreshToken(user *entity.User) string {
	claims := &jwtCustomClaim{
		Email:  user.Email,
		Role:   *user.Role,
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 3).Unix(),
			Issuer:    j.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (j *JwtService) validateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t_ *jwt.Token) (interface{}, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t_.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
}

func (j *JwtService) GetClaimsByToken(token string) (jwt.MapClaims, error) {
	aToken, err := j.validateToken(token)
	if err != nil {
		return nil, err
	}
	claims := aToken.Claims.(jwt.MapClaims)
	return claims, nil
}
