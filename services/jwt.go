package services

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type Payload struct {
	Email string `json:"email"`
}

type Claims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

//CreateJWT execute creation of token jwt
func CreateJWT(email string) (string, error) {
	SECRET := []byte(os.Getenv("SECRET"))
	expire := time.Now().Add(60 * time.Minute)
	claims := Claims{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire.Unix(),
		},
	}
	algoritm := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := algoritm.SignedString(SECRET)
	if err != nil {
		return "", err
	}
	return token, nil
}

//VerifyJWT execute verification for token
func VerifyJWT(token string) (*Payload, error) {
	SECRET := os.Getenv("SECRET")
	claims := Claims{}
	strToken, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SECRET), nil
	})
	if strToken == nil || claims.Email == "" {
		return nil, errors.New("token bad formated")
	}
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("token expired")
	}
	if !strToken.Valid {
		return nil, errors.New("token invalid")
	}
	if err != nil {
		return nil, errors.New("token invalid")
	}

	return &Payload{Email: claims.Email}, nil
}
