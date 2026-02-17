package jwt

import (
	"time"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(id int, username, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       id,
		"username": username,
		"exp":      time.Now().Add(5 * time.Minute).Unix(),
	})

	// generate token
	key := []byte(secretKey)
	tokenStr, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func ValidateToken(tokenStr, secretKey string) (int, string, error) {
	key := []byte(secretKey)
	claims := jwt.MapClaims{} //sebagai wadah pad saat unmarshall tokennya seperti id, username.

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
		return key, nil
	})
	if err != nil {
		return 0, "", err
	}

	// validasi jika token tidak valid
	if !token.Valid {
		return 0, "", constants.ErrInvalidToken
	}

	var id = int(claims["id"].(float64))
	var username = claims["username"].(string)

	return id, username, nil
}
