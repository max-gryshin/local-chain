package utils

import (
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("get_key_from_env")

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

const (
	hoursInDay  = 24
	daysInMonth = 30
)

// GenerateToken generate tokens used for auth
func GenerateToken(id int) (string, error) {
	nowTime := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		Audience:  "all",
		ExpiresAt: nowTime.Add(time.Hour * hoursInDay * daysInMonth).Unix(), // 1 month by develop
		Id:        strconv.Itoa(id),
		IssuedAt:  0,
		Issuer:    "local-chain",
		NotBefore: 0,
		Subject:   "",
	})

	return token.SignedString(jwtSecret)
}

// ParseToken parsing token
func ParseToken(token string) (jwt.StandardClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil && tokenClaims.Valid {
		if claims, ok := tokenClaims.Claims.(*jwt.StandardClaims); ok && tokenClaims.Valid {
			return *claims, nil
		}
	}

	return jwt.StandardClaims{}, err
}

func GetAuthenticatedUserID(token string) (string, error) {
	claims, err := ParseToken(token)
	if err != nil {
		return "", err
	}
	return claims.Id, nil
}
