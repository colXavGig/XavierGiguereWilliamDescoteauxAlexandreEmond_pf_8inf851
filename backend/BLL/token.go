package BLL

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claim struct {
	Email string

	jwt.StandardClaims
}

var secretKey = []byte("dafsedgfasgfaqerg")

func createToken(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		Claim{
			Email: email,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(), // Set token expiry to 1 hour
			},
		})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func verifyToken(tokenString string) (string, error) {

	token, err := jwt.ParseWithClaims(tokenString,
		&Claim{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid credentials")
	}

	claim, fine := token.Claims.(*Claim)
	if !fine {
		return "", err
	}

	return claim.Email, nil
}
