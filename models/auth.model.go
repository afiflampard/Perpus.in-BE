package models

import "github.com/dgrijalva/jwt-go"

type Claims struct {
	ID    uint64 `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}
