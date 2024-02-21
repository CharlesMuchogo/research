package models

import "github.com/dgrijalva/jwt-go"

type JWTClaim struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Phone     string `json:"phone"`
	jwt.StandardClaims
}
