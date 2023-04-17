package token

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = os.Getenv("jwt_key")

type JWTClaim struct {
	Name   string
	Email  string
	CodeID string
	jwt.StandardClaims
}

func GenerateToken(name string, email string, codeId string) (string, string, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claim := JWTClaim{
		Name:   name,
		Email:  email,
		CodeID: codeId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	refreshClaim := JWTClaim{
		Name:   name,
		Email:  email,
		CodeID: codeId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(jwtKey))
	if err != nil {
		return "", "", err
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaim).SignedString([]byte(jwtKey))
	if err != nil {
		return "", "", err
	}
	return token, refreshToken, nil

}

func ValidateToken(signedToken string) string {
	token, err := jwt.ParseWithClaims(signedToken, &JWTClaim{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return "unable to parse claims"
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return "couldnt parse claims"
	}
	if claims.ExpiresAt < time.Now().Unix() {
		return "expired token"
	}
	return ""
}

func SetTokenToExpired(signToken string) string {
	token, err := jwt.ParseWithClaims(signToken, &JWTClaim{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return "unable to parse claims"
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		return "unable to parse claims"
	}
	claims.ExpiresAt = time.Now().Unix()
	return ""
}
