package service

import (
	"crypto/rsa"
	"fmt"
	"gcw/dto"
	"gcw/entity"
	"io/ioutil"
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
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	issuer     string
}

func NewJwtService() *JwtService {
	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		issuer = "gcw"
	}

	privateKeyPath := os.Getenv("JWT_PRIVATE_KEY_PATH")
	if privateKeyPath == "" {
		privateKeyPath = "keys/private.pem"
	}

	publicKeyPath := os.Getenv("JWT_PUBLIC_KEY_PATH")
	if publicKeyPath == "" {
		publicKeyPath = "keys/public.pem"
	}

	// Load Private Key
	privateKeyBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		fmt.Printf("Warning: Failed to load JWT private key from %s: %v\n", privateKeyPath, err)
		return nil
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		fmt.Printf("Warning: Failed to parse JWT private key: %v\n", err)
		return nil
	}

	// Load Public Key
	publicKeyBytes, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		fmt.Printf("Warning: Failed to load JWT public key from %s: %v\n", publicKeyPath, err)
		return nil
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyBytes)
	if err != nil {
		fmt.Printf("Warning: Failed to parse JWT public key: %v\n", err)
		return nil
	}

	return &JwtService{
		privateKey: privateKey,
		publicKey:  publicKey,
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
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(j.privateKey)
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

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaims)
	signedToken, err := token.SignedString(j.privateKey)
	if err != nil {
		panic(err)
	}
	return signedToken
}

func (j *JwtService) validateToken(token string, isRefresh bool) (*jwt.Token, error) {
	return jwt.Parse(token, func(t_ *jwt.Token) (interface{}, error) {
		if _, ok := t_.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", t_.Header["alg"])
		}
		return j.publicKey, nil
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
