package service

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type jwtCustomClaim struct {
	Username string `json:"username"`
	jwt.StandardClaims
}
type jwtService struct {
	secretKey string
	issuer    string
}

type JwtService interface {
	GenerateToken(string) string
	GetClaimsByToken(string) (jwt.MapClaims, error)
}

func NewJwtService() JwtService {
	return &jwtService{
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

func (j *jwtService) GenerateToken(Username string) string {
	claims := &jwtCustomClaim{
		Username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(1, 0, 0).Unix(),
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

func (j *jwtService) validateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t_ *jwt.Token) (interface{}, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t_.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
}

func (j *jwtService) GetClaimsByToken(token string) (jwt.MapClaims, error) {
	aToken, err := j.validateToken(token)
	if err != nil {
		return nil, err
	}
	claims := aToken.Claims.(jwt.MapClaims)
	return claims, nil
}
