package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"math/rand"
	"time"
)

var secret = []byte(fmt.Sprintf("%x", rand.Int()))

func generateToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().AddDate(0, 0, 7).Unix(),
	})
	tokenStr, _ := token.SignedString(secret)
	return tokenStr
}

func validateToken(tokenStr string) bool {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return secret, nil
	})
	if err != nil {
		log.Println(err)
		return false
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if exp := int64(claims["exp"].(float64)); ok && exp > time.Now().Unix() {
			return true
		}
	}
	return false
}
