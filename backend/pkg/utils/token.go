package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

type Claims struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	UserType string `json:"userType"`
	jwt.StandardClaims
}

func CreateToken(id string, username string, userType string) (tokenString string, err error) {
	claims := &Claims{
		Id:       id,
		Username: username,
		UserType: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 20).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET"))); err != nil {
		return "", err
	} else {
		return signedToken, nil
	}
}

func VerifyToken(token string) (string, string, string, error) {
	if token == "" {
		return "", "", "", errors.New("token is empty")
	}
	claims := &Claims{}
	parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", "", "", errors.New("signature is invalid")
		}
		return "", "", "", errors.New("token is invalid")
	}
	if !parsedToken.Valid {
		return "", "", "", errors.New("parsed token is invalid")
	}
	if claims == nil {
		return "", "", "", errors.New("token claims are nil")
	}
	return claims.Id, claims.Username, claims.UserType, nil
}
