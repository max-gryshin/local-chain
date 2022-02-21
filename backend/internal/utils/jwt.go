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
func GenerateToken(username, password string, id int) (string, error) {
	nowTime := time.Now()
	// Create token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = strconv.Itoa(id)
	claims["name"] = username
	claims["Issuer"] = "local-chain"
	claims["admin"] = true
	claims["exp"] = nowTime.Add(time.Hour * hoursInDay * daysInMonth).Unix() // 1 month by develop
	return token.SignedString(jwtSecret)
}

// ParseToken parsing token
func ParseToken(token string) (*Claims, error) {
	// FIXME: error not checked, try setup golangci-lint
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		// FIXME: tokenClaims.Valid should be checked firstly:
		// when tokenClaims is not null and is not *Claims, this method returns nil, nil
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
