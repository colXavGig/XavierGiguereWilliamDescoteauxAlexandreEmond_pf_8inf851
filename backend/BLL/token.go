package BLL

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claim struct {
	email string
	exp int64
	jwt.Claims
}

var secretKey = []byte("dafsedgfasgfaqerg")

func createToken(email string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, 
        Claim{
			email: email,
			exp: time.Now().Add(time.Hour * 1).Unix(),
		})

    tokenString, err := token.SignedString(secretKey)
    if err != nil {
    return "", err
    }

 return tokenString, nil
}

func verifyToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString,
		Claim{},
		func(token *jwt.Token) (interface{}, error) {
	   		return secretKey, nil
		})
   
	if err != nil {
	   return "", err
	}
   
	if !token.Valid {
	   return "", errors.New("invalid credentials")
	}
	claim, fine := token.Claims.(Claim)
	if !fine {
		return "", err
	}
   
	return claim.email, nil
 }